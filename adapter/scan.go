package adapter

import (
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

func (s *Server) StartScan(ctx context.Context, req *ScanOptions) (*Id, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Start request")

	if id, worker, err := s.createWorker(req.Extra); err != nil {
		return nil, err
	} else {
		go func(req *ScanOptions, id *Id, worker *Worker) {
			defer worker.cancel()
			sync, resp := s.syncWorker(id, worker)
			scan(req, worker, sync, resp)
		}(req, id, worker)

		return id, nil
	}
}
