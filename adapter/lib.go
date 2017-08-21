package adapter

import (
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func validateApiKey(key string) error {
	if key != os.Getenv("ADAPTER_API_KEY") {
		err := grpc.Errorf(codes.PermissionDenied, "ADAPTER_API_KEY is invalid")
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Start request warn")
		return err
	} else {
		return nil
	}
}
