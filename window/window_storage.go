package window

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type rules struct {
	maxAPICalls int
	windowSize  time.Duration
}

// Storage is profile of API usage left and next time when usage will be reset
type Storage struct {
	APIUsageLeft int       `json:"api_usage_left"`
	NextWindow   time.Time `json:"next_window"`
}

func (s *Storage) createNewWindow(rule rules) {
	s.APIUsageLeft = rule.maxAPICalls
	s.NextWindow = time.Now().Add(rule.windowSize)
}

func (s *Storage) isValid(rule rules) (bool, bool) {
	if time.Now().After(s.NextWindow) {
		s.createNewWindow(rule)
		return true, true
	}

	return s.APIUsageLeft > 0, false
}

func (s *Storage) useAPI() {
	s.APIUsageLeft--
}

type ratelimit interface {
	New(int, time.Time) Memory
}

// Memory is used to store API usage and rules
type Memory struct {
	redisClient *redis.Client
	rules
}

// New will create new memory instance
func New(maxAPICalls int, windowSize time.Duration, redisClient *redis.Client) Memory {
	var newMemory Memory
	newMemory.maxAPICalls = maxAPICalls
	newMemory.windowSize = windowSize
	newMemory.redisClient = redisClient
	return newMemory
}

func (m *Memory) getStorage(identifier string) (Storage, bool) {
	var storage Storage
	serializedStorage, err := m.redisClient.Get(ctx, identifier).Result()
	json.Unmarshal([]byte(serializedStorage), &storage)
	if err != nil {
		return storage, false
	}
	return storage, true
}

func (m *Memory) setStorage(identifier string, storage Storage) bool {
	serializedStorage, _ := json.Marshal(storage)
	err := m.redisClient.Set(ctx, identifier, string(serializedStorage), 0).Err()
	if err != nil {
		return false
	}
	return true
}

// Use - will use token for uniqueIdentified and if no tokens are left then it will return with false.
func (m *Memory) Use(uniqueIdentifier string) (Storage, bool) {
	storageInstance, existInMemory := m.getStorage(uniqueIdentifier)
	if !existInMemory {
		newStorage := Storage{}
		newStorage.createNewWindow(m.rules)
		newStorage.useAPI()
		m.setStorage(uniqueIdentifier, newStorage)
		return newStorage, true
	}

	if valid, _ := storageInstance.isValid(m.rules); valid {
		storageInstance.useAPI()
		m.setStorage(uniqueIdentifier, storageInstance)
		return storageInstance, true
	}

	return storageInstance, false
}

// Status is used to get number of API usage left & next time to reset APIs
func (m *Memory) Status(uniqueIdentifier string) (Storage, bool) {
	instance, ok := m.getStorage(uniqueIdentifier)
	if !ok {
		return instance, false
	}

	_, changed := instance.isValid(m.rules)
	if changed {
		m.setStorage(uniqueIdentifier, instance)
	}
	return instance, ok
}
