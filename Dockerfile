# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.23-alpine

WORKDIR /usr/src/app

RUN apk update && apk upgrade && apk add --no-cache bash

# For Go Debugging
# RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY users-backend/go.mod users-backend/go.sum ./
RUN go mod download && go mod verify
COPY users-backend/ ./

RUN go clean -cache && go build -o main .

EXPOSE 8080
CMD ["./main"]

# For Go Debugging
# RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o main .
# EXPOSE 8080
# EXPOSE 2345
# CMD [ "/go/bin/dlv", "--listen=:2345", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/usr/src/app/main" ]
