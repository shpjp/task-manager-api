package realtime

import (
	"encoding/json"
	"sync"
)

// Event is a server-sent event describing a task change.
type Event struct {
	Type   string `json:"type"` // task.created | task.updated | task.deleted
	TaskID uint   `json:"task_id"`
	Task   any    `json:"task,omitempty"`
}

type subscriber struct {
	userID uint
	admin  bool
	ch     chan []byte
}

// Hub fans task events out to SSE subscribers. Events for a task are sent to
// the task owner's connections and to all admin connections.
type Hub struct {
	mu   sync.RWMutex
	subs map[*subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{subs: make(map[*subscriber]struct{})}
}

// Subscribe registers a connection and returns the event channel plus an
// unsubscribe function.
func (h *Hub) Subscribe(userID uint, admin bool) (<-chan []byte, func()) {
	sub := &subscriber{userID: userID, admin: admin, ch: make(chan []byte, 16)}

	h.mu.Lock()
	h.subs[sub] = struct{}{}
	h.mu.Unlock()

	return sub.ch, func() {
		h.mu.Lock()
		delete(h.subs, sub)
		h.mu.Unlock()
		close(sub.ch)
	}
}

// Publish delivers an event about ownerID's task to relevant subscribers.
// Slow subscribers are skipped rather than blocking the publisher.
func (h *Hub) Publish(ownerID uint, event Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for sub := range h.subs {
		if sub.userID != ownerID && !sub.admin {
			continue
		}
		select {
		case sub.ch <- payload:
		default:
		}
	}
}
