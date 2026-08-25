package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/irsyadjpp/tracking-diet-backend/internal/ai"
	"github.com/irsyadjpp/tracking-diet-backend/internal/config"
	"github.com/irsyadjpp/tracking-diet-backend/internal/repository"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/internal/worker"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	appLogger := logger.New(cfg.Env)

	appLogger.Info("Starting Tracking Diet Backend Worker", map[string]interface{}{
		"env": cfg.Env,
	})

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		appLogger.Error("Failed to connect to Redis", err, nil)
		log.Fatal("Failed to connect to Redis")
	}

	appLogger.Info("Connected to Redis successfully", nil)

	// Initialize database
	db, err := repository.NewPostgresDB()
	if err != nil {
		appLogger.Error("Failed to connect to database", err, nil)
		log.Fatal("Failed to connect to database")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	metabolicRepo := repository.NewMetabolicMeasurementRepository(db)
	bodyRepo := repository.NewBodyMeasurementRepository(db)
	nutritionRepo := repository.NewNutritionMeasurementRepository(db)

	// Initialize AI service
	aiService := ai.NewAIService(cfg.GoogleGenkitAPIKey)
	if cfg.GoogleGenkitAPIKey != "" {
		if err := aiService.Initialize(ctx); err != nil {
			appLogger.Error("Failed to initialize AI service", err, nil)
			log.Printf("Warning: AI service initialization failed: %v", err)
		} else {
			appLogger.Info("AI service initialized successfully", nil)
		}
	}

	// Initialize use cases
	aiUC := usecase.NewAIRecommendationUsecase(repository.NewAIRecommendationRepository(db), userRepo, aiService)

	// Initialize job queue
	queue := worker.NewQueue(redisClient, "tracking_diet")

	// Initialize job handlers
	aiJobHandler := worker.NewAIRecommendationJobHandler(userRepo, aiUC)
	healthScoreHandler := worker.NewHealthScoreJobHandler(metabolicRepo)
	dataCleanupHandler := worker.NewDataCleanupJobHandler(bodyRepo, nutritionRepo)
	nutritionReminderHandler := worker.NewNutritionReminderJobHandler(userRepo)

	// Initialize worker pool
	workerPool := worker.NewWorkerPool(queue, 5, 3) // 5 workers, max 3 retries
	workerPool.RegisterHandler(worker.JobTypeAIRecommendation, aiJobHandler)
	workerPool.RegisterHandler(worker.JobTypeHealthScore, healthScoreHandler)
	workerPool.RegisterHandler(worker.JobTypeDataCleanup, dataCleanupHandler)
	workerPool.RegisterHandler(worker.JobTypeNutritionReminder, nutritionReminderHandler)

	// Initialize scheduler
	scheduler := worker.NewScheduler(redisClient, queue, true) // true for distributed mode
	scheduler.RegisterHandler(worker.JobTypeAIRecommendation, aiJobHandler)
	scheduler.RegisterHandler(worker.JobTypeHealthScore, healthScoreHandler)
	scheduler.RegisterHandler(worker.JobTypeDataCleanup, dataCleanupHandler)
	scheduler.RegisterHandler(worker.JobTypeNutritionReminder, nutritionReminderHandler)

	// Start worker pool
	if err := workerPool.Start(); err != nil {
		appLogger.Error("Failed to start worker pool", err, nil)
		log.Fatal("Failed to start worker pool")
	}

	// Start scheduler
	if err := scheduler.Start(); err != nil {
		appLogger.Error("Failed to start scheduler", err, nil)
		log.Fatal("Failed to start scheduler")
	}

	appLogger.Info("Worker started successfully", map[string]interface{}{
		"worker_count": 5,
		"redis_addr":   cfg.RedisHost + ":" + cfg.RedisPort,
	})

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down worker...", nil)

	// Graceful shutdown
	scheduler.Stop()
	workerPool.Stop()

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		appLogger.Error("Failed to close Redis connection", err, nil)
	}

	appLogger.Info("Worker exited", nil)
}
