package workers

import (
	"github.com/sampiiiii-dev/anvil_server/anvil/workers/jobs"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type WorkerPool struct {
	jobs           chan jobs.Job
	shutdownSignal chan struct{}
	wg             sync.WaitGroup
	workers        int
	queue          *RedisJobQueue
	l              *zap.Logger
}

func NewWorkerPool(workers int, jobQueueSize int, queue *RedisJobQueue, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		jobs:           make(chan jobs.Job, jobQueueSize),
		workers:        workers,
		queue:          queue,
		l:              logger,
		shutdownSignal: make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	// Signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start Redis Pub/Sub
	pubsub := wp.queue.client.Subscribe(wp.queue.ctx, "job_channel")
	ch := pubsub.Channel()

	// Start workers
	wp.l.Info("Worker Pool: Starting workers")
	wp.wg.Add(wp.workers)
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}
	wp.l.Info("Worker Pool: Workers started")

	// Event loop for job handling
	for {
		select {
		case <-stop:
			err := pubsub.Close()
			if err != nil {
				return
			}
			wp.Stop()
			return
		case <-ch:
			job, err := wp.queue.Dequeue()
			if err == nil {
				wp.jobs <- job
			} else {
				wp.l.Error("Error dequeuing job", zap.Error(err))
			}
		}
	}
}

func (wp *WorkerPool) Stop() {
	wp.l.Info("Worker Pool: Initiating shutdown")
	close(wp.jobs)
	close(wp.shutdownSignal)
	wp.wg.Wait()
	wp.l.Info("Worker Pool: Workers stopped")
}

func (wp *WorkerPool) worker(workerID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	defer wg.Done()
	defer wp.wg.Done()
	wp.l.Info("Worker started", zap.Int("WorkerID", workerID))

	for {
		select {
		case <-wp.shutdownSignal:
			wp.l.Info("Worker shutting down", zap.Int("WorkerID", workerID))
			return
		case job, ok := <-wp.jobs:
			wp.l.Info("Worker received job", zap.Int("WorkerID", workerID))
			if ok {
				err := job.Execute()
				if err != nil {
					wp.l.Error("Error executing job", zap.Error(err))
				}
			}
		default:
			// Sleeps for 1-10 seconds
			time.Sleep(time.Second * time.Duration(1+(workerID%10)))
		}
	}
}
