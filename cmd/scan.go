package cmd

import (
	"context"
	"os"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/resin-io/adapter-base/adapter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Interact with the adapter-base scan endpoint",
}

var StartScanCmd = &cobra.Command{
	Use:   "start [mode]",
	Short: "Start a scan",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: startScanCmd,
}

var StatusScanCmd = &cobra.Command{
	Use:   "status",
	Short: "Get scan status",
	Run:   statusCmd,
}

var CancelScanCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Cancel a scan",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args, 1)
	},
	Run: cancelCmd,
}

func init() {
	ScanCmd.AddCommand(StartScanCmd)
	ScanCmd.AddCommand(StatusScanCmd)
	ScanCmd.AddCommand(CancelScanCmd)

	ScanCmd.PersistentFlags().IntVarP(&port, "port", "p", 8081, "API port")
	StartScanCmd.Flags().Int64VarP(&number, "number", "n", 0, "number of scans to run (default ∞)")
	StartScanCmd.Flags().Int64VarP(&delay, "delay", "d", 0, "millisecond pause between scans (default 0)")
	StartScanCmd.Flags().StringVarP(&application, "application", "a", "", "application")
	StartScanCmd.Flags().StringVarP(&mac, "mac", "m", "", "MAC address")
	StartScanCmd.Flags().Int64VarP(&timeout, "timeout", "t", 120000, "millisecond timeout")
	StatusScanCmd.Flags().StringVarP(&id, "id", "i", "", "job id")
}

func startScanCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	key := &adapter.ScanKey{
		AdapterApiKey: os.Getenv("ADAPTER_API_KEY"),
		Options: &adapter.ScanOptions{
			Mode:   args[0],
			Number: number,
			Delay:  delay,
			Extra:  make(map[string]*structpb.Value),
		},
	}
	key.Options.Extra["application"] = &structpb.Value{
		Kind: &structpb.Value_StringValue{StringValue: application},
	}
	key.Options.Extra["mac"] = &structpb.Value{
		Kind: &structpb.Value_StringValue{StringValue: mac},
	}
	key.Options.Extra["timeout"] = &structpb.Value{
		Kind: &structpb.Value_NumberValue{NumberValue: float64(timeout)},
	}

	client := adapter.NewScanClient(conn)
	resp, err := client.StartScan(context.Background(), key)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to start job")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Job status")
}
