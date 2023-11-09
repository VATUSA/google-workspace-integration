FROM golang:1.21.1-alpine:3.16 as build
WORKDIR /go/src/github.com/VATUSA/google-workspace-integration
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
RUN go build -o bin/core ./cmd/core/main.go

FROM alpine:latest as core
WORKDIR / app
COPY --from=build /go/src/github.com/VATUSA/google-workspace-integration/bin/core ./
CMD ["./core"]