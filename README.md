# status-check-notify-slack
Slack notification if the URL status code is other than 200.

## env
|  key  |  value  |
| ---- | ---- |
|  URL  |  https://google.com  |
|  SLACK_API_TOKEN  |  XXXXXXXXXXXXXX  |
|  SLACK_CHANNEL_ID  |  XXXXXXXXXXXXXX  |
|  BUCKET_NAME  |  XXXXXXXXXXXXXX  |
|  TOPIC  |  XXXXXXXXXXXXXX  |

### local env
```sh
set -gx URL https://google.com
set -gx SLACK_API_TOKEN XXXXXXXXXXXXXX
set -gx SLACK_CHANNEL_ID XXXXXXXXXXXXXX
set -gx BUCKET_NAME XXXXXXXXXXXXXX
```

## build

```sh
gcloud builds submit --config=cloudbuild.yaml --project==offiter-registry --substitutions=TAG_NAME="v0.0.1"
```