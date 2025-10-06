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
FROM opensuse/leap:42.3 as run
COPY --from=build /app/azure-dyndns2 /app/azure-dyndns2
RUN zypper refresh
RUN zypper --non-interactive update
WORKDIR /app
EXPOSE 8080
CMD ["/app/azure-dyndns2","serve"]
