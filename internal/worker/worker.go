package worker

import (
	"log"
	"time"
)

func StartScheduler() {
	log.Println("Worker scheduler started…")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			log.Println("Running scheduled job…")
		}
	}()
}
