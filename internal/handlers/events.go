package handlers

import (
	"fmt"
	"net/http"
	"time"

	"task-manager-api/internal/middleware"
	"task-manager-api/internal/realtime"

	"github.com/gin-gonic/gin"
)

type EventsHandler struct {
	hub *realtime.Hub
}

func NewEventsHandler(hub *realtime.Hub) *EventsHandler {
	return &EventsHandler{hub: hub}
}

// Stream serves task change events over Server-Sent Events. The browser's
// EventSource authenticates via the httpOnly auth cookie.
func (h *EventsHandler) Stream(c *gin.Context) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Streaming not supported")
		return
	}

	events, unsubscribe := h.hub.Subscribe(middleware.CurrentUserID(c), middleware.IsAdmin(c))
	defer unsubscribe()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	fmt.Fprint(c.Writer, ": connected\n\n")
	flusher.Flush()

	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case payload, open := <-events:
			if !open {
				return
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", payload)
			flusher.Flush()
		case <-heartbeat.C:
			fmt.Fprint(c.Writer, ": ping\n\n")
			flusher.Flush()
		}
	}
}
