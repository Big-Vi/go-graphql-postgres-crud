## Dev
FROM golang:1.18-alpine as development

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install live reloader
RUN go install github.com/cosmtrek/air@latest
RUN air init

# Copy app files
COPY . .

EXPOSE 8000

CMD air
