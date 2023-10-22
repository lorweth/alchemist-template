FROM golang:1.20.10-alpine3.17

ARG PROJECT_NAME=alchemist-template

RUN apk update

WORKDIR /${PROJECT_NAME}

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Install sqlboiler
RUN go install github.com/volatiletech/sqlboiler/v4@latest && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

# Install mockery
RUN go install github.com/vektra/mockery/v2@v2.36.0
