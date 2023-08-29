package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sampiiiii-dev/anvil_server/anvil/workers/jobs"
)

type JobWrapper struct {
	JobType string          `json:"job_type"`
	Payload json.RawMessage `json:"payload"`
}

type RedisJobQueue struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisJobQueue(client *redis.Client, ctx context.Context) *RedisJobQueue {
	return &RedisJobQueue{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisJobQueue) Enqueue(job jobs.Job, jobType string) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}
	wrapper := JobWrapper{
		JobType: jobType,
		Payload: payload,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	r.client.RPush(r.ctx, "job_queue", data)
	r.client.Publish(r.ctx, "job_channel", "new job")
	return nil
}

func (r *RedisJobQueue) Dequeue() (jobs.Job, error) {
	data, err := r.client.LPop(r.ctx, "job_queue").Bytes()
	if err != nil {
		return nil, err
	}

	var wrapper JobWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	var job jobs.Job
	switch wrapper.JobType {
	case "email":
		var eJob jobs.EmailJob
		if err := json.Unmarshal(wrapper.Payload, &eJob); err != nil {
			return nil, err
		}
		job = &eJob
	// Add other job types here as cases
	default:
		return nil, fmt.Errorf("unknown job type: %s", wrapper.JobType)
	}
	return job, nil
}
