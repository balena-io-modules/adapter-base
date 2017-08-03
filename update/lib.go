package update

import (
	"fmt"
	"time"
)

func update(worker *Worker, sync chan StatusResponse, resp StatusResponse) {
	// This is just an example which simulates increasing the progress percentage by one per second
	for i := 0; i < 100; i++ {
		select {
		case <-worker.ctx.Done():
			return
		case <-time.After(time.Second * 1):
			resp.State = StatusResponse_FLASHING
			resp.Progress = int32(i)
			resp.Message = fmt.Sprintf("message: %d", i)
			sync <- resp
		}
	}
}
