package handlers

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"docker-workshop-assesment-grader/internal/sse"
)

type EventsHandler struct {
	Hub *sse.Hub
}

func (h *EventsHandler) Stream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Flush headers + initial comment so EventSource.onopen fires immediately
	c.Writer.WriteString(": connected\n\n")
	c.Writer.Flush()

	events, unsubscribe := h.Hub.Subscribe()
	defer unsubscribe()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case evt, ok := <-events:
			if !ok {
				return false
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, evt.JSON())
			return true
		}
	})
}
