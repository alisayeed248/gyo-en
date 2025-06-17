package main

import (
	"fmt"
	"github.com/alisayeed248/gyo-en/internal/monitor"
	"time"
)

func main() {
	fmt.Println("gyo-en uptime monitor starting...")
	urls := []string{"https://www.google.com", "https://test-fake-website-12443.com", "https://httpstat.us/500"}

	for {
		fmt.Printf("\n--- Checking at %s ---\n", time.Now().Format("15:04:05"))

		for _, url := range urls {
			isUp, duration, err := monitor.CheckURL(url)
			if err != nil {
				fmt.Printf("Error checking %s: %v\n", url, err)
			} else if isUp {
				fmt.Printf("%s is UP (took %v)\n", url, duration)
			} else {
				fmt.Printf("%s is DOWN (took %v)\n", url, duration)
			}
		}

		fmt.Println("Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)
	}

}
