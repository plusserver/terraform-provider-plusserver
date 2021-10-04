package dns

import (
	"bytes"
	"context"
	"github.com/plusserver/terraform-provider-plusserver/api"
	"io"
)

type Client struct {
	api.Client
}

func NewDNSClient(credentials *api.OAuthConfig, baseURL string) (*Client, error) {
	client, err := api.NewHTTPClient(credentials, "dnsEntityService", baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{*client}, nil
}

func (c *Client) MakeRequest(ctx context.Context, method string, endpoint string, body *bytes.Buffer) (closer io.ReadCloser, err error) {
	return c.Client.MakeRequest(ctx, method, endpoint, body)
}


