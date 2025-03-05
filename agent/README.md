# for Developer

## Development
### make image
After changing code, build container image
```sh
make image/dev
```


### Run
```sh
GITHUB_TOKEN=$(gh auth token) \
ANTHROPIC_API_KEY="key" \
OPENAI_API_KEY="key" \
  go run cmd/runner/main.go create-pr clover0/example-repository/issues/123 \
     --base_branch main \
    --model claude-3-5-sonnet-latest \
     --language Japanese \
     --log_level debug
```


## Run using release image
```sh
# build image using Dockerfile on actual release
make image/dev-release

# run 
GITHUB_TOKEN=$(gh auth token) \
ANTHROPIC_API_KEY="key" \
OPENAI_API_KEY="key" \
  go run -ldflags "-X main.containerImageTag=dev-release" clover0/example-repository/issues/123 \
     --base_branch main \
     --model claude-3-5-sonnet-latest \
     --language Japanese \
     --log_level debug
```
