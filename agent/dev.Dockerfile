# This Dockerfile is used to build the development image for the agent.
# It is used to build the agent binary in this container and run at runner binary using this container.
# From CLI command to run this agent container as End to End test.

FROM golang:1.24.1-bookworm@sha256:d7d795d0a9f51b00d9c9bfd17388c2c626004a50c6ed7c581e095122507fe1ab AS development

WORKDIR /agent/src/

COPY go.mod go.sum /agent/src/

RUN go mod download

COPY . .

# TODO: SLSA
RUN cd /agent/src/cmd/agent && \
    CGO_ENABLED=0 go build \
      -ldflags "-X github.com/clover0/issue-agent/cli.version=dev-$(date -u +'%Y-%m-%dT%H%M%SZ')" \
      -o /agent/bin/agent


FROM gcr.io/distroless/static-debian12@sha256:3f2b64ef97bd285e36132c684e6b2ae8f2723293d09aae046196cca64251acac

ENV PATH="/agent/bin:$PATH"

COPY --from=development /agent/bin/agent /agent/bin/

ENTRYPOINT ["agent"]
