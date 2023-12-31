FROM golang:1.21.1-alpine3.18 as build
WORKDIR /go/src/github.com/VATUSA/google-workspace-integration
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o bin/core ./cmd/core/main.go

FROM alpine:3.18 as app
WORKDIR /app
COPY --from=build /go/src/github.com/VATUSA/google-workspace-integration/bin/* ./