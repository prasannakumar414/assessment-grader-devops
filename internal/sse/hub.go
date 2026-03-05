package sse

import (
	"encoding/json"
	"sync"
)

type Event struct {
	Type string
	Data any
}

func (e Event) JSON() []byte {
	b, _ := json.Marshal(e.Data)
	return b
}

type Hub struct {
	mu          sync.RWMutex
	subscribers map[chan Event]struct{}
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[chan Event]struct{}),
	}
}

func (h *Hub) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 64)

	h.mu.Lock()
	h.subscribers[ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		delete(h.subscribers, ch)
		h.mu.Unlock()
		close(ch)
	}

	return ch, unsubscribe
}

func (h *Hub) Broadcast(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.subscribers {
		select {
		case ch <- event:
		default:
			// drop if subscriber is too slow
		}
	}
}
