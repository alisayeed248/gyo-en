package main

import (
	"fmt"
	"github.com/alisayeed248/gyo-en/internal/monitor"
	"time"
	"os"
	"bufio"
	"strings"
	"net/http"
	"log"
)

func main() {
	fmt.Println("gyo-en uptime monitor starting...")
	fmt.Println("DEBUG: About to read URLs from file...")

	// Read URLs from file
	urls, err := readURLsFromFile("/config/urls.txt")
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("DEBUG: Successfully read URLs")  // ADD THIS

	fmt.Printf("Read %d URLs from file\n", len(urls))
	for _, url := range urls {
		fmt.Printf("- %s\n", url)
	}

	go startHealthServer()
	fmt.Println("HTTP server goroutine started!")

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

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}

func startHealthServer() {
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Health server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}