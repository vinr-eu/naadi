package event

import (
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
	Mutex           sync.Mutex
)
