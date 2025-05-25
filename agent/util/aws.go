package util

import "strings"

// IsAWSBedrockModel checks if the model is an AWS Bedrock model
// TODO: is this util function?
func IsAWSBedrockModel(model string) bool {
	return strings.Contains(model, "anthropic.claude-")
}
