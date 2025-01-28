FROM debian:bookworm-slim@sha256:f70dc8d6a8b6a06824c92471a1a258030836b26b043881358b967bf73de7c5ab AS release

RUN apt-get update && apt-get install -y git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Use binaly built by goreleaser
COPY agent /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
