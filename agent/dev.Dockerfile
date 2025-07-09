# This Dockerfile is used to build the development image for the agent.
# It is used to build the agent binary in this container and run at runner binary using this container.
# From CLI command to run this agent container as End to End test.

FROM golang:1.24.5-bookworm@sha256:69adc37c19ac6ef724b561b0dc675b27d8c719dfe848db7dd1092a7c9ac24bc6 AS development

WORKDIR /agent/src/

COPY go.mod go.sum /agent/src/

RUN go mod download

COPY . .

# TODO: SLSA
RUN cd /agent/src/cmd/agent && \
    CGO_ENABLED=0 go build \
      -ldflags "-X github.com/clover0/issue-agent/cli/command/version.version=dev-$(date -u +'%Y-%m-%dT%H%M%SZ')" \
      -o /agent/bin/agent


FROM gcr.io/distroless/static-debian12@sha256:b7b9a6953e7bed6baaf37329331051d7bdc1b99c885f6dbeb72d75b1baad54f9

ENV PATH="/agent/bin:$PATH"

COPY --from=development /agent/bin/agent /agent/bin/

ENTRYPOINT ["agent"]
