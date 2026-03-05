package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Client is the HTTP client for OpenDiscuz API
type Client struct {
	BaseURL     string
	AccessToken string
	HTTPClient  *http.Client
}

// Response wraps the standard API response
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// NewClient creates a new API client
func NewClient(baseURL, accessToken string) *Client {
	return &Client{
		BaseURL:     strings.TrimRight(baseURL, "/"),
		AccessToken: accessToken,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) do(method, path string, body interface{}) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AccessToken != "" {
		req.Header.Set("Cookie", "access_token="+c.AccessToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response (status=%d): %s", resp.StatusCode, string(data))
	}

	if apiResp.Code != 0 {
		return &apiResp, fmt.Errorf("API error: %s (code=%d)", apiResp.Message, apiResp.Code)
	}
	return &apiResp, nil
}

// GET performs a GET request
func (c *Client) GET(path string) (*Response, error) {
	return c.do("GET", path, nil)
}

// POST performs a POST request
func (c *Client) POST(path string, body interface{}) (*Response, error) {
	return c.do("POST", path, body)
}

// PUT performs a PUT request
func (c *Client) PUT(path string, body interface{}) (*Response, error) {
	return c.do("PUT", path, body)
}

// DELETE performs a DELETE request
func (c *Client) DELETE(path string) (*Response, error) {
	return c.do("DELETE", path, nil)
}

// UploadFile uploads a file via multipart form
func (c *Client) UploadFile(path, filePath string) (*Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)
	writer.Close()

	url := c.BaseURL + path
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.AccessToken != "" {
		req.Header.Set("Cookie", "access_token="+c.AccessToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var apiResp Response
	json.Unmarshal(data, &apiResp)
	if apiResp.Code != 0 {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Message)
	}
	return &apiResp, nil
}

// DataJSON returns pretty-printed JSON of the response data
func (r *Response) DataJSON() string {
	var out bytes.Buffer
	json.Indent(&out, r.Data, "", "  ")
	return out.String()
}
