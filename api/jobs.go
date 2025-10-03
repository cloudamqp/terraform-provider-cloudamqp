package api

import (
	"context"
	"fmt"
	"time"

	job "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
)

func (api *API) PollForJobCompleted(ctx context.Context, instanceID int, jobID string, sleep time.Duration) (job.JobResponse, error) {
	const interval = 5 * time.Second

	_, ok := ctx.Deadline()
	if !ok {
		return job.JobResponse{}, fmt.Errorf("context has no deadline")
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return job.JobResponse{}, fmt.Errorf("context cancelled while polling for job completed")
		case <-ticker.C:
			data, err := api.ReadJob(ctx, instanceID, jobID, sleep)
			if err != nil {
				return job.JobResponse{}, err
			}

			if *data.Status == "completed" {
				return data, nil
			}

			if *data.Status == "failed" {
				return job.JobResponse{}, fmt.Errorf("job failed: %s", *data.ErrorMessage)
			}

			continue
		}
	}
}

func (api *API) ReadJob(ctx context.Context, instanceID int, jobID string, sleep time.Duration) (job.JobResponse, error) {
	path := fmt.Sprintf("/api/instances/%d/jobs/%s", instanceID, jobID)
	var (
		data   job.JobResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Get(path), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobResponse{}, err
	}

	return data, nil
}
