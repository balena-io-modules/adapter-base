package update

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/resin-io/adapter-base/wifi"
)

func update(worker *Worker, sync chan StatusResponse, resp StatusResponse) {
	hosts, err := wifi.Scan(worker.ctx)
	if err != nil {
		if errors.Cause(err) == context.DeadlineExceeded {
			resp.State = StatusResponse_TIMED_OUT
		} else if errors.Cause(err) == context.Canceled {
			resp.State = StatusResponse_CANCELLED
		} else {
			resp.Message = err.Error()
			resp.State = StatusResponse_FAILED
		}
		sync <- resp
		return
	}

	ip := ""
	for _, host := range hosts {
		if strings.EqualFold(host.Mac, resp.StartRequest.Address) {
			ip = host.Ip
			break
		}
	}

	if ip == "" {
		resp.Message = "device is offline"
		resp.State = StatusResponse_FAILED
		sync <- resp
		return
	}

	if err := wifi.PostForm("http://"+ip+"/update", resp.StartRequest.Payload); err != nil {
		resp.Message = err.Error()
		resp.State = StatusResponse_FAILED
		sync <- resp
		return
	}

	resp.State = StatusResponse_COMPLETED
	resp.Completed = time.Now().UTC().Unix()
	sync <- resp
}
