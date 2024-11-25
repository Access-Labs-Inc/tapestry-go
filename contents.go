package tapestry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ContentProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type FindOrCreateContentRequest struct {
	ProfileID  string            `json:"profileId"`
	ID         string            `json:"id,omitempty"`
	Properties []ContentProperty `json:"properties"`
}

type UpdateContentRequest struct {
	Properties []ContentProperty `json:"properties"`
}

type Content struct {
	Namespace   string        `json:"namespace"`
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Title       string        `json:"title"`
	CreatedAt   UnixTimestamp `json:"created_at"`
}

type CreateOrUpdateContentResponse struct{ Content }
type GetContentResponse struct {
	Content      Content      `json:"content"`
	SocialCounts SocialCounts `json:"socialCounts"`
}

type AuthorProfile struct {
	ID        string        `json:"id"`
	Username  string        `json:"username"`
	Bio       string        `json:"bio"`
	Image     string        `json:"image"`
	CreatedAt UnixTimestamp `json:"created_at"`
}

type RequestingProfileSocialInfo struct {
	HasLiked bool `json:"hasLiked"`
}

type ContentListItem struct {
	AuthorProfile               AuthorProfile               `json:"authorProfile"`
	Content                     Content                     `json:"content"`
	SocialCounts                SocialCounts                `json:"socialCounts"`
	RequestingProfileSocialInfo RequestingProfileSocialInfo `json:"requestingProfileSocialInfo"`
}

type GetContentsResponse struct {
	Contents []ContentListItem `json:"contents"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

type BatchResponseContentListItem struct {
	Content      Content      `json:"content"`
	SocialCounts SocialCounts `json:"socialCounts"`
}

type GetContentsByBatchIDsResponse struct {
	Successful []BatchResponseContentListItem `json:"successful"`
	Failed     []struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	} `json:"failed"`
}

type SocialCounts struct {
	LikeCount    int `json:"likeCount"`
	CommentCount int `json:"commentCount"`
}

func (c *TapestryClient) FindOrCreateContent(ctx context.Context, profileId, id string, properties []ContentProperty) (*CreateOrUpdateContentResponse, error) {
	url := fmt.Sprintf("%s/contents/findOrCreate?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	jsonBody, err := json.Marshal(FindOrCreateContentRequest{
		ProfileID:  profileId,
		ID:         id,
		Properties: properties,
	})
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

	var contentResp CreateOrUpdateContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &contentResp, nil
}

func (c *TapestryClient) UpdateContent(ctx context.Context, contentId string, properties []ContentProperty) (*CreateOrUpdateContentResponse, error) {
	url := fmt.Sprintf("%s/contents/%s?apiKey=%s", c.tapestryApiBaseUrl, contentId, c.apiKey)

	jsonBody, err := json.Marshal(UpdateContentRequest{
		Properties: properties,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonBody)))
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

	var contentResp CreateOrUpdateContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &contentResp, nil
}

func (c *TapestryClient) DeleteContent(ctx context.Context, contentId string) error {
	url := fmt.Sprintf("%s/contents/%s?apiKey=%s", c.tapestryApiBaseUrl, contentId, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

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

func (c *TapestryClient) GetContentByID(ctx context.Context, contentId string) (*GetContentResponse, error) {
	url := fmt.Sprintf("%s/contents/%s?apiKey=%s", c.tapestryApiBaseUrl, contentId, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// unauthorized may also mean not found
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusUnauthorized {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var contentResp GetContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// not found
	if contentResp.Content.ID == "" {
		return nil, nil
	}

	return &contentResp, nil
}

func (c *TapestryClient) GetContentsByBatchIDs(ctx context.Context, batchIDs []string) (*GetContentsByBatchIDsResponse, error) {
	url := fmt.Sprintf("%s/contents/batch/read?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	jsonBody, err := json.Marshal(batchIDs)
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

	var batchResp GetContentsByBatchIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &batchResp, nil
}

func (c *TapestryClient) GetContents(ctx context.Context, opts ...GetContentsOption) (*GetContentsResponse, error) {
	params := &getContentsParams{
		apiKey: c.apiKey,
	}

	for _, opt := range opts {
		opt(params)
	}

	url := fmt.Sprintf("%s/contents/?%s", c.tapestryApiBaseUrl, params.encode())

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

	var contentsResp GetContentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentsResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &contentsResp, nil
}

type GetContentsSortDirection string

const (
	GetContentsSortDirectionAsc  GetContentsSortDirection = "ASC"
	GetContentsSortDirectionDesc GetContentsSortDirection = "DESC"
)

type getContentsParams struct {
	apiKey              string
	orderByField        string
	orderByDirection    GetContentsSortDirection
	page                string
	pageSize            string
	profileID           string
	requestingProfileID string
}

func (p *getContentsParams) encode() string {
	values := make([]string, 0)
	values = append(values, fmt.Sprintf("apiKey=%s", p.apiKey))

	if p.orderByField != "" {
		values = append(values, fmt.Sprintf("orderByField=%s", p.orderByField))
	}
	if p.orderByDirection != "" {
		values = append(values, fmt.Sprintf("orderByDirection=%s", p.orderByDirection))
	}
	if p.page != "" {
		values = append(values, fmt.Sprintf("page=%s", p.page))
	}
	if p.pageSize != "" {
		values = append(values, fmt.Sprintf("pageSize=%s", p.pageSize))
	}
	if p.profileID != "" {
		values = append(values, fmt.Sprintf("profileId=%s", p.profileID))
	}
	if p.requestingProfileID != "" {
		values = append(values, fmt.Sprintf("requestingProfileId=%s", p.requestingProfileID))
	}

	return strings.Join(values, "&")
}

// Option pattern for configurable parameters
type GetContentsOption func(*getContentsParams)

func WithOrderBy(field string, direction GetContentsSortDirection) GetContentsOption {
	return func(p *getContentsParams) {
		p.orderByField = field
		p.orderByDirection = direction
	}
}

func WithPagination(page, pageSize string) GetContentsOption {
	return func(p *getContentsParams) {
		p.page = page
		p.pageSize = pageSize
	}
}

func WithProfileID(profileID string) GetContentsOption {
	return func(p *getContentsParams) {
		p.profileID = profileID
	}
}

func WithRequestingProfileID(requestingProfileID string) GetContentsOption {
	return func(p *getContentsParams) {
		p.requestingProfileID = requestingProfileID
	}
}
