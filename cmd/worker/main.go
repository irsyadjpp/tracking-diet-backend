package main

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/config"
	"github.com/irsyadjpp/tracking-diet-backend/internal/worker"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New(cfg.Env)

	worker.StartScheduler(cfg, log)
}
