package update

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

	if req.Address == "" || req.Payload == "" || req.Timeout == 0 {
		err := grpc.Errorf(codes.InvalidArgument, "must specify address, payload and timeout")
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

		// Start of updating code
		// This is just an example which simulates increasing the progress percentage by one per second
		for i := 0; i < 100; i++ {
			select {
			case <-worker.ctx.Done():
				return
			case <-time.After(time.Second * 1):
				resp.State = StatusResponse_FLASHING
				resp.Progress = int32(i)
				resp.Message = fmt.Sprintf("message: %d", i)
				sync <- resp
			}
		}
		// End of updating code
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
