package cmd

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var verbose bool
var apiPort, rpcPort, port, concurrency int
var timeout int64
var name string

func validateArgs(args []string, required int) error {
	if len(args) != required {
		return errors.New("Incorrect arguments")
	}
	return nil
}

func openConnection() (*grpc.ClientConn, error) {
	serverAddr := fmt.Sprintf("localhost:%v", port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"Address": serverAddr,
		}).Fatal("Failed to open connection")
		return nil, err
	}

	return conn, nil
}
