package cmd

import (
	"context"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/resin-io/adapter-base/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Interact with the adapter-base update endpoint",
}

var StartUpdateCmd = &cobra.Command{
	Use:   "start [image]",
	Short: "Start an update",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, 1); err != nil {
			return err
		} else if err := validateFlags(cmd.Flags()); err != nil {
			return err
		} else {
			return nil
		}
	},
	Run: startUpdateCmd,
}

var StatusUpdateCmd = &cobra.Command{
	Use:   "status",
	Short: "Get update status",
	Run:   statusUpdateCmd,
}

var CancelUpdateCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Cancel an update",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: cancelUpdateCmd,
}

func init() {
	UpdateCmd.AddCommand(StartUpdateCmd)
	UpdateCmd.AddCommand(StatusUpdateCmd)
	UpdateCmd.AddCommand(CancelUpdateCmd)

	UpdateCmd.PersistentFlags().IntVarP(&port, "port", "p", 8081, "API port")
	StartUpdateCmd.Flags().StringArrayVarP(&destinations, "destinations", "d", nil, "update destinations")
	StartUpdateCmd.MarkFlagRequired("destinations")
	StartUpdateCmd.Flags().Int64VarP(&timeout, "timeout", "t", 120000, "millisecond timeout")
	StatusUpdateCmd.Flags().StringVarP(&id, "id", "i", "", "job id")
}

func startUpdateCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	options := &update.Options{
		Image: args[0],
		Extra: make(map[string]*structpb.Value),
	}
	for _, entry := range destinations {
		destination := &update.Destination{
			Id: entry,
		}
		options.Destinations = append(options.Destinations, destination)
	}
	options.Extra["timeout"] = &structpb.Value{
		Kind: &structpb.Value_NumberValue{NumberValue: float64(timeout)},
	}

	client := update.NewUpdateClient(conn)
	resp, err := client.Start(context.Background(), options)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to start job")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Job status")
}

func statusUpdateCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := update.NewUpdateClient(conn)
	resp, err := client.Status(context.Background(), &update.Id{Id: id})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get job status")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Job status")
}

func cancelUpdateCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := update.NewUpdateClient(conn)
	resp, err := client.Cancel(context.Background(), &update.Id{Id: args[0]})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to cancel job")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Job status")
}
