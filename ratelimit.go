package ratelimit

import (
	"time"
)

type rules struct {
	maxAPICalls int
	windowSize  time.Duration
}

// Storage is profile of API usage left and next time when usage will be reset
type Storage struct {
	APIUsageLeft int
	NextWindow   time.Time
}

func (s *Storage) createNewWindow(rule rules) {
	s.APIUsageLeft = rule.maxAPICalls
	s.NextWindow = time.Now().Add(rule.windowSize)
}

func (s *Storage) isValid(rule rules) bool {
	if time.Now().After(s.NextWindow) {
		s.createNewWindow(rule)
		return true
	}

	return s.APIUsageLeft > 0
}

func (s *Storage) useAPI() {
	s.APIUsageLeft--
}

type ratelimit interface {
	New(int, time.Time) Memory
}

// Memory is used to store API usage and rules
type Memory struct {
	memory map[string](*Storage)
	rules
}

// New will create new memory instance
func New(maxAPICalls int, windowSize time.Duration) Memory {
	var newMemory Memory
	newMemory.maxAPICalls = maxAPICalls
	newMemory.windowSize = windowSize
	newMemory.memory = make(map[string](*Storage))
	return newMemory
}

// Use - will use token for uniqueIdentified and if no tokens are left then it will return with false.
func (m *Memory) Use(uniqueIdentifier string) (Storage, bool) {
	storageInstance, existInMemory := m.memory[uniqueIdentifier]
	if !existInMemory {
		newStorage := Storage{}
		newStorage.createNewWindow(m.rules)
		m.memory[uniqueIdentifier] = &newStorage
		newStorage.useAPI()
		return newStorage, true
	}

	if storageInstance.isValid(m.rules) {
		storageInstance.useAPI()
		return *storageInstance, true
	}

	return Storage{}, false
}

// Status is used to get number of API usage left & next time to reset APIs
func (m *Memory) Status(uniqueIdentifier string) (Storage, bool) {
	instance, ok := m.memory[uniqueIdentifier]
	if !ok {
		return Storage{}, false
	}

	instance.isValid(m.rules)
	return *instance, ok
}
