# For development with docker compose.

FROM golang:1.23.6-bookworm@sha256:441f59f8a2104b99320e1f5aaf59a81baabbc36c81f4e792d5715ef09dd29355 AS devlopment

WORKDIR /usr/local/agent
