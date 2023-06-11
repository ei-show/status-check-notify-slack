package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

// Config is the configuration for the job.
type Config struct {
	url            string // URL
	slackApiToken  string // SLACK_API_TOKEN
	slackChannelId string // SLACK_CHANNEL_ID
	bucketName     string // BUCKET_NAME
}

// get Environment Variables
func configFromEnv() (Config, error) {
	url := os.Getenv("URL")
	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	slackChannelId := os.Getenv("SLACK_CHANNEL_ID")
	bucketName := os.Getenv("BUCKET_NAME")

	// set config from Environment Variables
	config := Config{
		url:            url,
		slackApiToken:  slackApiToken,
		slackChannelId: slackChannelId,
		bucketName:     bucketName,
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
	// if config.bucketName == "" {
	// 	return config, errors.New("BUCKET_NAME is not set")
	// }

	return config, nil
}

// URLをjsonのファイル名に変換する
func convertURLToFileName(inputURL string) string {
	// URLを解析
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		// エラーハンドリング
		panic(err)
	}

	// ホスト部分を変換
	host := strings.ReplaceAll(parsedURL.Host, ".", "_")

	// パス部分を変換
	path := strings.ReplaceAll(parsedURL.Path, "/", "_")

	// ファイル名を組み立て
	fileName := fmt.Sprintf("%s%s.json", host, path)

	return fileName
}

// オブジェクトを作成する関数
func createObject(ctx context.Context, obj *storage.ObjectHandle, data []byte) error {
	writer := obj.NewWriter(ctx)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

// オブジェクトの内容を読み込む関数
func readObject(ctx context.Context, obj *storage.ObjectHandle) ([]byte, error) {
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// notify to slack
func notifyToSlack(config Config, title string, color string) {
	// notify to slack
	api := slack.New(config.slackApiToken)
	// Notification content
	attachment := slack.Attachment{
		Title: title,
		Color: color,
		Text:  config.url,
	}
	channelID, timestamp, err := api.PostMessage(
		config.slackChannelId,
		// slack.MsgOptionText("Hello world!", false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}

func main() {
	// check Environment Variables
	config, err := configFromEnv()
	if err != nil {
		log.Fatalf("{Environment variables error: %v }\n", err)
	}

	// config.bucketName が存在する場合は、GCSにアクセスする
	if config.bucketName != "" {
		// GCS clientの作成
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		// GCSで使用するオブジェクト名を生成する
		objectName := convertURLToFileName(config.url)
		obj := client.Bucket(config.bucketName).Object(objectName)

		// objectが存在するか確認する
		_, err = obj.Attrs(ctx)
		if err != nil {
			if err == storage.ErrObjectNotExist {
				// オブジェクトが存在しない場合は、作成する
				err := createObject(ctx, obj, nil)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("オブジェクトを作成しました")
			} else {
				log.Fatal(err)
			}
		} else {
			// オブジェクトが存在する場合は、読み込む
			data, err := readObject(ctx, obj)
			fmt.Printf("オブジェクトの内容: %s\n", data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// check status code
	resp, err := http.Get(config.url)
	if err != nil {
		log.Fatalf("{Request error: {message: %v, url: %v }}\n", err, config.url)
	}
	defer resp.Body.Close()

	fmt.Printf("{Response: {url: %v, status: %v }}\n", config.url, resp.Status)

	// status code not 200
	if resp.StatusCode != http.StatusOK {
		notifyToSlack(config, "サーバがダウンしています", "danger")
	} else {
		notifyToSlack(config, "サーバが復帰しました", "good")
	}

}
