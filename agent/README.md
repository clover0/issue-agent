# for Developer

## Test using release image
```
# change var to
# TODO: changeable from outside
var containerImageTag = "dev-release

# build image using Dockerfile on asctual release
make image/release

# run 
GITHUB_TOKEN=$(gh auth token) \
ANTHROPIC_API_KEY="key" \
OPENAI_API_KEY="key" \
  go run cmd/runner/main.go issue \
     --base_branch main \
     --github_issue_number 10 \
     --work_repository hobby \
     --github_owner clover0 \
     --model "gpt-4o" \
     --language Japanese \
     --log_level debug \
     --aws_region us-east-1
```
