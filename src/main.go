package main

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"net/http"
	"os"
)

// Config is the configuration for the job.
type Config struct {
	url            string // URL
	slackApiToken  string // SLACK_API_TOKEN
	slackChannelId string // SLACK_CHANNEL_ID
}

// get Environment Variables
func configFromEnv() (Config, error) {
	url := os.Getenv("URL")
	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	slackChannelId := os.Getenv("SLACK_CHANNEL_ID")

	// set config from Environment Variables
	config := Config{
		url:            url,
		slackApiToken:  slackApiToken,
		slackChannelId: slackChannelId,
	}

	// check Environment Variables
	if config.url == "" {
		return config, errors.New("URL is not set")
	}
	if config.slackApiToken == "" {
		return config, errors.New("SLACK_API_TOKEN is not set")
	}
	if config.slackChannelId == "" {
		return config, errors.New("SLACK_CHANNEL_ID is not set")
	}

	return config, nil
}

func main() {
	// check Environment Variables
	config, err := configFromEnv()
	if err != nil {
		log.Fatalf("{Environment variables error: %v }", err)
	}

	// check status code
	resp, err := http.Get(config.url)
	if err != nil {
		log.Fatalf("{Request error: {message: %v, url: %v }}", err, config.url)
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %v\n", resp.Status)

	// status code not 200
	if resp.StatusCode != http.StatusOK {
		// notiry to slack
		api := slack.New(config.slackApiToken)
		// Notification content
		attachment := slack.Attachment{
			Title: "サーバがダウンしました",
			Color: "danger",
			Text:  config.url,
		}
		channelID, timestamp, err := api.PostMessage(
			config.slackChannelId,
			// slack.MsgOptionText("Hello world!", false),
			slack.MsgOptionAttachments(attachment),
		)
		if err != nil {
			fmt.Printf("Slack Post Message Error %s\n", err)
			return
		}

		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	}

}
