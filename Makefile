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
	sam package --template-file template.yaml --s3-bucket coffee.myunidays.dev --output-template-file package.yaml --profile dev
	sam deploy \
		--template-file package.yaml \
		--stack-name coffee-stack \
		--capabilities CAPABILITY_IAM \
		--profile dev \
		--parameter-overrides \
			SlackTokenParameter=${SLACK_TOKEN} \
			SlackChannelParameter=${SLACK_CHANNEL}
