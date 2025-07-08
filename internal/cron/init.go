package cron

import (
	"peanut/internal/cache"
	"time"
)

func Init() error {
	go func() {
		for {
			time.Sleep(5 * time.Minute)

			cache.Refresh()
		}
	}()

	return nil
}
