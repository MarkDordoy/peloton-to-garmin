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

func (c *Client) GetWorkouts(instances int) ([]WorkoutData, error) {
	workoutData := []WorkoutData{}
	instanceCount := 0
	page := 0
	for instanceCount < instances {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/user/%s/workouts?joins=peloton.ride&limit=%d&page=%d&sort_by=-created", c.Host, c.UserID, instances, page), nil)
		if err != nil {
			return workoutData, errors.Wrap(err, "failed to bulid users workouts request")
		}

		req.Header.Add("Content-Type", "application/json")
		req.AddCookie(c.authCookie)

		resp, err := c.httpClient.Do(req)

		if err != nil {
			return workoutData, errors.Wrapf(err, "failed to get user workouts response. StatusCode: %d", resp.StatusCode)
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return workoutData, errors.New(fmt.Sprintf("API returned an unxpected status code: %d", resp.StatusCode))
		}

		defer resp.Body.Close()
		workouts := Workouts{}
		err = json.NewDecoder(resp.Body).Decode(&workouts)
		if err != nil {
			return workoutData, errors.Wrap(err, "failed to decode response for user workouts")
		}

		for _, data := range workouts.Data {
			if len(workoutData) <= instances {
				workoutData = append(workoutData, data)
				instanceCount++
			} else {
				break
			}
		}

		if page == workouts.PageCount {
			break
		}

		page++
	}
	return workoutData, nil
}

func (c *Client) GetWorkoutDetails(detail WorkoutData, dataFrequency int) (WorkoutDetail, error) {
	workoutDetails := WorkoutDetail{
		ID:                       detail.ID,
		Title:                    detail.Peloton.Ride.Title,
		Description:              detail.Peloton.Ride.Description,
		FitnessDiscipline:        detail.FitnessDiscipline,
		DataGranularityInSeconds: dataFrequency,
		StartTime:                time.Unix(int64(detail.StartTime), 0),
		EndTime:                  time.Unix(int64(detail.EndTime), 0),
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/workout/%s/performance_graph?every_n=%d", c.Host, detail.ID, dataFrequency), nil)
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
