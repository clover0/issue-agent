# This Dockerfile is used to build the development image for the agent.
# It is used to build the agent binary in this container and run at runner binary using this container.
# From CLI command to run this agent container as End to End test.

FROM golang:1.24.2-bookworm@sha256:00eccd446e023d3cd9566c25a6e6a02b90db3e1e0bbe26a48fc29cd96e800901 AS development

WORKDIR /agent/src/

COPY go.mod go.sum /agent/src/

RUN go mod download

COPY . .

# TODO: SLSA
RUN cd /agent/src/cmd/agent && \
    CGO_ENABLED=0 go build \
      -ldflags "-X github.com/clover0/issue-agent/cli/command/version.version=dev-$(date -u +'%Y-%m-%dT%H%M%SZ')" \
      -o /agent/bin/agent


FROM gcr.io/distroless/static-debian12@sha256:3d0f463de06b7ddff27684ec3bfd0b54a425149d0f8685308b1fdf297b0265e9

ENV PATH="/agent/bin:$PATH"

COPY --from=development /agent/bin/agent /agent/bin/

ENTRYPOINT ["agent"]
