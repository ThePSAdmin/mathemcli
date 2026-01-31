package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	BaseURL     = "https://www.mathem.se/tienda-web-api/v1"
	WebBaseURL  = "https://www.mathem.se"
	ContentType = "application/json"

	// Browser-like User-Agent
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36"
)

// Client handles communication with the Mathem API
type Client struct {
	httpClient *http.Client
	baseURL    string
	sessionID  string
	csrfToken  string
}

// NewClient creates a new API client
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		baseURL: BaseURL,
	}
}

// NewClientWithSession creates a client with an existing session
func NewClientWithSession(sessionID, csrfToken string) *Client {
	client := NewClient()
	client.sessionID = sessionID
	client.csrfToken = csrfToken
	return client
}

// SessionID returns the current session ID
func (c *Client) SessionID() string {
	return c.sessionID
}

// CSRFToken returns the current CSRF token
func (c *Client) CSRFToken() string {
	return c.csrfToken
}

// setBrowserHeaders adds browser-like headers to mimic Chrome
func (c *Client) setBrowserHeaders(req *http.Request, referer string) {
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,sv;q=0.8")
	req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	if referer != "" {
		req.Header.Set("Referer", referer)
	}
}

// doRequest performs an HTTP request with proper headers
func (c *Client) doRequest(method, endpoint string, body any, referer string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	reqURL := c.baseURL + endpoint
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setBrowserHeaders(req, referer)

	if body != nil {
		req.Header.Set("Content-Type", ContentType)
	}

	// Add session cookie if available
	if c.sessionID != "" {
		req.AddCookie(&http.Cookie{
			Name:  "sessionid",
			Value: c.sessionID,
		})
	}

	// Add CSRF token if available
	if c.csrfToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  "csrftoken",
			Value: c.csrfToken,
		})
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Extract cookies from response
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "sessionid":
			c.sessionID = cookie.Value
		case "csrftoken":
			c.csrfToken = cookie.Value
		}
	}

	return resp, nil
}

// decodeResponse decodes a JSON response into the given target
func decodeResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// initSession visits the login page to initialize cookies (csrftoken)
func (c *Client) initSession() error {
	req, err := http.NewRequest(http.MethodGet, WebBaseURL+"/se/user/login/", nil)
	if err != nil {
		return fmt.Errorf("failed to create init request: %w", err)
	}

	c.setBrowserHeaders(req, WebBaseURL+"/se/")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to init session: %w", err)
	}
	defer resp.Body.Close()

	// Extract cookies
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrftoken" {
			c.csrfToken = cookie.Value
		}
	}

	return nil
}

// Login authenticates with email and password
func (c *Client) Login(email, password string) error {
	// First, visit the login page to get CSRF token and initial cookies
	if err := c.initSession(); err != nil {
		return fmt.Errorf("failed to initialize session: %w", err)
	}

	payload := map[string]string{
		"username": email,
		"password": password,
	}

	resp, err := c.doRequest(http.MethodPost, "/user/login/", payload, WebBaseURL+"/se/user/login/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Session ID is extracted in doRequest from cookies
	if c.sessionID == "" {
		return fmt.Errorf("login succeeded but no session cookie received")
	}

	return nil
}

// Search searches for products
func (c *Client) Search(query string, page int) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("/search/mixed/?q=%s&type=product&page=%d",
		url.QueryEscape(query), page)

	resp, err := c.doRequest(http.MethodGet, endpoint, nil, WebBaseURL+"/se/")
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCart retrieves the current cart
func (c *Client) GetCart() (*Cart, error) {
	resp, err := c.doRequest(http.MethodGet, "/cart/?group_by=recipes", nil, WebBaseURL+"/se/")
	if err != nil {
		return nil, err
	}

	var cart Cart
	if err := decodeResponse(resp, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

// AddToCart adds items to the cart
func (c *Client) AddToCart(items []CartItem) (*Cart, error) {
	payload := map[string][]CartItem{
		"items": items,
	}

	resp, err := c.doRequest(http.MethodPost, "/cart/items/", payload, WebBaseURL+"/se/")
	if err != nil {
		return nil, err
	}

	var cart Cart
	if err := decodeResponse(resp, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

// ClearCart removes all items from the cart
func (c *Client) ClearCart() (*Cart, error) {
	resp, err := c.doRequest(http.MethodPost, "/cart/clear/", nil, WebBaseURL+"/se/cart/")
	if err != nil {
		return nil, err
	}

	var cart Cart
	if err := decodeResponse(resp, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}
