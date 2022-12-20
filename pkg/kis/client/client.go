package pkgkisclient

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	baseURL = "https://openapi.koreainvestment.com:9443"
	testURL = "https://openapivts.koreainvestment.com:29443"
)

type Credential struct {
	APIKey    string `json:"APIKey" validate:"required"`
	APISecret string `json:"APISecret" validate:"required"`
}

type Client struct {
	http.Client

	cred  Credential
	debug bool
}

type accessToken struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
}

func NewClient(cred Credential, debug bool) Client {
	c := Client{
		cred:  cred,
		debug: debug,
	}

	c.Timeout = 10 * time.Second

	return c
}

func (c *Client) authenticate() (accessToken, error) {
	urlStr := baseURL
	if c.debug {
		urlStr = testURL
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return accessToken{}, err
	}

	u.Path = "/oauth2/tokenP"

	b := new(bytes.Buffer)

	if err = json.NewEncoder(b).Encode(map[string]string{
		"grant_type": "client_credentials",
		"appkey":     c.cred.APIKey,
		"appsecret":  c.cred.APISecret,
	}); err != nil {
		return accessToken{}, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), b)
	if err != nil {
		return accessToken{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := c.Do(req)
	if err != nil {
		return accessToken{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return accessToken{}, errors.Errorf("http response error [code=%d]", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return accessToken{}, err
	}

	token := accessToken{}
	if err = json.Unmarshal(body, &token); err != nil {
		return accessToken{}, err
	}

	return token, nil
}
