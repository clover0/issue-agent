.PHONY: test
test:
	go test ./... -v

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
