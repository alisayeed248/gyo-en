package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/alisayeed248/gyo-en/internal/monitor"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("🚀 NEW VERSION v2 - Testing simple deployment!")
	fmt.Println("DEBUG: About to read URLs from file...")

	// Read URLs from file
	urls, err := readURLsFromFile("/config/urls.txt")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("DEBUG: Successfully read URLs") // ADD THIS

	fmt.Printf("Read %d URLs from file\n", len(urls))
	for _, url := range urls {
		fmt.Printf("- %s\n", url)
	}

	rdb := connectRedis()
	if rdb == nil {
		fmt.Println("Cannot continue without Redis")
		return
	}

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

	rdb.LTrim(ctx, key, 0, 99)

	return nil
}
