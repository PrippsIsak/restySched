package n8n

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/isak/restySched/internal/domain"
)

// Client defines the interface for n8n webhook client
type Client interface {
	SendSchedule(ctx context.Context, payload domain.N8NSchedulePayload) error
}

type client struct {
	webhookURL string
	httpClient *http.Client
}

// NewClient creates a new n8n webhook client
func NewClient(webhookURL string) Client {
	return &client{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendSchedule sends schedule data to n8n webhook
func (c *client) SendSchedule(ctx context.Context, payload domain.N8NSchedulePayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("n8n webhook returned non-success status: %d", resp.StatusCode)
	}

	return nil
}
