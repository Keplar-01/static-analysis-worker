FROM --platform=linux/amd64 golang:1.22-alpine AS go-builder

WORKDIR /app

COPY worker-static-analyzer/go.mod worker-static-analyzer/go.sum ./
RUN go mod download

COPY worker-static-analyzer/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /worker-static ./cmd

FROM --platform=linux/amd64 keplar01/static-analyzer:latest

COPY --from=go-builder /worker-static /usr/local/bin/worker-static

WORKDIR /app

ENTRYPOINT []
CMD ["worker-static"]
