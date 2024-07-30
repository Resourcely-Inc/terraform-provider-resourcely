package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	ProjectURL     = "github.com/Resourcely-Inc/terraform-provider-resourcely"
	ProjectVersion = "0.1.0"

	DefaultScheme = "https"

	HeaderToken       = "Authorization"
	HeaderTokenFormat = "Bearer %s"

	MediaTypeJSON = "application/json"

	BasePath = "api/v1"
)

var (
	DefaultUserAgent = fmt.Sprintf(
		"terraform-provider-resourcely/%s (+%s; %s)",
		ProjectVersion, ProjectURL, runtime.Version(),
	)
	WantAcceptHeaders      = []string{MediaTypeJSON}
	WantContentTypeHeaders = []string{MediaTypeJSON}
)

type Client struct {
	Client *retryablehttp.Client

	BasePath string

	BaseURL   *url.URL
	UserAgent string
	AuthToken string

	// Requse a single struct instead of allocating one for each service on the heap
	common service

	// Services
	Blueprints       *BlueprintsService
	ContextQuestions *ContextQuestionsService
	System           *SystemService
	Guardrails       *GuardrailsService
}

type service struct {
	Client *Client
}

// All the exported methods in this file are designed to be
// general-purpose HTTP helpers. These methods will accept any request
// struct and deserialize responses into any struct type.

func (c *Client) Get(ctx context.Context, path string, queryParameters url.Values, respBody interface{}) (interface{}, *http.Response, error) {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add the query parameters to the request URL.
	req.URL.RawQuery = queryParameters.Encode()

	return c.MakeRequest(ctx, req, respBody)
}

// post makes post requests to the given path with the given fields and stores the response in the given body object.
func (c *Client) Post(ctx context.Context, path string, fields interface{}, respBody interface{}) (interface{}, *http.Response, error) {
	req, err := c.NewRequest("POST", path, fields)
	if err != nil {
		return nil, nil, err
	}
	return c.MakeRequest(ctx, req, respBody)
}

// put makes put requests to the given path with the given fields and stores the response in the given body object.
func (c *Client) Put(ctx context.Context, path string, fields interface{}, respBody interface{}) (interface{}, *http.Response, error) {
	req, err := c.NewRequest("PUT", path, fields)
	if err != nil {
		return nil, nil, err
	}
	return c.MakeRequest(ctx, req, respBody)
}

// delete makes delete requests to the given path.
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, nil)
}

// makeRequest makes the given request, and stores the result into the given body.
func (c *Client) MakeRequest(ctx context.Context, req *http.Request, body interface{}) (interface{}, *http.Response, error) {
	resp, err := c.Do(ctx, req, body)
	return body, resp, err
}

// NewClient returns a new Resourcely API client.
func NewClient(httpClient *retryablehttp.Client, host string, authToken string) (*Client, error) {
	if httpClient == nil {
		httpClient = retryablehttp.NewClient()
		httpClient.HTTPClient = &http.Client{}
	}

	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	if baseURL.Scheme == "" {
		baseURL.Scheme = DefaultScheme
	}

	c := &Client{Client: httpClient, BasePath: BasePath, BaseURL: baseURL, UserAgent: DefaultUserAgent, AuthToken: authToken}
	c.common.Client = c
	c.Blueprints = (*BlueprintsService)(&c.common)
	c.ContextQuestions = (*ContextQuestionsService)(&c.common)
	c.Guardrails = (*GuardrailsService)(&c.common)
	c.System = (*SystemService)(&c.common)
	return c, nil
}

// NewRequest creates a new HTTP request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	if c.AuthToken != "" {
		req.Header.Set(HeaderToken, fmt.Sprintf(HeaderTokenFormat, c.AuthToken))
	}

	for _, accept := range WantAcceptHeaders {
		req.Header.Set("Accept", accept)
	}

	if body != nil {
		for _, contentType := range WantContentTypeHeaders {
			req.Header.Set("Content-Type", contentType)
		}
	}

	return req, nil
}

// Do executes an HTTP request.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	retryableReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(retryableReq)
	if err != nil {
		// If we got an error, and the context has been canceled,.
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, _ = io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

// ErrorResponse represents the error response from the API.
type Err struct {
	Status      uint32   `json:"status"`
	RequestId   string   `json:"request_id"`
	Errors      []string `json:"errors"`
	RequestPath string   `json:"request_path"`
	AppVersion  string   `json:"app_version"`
}

type ErrorResponse struct {
	Response *http.Response
	Err      Err
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: [%d] %v - %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Err.RequestId, strings.Join(r.Err.Errors, ", "))
}

// CheckResponse checks the HTTP response for an error.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// Authorization errors do not follow the standard ErrorResponse
	// format
	if r.StatusCode == 401 {
		return &ErrorResponse{
			Response: r,
			Err: Err{
				Status: 401,
				Errors: []string{"Unauthorized"},
			},
		}
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		_ = json.Unmarshal(data, &errorResponse.Err)
	}

	return errorResponse
}

type SystemHealthResponse struct {
	Status string `json:"status"`
}

func (c *Client) Check() error {
	shr, _, err := c.System.GetHealth(context.Background())
	if err != nil {
		return err
	}

	if !(shr.Status == "ok") {
		return fmt.Errorf("system/health: available: %s, body: %s", shr.Status, shr)
	}

	return nil
}

// Model of Resourcley Auth Token claims
type ResourcelyClaims struct {
	// We currently only care about the tenant claim
	Tenant string `json:"@resourcely/tenant"`
	jwt.RegisteredClaims
}

// Returns the tenant name from the auth token
func (c *Client) Tenant() (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(c.AuthToken, &ResourcelyClaims{})
	if err != nil {
		return "", fmt.Errorf("Error parsing Resourcely auth token: %w", err)
	}

	claims, ok := token.Claims.(*ResourcelyClaims)
	if !ok {
		return "", fmt.Errorf("Error parsing Resourcely auth token: invalid claims")
	}

	return claims.Tenant, nil
}
