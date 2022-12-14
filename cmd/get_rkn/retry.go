package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// retry tries to call the given function with `sleep` interval while the function returns a non-nil error,
// maximum `attempts` times
func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := attempts - 1; i > 0; i-- {
		if err = f(); err == nil {
			return
		}
		log.Infof("retry, %d left, sleep for %v, err: %v", i, sleep, err)
		time.Sleep(sleep)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
