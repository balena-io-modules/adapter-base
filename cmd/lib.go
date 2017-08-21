package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var verbose bool
var apiPort, rpcPort, port, concurrency int
var number, delay, timeout int64
var application, mac, id string
var destinations []string

func validateArgs(args []string, n int) error {
	if len(args) != n {
		return fmt.Errorf("accepts exactly %d arg(s), received %d", n, len(args))
	}

	return nil
}

func validateFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	})

	if requiredError {
		return fmt.Errorf("required flag \"%s\" has not been set", flagName)
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
