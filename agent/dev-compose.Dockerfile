# For development with docker compose.

FROM golang:1.24.1-bookworm@sha256:fa1a01d362a7b9df68b021d59a124d28cae6d99ebd1a876e3557c4dd092f1b1d AS devlopment

WORKDIR /usr/local/agent
