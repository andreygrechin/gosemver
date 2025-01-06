# Build arguments
ARG BUILD_IMAGE_PATH=golang
ARG BUILD_IMAGE_TAG=1.23.4-alpine3.21

# Build-time metadata
ARG VERSION=unknown
ARG COMMIT=unknown
ARG BUILDTIME=unknown
ARG APP_NAME=gosemver
ARG MOD_PATH=github.com/andreygrechin/gosemver

# Build platform
ARG GOARCH=amd64
ARG GOOS=linux

# Build stage
FROM ${BUILD_IMAGE_PATH}:${BUILD_IMAGE_TAG} AS build-stage
ARG VERSION
ARG COMMIT
ARG BUILDTIME
ARG APP_NAME
ARG MOD_PATH
ARG GOARCH=amd64
ARG GOOS=linux

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY *.go ./

# Build application
RUN CGO_ENABLED=0 GOARCH=${GOARCH} GOOS=${GOOS} go build \
    -ldflags "-s -w \
    -X ${MOD_PATH}/internal/config.Version=${VERSION} \
    -X ${MOD_PATH}/internal/config.BuildTime=${BUILDTIME} \
    -X ${MOD_PATH}/internal/config.Commit=${COMMIT}" \
    -o "/app/${APP_NAME}"

# Final stage
FROM alpine:3.21.0
ARG APP_NAME

COPY --from=build-stage /app/${APP_NAME} /app/${APP_NAME}
ENV PATH="/app:$PATH"

# Setup non-root user
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
USER nonroot:nonroot

CMD [ "gosemver", "--help" ]
