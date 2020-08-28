package window_test

import (
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis"

	"github.com/Marvin9/ratelimit/window"
)

func TestRateLimit(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	err := redisClient.Ping().Err()
	if err != nil {
		log.Fatalf("Start redis-server on port 6379.\nCommand: redis-server")
	}
	rt := window.New(5, time.Second*5, redisClient)

	// User can use all APIs in one window without any issue
	for i := 0; i < 5; i++ {
		_, ok := rt.Use("unique")
		if !ok {
			t.Errorf("Could not use %v times, limit exceed at %v", 5, i+1)
		}
	}
	redisClient.Del("unique")

	rt = window.New(1, time.Second*5, redisClient)
	// User cannot consume APIs more than limit and in given time window
	rt.Use("unique")
	instance, ok := rt.Use("unique")
	if ok {
		t.Errorf("Not stopping even after exceeding limit.\nInstance: %v", instance)
	}
	redisClient.Del("unique")

	rt = window.New(2, time.Second*2, redisClient)
	report, _ := rt.Use("unique")
	if report.APIUsageLeft != 1 {
		t.Errorf("Expected %v usage left, got %v.\n", 1, report.APIUsageLeft)
	}

	time.Sleep(time.Second * 3)
	report, _ = rt.Status("unique")
	if report.APIUsageLeft != 2 {
		t.Errorf("It must reset usage after each given window time. Expected usage left %v, got %v", 2, report.APIUsageLeft)
	}
	redisClient.Del("unique")
}
