package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "adapter-base",
	Short: "API for interacting with adapter-base",
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(UpdateCmd)
	RootCmd.AddCommand(ScanCmd)
}
