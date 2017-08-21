package cmd

import (
	"context"

	"github.com/resin-io/adapter-base/adapter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func statusCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	parent := cmd.Parent().Name()

	var resp *adapter.Jobs
	if parent == "scan" {
		resp, err = adapter.NewScanClient(conn).Status(
			context.Background(),
			&adapter.Id{Id: id},
		)
	} else if parent == "update" {
		resp, err = adapter.NewUpdateClient(conn).Status(
			context.Background(),
			&adapter.Id{Id: id},
		)
	} else {
		log.WithFields(log.Fields{
			"parent command": parent,
		}).Fatal("Command not found")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get job status")
	} else {
		log.WithFields(log.Fields{
			"response": resp,
		}).Info("Job status")
	}
}

func cancelCmd(cmd *cobra.Command, args []string) {
	conn, err := openConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	parent := cmd.Parent().Name()

	var resp *adapter.Job
	if parent == "scan" {
		resp, err = adapter.NewScanClient(conn).Cancel(
			context.Background(),
			&adapter.Id{Id: args[0]},
		)
	} else if parent == "update" {
		resp, err = adapter.NewUpdateClient(conn).Cancel(
			context.Background(),
			&adapter.Id{Id: args[0]},
		)
	} else {
		log.WithFields(log.Fields{
			"parent command": parent,
		}).Fatal("Command not found")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to cancel job status")
	} else {
		log.WithFields(log.Fields{
			"response": resp,
		}).Info("Job status")
	}
}
