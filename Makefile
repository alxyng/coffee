GOCMD=go
GOBUILD=GOOS=linux $(GOCMD) build
GOTEST=$(GOCMD) test

all: test build
build:
	(cd cmd/coffee; $(GOBUILD) -o main -v)
test:
	$(GOTEST) ./...
deploy:
	sam package \
		--profile dev \
		--template-file template.yaml \
		--s3-bucket coffee.myunidays.dev \
		--output-template-file package.yaml
	sam deploy \
		--profile dev \
		--template-file package.yaml \
		--stack-name coffee-stack \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			SlackTokenParameter=${SLACK_TOKEN} \
			SlackChannelParameter=${SLACK_CHANNEL}
