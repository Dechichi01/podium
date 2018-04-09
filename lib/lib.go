package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// Podium is a struct that represents a podium API application
type Podium struct {
	Config            *viper.Viper
	URL               string
	baseLeaderboard   string
	localeLeaderboard string
}

// NewPodium returns a new podium API application
func NewPodium(config *viper.Viper) *Podium {
	app := &podium{
		Config:            config,
		URL:               config.Get("podium.url").(string),
		baseLeaderboard:   config.Get("leaderboards.globalLeaderboard").(string),
		localeLeaderboard: config.Get("leaderboards.localeLeaderboard").(string),
	}
	return app
}

func sendTo(method, url string, payload map[string]interface{}) (int, string, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return -1, "", err
	}

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
		if err != nil {
			return -1, "", err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return -1, "", err
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return -1, "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, string(body), nil
}

func (app *podium) buildURL(pathname string) string {
	return fmt.Sprintf("%s%s", app.URL, pathname)
}

func (app *podium) buildDeleteLeaderboardURL(leaderboard string) string {
	var pathname = fmt.Sprintf("/l/%s", leaderboard)
	return app.buildURL(pathname)
}

func (app *podium) buildGetTopPercentURL(leaderboard string, percentage int) string {
	var pathname = fmt.Sprintf("/l/%s/top-percent/%d", leaderboard, percentage)
	return app.buildURL(pathname)
}

func (app *podium) buildUpdateScoreURL(leaderboard string, playerID string) string {
	var pathname = fmt.Sprintf("/l/%s/members/%s/score", leaderboard, playerID)
	return app.buildURL(pathname)
}

func (app *podium) buildIncrementScoreURL(leaderboard string, playerID string) string {
	return app.buildUpdateScoreURL(leaderboard, playerID)
}

func (app *podium) buildUpdateScoresURL(playerID string) string {
	var pathname = fmt.Sprintf("/m/%s/scores", playerID)
	return app.buildURL(pathname)
}

func (app *podium) buildRemoveMemberFromLeaderboardURL(leaderboard string, member string) string {
	var pathname = fmt.Sprintf("/l/%s/members/%s", leaderboard, member)
	return app.buildURL(pathname)
}

// page is 1-based
func (app *podium) buildGetTopURL(leaderboard string, page int, pageSize int) string {
	var pathname = fmt.Sprintf("/l/%s/top/%d?pageSize=%d", leaderboard, page, pageSize)
	return app.buildURL(pathname)
}

func (app *podium) buildGetPlayerURL(leaderboard string, playerID string) string {
	var pathname = fmt.Sprintf("/l/%s/members/%s", leaderboard, playerID)
	return app.buildURL(pathname)
}

func (app *podium) buildHealthcheckURL() string {
	var pathname = "/healthcheck"
	return app.buildURL(pathname)
}

// external functions:
func (app *podium) GetBaseLeaderboards() string {
	return app.baseLeaderboard
}

func (app *podium) GetLocalizedLeaderboard(locale string) string {
	localeLeaderboard := app.localeLeaderboard
	result := strings.Replace(localeLeaderboard, "%{locale}", locale, -1)
	return result
}

func (app *podium) GetTop(leaderboard string, page int, pageSize int) (int, string, error) {
	route := app.buildGetTopURL(leaderboard, page, pageSize)
	status, body, err := sendTo("GET", route, nil)
	return status, body, err
}

func (app *podium) GetTopPercent(leaderboard string, percentage int) (int, string, error) {
	route := app.buildGetTopPercentURL(leaderboard, percentage)
	status, body, err := sendTo("GET", route, nil)
	return status, body, err
}

func (app *podium) UpdateScore(leaderboard string, playerID string, score int) (int, string, error) {
	route := app.buildUpdateScoreURL(leaderboard, playerID)
	payload := map[string]interface{}{
		"score": score,
	}
	status, body, err := sendTo("PUT", route, payload)
	return status, body, err
}

func (app *podium) IncrementScore(leaderboard string, playerID string, increment int) (int, string, error) {
	route := app.buildIncrementScoreURL(leaderboard, playerID)
	payload := map[string]interface{}{
		"increment": increment,
	}
	status, body, err := sendTo("PATCH", route, payload)
	return status, body, err
}

func (app *podium) UpdateScores(leaderboards []string, playerID string, score int) (int, string, error) {
	route := app.buildUpdateScoresURL(playerID)
	payload := map[string]interface{}{
		"score":        score,
		"leaderboards": leaderboards,
	}
	status, body, err := sendTo("PUT", route, payload)
	return status, body, err
}

func (app *podium) RemoveMemberFromLeaderboard(leaderboard string, member string) (int, string, error) {
	route := app.buildRemoveMemberFromLeaderboardURL(leaderboard, member)
	status, body, err := sendTo("DELETE", route, nil)
	return status, body, err
}

func (app *podium) GetPlayer(leaderboard string, playerID string) (int, string, error) {
	route := app.buildGetPlayerURL(leaderboard, playerID)
	status, body, err := sendTo("GET", route, nil)
	return status, body, err
}

func (app *podium) Healthcheck(leaderboard string, playerID string) (int, string, error) {
	route := app.buildHealthcheckURL()
	status, body, err := sendTo("GET", route, nil)
	return status, body, err
}

func (app *podium) DeleteLeaderboard(leaderboard string) (int, string, error) {
	route := app.buildDeleteLeaderboardURL(leaderboard)
	status, body, err := sendTo("DELETE", route, nil)
	return status, body, err
}
