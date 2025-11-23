package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
)

type WebhookEvent struct {
	Event     string      `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type WebhookService struct {
	endpoints []string
	client    *http.Client
}

func NewWebhookService(endpoints []string) *WebhookService {
	return &WebhookService{
		endpoints: endpoints,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendWebhook sends a webhook to all registered endpoints
func (s *WebhookService) SendWebhook(event string, data interface{}) error {
	payload := WebhookEvent{
		Event:     event,
		Timestamp: time.Now(),
		Data:      data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	for _, endpoint := range s.endpoints {
		go s.sendToEndpoint(endpoint, jsonData)
	}

	return nil
}

func (s *WebhookService) sendToEndpoint(endpoint string, payload []byte) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to create webhook request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Community-Issue-Mapper-Webhook/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send webhook to %s: %v\n", endpoint, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("Webhook endpoint %s returned error: %d\n", endpoint, resp.StatusCode)
	}
}

// Convenience methods for different events
func (s *WebhookService) IssueCreated(issue *models.Issue) error {
	return s.SendWebhook("issue.created", issue)
}

func (s *WebhookService) IssueUpdated(issue *models.Issue) error {
	return s.SendWebhook("issue.updated", issue)
}

func (s *WebhookService) IssueResolved(issue *models.Issue) error {
	return s.SendWebhook("issue.resolved", issue)
}

func (s *WebhookService) CommentAdded(comment *models.Comment) error {
	return s.SendWebhook("comment.added", comment)
}
