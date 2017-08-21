package adapter

import (
	"context"
	"fmt"
	"time"

	structpb "github.com/golang/protobuf/ptypes/struct"
)

func simulateScan(req *ScanOptions, worker *Worker, sync chan Job, resp Job) {
	// This is just an example which simulates finding one device per second
	for i := 0; i < 100; i++ {
		select {
		case <-worker.ctx.Done():
			if worker.ctx.Err() == context.DeadlineExceeded {
				resp.State = State_TIMED_OUT
				sync <- resp
			}
			return
		case <-time.After(time.Second * 1):
			extra := make(map[string]*structpb.Value)
			extra["application"] = &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("Application: %d", i)},
			}
			extra["MAC"] = &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("MAC: %d", i)},
			}
			destination := &Destination{
				Id:    fmt.Sprintf("IP: %d", i),
				Extra: extra,
			}

			resp.Destinations = append(resp.Destinations, destination)
			sync <- resp
		}
	}

	resp.State = State_COMPLETED
	sync <- resp
}

func simulateUpdate(req *UpdateOptions, worker *Worker, sync chan Job, resp Job) {
	// This is just an example which simulates increasing the progress percentage by one per second
	for i := 0; i < 100; i++ {
		select {
		case <-worker.ctx.Done():
			if worker.ctx.Err() == context.DeadlineExceeded {
				resp.State = State_TIMED_OUT
				sync <- resp
			}
			return
		case <-time.After(time.Second * 1):
			resp.Progress.Percentage = int64(i)
			resp.Progress.Transferred = int64(i) * 3
			resp.Progress.Length = 300
			resp.Progress.Remaining = resp.Progress.Length - resp.Progress.Transferred
			resp.Progress.Eta = 10000000000
			resp.Progress.Runtime = int64(i)
			resp.Progress.Speed = 3
			sync <- resp
		}
	}

	resp.State = State_COMPLETED
	sync <- resp
}
