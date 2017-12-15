# devcoffee

Setup, build and run by running the following:
```
export SLACK_TOKEN=<YOUR_SLACK_TOKEN>
export SLACK_CHANNEL=<YOUR_SLACK_CHANNEL>
docker build -t devcoffee:latest .
docker run -e SLACK_TOKEN -e SLACK_CHANNEL -p "3000:3000" devcoffee:latest
```

## Todo

- Persist devcoffee.db on the host using a docker volume
- Show member names in stats without mentioning members
