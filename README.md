# devcoffee

[![Build Status](https://travis-ci.org/nullseed/devcoffee.svg?branch=master)](https://travis-ci.org/nullseed/devcoffee)

An endpoint for a Slack command that chooses a random member from a Slack
channel to make coffee.

## Testing locally

### Setup

Copy `env.example.json` to `env.json` and replace `<YOUR_SLACK_TOKEN>` with your slack
token and `<YOUR_SLACK_CHANNEL>` with the channel to pull members from for the
coffee draw.

### Build & Run

From the root of the project, run:
```
docker-compose up -d
make run
```

### Stop

From the root of the project, run:
```
docker-compose down
```
