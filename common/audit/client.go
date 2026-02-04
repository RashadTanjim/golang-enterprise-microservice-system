package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"enterprise-microservice-system/common/logger"

	"go.uber.org/zap"
)

// Config holds configuration for the audit log client.
type Config struct {
	Enabled bool
	BaseURL string
	Timeout time.Duration
}

// Event represents an audit log event to record.
type Event struct {
	Actor        string `json:"actor"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Description  string `json:"description,omitempty"`
	Metadata     string `json:"metadata,omitempty"`
}

// Client sends audit log events to the audit-log-service.
type Client struct {
	enabled bool
	baseURL string
	client  *http.Client
	logger  *logger.Logger
}

// NewClient creates a new audit log client.
func NewClient(cfg Config, log *logger.Logger) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 3 * time.Second
	}

	return &Client{
		enabled: cfg.Enabled,
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		client: &http.Client{
			Timeout: timeout,
		},
		logger: log,
	}
}

// Track sends the event to audit-log-service using the provided bearer token.
// This is best-effort and does not return errors to callers.
func (c *Client) Track(ctx context.Context, event Event, bearerToken string) {
	if c == nil || !c.enabled || c.baseURL == "" {
		return
	}

	token := normalizeBearerToken(bearerToken)
	if token == "" {
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		c.warn("failed to marshal audit log payload", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/audit-logs", bytes.NewReader(payload))
	if err != nil {
		c.warn("failed to create audit log request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		c.warn("failed to send audit log event", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		c.warn("audit log request returned non-2xx status", nil, zap.Int("status", resp.StatusCode))
	}
}

func normalizeBearerToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}

	return value
}

func (c *Client) warn(message string, err error, fields ...zap.Field) {
	if c.logger == nil {
		return
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	c.logger.Warn(message, fields...)
}
