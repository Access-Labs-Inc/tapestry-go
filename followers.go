package tapestry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FollowRequest struct {
	StartID string `json:"startId"`
	EndID   string `json:"endId"`
}

func (c *TapestryClient) AddFollower(startID, endID string) error {
	url := fmt.Sprintf("%s/followers/add?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	jsonBody, err := json.Marshal(FollowRequest{
		StartID: startID,
		EndID:   endID,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *TapestryClient) RemoveFollower(startID, endID string) error {
	url := fmt.Sprintf("%s/followers/remove?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	jsonBody, err := json.Marshal(FollowRequest{
		StartID: startID,
		EndID:   endID,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
