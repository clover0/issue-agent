# For development with docker compose.

FROM golang:1.24.2-bookworm@sha256:00eccd446e023d3cd9566c25a6e6a02b90db3e1e0bbe26a48fc29cd96e800901 AS devlopment

WORKDIR /usr/local/agent
