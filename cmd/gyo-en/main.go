package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alisayeed248/gyo-en/internal/monitor"
	"github.com/redis/go-redis/v9"
	"github.com/alisayeed248/gyo-en/backend/database"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var urls []string
var rdb *redis.Client

func main() {
	fmt.Println("🚀 gyo-en starting...")

	database.InitDatabase()
	
	// Environment-aware configuration
	environment := getEnv("ENVIRONMENT", "development")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	port := getEnv("PORT", "8080")
	urlsFile := getEnv("URLS_FILE", "test-urls.txt")
	
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Redis: %s\n", redisAddr)
	fmt.Printf("Port: %s\n", port)

	// Try to read URLs from file, fallback to hardcoded
	var err error
	urls, err = readURLsFromFile(urlsFile)
	if err != nil {
		fmt.Printf("Failed to read URLs from file (%s), using defaults: %v\n", urlsFile, err)
		// Fallback URLs for local development
		urls = []string{
			"https://www.google.com",
			"https://github.com",
			"https://httpstat.us/500",
		}
	}
	fmt.Printf("Monitoring %d URLs\n", len(urls))

	// Try to connect to Redis, but don't fail if it's not available
	rdb = connectRedis(redisAddr)
	redisAvailable := rdb != nil

	// Set up API endpoints
	http.HandleFunc("/api/status", apiStatusHandler)
	http.HandleFunc("/health", healthHandler)

	// Start HTTP server in background
	go func() {
		fmt.Printf("HTTP server starting on port %s...\n", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	// Main monitoring loop
	for {
		fmt.Printf("\n--- Checking at %s ---\n", time.Now().Format("15:04:05"))

		for _, url := range urls {
			isUp, duration, err := monitor.CheckURL(url)

			if err != nil {
				fmt.Printf("Error checking %s: %v\n", url, err)
				isUp = false
			}

			// Only do Redis operations if Redis is available
			if redisAvailable {
				hasChanged, changeType, detectErr := detectStatusChange(rdb, url, isUp)
				if detectErr != nil {
					fmt.Printf("Failed to detect changes for %s: %v\n", url, detectErr)
				}

				storeCheckInDatabase(url, isUp, duration, 200, "")
				if storeErr != nil {
					fmt.Printf("Failed to store result for %s: %v\n", url, storeErr)
				}

				if hasChanged {
					fmt.Printf("🚨 ALERT: %s changed status: %s\n", url, changeType)
				}
			}

			// Always show status in console
			if isUp {
				fmt.Printf("✅ %s is UP (took %v)\n", url, duration)
			} else {
				fmt.Printf("❌ %s is DOWN (took %v)\n", url, duration)
			}
		}

		fmt.Println("💤 Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)
	}
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func connectRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("⚠️  Redis not available at %s: %v\n", addr, err)
		fmt.Println("   Running without Redis (status won't be stored)")
		return nil
	}

	fmt.Printf("✅ Connected to Redis at %s\n", addr)
	return rdb
}

func storeCheckResult(rdb *redis.Client, url string, isUp bool, duration time.Duration) error {
	if rdb == nil {
		return nil // Skip if Redis not available
	}
	
	ctx := context.Background()
	timestamp := time.Now().Format("2006-01-02T15:04:05")
	
	status := "DOWN"
	if isUp {
		status = "UP"
	}

	result := fmt.Sprintf("%s|%s|%v", timestamp, status, duration)
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
	if rdb == nil {
		return false, "NO_REDIS", nil
	}
	
	ctx := context.Background()
	key := fmt.Sprintf("checks:%s", url)

	lastResult, err := rdb.LIndex(ctx, key, 0).Result()
	if err != nil {
		if err == redis.Nil {
			return true, "NEW", nil
		}
		return false, "", err
	}

	parts := strings.Split(lastResult, "|")
	if len(parts) < 2 {
		return false, "", fmt.Errorf("invalid stored result format")
	}

	previousStatus := parts[1] == "UP"

	if previousStatus && !currentStatus {
		return true, "UP->DOWN", nil
	} else if !previousStatus && currentStatus {
		return true, "DOWN->UP", nil
	}

	return false, "NO_CHANGE", nil
}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // For local development

	var statusList []map[string]interface{}

	for _, url := range urls {
		status := map[string]interface{}{
			"url":       url,
			"status":    "UNKNOWN",
			"lastCheck": "",
		}

		// Only get from Redis if available
		if rdb != nil {
			ctx := context.Background()
			key := fmt.Sprintf("checks:%s", url)
			lastResult, err := rdb.LIndex(ctx, key, 0).Result()

			if err == nil && lastResult != "" {
				parts := strings.Split(lastResult, "|")
				if len(parts) >= 2 {
					status["status"] = parts[1]
					status["lastCheck"] = parts[0]
				}
			}
		}

		statusList = append(statusList, status)
	}

	response := map[string]interface{}{
		"urls":      statusList,
		"timestamp": time.Now().Format(time.RFC3339),
		"redis":     rdb != nil,
	}
	
	json.NewEncoder(w).Encode(response)
}

func storeCheckInDatabase(url string, isUp bool, duration time.Duration, statusCode int, errorMsg string) {
	checkResult := database.CheckResult{
		URL:          url,
		IsUp:         isUp,
		ResponseTime: duration,
		StatusCode:   statusCode,
		ErrorMessage: errorMsg,
		CheckedAt:    time.Now(),
	}
	
	database.DB.Create(&checkResult)
}