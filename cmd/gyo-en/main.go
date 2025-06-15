package main

import (
	"fmt"
	"github.com/alisayeed248/gyo-en/internal/monitor"
)

func main() {
	fmt.Println("gyo-en uptime monitor starting...")
	urls := []string{"https://www.google.com", "https://test-fake-website-12443.com", "https://httpstat.us/500",}
	
	
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
	
}