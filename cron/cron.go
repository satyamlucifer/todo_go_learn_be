package cron

import (
    "github.com/robfig/cron/v3"
    "log"
)

func StartCron() {
    c := cron.New()
    c.AddFunc("@every 1m", func() {
        log.Println("Running cron job every 1 minute")
        // your scheduled task here, e.g. clean up todos, send email reminders, etc.
    })
    c.Start()
}
