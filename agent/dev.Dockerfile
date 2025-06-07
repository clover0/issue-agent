# This Dockerfile is used to build the development image for the agent.
# It is used to build the agent binary in this container and run at runner binary using this container.
# From CLI command to run this agent container as End to End test.

FROM golang:1.24.4-bookworm@sha256:c83619bb18b0207412fffdaf310f57ee3dd02f586ac7a5b44b9c36a29a9d5122 AS development

WORKDIR /agent/src/

COPY go.mod go.sum /agent/src/

RUN go mod download

COPY . .

# TODO: SLSA
RUN cd /agent/src/cmd/agent && \
    CGO_ENABLED=0 go build \
      -ldflags "-X github.com/clover0/issue-agent/cli/command/version.version=dev-$(date -u +'%Y-%m-%dT%H%M%SZ')" \
      -o /agent/bin/agent


FROM gcr.io/distroless/static-debian12@sha256:d9f9472a8f4541368192d714a995eb1a99bab1f7071fc8bde261d7eda3b667d8

ENV PATH="/agent/bin:$PATH"

COPY --from=development /agent/bin/agent /agent/bin/

ENTRYPOINT ["agent"]
