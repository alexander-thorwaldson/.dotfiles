package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ICEClient talks to the ice prompt injection classifier daemon.
type ICEClient struct {
	Endpoint string
	HTTP     *http.Client
}

// NewICEClient creates an ICE client pointing at the given base URL.
func NewICEClient(endpoint string) *ICEClient {
	return &ICEClient{
		Endpoint: endpoint,
		HTTP:     http.DefaultClient,
	}
}

// Classification is the response from the ice /classify endpoint.
type Classification struct {
	Label  string             `json:"label"`
	Score  float64            `json:"score"`
	Scores map[string]float64 `json:"scores"`
}

// Classify sends text to ice for prompt injection classification.
func (c *ICEClient) Classify(text string) (*Classification, error) {
	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}
	resp, err := c.HTTP.Post(c.Endpoint+"/classify", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ice request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading ice response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ice returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result Classification
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding ice response: %w", err)
	}
	return &result, nil
}

// IsInjection returns true if the text is classified as a prompt injection.
func (c *ICEClient) IsInjection(text string) (bool, error) {
	result, err := c.Classify(text)
	if err != nil {
		return false, err
	}
	return result.Label == "INJECTION", nil
}
