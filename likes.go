package tapestry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CreateLikeRequest struct {
	StartId   string `json:"startId"`
	Execution string `json:"execution"`
}
type DeleteLikeRequest struct {
	StartId string `json:"startId"`
}

func (c *TapestryClient) CreateLike(contentID string, profile Profile) error {
	url := fmt.Sprintf("%s/likes/%s?apiKey=%s&username=%s", c.tapestryApiBaseUrl, contentID, c.apiKey, url.QueryEscape(profile.Username))

	reqBody, err := json.Marshal(CreateLikeRequest{StartId: profile.ID, Execution: "FAST_UNCONFIRMED"})
	if err != nil {
		return fmt.Errorf("error encoding body: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *TapestryClient) DeleteLike(contentID string, profile Profile) error {
	url := fmt.Sprintf("%s/likes/%s?apiKey=%s&username=%s", c.tapestryApiBaseUrl, contentID, c.apiKey, url.QueryEscape(profile.Username))

	reqBody, err := json.Marshal(DeleteLikeRequest{StartId: profile.ID})
	if err != nil {
		return fmt.Errorf("error encoding body: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(reqBody))
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
