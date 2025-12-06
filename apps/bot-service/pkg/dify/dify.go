package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
)

type Request struct {
	Inputs         map[string]any `json:"inputs"`
	Query          string         `json:"query"`
	ResponseMode   string         `json:"response_mode"` // "streaming" or "blocking"
	ConversationID string         `json:"conversation_id"`
	User           string         `json:"user"`
	Files          []any          `json:"files"`
}

type Response struct {
	Event          string         `json:"event"`
	TaskID         string         `json:"task_id"`
	ID             string         `json:"id"`
	MessageID      string         `json:"message_id"`
	ConversationID string         `json:"conversation_id"`
	Mode           string         `json:"mode"`
	Answer         string         `json:"answer"`
	Metadata       map[string]any `json:"metadata"`
	CreatedAt      int64          `json:"created_at"`
}

type CustomDifyInterface interface {
	ChatMessages(ctx context.Context, req *Request) (*Response, error)
}

type CustomDifyStruct struct {
	DifyAPIURL string
	DifyAPIKey string
}

func getDify() CustomDifyInterface {
	difyAPIURL := env.AppEnv.DifyAPIURL
	difyAPIKey := env.AppEnv.DifyAPIKey

	return &CustomDifyStruct{
		DifyAPIURL: difyAPIURL,
		DifyAPIKey: difyAPIKey,
	}
}

var Dify = getDify()

// ChatMessages sends a chat message request to the Dify API and returns the response.
func (o *CustomDifyStruct) ChatMessages(ctx context.Context, req *Request) (*Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.DifyAPIURL+"/chat-messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.DifyAPIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var res Response
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &res, nil
}
