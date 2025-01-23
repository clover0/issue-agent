TEMPDIR ?= .tmp

$(TEMPDIR):
	mkdir -p $(TEMPDIR)

build: $(TEMPDIR)
	cd cmd/runner && go build -o ../../$(TEMPDIR)/build/agent

build/agent:
	cd cmd/agent && go build -o $(DIST)

.PHONY: test
run/test:
	go test ./... -v -parallel 4 -race

image/dev:
	docker build -t ghcr.io/clover0/issue-agent:dev -f dev.Dockerfile .

#image/release:
#	echo "See .github/workflows/release.yml"
