package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	baseURL string
	apiService string
	httpClient *http.Client
}

type IClient interface {
	MakeRequest(ctx context.Context, method string, endpoint string, body *bytes.Buffer) (closer io.ReadCloser, err error)
}

func NewHTTPClient(credentials *OAuthConfig, apiService string, baseURL string) (*Client, error) {
	client, err := NewClient(credentials)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil, err
	}

	return &Client{httpClient: client, baseURL: baseURL,
		apiService: apiService}, nil
}

func (c *Client) makeURL(endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", c.baseURL, c.apiService, endpoint)
}

func (c *Client) MakeRequest(ctx context.Context, method string, endpoint string, body *bytes.Buffer) (closer io.ReadCloser, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, c.makeURL(endpoint), body)
	if err != nil {
		return nil, err
	}
	switch method {
		case "GET":
		case "DELETE":
		default:
			req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Set("X-Message-Id", uuid.New().String())
	req.Header.Set("X-Transaction-Id", uuid.New().String())
	req.Header.Set("X-Transaction-Caller", "terraform-provider-plusserver")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s URL: %s", resp.StatusCode, string(bodyBytes), resp.Request.URL)
	}
	return resp.Body, nil
}