package main

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/config"
	"github.com/irsyadjpp/tracking-diet-backend/internal/worker"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

func main() {
	// load config (jika dibutuhkan untuk hal lain)
	cfg := config.LoadConfig()

	// inisialisasi logger tanpa argumen
	logg := logger.New()
	logg.Info("Worker is starting… using DB: " + cfg.DBUrl)

	// jalankan worker scheduler tanpa argumen
	worker.StartScheduler()

	// blokir agar tidak langsung keluar
	select {}
}
