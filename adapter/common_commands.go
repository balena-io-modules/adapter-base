package adapter

import (
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *Server) Status(ctx context.Context, req *Id) (*Jobs, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Status request")

	if err := validateApiKey(req.AdapterApiKey); err != nil {
		return nil, err
	}

	resp := &Jobs{}
	if req.Id == "" {
		if len(s.workers) == 0 {
			err := grpc.Errorf(codes.NotFound, "no workers found")
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("Status request warn")
			return nil, err
		}

		for _, worker := range s.workers {
			worker.input <- req
			r := <-worker.output
			resp.Jobs = append(resp.Jobs, &r)
		}
	} else {
		if worker, err := s.findWorker(req); err != nil {
			return nil, err
		} else {
			worker.input <- req
			r := <-worker.output
			s.cleanupWorker(req, worker)
			resp.Jobs = append(resp.Jobs, &r)
		}
	}

	return resp, nil
}

func (s *Server) Cancel(ctx context.Context, req *Id) (*Job, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Cancel request")

	if err := validateApiKey(req.AdapterApiKey); err != nil {
		return nil, err
	}

	worker, err := s.findWorker(req)
	if err != nil {
		return nil, err
	}

	worker.cancel()

	worker.input <- req
	resp := <-worker.output

	s.cleanupWorker(req, worker)

	return &resp, nil
}
