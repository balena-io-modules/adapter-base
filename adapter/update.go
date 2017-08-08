package adapter

import (
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

func (s *Server) StartUpdate(ctx context.Context, req *UpdateOptions) (*Id, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Start request")

	if id, worker, err := s.createWorker(req.Extra); err != nil {
		return nil, err
	} else {
		go func(req *UpdateOptions, id *Id, worker *Worker) {
			defer worker.cancel()
			sync, resp := s.syncWorker(id, worker)
			update(req, worker, sync, resp)
		}(req, id, worker)

		return id, nil
	}
}
