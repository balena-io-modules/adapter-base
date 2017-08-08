package adapter

import (
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

func (s *Server) Status(ctx context.Context, req *Id) (*Jobs, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Status request")

	resp := &Jobs{}
	if req.Id == "" {
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
