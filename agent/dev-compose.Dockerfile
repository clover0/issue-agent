# For development with docker compose.

FROM golang:1.24.0-bookworm@sha256:6260304a09fb81a1983db97c9e6bfc1779ebce33d39581979a511b3c7991f076 AS devlopment

WORKDIR /usr/local/agent
