FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /lua-tracing-poc-client
CMD ["/lua-tracing-poc-client"]
