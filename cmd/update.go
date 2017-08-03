package cmd

import (
	"context"

	"github.com/resin-io/adapter-base/update"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Interact with the adapter-base update endpoint",
}

var StartUpdateCmd = &cobra.Command{
	Use:   "start [address] [payload]",
	Short: "Start an update",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 2)
	},
	Run: startUpdateCmd,
}

var StatusUpdateCmd = &cobra.Command{
	Use:   "status [id]",
	Short: "Get the status of an update",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: statusUpdateCmd,
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
	StartUpdateCmd.Flags().Int64VarP(&timeout, "timeout", "t", 120, "timeout")
}

func startUpdateCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := update.NewUpdateClient(conn)
	resp, err := client.Start(context.Background(), &update.StartRequest{Address: args[0], Payload: args[1], Timeout: timeout})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to start update")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Update status")
}

func statusUpdateCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := update.NewUpdateClient(conn)
	resp, err := client.Status(context.Background(), &update.StatusRequest{Id: args[0]})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get update status")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Update status")
}

func cancelUpdateCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := update.NewUpdateClient(conn)
	resp, err := client.Cancel(context.Background(), &update.StatusRequest{Id: args[0]})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to cancel update")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Update status")
}
