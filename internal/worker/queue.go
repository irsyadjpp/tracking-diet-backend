package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Job represents a background job
type Job struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	RetryCount int                   `json:"retry_count"`
	CreatedAt time.Time              `json:"created_at"`
}

// Queue handles Redis-based job queue operations
type Queue struct {
	client     *redis.Client
	queueName  string
	processing string
	deadLetter string
}

// NewQueue creates a new Redis queue
func NewQueue(client *redis.Client, queueName string) *Queue {
	return &Queue{
		client:     client,
		queueName:  queueName + ":queue",
		processing: queueName + ":processing",
		deadLetter: queueName + ":dead_letter",
	}
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(ctx context.Context, job *Job) error {
	job.CreatedAt = time.Now()
	job.RetryCount = 0

	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Add to Redis list (left push)
	if err := q.client.LPush(ctx, q.queueName, jobData).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

// Dequeue retrieves a job from the queue (blocking with timeout)
func (q *Queue) Dequeue(ctx context.Context, timeout time.Duration) (*Job, error) {
	// Use BRPOPLPUSH to atomically move job from queue to processing
	result, err := q.client.BRPopLPush(ctx, q.queueName, q.processing, timeout).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No job available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(result), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

// Complete marks a job as completed and removes it from processing
func (q *Queue) Complete(ctx context.Context, jobID string) error {
	// Remove from processing queue
	pattern := fmt.Sprintf("*\"id\":\"%s\"*", jobID)
	
	// Get all jobs in processing queue
	jobs, err := q.client.LRange(ctx, q.processing, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get processing jobs: %w", err)
	}

	// Find and remove the specific job
	for _, jobData := range jobs {
		var job Job
		if err := json.Unmarshal([]byte(jobData), &job); err != nil {
			continue
		}
		if job.ID == jobID {
			if err := q.client.LRem(ctx, q.processing, 1, jobData).Err(); err != nil {
				return fmt.Errorf("failed to remove completed job: %w", err)
			}
			break
		}
	}

	return nil
}

// Retry moves a job back to the queue with retry count increment
func (q *Queue) Retry(ctx context.Context, job *Job, maxRetries int) error {
	job.RetryCount++

	if job.RetryCount > maxRetries {
		// Move to dead letter queue
		return q.MoveToDeadLetter(ctx, job)
	}

	// Remove from processing
	pattern := fmt.Sprintf("*\"id\":\"%s\"*", job.ID)
	jobs, err := q.client.LRange(ctx, q.processing, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get processing jobs: %w", err)
	}

	for _, jobData := range jobs {
		var existingJob Job
		if err := json.Unmarshal([]byte(jobData), &existingJob); err != nil {
			continue
		}
		if existingJob.ID == job.ID {
			if err := q.client.LRem(ctx, q.processing, 1, jobData).Err(); err != nil {
				return fmt.Errorf("failed to remove job for retry: %w", err)
			}
			break
		}
	}

	// Re-enqueue with exponential backoff
	backoff := time.Duration(job.RetryCount) * time.Second
	time.Sleep(backoff)

	return q.Enqueue(ctx, job)
}

// MoveToDeadLetter moves a failed job to the dead letter queue
func (q *Queue) MoveToDeadLetter(ctx context.Context, job *Job) error {
	// Remove from processing
	pattern := fmt.Sprintf("*\"id\":\"%s\"*", job.ID)
	jobs, err := q.client.LRange(ctx, q.processing, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get processing jobs: %w", err)
	}

	for _, jobData := range jobs {
		var existingJob Job
		if err := json.Unmarshal([]byte(jobData), &existingJob); err != nil {
			continue
		}
		if existingJob.ID == job.ID {
			if err := q.client.LRem(ctx, q.processing, 1, jobData).Err(); err != nil {
				return fmt.Errorf("failed to remove job for dead letter: %w", err)
			}
			break
		}
	}

	// Add to dead letter queue
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job for dead letter: %w", err)
	}

	return q.client.LPush(ctx, q.deadLetter, jobData).Err()
}

// GetQueueSize returns the current queue size
func (q *Queue) GetQueueSize(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, q.queueName).Result()
}

// GetProcessingSize returns the current processing queue size
func (q *Queue) GetProcessingSize(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, q.processing).Result()
}

// GetDeadLetterSize returns the current dead letter queue size
func (q *Queue) GetDeadLetterSize(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, q.deadLetter).Result()
}