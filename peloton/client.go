package peloton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type authResponse struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}

type Client struct {
	httpClient http.Client
	UserID     string
	authCookie *http.Cookie
	Host       string
}

func NewClient(username, password, host string) (Client, error) {

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	pelotonClient := Client{
		httpClient: httpClient,
		Host:       host,
	}

	err := pelotonClient.getSessionCookie(username, password)
	if err != nil {
		return pelotonClient, errors.Wrap(err, "failed to authenticate to peloton")
	}

	return pelotonClient, nil
}

func (c *Client) getSessionCookie(username, password string) error {

	postData := bytes.NewBuffer([]byte(fmt.Sprintf("{\"username_or_email\": \"%v\", \"password\": \"%v\"}", username, password)))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/auth/login", c.Host), postData)
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
	c.UserID = authResp.UserID
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "peloton_session_id" {
			c.authCookie = cookie
		}
	}

	return nil
}

func (c *Client) GetWorkouts() (Workouts, error) {
	workouts := Workouts{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/user/%s/workouts", c.Host, c.UserID), nil)
	if err != nil {
		return workouts, errors.Wrap(err, "failed to bulid users workouts request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(c.authCookie)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return workouts, errors.Wrapf(err, "failed to get user workouts response. StatusCode: %d", resp.StatusCode)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return workouts, errors.New(fmt.Sprintf("API returned an unxpected status code: %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&workouts)
	if err != nil {
		return workouts, errors.Wrap(err, "failed to decode response for user workouts")
	}
	return workouts, nil
}

func (c *Client) GetWorkoutDetails(id string, dataFrequency int) (WorkoutDetail, error) {
	workoutDetails := WorkoutDetail{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/workout/%s/performance_graph?every_n=%d", c.Host, id, dataFrequency), nil)
	if err != nil {
		return workoutDetails, errors.Wrap(err, "failed to bulid workout details request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(c.authCookie)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return workoutDetails, errors.Wrap(err, "failed to get workout detail response")
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&workoutDetails)
	if err != nil {
		return workoutDetails, errors.Wrap(err, "failed to decode response for workout details")
	}
	return workoutDetails, nil
}
