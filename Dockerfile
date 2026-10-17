# syntax=docker/dockerfile:1

FROM golang:1.27.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/main .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app -u 10001 app

WORKDIR /app

COPY --from=builder /out/main /main
COPY --from=builder /src/config.yaml /app/config.yaml

USER app
EXPOSE 3333

ENTRYPOINT ["/main"]
