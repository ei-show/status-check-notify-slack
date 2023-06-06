package main

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"net/http"
	"os"
)

// get Environment Variables URL
func getUrl() (string, error) {
	url := os.Getenv("URL")
	if url == "" {
		return "", errors.New("URL is not set")
	}
	return url, nil
}

// get Environment Variables SLACK_API_TOKEN
func getSlackApiToken() (string, error) {
	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	if slackApiToken == "" {
		return "", errors.New("SLACK_API_TOKEN is not set")
	}
	return slackApiToken, nil
}

// get Environment Variables SLACK_API_TOKEN
func getSlackChannelId() (string, error) {
	slackChannelId := os.Getenv("SLACK_CHANNEL_ID")
	if slackChannelId == "" {
		return "", errors.New("SLACK_CHANNEL_ID is not set")
	}
	return slackChannelId, nil
}

func main() {
	// check Environment Variables
	url, err := getUrl()
	if err != nil {
		fmt.Printf("get Environment Variables: %s\n", err)
		return
	}

	slackApiToken, err := getSlackApiToken()
	if err != nil {
		fmt.Printf("get Environment Variables: %s\n", err)
		return
	}

	slackChannelId, err := getSlackChannelId()
	if err != nil {
		fmt.Printf("get Environment Variables: %s\n", err)
		return
	}

	// check status code
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %v\n", resp.Status)

	// status code not 200
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response error: %v\n", resp.Status)

		// nortiry to slack
		api := slack.New(slackApiToken)
		channelID, timestamp, err := api.PostMessage(
			slackChannelId,
			slack.MsgOptionText("Hello world!", false),
		)
		if err != nil {
			fmt.Printf("Slack Post Message Error %s\n", err)
			return
		}

		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	}

}
