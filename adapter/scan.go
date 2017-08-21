package adapter

import (
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *Server) StartScan(ctx context.Context, req *ScanKey) (*Id, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Start request")

	if mode, err := validateScan(req); err != nil {
		return nil, err
	} else if id, worker, err := s.createWorker(req.Options.Extra); err != nil {
		return nil, err
	} else {
		go func(mode scanMode, req *ScanOptions, id *Id, worker *Worker) {
			defer worker.cancel()
			sync, resp := s.syncWorker(id, worker)
			mode(req, worker, sync, resp)
		}(mode, req.Options, id, worker)

		return id, nil
	}
}

func validateScan(req *ScanKey) (scanMode, error) {
	var err error
	if req == nil {
		err = grpc.Errorf(codes.InvalidArgument, "request not found")
	} else if err = validateApiKey(req.AdapterApiKey); err != nil {
		return nil, err
	} else if req.Options == nil {
		err = grpc.Errorf(codes.InvalidArgument, "request options not found")
	} else if req.Options.Mode == "" {
		err = grpc.Errorf(codes.InvalidArgument, "request mode not found")
	} else if mode, ok := scanModes[req.Options.Mode]; !ok {
		err = grpc.Errorf(codes.Unimplemented, "request mode not implemented")
	} else {
		return mode, nil
	}

	log.WithFields(log.Fields{
		"error": err,
	}).Warn("Start request warn")

	return nil, err
}
