package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/redis/go-redis/v9"
)

// JobType represents different types of scheduled jobs
type JobType string

const (
	JobTypeAIRecommendation JobType = "ai_recommendation"
	JobTypeHealthScore       JobType = "health_score"
	JobTypeDataCleanup       JobType = "data_cleanup"
	JobTypeNutritionReminder JobType = "nutrition_reminder"
)

// Scheduler handles cron-based job scheduling
type Scheduler struct {
	cron         *cron.Cron
	redisClient  *redis.Client
	queue        *Queue
	jobHandlers  map[JobType]JobHandler
	distributed  bool
	lockKey      string
	mu           sync.Mutex
}

// JobHandler defines the interface for job handlers
type JobHandler interface {
	Handle(ctx context.Context, payload map[string]interface{}) error
}

// NewScheduler creates a new cron scheduler
func NewScheduler(redisClient *redis.Client, queue *Queue, distributed bool) *Scheduler {
	options := cron.Options{
		Seconds: true, // Enable second-level precision
	}

	if distributed {
		options = cron.Options{
			Seconds: true,
		}
	}

	return &Scheduler{
		cron:        cron.New(options),
		redisClient: redisClient,
		queue:       queue,
		jobHandlers: make(map[JobType]JobHandler),
		distributed: distributed,
		lockKey:     "scheduler:lock",
	}
}

// RegisterHandler registers a job handler for a specific job type
func (s *Scheduler) RegisterHandler(jobType JobType, handler JobHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobHandlers[jobType] = handler
}

// Start starts the scheduler
func (s *Scheduler) Start() error {
	// Register scheduled jobs
	if err := s.registerScheduledJobs(); err != nil {
		return fmt.Errorf("failed to register scheduled jobs: %w", err)
	}

	s.cron.Start()
	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// registerScheduledJobs registers all scheduled jobs
func (s *Scheduler) registerScheduledJobs() error {
	// Daily AI recommendation generation at 9 AM
	if _, err := s.cron.AddFunc("0 0 9 * * *", func() {
		s.executeJob(JobTypeAIRecommendation, map[string]interface{}{
			"trigger": "scheduled",
			"time":    "daily_9am",
		})
	}); err != nil {
		return fmt.Errorf("failed to schedule AI recommendation job: %w", err)
	}

	// Weekly health score calculation on Sunday at 8 AM
	if _, err := s.cron.AddFunc("0 0 8 * * 0", func() {
		s.executeJob(JobTypeHealthScore, map[string]interface{}{
			"trigger": "scheduled",
			"time":    "weekly_sunday_8am",
		})
	}); err != nil {
		return fmt.Errorf("failed to schedule health score job: %w", err)
	}

	// Monthly data cleanup on the 1st of each month at 2 AM
	if _, err := s.cron.AddFunc("0 0 2 1 * *", func() {
		s.executeJob(JobTypeDataCleanup, map[string]interface{}{
			"trigger": "scheduled",
			"time":    "monthly_1st_2am",
		})
	}); err != nil {
		return fmt.Errorf("failed to schedule data cleanup job: %w", err)
	}

	// Nutrition reminders at 8 AM, 12 PM, and 6 PM
	mealTimes := []string{"0 0 8 * * *", "0 0 12 * * *", "0 0 18 * * *"}
	for _, cronExpr := range mealTimes {
		if _, err := s.cron.AddFunc(cronExpr, func() {
			s.executeJob(JobTypeNutritionReminder, map[string]interface{}{
				"trigger": "scheduled",
				"type":    "meal_reminder",
			})
		}); err != nil {
			return fmt.Errorf("failed to schedule nutrition reminder job: %w", err)
		}
	}

	return nil
}

// executeJob executes a job by enqueuing it to the queue
func (s *Scheduler) executeJob(jobType JobType, payload map[string]interface{}) {
	ctx := context.Background()

	// For distributed mode, acquire lock to prevent duplicate execution
	if s.distributed {
		if !s.acquireLock(ctx, string(jobType)) {
			return // Another instance is handling this job
		}
		defer s.releaseLock(ctx, string(jobType))
	}

	job := &Job{
		ID:      generateJobID(),
		Type:    string(jobType),
		Payload: payload,
	}

	if err := s.queue.Enqueue(ctx, job); err != nil {
		fmt.Printf("Failed to enqueue job %s: %v\n", jobType, err)
	}
}

// acquireLock acquires a distributed lock for the job
func (s *Scheduler) acquireLock(ctx context.Context, key string) bool {
	lockKey := fmt.Sprintf("%s:%s", s.lockKey, key)
	
	// Try to set lock with expiration
	result, err := s.redisClient.SetNX(ctx, lockKey, "1", 5*time.Minute).Result()
	if err != nil {
		return false
	}
	
	return result
}

// releaseLock releases a distributed lock
func (s *Scheduler) releaseLock(ctx context.Context, key string) {
	lockKey := fmt.Sprintf("%s:%s", s.lockKey, key)
	s.redisClient.Del(ctx, lockKey)
}

// generateJobID generates a unique job ID
func generateJobID() string {
	return fmt.Sprintf("job_%d", time.Now().UnixNano())
}