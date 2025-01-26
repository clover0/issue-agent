package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util/pointer"
)

type BedrockClient struct {
	client *bedrockruntime.Client
	logger logger.Logger

	// services
	Messages *BedrockMessageService
}

func NewBedrock(logger logger.Logger) (BedrockClient, error) {
	ctx := context.Background()
	var opts []func(*config.LoadOptions) error
	sdkConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return BedrockClient{}, fmt.Errorf("failed to load AWS SDK config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(sdkConfig)

	c := BedrockClient{
		logger: logger,
		client: client,
	}

	c.Messages = &BedrockMessageService{client: &c}

	return c, nil
}

type BedrockMessageService struct {
	client *BedrockClient
}

type BedrockConverseMessageResponse struct {
	Value string
	Role  MessageRole
}

func (s *BedrockMessageService) Create(
	ctx context.Context,
	modelID string,
	systemMessage string,
	messages []types.Message,
	toolSpecs []*types.ToolMemberToolSpec) (response *bedrockruntime.ConverseOutput, _ error) {
	input := &bedrockruntime.ConverseInput{
		// todo: changeable models
		ModelId: aws.String(modelID),
		InferenceConfig: &types.InferenceConfiguration{
			Temperature: pointer.Float32(0),
		},
		System:     []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: systemMessage}},
		Messages:   messages,
		ToolConfig: &types.ToolConfiguration{},
	}
	for _, tool := range toolSpecs {
		input.ToolConfig.Tools = append(input.ToolConfig.Tools, tool)
	}

	result, err := s.client.client.Converse(ctx, input)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "provided model identifier is invalid") {
			return response, fmt.Errorf("failed to invoke model: %w: hint - check whether enabled the model and in the AWS region", err)
		}
		return response, fmt.Errorf("failed to invoke model: %w", err)
	}

	return result, nil
}
