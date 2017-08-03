package scan

import (
	"fmt"
	"time"
)

func scan(worker *Worker, sync chan StatusResponse, resp StatusResponse) {
	// This is just an example which simulates finding one device per second
	for i := 0; i < 100; i++ {
		select {
		case <-worker.ctx.Done():
			return
		case <-time.After(time.Second * 1):
			result := &StatusResponse_Result{
				Address: fmt.Sprintf("address: %d", i),
				Name:    fmt.Sprintf("name: %d", i),
			}
			resp.Results = append(resp.Results, result)
			sync <- resp
		}
	}
}
