# devcoffee

An endpoint for a Slack command that chooses a random member from a Slack
channel to make coffee.

## Setup

Copy `.env.example` to `.env` and replace `<YOUR_SLACK_TOKEN>` with your slack
token and `<YOUR_SLACK_CHANNEL>` with the channel to pull members from for the
coffee draw.

## Build & Run

From the root of the project, run:
```
docker-compose up --build -d
```

## Stop

From the root of the project, run:
```
docker-compose down
```

## Lambda

Build for AWS Lambda using:

```
GOOS=linux go build -o main
zip deployment.zip main
```
