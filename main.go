package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/JonathanHeinz/slackify/spotify"
	"github.com/getlantern/systray"
)

var slackAPI = "https://slack.com/api/users.profile.set"
var slackAPIKey = ""
var infoButton *systray.MenuItem

type slackStatus struct {
	Profile profile `json:"profile"`
}

type profile struct {
	StatusText       string `json:"status_text"`
	StatusEmoji      string `json:"status_emoji"`
	StatusExpiration int    `json:"status_expiration"`
}

func main() {
	slackAPIKey = os.Getenv("SLACK_TOKEN")
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Slackify")
	infoButton = systray.AddMenuItem("", "")
	infoButton.Disable()
	exitButton := systray.AddMenuItem("Quit", "Quit Slackify")

	err := setStatus()
	if err != nil {
		return
	}

	timer := time.NewTicker(time.Second * 30).C

	for {
		select {
		case <-timer:
			err := setStatus()
			if err != nil {
				systray.Quit()
				return
			}
		case <-exitButton.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onExit() {
	fmt.Println("quitting...")
}

func setStatus() error {
	song, err := spotify.GetCurrentTitle()
	if err != nil {
		return err
	}

	status := song.Item.Artists + " - " + song.Item.Name
	infoButton.SetTitle(status)

	err = setSlackStatus(status)
	if err != nil {
		return err
	}

	return nil
}

func setSlackStatus(status string) error {
	if slackAPIKey == "" {
		return errors.New("token is not defined")
	}

	client := http.Client{}
	newStatus := slackStatus{
		profile{
			StatusText:       status,
			StatusEmoji:      ":notes:",
			StatusExpiration: 0,
		},
	}

	JSONStatus, err := json.Marshal(newStatus)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", slackAPI, bytes.NewBuffer(JSONStatus))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+slackAPIKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
