package tapestry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Profile struct {
	Namespace  string `json:"namespace"`
	ID         string `json:"id"`
	Blockchain string `json:"blockchain"`
	Username   string `json:"username"`
}

type ProfileResponse struct {
	Profile       Profile `json:"profile"`
	WalletAddress string  `json:"walletAddress"`
}

type FindOrCreateProfileParameters struct {
	WalletAddress string `json:"walletAddress"`
	Username      string `json:"username"`
	Bio           string `json:"bio,omitempty"`
	Image         string `json:"image,omitempty"`
	ID            string `json:"id,omitempty"`
}

type FindOrCreateProfileRequest struct {
	FindOrCreateProfileParameters
	Execution  string `json:"execution,omitempty"`
	Blockchain string `json:"blockchain,omitempty"`
}

type UpdateProfileParameters struct {
	Username string `json:"username"`
	Bio      string `json:"bio,omitempty"`
	Image    string `json:"image,omitempty"`
}

type UpdateProfileRequest struct {
	UpdateProfileParameters
	Execution string `json:"execution,omitempty"`
}

func (c *TapestryClient) FindOrCreateProfile(params FindOrCreateProfileParameters) (*ProfileResponse, error) {
	req := FindOrCreateProfileRequest{
		FindOrCreateProfileParameters: params,
		Execution:                     string(c.execution),
		Blockchain:                    c.blockchain,
	}

	url := fmt.Sprintf("%s/profiles/findOrCreate?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var profileResp ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profileResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &profileResp, nil
}

func (c *TapestryClient) UpdateProfile(id string, req UpdateProfileParameters) error {
	url := fmt.Sprintf("%s/profiles/%s?apiKey=%s", c.tapestryApiBaseUrl, id, c.apiKey)

	jsonBody, err := json.Marshal(UpdateProfileRequest{
		UpdateProfileParameters: req,
		Execution:               string(c.execution),
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *TapestryClient) GetProfileByID(id string) (*ProfileResponse, error) {
	url := fmt.Sprintf("%s/profiles/%s?apiKey=%s", c.tapestryApiBaseUrl, id, c.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil // ok but not found
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var profileResp ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profileResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &profileResp, nil
}
