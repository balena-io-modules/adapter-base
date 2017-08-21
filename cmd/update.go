package cmd

import (
	"context"
	"os"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/resin-io/adapter-base/adapter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Interact with the adapter-base update endpoint",
}

var StartUpdateCmd = &cobra.Command{
	Use:   "start [mode] [image]",
	Short: "Start an update",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, 2); err != nil {
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
	Run:   statusCmd,
}

var CancelUpdateCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Cancel an update",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: cancelCmd,
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

	key := &adapter.UpdateKey{
		AdapterApiKey: os.Getenv("ADAPTER_API_KEY"),
		Options: &adapter.UpdateOptions{
			Mode:  args[0],
			Image: args[1],
			Extra: make(map[string]*structpb.Value),
		},
	}
	for _, entry := range destinations {
		destination := &adapter.Destination{
			Id: entry,
		}
		key.Options.Destinations = append(key.Options.Destinations, destination)
	}
	key.Options.Extra["timeout"] = &structpb.Value{
		Kind: &structpb.Value_NumberValue{NumberValue: float64(timeout)},
	}

	client := adapter.NewUpdateClient(conn)
	resp, err := client.StartUpdate(context.Background(), key)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to start job")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Job status")
}
