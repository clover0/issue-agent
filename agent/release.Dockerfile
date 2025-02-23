FROM gcr.io/distroless/static-debian12

# Use binaly built by goreleaser
COPY agent /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
