.PHONY: deps clean test build run package deploy

deps:
	dep ensure

clean:
	rm -rf main

test:
	go test -v ./...

build:
	GOOS=linux GOARCH=amd64 go build -o devcoffee

package:
	aws cloudformation package \
		--template-file template.yaml \
		--s3-bucket devcoffee.myunidays.com \
		--output-template-file packaged.yaml

deploy:
	aws cloudformation deploy \
		--template-file ./packaged.yaml \
		--stack-name coffee \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			SlackTokenParameter=${SLACK_TOKEN} \
			SlackChannelParameter=${SLACK_CHANNEL}

run:
	sam local start-api \
		--env-vars env.json \
		--docker-network devcoffee_default
