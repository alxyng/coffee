GOCMD=go
GOBUILD=GOOS=linux $(GOCMD) build
GOTEST=$(GOCMD) test

BINARY_NAME=main

all: test build
build:
	(cd handlers/coffee; $(GOBUILD) -o $(BINARY_NAME) -v)
test:
	$(GOTEST) -v ./...
deploy:
	sam package --template-file template.yaml --s3-bucket coffee-storage --output-template-file package.yaml
	sam deploy \
		--template-file package.yaml \
		--stack-name coffee-stack \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			SlackTokenParameter=${SLACK_TOKEN} \
			SlackChannelParameter=${SLACK_CHANNEL}
