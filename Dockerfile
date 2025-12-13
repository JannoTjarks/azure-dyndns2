# ---------- Stage 1: Build ----------
FROM golang:1.25 as build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o azure-dyndns2 .
 
# ---------- Stage 2: Final ----------
FROM opensuse/leap:16.0 as run
LABEL org.opencontainers.image.source=https://github.com/JannoTjarks/azure-dyndns2
LABEL org.opencontainers.image.description="Simple dyndns2-compatible web api for Azure DNS"
LABEL org.opencontainers.image.licenses=Apache-2.0

RUN groupadd nonroot \
    && useradd nonroot -g nonroot
COPY --from=build /app/azure-dyndns2 /app/azure-dyndns2
RUN zypper refresh \
    && zypper --non-interactive update
USER nonroot
WORKDIR /app
EXPOSE 8080
CMD ["/app/azure-dyndns2","serve"]
