# This Dockerfile is used to build the development image for the agent.
# It is used to build the agent binary in this container and run at runner binary using this container.
# From CLI command to run this agent container as End to End test.

FROM golang:1.23.5-bookworm@sha256:3149bc5043fa58cf127fd8db1fdd4e533b6aed5a40d663d4f4ae43d20386665f AS development

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
