# This Dockerfile is used to build the development image for the agent.
# It is used to build the agent binary in this container and run at runner binary using this container.
# From CLI command to run this agent container as End to End test.

FROM golang:1.23.5-bookworm@sha256:9c79a16e024bcfb856b6d063cf7ed9a6257f554466761f5a99b09787d2b55fbd AS development

WORKDIR /agent/src/

COPY go.mod go.sum /agent/src/

RUN go mod download

COPY . .

# TODO: SLSA
RUN cd /agent/src/cmd/agent && \
    go build \
      -ldflags "-X github.com/clover0/issue-agent/cli.version=dev-$(date -u +'%Y-%m-%dT%H%M%SZ')" \
      -o /agent/bin/agent

ENV PATH="/agent/bin:$PATH"

ENTRYPOINT ["/agent/bin/agent"]
