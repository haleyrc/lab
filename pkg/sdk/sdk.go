package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// TODO (RCH): Change this to the production URL
const DefaultURL = `http://localhost:8080/api`

type Config struct {
	Debug bool
	URL   string
}

func NewClient(cfg Config) *Client {
	if cfg.URL == "" {
		cfg.URL = DefaultURL
	}
	return &Client{
		BaseURL: cfg.URL,
		Debug:   cfg.Debug,
		c:       &http.Client{Timeout: 5 * time.Second},
	}
}

type Client struct {
	BaseURL string
	Debug   bool

	bearerToken string
	c           *http.Client
}

type response struct {
	Data  json.RawMessage `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) post(ctx context.Context, path string, body interface{}) (json.RawMessage, error) {
	return c.request(ctx, http.MethodPost, c.BaseURL+path, body)
}

func (c *Client) request(ctx context.Context, method, path string, body interface{}) (json.RawMessage, error) {
	bb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(bb))
	if err != nil {
		return nil, err
	}

	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}

	if c.Debug {
		fmt.Printf("--> %s %s %t\n", req.Method, path, req.Header.Get("Authorization") != "")
		fmt.Println("   ", string(bb))
	}

	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		dumpResponseBody(res)
	}

	var resp response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Error.Message != "" {
		return nil, fmt.Errorf("request failed: %s", resp.Error.Message)
	}

	return resp.Data, nil
}

func dumpResponseBody(resp *http.Response) {
	body := getBody(resp)
	fmt.Printf("<-- %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	fmt.Printf("    %s\n", body)
}

func getBody(resp *http.Response) string {
	b, _ := httputil.DumpResponse(resp, true)
	lines := strings.Split(string(b), "\n")
	idx := findBlankLine(lines)
	if idx < 0 {
		return "malformed response"
	}

	return join(lines[idx:])
}

func findBlankLine(lines []string) int {
	for i := range lines {
		if strings.TrimSpace(lines[i]) != "" {
			continue
		}
		return i
	}
	return -1
}

func join(lines []string) string {
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "")
}
