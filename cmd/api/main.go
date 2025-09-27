package api

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/config"
	"github.com/irsyadjpp/tracking-diet-backend/internal/delivery"
	"github.com/irsyadjpp/tracking-diet-backend/internal/repository"
	"github.com/irsyadjpp/tracking-diet-backend/internal/usecase"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New(cfg.Env)

	db := repository.NewPostgresDB(cfg.DatabaseURL)

	userRepo := repository.NewUserRepository(db)
	bodyRepo := repository.NewBodyMeasurementRepository(db)

	userUC := usecase.NewUserUsecase(userRepo)
	bodyUC := usecase.NewBodyMeasurementUsecase(bodyRepo)

	srv := delivery.NewHTTPServer(cfg, log, userUC, bodyUC)
	srv.Run()
}
