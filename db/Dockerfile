FROM golang:1.23.3-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# make goose available in the path
ENV PATH=$PATH:/go/bin

COPY . /db