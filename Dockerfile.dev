FROM golang:1.27-alpine

WORKDIR /app

RUN apk add --no-cache git bash inotify-tools

# Air (hot reload) — requires Go 1.25+
RUN go install github.com/air-verse/air@latest

# Cache deps (source code comes via volume)
COPY go.mod go.sum ./
RUN go mod download

COPY .air.toml ./

EXPOSE 9090

CMD ["air", "-c", ".air.toml"]
