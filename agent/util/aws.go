package util

import "strings"

// IsAWSBedrockModel checks if the model is an AWS Bedrock model
// Currently, support for Claude 3.5 or 3.7
// TODO: is this util function?
func IsAWSBedrockModel(model string) bool {
	return strings.Contains(model, "anthropic.claude-3-")
}
