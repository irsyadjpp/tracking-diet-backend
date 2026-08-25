package worker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkerPool manages a pool of worker goroutines
type WorkerPool struct {
	workers     []*Worker
	queue       *Queue
	jobHandlers map[JobType]JobHandler
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// Worker represents a single worker
type Worker struct {
	id         int
	pool       *WorkerPool
	maxRetries int
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(queue *Queue, workerCount int, maxRetries int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		queue:       queue,
		jobHandlers: make(map[JobType]JobHandler),
		ctx:         ctx,
		cancel:      cancel,
		workers:     make([]*Worker, workerCount),
	}
}

// RegisterHandler registers a job handler for a specific job type
func (p *WorkerPool) RegisterHandler(jobType JobType, handler JobHandler) {
	p.jobHandlers[jobType] = handler
}

// Start starts the worker pool
func (p *WorkerPool) Start() error {
	for i := 0; i < len(p.workers); i++ {
		worker := &Worker{
			id:         i,
			pool:       p,
			maxRetries: 3,
		}
		p.workers[i] = worker

		p.wg.Add(1)
		go worker.run()
	}

	return nil
}

// Stop stops the worker pool gracefully
func (p *WorkerPool) Stop() {
	p.cancel() // Cancel context to signal workers to stop
	
	// Wait for all workers to finish their current jobs
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		// All workers stopped gracefully
	case <-time.After(30 * time.Second):
		// Force shutdown after timeout
	}
}

// run is the main worker loop
func (w *Worker) run() {
	defer w.pool.wg.Done()

	for {
		select {
		case <-w.pool.ctx.Done():
			// Context cancelled, stop worker
			return
		default:
			// Try to dequeue a job
			job, err := w.pool.queue.Dequeue(w.pool.ctx, 5*time.Second)
			if err != nil {
				fmt.Printf("Worker %d: error dequeuing job: %v\n", w.id, err)
				continue
			}

			if job == nil {
				// No job available, continue loop
				continue
			}

			// Process the job
			w.processJob(job)
		}
	}
}

// processJob processes a single job
func (w *Worker) processJob(job *Job) {
	handler, exists := w.pool.jobHandlers[JobType(job.Type)]
	if !exists {
		fmt.Printf("Worker %d: no handler for job type %s\n", w.id, job.Type)
		w.pool.queue.MoveToDeadLetter(w.pool.ctx, job)
		return
	}

	// Execute the job handler
	err := handler.Handle(w.pool.ctx, job.Payload)
	if err != nil {
		fmt.Printf("Worker %d: job %s failed: %v\n", w.id, job.ID, err)
		
		// Retry logic
		if job.RetryCount < w.maxRetries {
			if retryErr := w.pool.queue.Retry(w.pool.ctx, job, w.maxRetries); retryErr != nil {
				fmt.Printf("Worker %d: failed to retry job %s: %v\n", w.id, job.ID, retryErr)
			} else {
				fmt.Printf("Worker %d: job %s queued for retry (attempt %d)\n", w.id, job.ID, job.RetryCount)
			}
		} else {
			fmt.Printf("Worker %d: job %s exceeded max retries, moving to dead letter\n", w.id, job.ID)
			w.pool.queue.MoveToDeadLetter(w.pool.ctx, job)
		}
		return
	}

	// Job completed successfully
	if err := w.pool.queue.Complete(w.pool.ctx, job.ID); err != nil {
		fmt.Printf("Worker %d: failed to mark job %s as complete: %v\n", w.id, job.ID, err)
	}

	fmt.Printf("Worker %d: job %s completed successfully\n", w.id, job.ID)
}

// GetStats returns worker pool statistics
func (p *WorkerPool) GetStats(ctx context.Context) map[string]interface{} {
	queueSize, _ := p.queue.GetQueueSize(ctx)
	processingSize, _ := p.queue.GetProcessingSize(ctx)
	deadLetterSize, _ := p.queue.GetDeadLetterSize(ctx)

	return map[string]interface{}{
		"worker_count":     len(p.workers),
		"queue_size":       queueSize,
		"processing_size":  processingSize,
		"dead_letter_size": deadLetterSize,
	}
}