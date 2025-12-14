# https://goreleaser.com/errors/docker-build/#dont
FROM opensuse/leap:16.0

ARG TARGETPLATFORM

EXPOSE 8080

RUN groupadd nonroot \
    && useradd nonroot -g nonroot \
    && zypper refresh \
    && zypper --non-interactive update

USER nonroot
COPY ${TARGETPLATFORM}/azure-dyndns2 /app/azure-dyndns2
CMD ["/app/azure-dyndns2","serve"]
