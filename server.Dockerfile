# docker build -t brage1t --build-arg VERSION="v0.0.1" -f server.Dockerfile .
# ---- BUILD FRONTEND
FROM node:22-alpine AS fe-build

ARG VERSION
RUN test -n "$VERSION" || (echo "VERSION is required" && exit 1)

RUN apk add --no-cache git

WORKDIR /fe-build

COPY . .

RUN mkdir -p ./assets/frontend
RUN VITE_APP_VERSION=${VERSION} npm install
RUN npm install --prefix web/
RUN npm run build --prefix web/
RUN cp -r ./web/dist/* ./assets/frontend
RUN cp ./web/static/* ./assets/frontend/assets

# ---- BUILD BACKEND
FROM golang:1.25-alpine AS go-build

ARG VERSION
RUN test -n "$VERSION" || (echo "VERSION is required" && exit 1)

WORKDIR /build

COPY . .

COPY --from=fe-build /fe-build/assets/frontend ./assets/frontend

RUN apk add --no-cache \
    git \
    build-base \
    portaudio-dev \
    alsa-lib-dev \
    jack-dev

RUN go mod download
RUN CGO_ENABLED=1 go build -tags fts5 -o bragemusic-server --ldflags="-X 'github.com/bragemusic/bragemusic/internal/vars.VERSION=${VERSION}'" cmd/server/main.go

# ---- BUILD RELEASE STAGE
FROM alpine:latest AS release

LABEL org.opencontainers.image.source="https://github.com/bragemusic/bragemusic"

WORKDIR /bragemusic-server
COPY --from=go-build /build/db/migrations ./migrations

COPY --from=amacneil/dbmate:latest /usr/local/bin/dbmate /usr/local/bin/dbmate
COPY --from=ghcr.io/lucas-ingemar/hermox:latest /usr/local/bin/hermox /usr/local/bin/hermox

WORKDIR /bragemusic-server
COPY --from=go-build /build/bragemusic-server ./bragemusic-server

RUN apk add --no-cache \
    libstdc++ \
    portaudio \
    alsa-lib \
    jack \
    chromaprint \
    imagemagick \
    imagemagick-libs \
    libjpeg-turbo \
    libpng \
    libwebp \
    libheif \
    libde265

RUN apk add --no-cache tini

ARG BM_PATHS_MUSIC_DIR="/bragemusic/data/music"
ARG BM_PATHS_IMAGE_DIR="/bragemusic/data/img"
ARG BM_PATHS_CONFIG_DIR="/bragemusic/data/config"
ARG BM_PATHS_IMPORT_DIR="/bragemusic/importDir"
ARG BM_PATHS_MANUAL_IMPORT_DIR="/bragemusic/manualImportDir"
ARG BM_PATHS_BACKUP_IMPORT_DIR="/bragemusic/doneImportsDir"
ARG BM_ADMIN_USERNAME="Admin"
ARG BM_ADMIN_EMAIL="admin@example.com"
ARG BM_ADMIN_PASSWORD="password"

ENV BM_PATHS_MUSIC_DIR=$BM_PATHS_MUSIC_DIR
ENV BM_PATHS_IMAGE_DIR=$BM_PATHS_IMAGE_DIR
ENV BM_PATHS_CONFIG_DIR=$BM_PATHS_CONFIG_DIR
ENV BM_PATHS_IMPORT_DIR=$BM_PATHS_IMPORT_DIR
ENV BM_PATHS_MANUAL_IMPORT_DIR=$BM_PATHS_MANUAL_IMPORT_DIR
ENV BM_PATHS_BACKUP_IMPORT_DIR=$BM_PATHS_BACKUP_IMPORT_DIR
ENV BM_ADMIN_USERNAME=$BM_ADMIN_USERNAME
ENV BM_ADMIN_EMAIL=$BM_ADMIN_EMAIL
ENV BM_ADMIN_PASSWORD=$BM_ADMIN_PASSWORD

COPY --chmod=555 entrypoint.sh ./entrypoint.sh

ENTRYPOINT ["/sbin/tini", "--", "./entrypoint.sh"]
