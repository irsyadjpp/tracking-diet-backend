package api

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/irsyadjpp/tracking-diet-backend/internal/auth"
	"github.com/irsyadjpp/tracking-diet-backend/internal/config"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery"
	"github.com/irsyadjpp/tracking-diet-backend/internal/repository"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	appLogger := logger.New(cfg.Env)

	appLogger.Info("Starting Tracking Diet Backend", map[string]interface{}{
		"env":  cfg.Env,
		"port": cfg.Port,
	})

	db, err := repository.NewPostgresDB()
	if err != nil {
		appLogger.Error("Failed to connect to database", err, nil)
		log.Fatal("Failed to connect to database")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bodyRepo := repository.NewBodyMeasurementRepository(db)
	nutritionRepo := repository.NewNutritionMeasurementRepository(db)
	metabolicRepo := repository.NewMetabolicMeasurementRepository(db)
	fitnessRepo := repository.NewFitnessMeasurementRepository(db)
	wellbeingRepo := repository.NewWellbeingMeasurementRepository(db)
	labTestRepo := repository.NewLabTestRepository(db)
	aiRecommendationRepo := repository.NewAIRecommendationRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, time.Duration(cfg.JWTExpiration)*time.Hour)

	// Initialize use cases
	userUC := usecase.NewUserUsecase(userRepo)
	bodyUC := usecase.NewBodyMeasurementUsecase(bodyRepo, userRepo)
	nutritionUC := usecase.NewNutritionUsecase(nutritionRepo, userRepo)
	metabolicUC := usecase.NewMetabolicUsecase(metabolicRepo, userRepo)
	fitnessUC := usecase.NewFitnessUsecase(fitnessRepo, userRepo)
	wellbeingUC := usecase.NewWellbeingUsecase(wellbeingRepo, userRepo)
	labTestUC := usecase.NewLabTestUsecase(labTestRepo, userRepo)
	aiRecommendationUC := usecase.NewAIRecommendationUsecase(aiRecommendationRepo, userRepo, aiService)
	dashboardUC := usecase.NewDashboardUsecase(dashboardRepo)

	// Initialize HTTP server
	server := delivery.NewHTTPServer(cfg, appLogger, userRepo, jwtManager, userUC, bodyUC, nutritionUC, metabolicUC, fitnessUC, wellbeingUC, labTestUC, aiRecommendationUC, dashboardUC)

	// Start server in a goroutine
	go func() {
		if err := server.Run(); err != nil {
			appLogger.Error("Server failed to start", err, nil)
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...", nil)

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", err, nil)
		log.Fatal(err)
	}

	appLogger.Info("Server exited", nil)
}
