package notifier

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vinr-eu/naadi/internal/domain/event"
)

func NotifyEvent(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header missing", http.StatusBadRequest)
		return
	}

	var entity event.Event
	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	now := time.Now()
	entity.NotificationTime = &now

	event.Mutex.Lock()
	defer event.Mutex.Unlock()

	// Check for idempotency using the IdempotencyKeys map
	if _, exists := event.IdempotencyKeys[idempotencyKey]; exists {
		http.Error(w, "Duplicate event", http.StatusConflict)
		return
	}

	// Generate a new ID by finding the highest ID and incrementing it
	newOffset := int64(len(event.Store))

	// Store the new event and idempotency key
	event.Store[newOffset] = entity
	event.IdempotencyKeys[idempotencyKey] = true

	// Save store to file
	if err := event.SaveStore(); err != nil {
		http.Error(w, "Failed to save store", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
