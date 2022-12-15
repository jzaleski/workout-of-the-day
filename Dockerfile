FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./main.go
RUN go build -o ./workout-of-the-day

FROM alpine:3.15

RUN apk add ca-certificates

WORKDIR /app

COPY assets/*.css ./assets
COPY assets/*.ico ./assets
COPY templates/*.html ./templates
COPY --from=builder /app/workout-of-the-day ./bin/workout-of-the-day

CMD ["/app/bin/workout-of-the-day"]
