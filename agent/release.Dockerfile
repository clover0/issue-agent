FROM gcr.io/distroless/static-debian12@sha256:d75cdd72874d4790092fcb1b058493ecf6bb5bf2b2b897045b00ff01d91843f2

# Use binaly built by goreleaser
COPY agent /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
