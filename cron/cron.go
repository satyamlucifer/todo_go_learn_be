package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

func StartCron() {
	croner := cron.New()
	croner.AddFunc("@every 1m", func() {
		log.Println("Running cron job every 1 minute")

	})
	croner.Start()
}
