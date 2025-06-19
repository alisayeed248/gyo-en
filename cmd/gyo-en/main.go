package main

import (
	"bufio"
	"context"
	"encoding/json" // need this for JSON responses
	"fmt"
	"github.com/alisayeed248/gyo-en/internal/monitor"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var urls []string     // accessed from API handlers
var rdb *redis.Client // moved to outside main funct

func main() {
	fmt.Println("🚀 NEW VERSION v3 - Adding API endpoints!")

	// Read URLs from file
	var err error
	urls, err := readURLsFromFile("/config/urls.txt")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Read %d URLs from file\n", len(urls))

	rdb = connectRedis()
	if rdb == nil {
		fmt.Println("Cannot continue without Redis")
		return
	}

	http.HandleFunc("/api/status", apiStatusHandler)

	ctx := context.Background()
	testKey := fmt.Sprintf("test:%d", time.Now().Unix())
	err = rdb.Set(ctx, testKey, "Hello Redis!", 0).Err()
	if err != nil {
		fmt.Printf("❌ Redis SET failed: %v\n", err)
	} else {
		val, err := rdb.Get(ctx, testKey).Result()
		if err != nil {
			fmt.Printf("❌ Redis GET failed: %v\n", err)
		} else {
			fmt.Printf("🚀 Redis test SUCCESS: %s\n", val)
		}
	}

	go startHealthServer()
	fmt.Println("HTTP server goroutine started!")

	for {
		fmt.Printf("\n--- Checking at %s ---\n", time.Now().Format("15:04:05"))

		for _, url := range urls {
			isUp, duration, err := monitor.CheckURL(url)

			if err != nil {
				fmt.Printf("Error checking %s: %v\n", url, err)
				isUp = false // Treat errors as DOWN
			}

			hasChanged, changeType, detectErr := detectStatusChange(rdb, url, isUp)
			if detectErr != nil {
				fmt.Printf("Failed to detect changes for %s: %v\n", url, detectErr)
			}

			// Store the new result
			storeErr := storeCheckResult(rdb, url, isUp, duration)
			if storeErr != nil {
				fmt.Printf("Failed to store result for %s: %v\n", url, storeErr)
			}

			if hasChanged {
				fmt.Printf("🚨 ALERT: %s changed status: %s\n", url, changeType)
			}

			if isUp {
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

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-service:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return nil
	}

	fmt.Println("✅ Connected to Redis!")
	return rdb
}

// rdb is pointer to our redis connection, isUp is our true/false result from the check
func storeCheckResult(rdb *redis.Client, url string, isUp bool, duration time.Duration) error {
	// run normally
	ctx := context.Background()

	timestamp := time.Now().Format("2006-01-02T15:04:05")

	// Default status is down, we can move to up
	status := "DOWN"
	if isUp {
		status = "UP"
	}

	// Sprintf is the fmt command that lets us build a string. this is string string value
	result := fmt.Sprintf("%s|%s|%v", timestamp, status, duration)

	// store in Redis list (most recent first)
	key := fmt.Sprintf("checks:%s", url)
	err := rdb.LPush(ctx, key, result).Err()
	if err != nil {
		return err
	}

	err = rdb.LTrim(ctx, key, 0, 99).Err()
	if err != nil {
		fmt.Printf("Warning: Failed to trim old data for %s: %v\n", url, err)
	}

	return nil
}

func detectStatusChange(rdb *redis.Client, url string, currentStatus bool) (bool, string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("checks:%s", url)

	// get most recent stored result (index 0)
	lastResult, err := rdb.LIndex(ctx, key, 0).Result()
	if err != nil {
		// if no key, not error
		if err == redis.Nil {
			return true, "NEW", nil
		}
		return false, "", err
	}

	// parse the last result to get previous status
	// original format is like timestamp:status:ms
	parts := strings.Split(lastResult, "|")
	if len(parts) < 2 {
		return false, "", fmt.Errorf("invalid stored result format")
	}

	previousStatus := parts[1] == "UP"

	// Detect changes
	if previousStatus && !currentStatus {
		return true, "UP->DOWN", nil
	} else if !previousStatus && currentStatus {
		return true, "DOWN->UP", nil
	}

	return false, "NO_CHANGE", nil
}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Create context for Redis operations
	ctx := context.Background()

	// We hold all the URL statuses in a slice
	var statusList []map[string]interface{}

	// Go through each URL and get the latest status from Redis cache
	for _, url := range urls {
		// build our redis key
		key := fmt.Sprintf("checks:%s", url)

		// use this key and get the most recent result
		lastResult, err := rdb.LIndex(ctx, key, 0).Result()

		// create a status object for the URL
		status := map[string]interface{}{
			"url":       url,
			"status":    "UNKNOWN",
			"lastCheck": "",
		}

		// if there's data, we parse it and add the URl status to our list
		if err == nil && lastResult != "" {
			parts := strings.Split(lastResult, "|")
			if len(parts) >= 2 {
				status["status"] = parts[1] // "UP" or "DOWN"
				status["lastCheck"] = parts[0]
			}
		}

		statusList = append(statusList, status)

		// create final response

		// we send JSON response

	}
	response := map[string]interface{}{
		"urls":      statusList,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)

}
