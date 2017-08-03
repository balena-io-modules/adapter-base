package scan

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/resin-io/adapter-base/wifi"
)

func scan(worker *Worker, sync chan StatusResponse, resp StatusResponse) {
	hosts, err := wifi.Scan(worker.ctx)
	if err != nil {
		if errors.Cause(err) == context.DeadlineExceeded {
			resp.State = StatusResponse_TIMED_OUT
		} else if errors.Cause(err) == context.Canceled {
			resp.State = StatusResponse_CANCELLED
		} else {
			resp.Message = err.Error()
			resp.State = StatusResponse_SCAN_FALURE
		}
		sync <- resp
		return
	}

	for _, host := range hosts {
		if resp.StartRequest.Address != "" && !strings.EqualFold(host.Mac, resp.StartRequest.Address) {
			break
		}

		result := &StatusResponse_Result{
			Address: host.Mac,
			Name:    host.Name,
		}
		resp.Results = append(resp.Results, result)
		resp.State = StatusResponse_COMPLETED
		resp.Completed = time.Now().UTC().Unix()
		sync <- resp
	}
}
