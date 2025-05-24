# For development with docker compose.

FROM golang:1.24.3-bookworm@sha256:29d97266c1d341b7424e2f5085440b74654ae0b61ecdba206bc12d6264844e21 AS devlopment

WORKDIR /usr/local/agent
