# devcoffee

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

## Todo

- Show member names in stats without mentioning members
