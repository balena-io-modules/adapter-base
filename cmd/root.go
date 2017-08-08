package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:              "adapter-base",
	Short:            "API for interacting with adapter-base",
	PersistentPreRun: setVerbosityCmd,
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(UpdateCmd)
	RootCmd.AddCommand(ScanCmd)
}

func setVerbosityCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
