package switchbot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const defaultEndpoint = "https://api.switch-bot.com/v1.1"

type Client struct {
	httpClient *http.Client
	endpoint   string
	token      string
	secretKey  string
}

func New(token, secretKey string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		endpoint:   defaultEndpoint,
		token:      token,
		secretKey:  secretKey,
	}
}

type response struct {
	Body io.ReadCloser
}

func (r *response) Read(p []byte) (n int, err error) {
	return r.Body.Read(p)
}

func (r *response) Close() error {
	return r.Body.Close()
}

func (c *Client) do(ctx context.Context, method, path string) (io.ReadCloser, error) {
	nonce := uuid.New().String()
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	signer := hmac.New(sha256.New, []byte(c.secretKey))
	signer.Write([]byte(fmt.Sprintf("%s%s%s", c.token, t, nonce)))
	sign := base64.StdEncoding.EncodeToString(signer.Sum(nil))

	p, err := url.JoinPath(c.endpoint, path)
	if err != nil {
		return nil, fmt.Errorf("failed to join URL paths: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, p, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the HTTP request: %w", err)
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Add("t", t)
	req.Header.Add("sign", sign)
	req.Header.Add("nonce", nonce)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do the HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return &response{Body: resp.Body}, nil
}
