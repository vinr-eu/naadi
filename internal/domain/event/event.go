package event

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Event struct {
	Name             string     `json:"name"`
	EntityID         string     `json:"entityId"`
	NotificationTime *time.Time `json:"notificationTime"`
}

// Exported in-memory storage for events
var (
	Store           = make(map[int64]Event)
	IdempotencyKeys = make(map[string]bool)
	ServiceOffsets  = make(map[string]int64)
	storeFile       = "event_store.json"
	offsetFile      = "service_offsets.json"
	Mutex           sync.Mutex
)

// SaveStore saves the event store to a file.
func SaveStore() error {
	file, err := os.Create(storeFile)
	if err != nil {
		return fmt.Errorf("failed to create store file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(Store); err != nil {
		return fmt.Errorf("failed to encode store: %v", err)
	}
	return nil
}

func Init() error {
	storeErr := LoadStore()
	if storeErr != nil {
		return storeErr
	}
	offsetErr := LoadOffsets()
	if offsetErr != nil {
		return offsetErr
	}
	return nil
}

// LoadStore loads the event store from a file.
func LoadStore() error {
	file, err := os.Open(storeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing store file, nothing to load
		}
		return fmt.Errorf("failed to open store file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&Store); err != nil {
		return fmt.Errorf("failed to decode store: %v", err)
	}
	return nil
}

// SaveOffsets saves the service offsets to a file.
func SaveOffsets() error {
	file, err := os.Create(offsetFile)
	if err != nil {
		return fmt.Errorf("failed to create offset file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(ServiceOffsets); err != nil {
		return fmt.Errorf("failed to encode offsets: %v", err)
	}
	return nil
}

// LoadOffsets loads the service offsets from a file.
func LoadOffsets() error {
	file, err := os.Open(offsetFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing offset file, nothing to load
		}
		return fmt.Errorf("failed to open offset file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&ServiceOffsets); err != nil {
		return fmt.Errorf("failed to decode offsets: %v", err)
	}
	return nil
}
