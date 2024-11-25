package tapestry

import (
	"context"
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

type GetFollowersResponse struct {
	Profiles []ProfileDetails `json:"profiles"`
}

type GetFollowingResponse struct {
	Profiles []ProfileDetails `json:"profiles"`
}

type ProfileDetails struct {
	ID        string        `json:"id"`
	Username  string        `json:"username"`
	Bio       string        `json:"bio,omitempty"`
	Image     string        `json:"image,omitempty"`
	CreatedAt UnixTimestamp `json:"created_at"`
}

type GetFollowingWhoFollowResponse struct {
	Profiles []ProfileDetails `json:"profiles"`
}

type SuggestedProfileValue struct {
	Namespaces []struct {
		Name         string `json:"name"`
		ReadableName string `json:"readableName,omitempty"`
		FaviconURL   string `json:"faviconURL,omitempty"`
	} `json:"namespaces"`
	Profile ProfileDetails `json:"profile"`
	Wallet  struct {
		Address string `json:"address"`
	} `json:"wallet"`
}

type GetSuggestedProfilesResponse struct {
	Profiles map[string]SuggestedProfileValue
}

func (c *TapestryClient) FindOrCreateProfile(ctx context.Context, params FindOrCreateProfileParameters) (*ProfileResponse, error) {
	reqData := FindOrCreateProfileRequest{
		FindOrCreateProfileParameters: params,
		Execution:                     string(c.execution),
		Blockchain:                    c.blockchain,
	}

	url := fmt.Sprintf("%s/profiles/findOrCreate?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	jsonBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func (c *TapestryClient) UpdateProfile(ctx context.Context, id string, reqData UpdateProfileParameters) error {
	url := fmt.Sprintf("%s/profiles/%s?apiKey=%s", c.tapestryApiBaseUrl, id, c.apiKey)

	jsonBody, err := json.Marshal(UpdateProfileRequest{
		UpdateProfileParameters: reqData,
		Execution:               string(c.execution),
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonBody)))
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

func (c *TapestryClient) GetProfileByID(ctx context.Context, id string) (*ProfileResponse, error) {
	url := fmt.Sprintf("%s/profiles/%s?apiKey=%s", c.tapestryApiBaseUrl, id, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
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

func (c *TapestryClient) GetFollowers(ctx context.Context, profileID string) (*GetFollowersResponse, error) {
	url := fmt.Sprintf("%s/profiles/%s/followers?apiKey=%s", c.tapestryApiBaseUrl, profileID, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var followersResp GetFollowersResponse
	if err := json.NewDecoder(resp.Body).Decode(&followersResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &followersResp, nil
}

func (c *TapestryClient) GetFollowing(ctx context.Context, profileID string) (*GetFollowingResponse, error) {
	url := fmt.Sprintf("%s/profiles/%s/following?apiKey=%s", c.tapestryApiBaseUrl, profileID, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var followingResp GetFollowingResponse
	if err := json.NewDecoder(resp.Body).Decode(&followingResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &followingResp, nil
}

func (c *TapestryClient) GetFollowingWhoFollow(ctx context.Context, profileID string, requestorID string) (*GetFollowingWhoFollowResponse, error) {
	url := fmt.Sprintf("%s/profiles/%s/following-who-follow?apiKey=%s&requestorId=%s",
		c.tapestryApiBaseUrl, profileID, c.apiKey, requestorID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var followingWhoFollowResp GetFollowingWhoFollowResponse
	if err := json.NewDecoder(resp.Body).Decode(&followingWhoFollowResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &followingWhoFollowResp, nil
}

func (c *TapestryClient) GetSuggestedProfiles(ctx context.Context, address string, ownAppOnly bool) (*GetSuggestedProfilesResponse, error) {
	url := fmt.Sprintf("%s/profiles/suggested/%s?apiKey=%s&ownAppOnly=%t",
		c.tapestryApiBaseUrl, address, c.apiKey, ownAppOnly)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawResponse map[string]SuggestedProfileValue
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &GetSuggestedProfilesResponse{Profiles: rawResponse}, nil
}
