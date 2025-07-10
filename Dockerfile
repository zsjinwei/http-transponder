FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o http-transponder ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/http-transponder .
COPY config.yaml .
EXPOSE 8080
ENTRYPOINT ["./http-transponder", "-config", "config.yaml"]
