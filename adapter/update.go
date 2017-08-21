package adapter

import (
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *Server) StartUpdate(ctx context.Context, req *UpdateKey) (*Id, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Debug("Start request")

	if mode, err := validateUpdate(req); err != nil {
		return nil, err
	} else if id, worker, err := s.createWorker(req.Options.Extra); err != nil {
		return nil, err
	} else {
		go func(mode updateMode, req *UpdateOptions, id *Id, worker *Worker) {
			defer worker.cancel()
			sync, resp := s.syncWorker(id, worker)
			mode(req, worker, sync, resp)
		}(mode, req.Options, id, worker)

		return id, nil
	}
}

func validateUpdate(req *UpdateKey) (updateMode, error) {
	var err error
	if req == nil {
		err = grpc.Errorf(codes.InvalidArgument, "request not found")
	} else if err = validateApiKey(req.AdapterApiKey); err != nil {
		return nil, err
	} else if req.Options == nil {
		err = grpc.Errorf(codes.InvalidArgument, "request options not found")
	} else if req.Options.Mode == "" {
		err = grpc.Errorf(codes.InvalidArgument, "request mode not found")
	} else if mode, ok := updateModes[req.Options.Mode]; !ok {
		err = grpc.Errorf(codes.Unimplemented, "request mode not implemented")
	} else if req.Options.Image == "" {
		err = grpc.Errorf(codes.Unimplemented, "request image not found")
	} else if req.Options.Destinations == nil {
		err = grpc.Errorf(codes.InvalidArgument, "update destinations not found")
	} else if len(req.Options.Destinations) == 0 {
		err = grpc.Errorf(codes.InvalidArgument, "at least one update destination required")
	} else {
		return mode, nil
	}

	log.WithFields(log.Fields{
		"error": err,
	}).Warn("Start request warn")

	return nil, err
}
