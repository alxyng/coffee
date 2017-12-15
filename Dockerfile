FROM golang:1.9 as build
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/nullseed/devcoffee
COPY handler.go .
COPY main.go .
COPY services services
COPY Gopkg.lock .
COPY Gopkg.toml .
RUN dep ensure --vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o devcoffee .

FROM debian:latest
RUN apt-get update
RUN apt-get install -y ca-certificates
WORKDIR /app
COPY --from=build /go/src/github.com/nullseed/devcoffee/devcoffee .

EXPOSE 3000

CMD ["./devcoffee"]
