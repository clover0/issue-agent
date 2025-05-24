FROM gcr.io/distroless/static-debian12@sha256:d9f9472a8f4541368192d714a995eb1a99bab1f7071fc8bde261d7eda3b667d8

# Use binaly built by goreleaser
COPY agent /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
