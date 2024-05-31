package receiver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vinr-eu/naadi/internal/domain/event"
)

func ReceiveEvents(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		http.Error(w, "Service name missing", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "1" // default limit to 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	event.Mutex.Lock()
	defer event.Mutex.Unlock()

	// Get the current offset (number of events processed) for the service
	offset, exists := event.ServiceOffsets[serviceName]
	if !exists {
		offset = 0 // Initialize to 0 for the first time
	}

	events := make([]event.Event, 0)
	count := 0

	// Iterate through the events starting from the offset
	for id := offset; count < limit; id++ {
		if e, ok := event.Store[id]; ok {
			events = append(events, e)
			count++
		} else {
			break // No more events available
		}
	}

	// Update the offset for the service
	if len(events) > 0 {
		event.ServiceOffsets[serviceName] = offset + int64(len(events))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(events)
}
