package update

import (
	"fmt"
	"time"

	"github.com/currantlabs/ble"
	"github.com/resin-io/adapter-base/bluetooth"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	nRF51822DK = "nRF51822-DK"
	microbit   = "micro:bit"
)

type Worker struct {
	input  chan *StatusRequest
	output chan StatusResponse
	ctx    context.Context
	cancel context.CancelFunc
}

type Server struct {
	concurrency int
	workers     map[string]*Worker
}

func NewServer(concurrency int, verbose bool) *Server {
	server := &Server{
		concurrency: concurrency,
		workers:     make(map[string]*Worker),
	}

	return server
}

func (s *Server) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Start request received")

	if req.Address == "" || req.Payload == "" || req.Device == "" || req.Timeout == 0 {
		err := grpc.Errorf(codes.InvalidArgument, "must specify address, payload, device and timeout")
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Start request error")
		return nil, err
	} else if req.Device != nRF51822DK && req.Device != microbit {
		err := grpc.Errorf(codes.InvalidArgument, "device must be one of: %s, %s", nRF51822DK, microbit)
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Start request error")
		return nil, err
	}

	if len(s.workers) >= s.concurrency {
		err := grpc.Errorf(codes.ResourceExhausted, "no workers available")
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Start request warn")
		return nil, err
	}

	id := uuid.NewV4().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.Timeout)*time.Second)
	worker := &Worker{
		input:  make(chan *StatusRequest),
		output: make(chan StatusResponse),
		ctx:    ctx,
		cancel: cancel,
	}
	s.workers[id] = worker

	s.update(req, id, worker)

	return &StartResponse{Id: id}, nil
}

func (s *Server) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Status request")

	worker, err := s.findWorker(req)
	if err != nil {
		return nil, err
	}

	worker.input <- req
	resp := <-worker.output

	s.cleanup(worker, req)

	return &resp, nil
}

func (s *Server) Cancel(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Cancel request")

	worker, err := s.findWorker(req)
	if err != nil {
		return nil, err
	}

	worker.cancel()

	worker.input <- req
	resp := <-worker.output

	s.cleanup(worker, req)

	return &resp, nil
}

//TODO pipe errors back via message - can we do this with a nice handler
// pass around sync, response and get back resp
// check ctx somehow - maybe something like u cannot cancel until the op has finished
func (s *Server) update(req *StartRequest, id string, worker *Worker) {
	go func(req *StartRequest, id string, worker *Worker) {
		defer worker.cancel()

		resp := StatusResponse{
			Id:           id,
			StartRequest: req,
			State:        StatusResponse_STARTED,
			Started:      time.Now().UTC().Unix(),
		}

		sync := make(chan StatusResponse)

		go func(worker *Worker, sync chan StatusResponse, resp StatusResponse) {
			for {
				select {
				case resp = <-sync:
				case <-worker.input:
					worker.output <- resp
					select {
					case <-worker.ctx.Done():
						return
					default:
					}
				default:
				}
			}
		}(worker, sync, resp)

		// TODO check ctx done regularly
		if err := bluetooth.OpenDevice(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error opening device")
			return
		}
		defer bluetooth.CloseDevice()

		micro := &Nrf51822{
			LocalUUID:           resp.StartRequest.Address,
			Firmware:            Firmware{},
			NotificationChannel: make(chan []byte),
		}

		if resp, err := startBootloader(sync, resp); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error starting bootloader")
			return
		} else if resp, err := micro.ExtractPayload(sync, resp); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error extracting payload")
			return
		} else if _, err := micro.Update(sync, resp); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error updating")
			return
		}
	}(req, id, worker)
}

func (s *Server) findWorker(req *StatusRequest) (*Worker, error) {
	worker, ok := s.workers[req.Id]
	if !ok {
		err := grpc.Errorf(codes.NotFound, "worker not found")
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Status request warn")
		return nil, err
	}

	return worker, nil
}

func (s *Server) cleanup(worker *Worker, req *StatusRequest) {
	select {
	case <-worker.ctx.Done():
		delete(s.workers, req.Id)
	default:
	}
}

func startBootloader(sync chan StatusResponse, resp StatusResponse) (StatusResponse, error) {
	if name, err := bluetooth.GetName(resp.StartRequest.Address, 10); err != nil {
		return resp, err
	} else if name != "DfuTarg" {
		resp.Message = "Starting bootloader"
		sync <- resp

		if client, err := bluetooth.Connect(resp.StartRequest.Address, 10); err != nil {
			return resp, err
		} else {
			switch resp.StartRequest.Device {
			case nRF51822DK:
				if dfu, err := bluetooth.GetCharacteristic("000015311212efde1523785feabcd123", ble.CharWrite+ble.CharNotify, 0x0F, 0x10); err != nil {
					return resp, err
				} else if dfu.CCCD, err = bluetooth.GetDescriptor("2902", 0x11); err != nil {
					return resp, err
				} else if err := bluetooth.WriteDescriptor(client, dfu.CCCD, []byte{0x001}, 1); err != nil {
					return resp, err
				} else {
					bluetooth.WriteCharacteristic(client, dfu, []byte{Start, 0x04}, false, 1)
				}
			case microbit:
				if dfu, err := bluetooth.GetCharacteristic("e95d93b1251d470aa062fa1922dfa9a8", ble.CharRead+ble.CharWrite, 0x0D, 0x0E); err != nil {
					return resp, err
				} else {
					bluetooth.WriteCharacteristic(client, dfu, []byte{Start}, false, 1)
				}
			default:
				return resp, fmt.Errorf("Device not supported")
			}

			time.Sleep(time.Duration(1) * time.Second)

			resp.Message = "Started bootloader"
			sync <- resp
		}
	} else {
		resp.Message = "Bootloader already started"
		sync <- resp
	}

	return resp, nil
}
