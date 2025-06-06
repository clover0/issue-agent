TEMPDIR ?= .tmp

DIST ?= $(TEMPDIR)/build/agent

DEV_IMAGE ?= ghcr.io/clover0/issue-agent
DEV_IMAGE_TAG ?= dev
DEV_RELEASE_IMAGE_TAG ?= dev-release

$(TEMPDIR):
	mkdir -p $(TEMPDIR)

$(DIST):
	mkdir -p $(dir $(DIST))

clean:
	rm -rf $(TEMPDIR)

build: $(TEMPDIR)
	cd cmd/runner && CGO_ENABLED=0 go build -o ../../$(TEMPDIR)/build/agent

build/agent:
	mkdir -p $(dir $(DIST))
	CGO_ENABLED=0 go build -o $(DIST) ./cmd/agent/

run/test:
	go test ./... -v -race

image/dev:
	docker build -t $(DEV_IMAGE):$(DEV_IMAGE_TAG) -f dev.Dockerfile .

# for testing release image on development
image/dev-release: $(DIST) image/dev
	$(eval CONTAINER_NAME = issu-agent-dev-built-container)
	docker rm -f $(CONTAINER_NAME) || true
	docker run -it --name $(CONTAINER_NAME) $(DEV_IMAGE):$(DEV_IMAGE_TAG) version
	docker cp $(CONTAINER_NAME):/agent/bin/agent $(DIST)
	cd $(dir $(DIST)) && \
	docker build -t $(DEV_IMAGE):$(DEV_RELEASE_IMAGE_TAG) . -f ../../release.Dockerfile
