package consumer

import (
	"bufio"
	"channel-test/internal/store"
	"channel-test/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// SSEConsumer consumes Server-Sent Events from the test scores endpoint
type SSEConsumer struct {
	url    string
	store  store.Store
	client *http.Client
}

// NewSSEConsumer creates a new SSE consumer
func NewSSEConsumer(url string, store store.Store) *SSEConsumer {
	return &SSEConsumer{
		url:   url,
		store: store,
		client: &http.Client{
			Timeout: 0, // No timeout for SSE connections
		},
	}
}

// Start begins consuming events from the SSE endpoint
func (c *SSEConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("SSE consumer shutting down")
			return ctx.Err()
		default:
			if err := c.connect(ctx); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				log.Printf("Connection error: %v. Reconnecting in 5 seconds...", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (c *SSEConsumer) connect(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Println("Connected to SSE endpoint")

	return c.readEvents(ctx, resp.Body)
}

func (c *SSEConsumer) readEvents(ctx context.Context, body io.ReadCloser) error {
	reader := bufio.NewReader(body)
	var eventType string
	var data string

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("connection closed")
			}
			return fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)

		// Empty line indicates end of event
		if line == "" {
			if eventType == "score" && data != "" {
				c.processScoreEvent(data)
			}
			eventType = ""
			data = ""
			continue
		}

		// Parse SSE field
		if strings.HasPrefix(line, "event:") {
			eventType = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			data = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}
}

func (c *SSEConsumer) processScoreEvent(data string) {
	var event models.ScoreEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		log.Printf("Failed to parse score event: %v", err)
		return
	}

	// Validate the event
	if event.StudentID == "" {
		log.Printf("Invalid event: missing student ID")
		return
	}

	if event.Score < 0 || event.Score > 1 {
		log.Printf("Invalid event: score out of range [0,1]: %f", event.Score)
		return
	}

	if err := c.store.AddScore(event); err != nil {
		log.Printf("Failed to store score: %v", err)
		return
	}

	log.Printf("Stored score: student=%s, exam=%d, score=%.3f",
		event.StudentID, event.Exam, event.Score)
}
