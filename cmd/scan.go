package cmd

import (
	"context"

	"github.com/resin-io/adapter-base/scan"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Interact with the adapter-base scan endpoint",
}

var StartScanCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a scan",
	Run:   startScanCmd,
}

var StatusScanCmd = &cobra.Command{
	Use:   "status [id]",
	Short: "Get the status of a scan",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: statusScanCmd,
}

var CancelScanCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Cancel a scan",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: cancelScanCmd,
}

func init() {
	ScanCmd.AddCommand(StartScanCmd)
	ScanCmd.AddCommand(StatusScanCmd)
	ScanCmd.AddCommand(CancelScanCmd)

	ScanCmd.PersistentFlags().IntVarP(&port, "port", "p", 8081, "API port")
	StartScanCmd.Flags().StringVarP(&address, "address", "a", "", "address")
	StartScanCmd.Flags().Int64VarP(&timeout, "timeout", "t", 120, "timeout")
}

func startScanCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := scan.NewScanClient(conn)
	resp, err := client.Start(context.Background(), &scan.StartRequest{Address: address, Timeout: timeout})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to start scan")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Scan status")
}

func statusScanCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := scan.NewScanClient(conn)
	resp, err := client.Status(context.Background(), &scan.StatusRequest{Id: args[0]})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get scan status")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Scan status")
}

func cancelScanCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	client := scan.NewScanClient(conn)
	resp, err := client.Cancel(context.Background(), &scan.StatusRequest{Id: args[0]})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to cancel scan")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Scan status")
}
