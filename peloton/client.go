package peloton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

type authResponse struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}

type Client struct {
	httpClient http.Client
	UserID     string
	sessionID  string
}

func NewClient(username, password string) (Client, error) {

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	standardClient := retryClient.StandardClient()
	standardClient.Timeout = time.Second * 10

	pelotonClient := Client{
		httpClient: *standardClient,
	}

	err := pelotonClient.getSessionToken(username, password)
	if err != nil {
		return pelotonClient, errors.Wrap(err, "failed to authenticate to peloton")
	}

	return pelotonClient, nil
}

func (c *Client) getSessionToken(username, password string) error {

	postData := bytes.NewBuffer([]byte(fmt.Sprintf("{\"username_or_email\": \"%v\", \"password\": \"%v\"}", username, password)))
	req, err := http.NewRequest("POST", "https://api.onepeloton.com/auth/login", postData)
	if err != nil {
		return errors.Wrap(err, "failed to build http request")
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}

	defer resp.Body.Close()

	authResp := authResponse{}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return errors.Wrap(err, "failed to decode peloton auth response")
	}
	c.sessionID = authResp.SessionID
	c.UserID = authResp.SessionID
	return nil
}
