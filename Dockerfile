# ---------- Stage 1: Build ----------
FROM golang:1.25.6 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./cmd ./cmd
COPY ./docs ./docs
COPY ./internal ./internal
COPY ./vendor ./vendor
COPY ./main.go .
RUN go build -trimpath -ldflags="-s -w" -o azure-dyndns2 .
 
# ---------- Stage 2: Final ----------
FROM opensuse/leap:16.0 AS run
LABEL org.opencontainers.image.source=https://github.com/JannoTjarks/azure-dyndns2
LABEL org.opencontainers.image.description="Simple dyndns2-compatible web api for Azure DNS"
LABEL org.opencontainers.image.licenses=AGPLv3

EXPOSE 8080
WORKDIR /app

RUN groupadd nonroot \
    && useradd nonroot -g nonroot \
    && zypper refresh \
    && zypper --non-interactive update

USER nonroot
COPY --from=build /app/azure-dyndns2 /app/azure-dyndns2
CMD ["/app/azure-dyndns2","serve"]
