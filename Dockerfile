# syntax=docker/dockerfile:1
ARG GO_VERSION=1.26.3
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go mod download
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gateway ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=build /gateway /app/gateway
COPY web /app/web
ENV HTTP_ADDR=:8080
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/v1/health || exit 1
CMD ["/app/gateway"]
