package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Request struct {
	SourceFormat string `json:"source_format"`
	TargetFormat string `json:"target_format"`
	Payload      string `json:"payload"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

type fastAPIError struct {
	Detail any `json:"detail"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Transform(source, target, payload string) (string, error) {
	body, err := json.Marshal(Request{
		SourceFormat: source,
		TargetFormat: target,
		Payload:      payload,
	})
	if err != nil {
		return "", err
	}
	url := c.baseURL + "/transform"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("python service: %w", err)
	}
	raw, readErr := io.ReadAll(resp.Body)
	closeErr := resp.Body.Close()
	if readErr != nil {
		return "", readErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	if resp.StatusCode >= 400 {
		var fe fastAPIError
		if json.Unmarshal(raw, &fe) == nil && fe.Detail != nil {
			return "", fmt.Errorf("%v", fe.Detail)
		}
		return "", fmt.Errorf("transform failed: status %d", resp.StatusCode)
	}
	var out Response
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", err
	}
	if out.Error != "" {
		return "", fmt.Errorf("%s", out.Error)
	}
	return out.Result, nil
}
