package monitor

import (
	"net/http"
	"time"
)

func CheckURL(url string) (bool, time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return false, duration, err
	}

	defer resp.Body.Close()
	isUp := resp.StatusCode >= 200 && resp.StatusCode < 300

	return isUp, duration, nil
}