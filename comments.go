package tapestry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type CommentProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type CreateCommentOptions struct {
	ContentID  string            `json:"contentId"`
	ProfileID  string            `json:"profileId"`
	Text       string            `json:"text"`
	CommentID  string            `json:"commentId"`
	Properties []CommentProperty `json:"properties"`
}

type CreateCommentRequest struct {
	CreateCommentOptions
	Execution string `json:"execution"`
}

type CreateCommentResponse struct {
	Comment
}

type GetCommentByIdResponse struct {
	CommentData
}

type GetCommentsResponse struct {
	Comments []CommentData `json:"comments"`
}

type CommentData struct {
	Comment                     Comment                `json:"comment"`
	Author                      Author                 `json:"author"`
	SocialCounts                SocialCounts           `json:"socialCounts"`
	RequestingProfileSocialInfo map[string]interface{} `json:"requestingProfileSocialInfo"`
}

type Comment struct {
	Namespace string        `json:"namespace"`
	CreatedAt UnixTimestamp `json:"created_at"`
	Text      string        `json:"text"`
	ID        string        `json:"id"`
}

type Author struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
}

type GetCommentsOptions struct {
	ContentID           string
	CommentID           string
	ProfileID           string
	RequestingProfileID string
	Page                int
	PageSize            int
}

type UpdateCommentRequest struct {
	Properties []CommentProperty `json:"properties"`
}

type UpdateCommentResponse struct {
	Comment
}

func (c *TapestryClient) CreateComment(ctx context.Context, options CreateCommentOptions) (*CreateCommentResponse, error) {
	uri := fmt.Sprintf("%s/comments?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)
	if options.Properties == nil {
		options.Properties = []CommentProperty{}
	}
	reqData := CreateCommentRequest{
		CreateCommentOptions: options,
		Execution:            string(c.execution),
	}
	if options.CommentID != "" {
		reqData.CommentID = options.CommentID
	}

	jsonBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(string(jsonBody)))
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
		// read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var commentResp CreateCommentResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &commentResp, nil
}

func (c *TapestryClient) GetComments(ctx context.Context, options GetCommentsOptions) (*GetCommentsResponse, error) {
	baseURL := fmt.Sprintf("%s/comments?apiKey=%s", c.tapestryApiBaseUrl, c.apiKey)

	params := url.Values{}
	if options.ContentID != "" {
		params.Add("contentId", options.ContentID)
	}
	if options.CommentID != "" {
		params.Add("commentId", options.CommentID)
	}
	if options.ProfileID != "" {
		params.Add("profileId", options.ProfileID)
	}
	if options.RequestingProfileID != "" {
		params.Add("requestingProfileId", options.RequestingProfileID)
	}

	if options.Page > 0 {
		params.Add("page", strconv.Itoa(options.Page))
	}
	if options.PageSize > 0 {
		params.Add("pageSize", strconv.Itoa(options.PageSize))
	}

	uri := baseURL
	if len(params) > 0 {
		uri += "&" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
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
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var comments GetCommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &comments, nil
}

func (c *TapestryClient) GetCommentByID(ctx context.Context, commentID string, requestingProfileID string) (*GetCommentByIdResponse, error) {
	uri := fmt.Sprintf("%s/comments/%s?apiKey=%s", c.tapestryApiBaseUrl, commentID, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
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

	var commentResp GetCommentByIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &commentResp, nil
}

func (c *TapestryClient) DeleteComment(ctx context.Context, commentID string) error {
	uri := fmt.Sprintf("%s/comments/%s?apiKey=%s", c.tapestryApiBaseUrl, commentID, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
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

func (c *TapestryClient) UpdateComment(ctx context.Context, commentID string, properties []CommentProperty) (*UpdateCommentResponse, error) {
	uri := fmt.Sprintf("%s/comments/%s?apiKey=%s", c.tapestryApiBaseUrl, commentID, c.apiKey)

	reqData := UpdateCommentRequest{
		Properties: properties,
	}

	jsonBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, strings.NewReader(string(jsonBody)))
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var commentResp UpdateCommentResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &commentResp, nil
}
