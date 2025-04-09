# For development with docker compose.

FROM golang:1.24.2-bookworm@sha256:e719692f259f78b4496dbfe80628fbbef542da15314a24ddb98f26bac39833cf AS devlopment

WORKDIR /usr/local/agent
