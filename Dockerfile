# --platform=$BUILDPLATFORM keeps the build stage on the runner's native arch and
# lets Go cross-compile to $TARGETARCH instead of emulating the whole toolchain
# under QEMU. For a pure-Go binary this is the same output at a fraction of the
# build time, so multi-arch costs us almost nothing here.
FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18 as build
WORKDIR /go/src/github.com/VATUSA/google-workspace-integration
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal
ARG TARGETOS
ARG TARGETARCH
# CGO_ENABLED=0 is required for cross-compilation without a target C toolchain,
# and is safe here: the binary is pure Go.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/core ./cmd/core/main.go

FROM alpine:3.18 as app
WORKDIR /app
COPY --from=build /go/src/github.com/VATUSA/google-workspace-integration/bin/* ./
