package adapter

import (
	"time"

	structpb "github.com/golang/protobuf/ptypes/struct"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Worker struct {
	input  chan *Id
	output chan Job
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *Server) createWorker(extra map[string]*structpb.Value) (*Id, *Worker, error) {
	if len(s.workers) >= s.concurrency {
		err := grpc.Errorf(codes.ResourceExhausted, "no workers available")
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Start request warn")
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	if value, ok := extra["timeout"]; ok {
		if timeout := value.GetNumberValue(); timeout != 0 {
			ctx, cancel = context.WithTimeout(ctx, time.Duration(int64(timeout))*time.Second)
		}
	}

	id := uuid.NewV4().String()
	worker := &Worker{
		input:  make(chan *Id),
		output: make(chan Job),
		ctx:    ctx,
		cancel: cancel,
	}
	s.workers[id] = worker

	return &Id{Id: id}, worker, nil
}

func (s *Server) syncWorker(id *Id, worker *Worker) (chan Job, Job) {
	sync := make(chan Job)

	resp := Job{
		Id:    id.Id,
		State: State_STARTED,
		Progress: &Progress{
			Started: time.Now().UTC().Unix(),
		},
	}

	go func(worker *Worker, sync chan Job, resp Job) {
		for {
			select {
			case resp = <-sync:
			case <-worker.input:
				if resp.Progress.Completed != 0 {
					resp.Progress.Duration = resp.Progress.Completed - resp.Progress.Started
				} else {
					resp.Progress.Duration = time.Now().UTC().Unix() - resp.Progress.Started
				}
				worker.output <- resp
				select {
				case <-worker.ctx.Done():
					return
				default:
				}
			}
		}
	}(worker, sync, resp)

	return sync, resp
}

func (s *Server) findWorker(req *Id) (*Worker, error) {
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

func (s *Server) cleanupWorker(req *Id, worker *Worker) {
	select {
	case <-worker.ctx.Done():
		delete(s.workers, req.Id)
	default:
	}
}
