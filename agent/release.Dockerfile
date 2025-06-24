FROM gcr.io/distroless/static-debian12@sha256:b7b9a6953e7bed6baaf37329331051d7bdc1b99c885f6dbeb72d75b1baad54f9

# Use binaly built by goreleaser
COPY agent /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
