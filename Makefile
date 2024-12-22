TEST_DIR := ./agent/functions

.PHONY: test

test:
	cd $(TEST_DIR) && go test -v
