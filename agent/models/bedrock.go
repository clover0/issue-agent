package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	"github.com/clover0/issue-agent/core"
	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/util/pointer"
)

type BedrockClient struct {
	client *bedrockruntime.Client
	logger logger.Logger

	// services
	Messages *BedrockMessageService
}

type awsCustomRetryer struct {
	*retry.Standard
	logger logger.Logger
}

func (r *awsCustomRetryer) IsErrorRetryable(err error) bool {
	var v interface{ HTTPStatusCode() int }

	// custom error handling
	if errors.As(err, &v) {
		if v.HTTPStatusCode() == http.StatusTooManyRequests {
			r.logger.Info(fmt.Sprintf("%s\nRate limited, retrying after 60 seconds...\n", err))
			time.Sleep(60 * time.Second)
			return true
		}
	}

	return r.Standard.IsErrorRetryable(err)
}

func (r *awsCustomRetryer) MaxAttempts() int {
	return r.Standard.MaxAttempts()
}

func (r *awsCustomRetryer) RetryDelay(attempt int, opErr error) (time.Duration, error) {
	return r.Standard.RetryDelay(attempt, opErr)
}

func (r *awsCustomRetryer) GetRetryToken(ctx context.Context, opErr error) (releaseToken func(error) error, err error) {
	return r.Standard.GetRetryToken(ctx, opErr)
}

func (r *awsCustomRetryer) GetAttemptToken(ctx context.Context) (func(error) error, error) {
	return r.Standard.GetAttemptToken(ctx)
}

func NewBedrock(logger logger.Logger) (BedrockClient, error) {
	ctx := context.Background()
	stdRetryer := retry.NewStandard()
	sdkConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRetryer(func() aws.Retryer {
			return &awsCustomRetryer{Standard: stdRetryer, logger: logger}
		}))
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
	Role  core.MessageRole
}

func (s *BedrockMessageService) Create(
	ctx context.Context,
	modelID string,
	systemMessage string,
	messages []types.Message,
	toolSpecs []*types.ToolMemberToolSpec,
) (response *bedrockruntime.ConverseOutput, _ error) {
	input := &bedrockruntime.ConverseInput{
		// todo: changeable models
		ModelId: aws.String(modelID),
		InferenceConfig: &types.InferenceConfiguration{
			Temperature: pointer.Float32(0),
			MaxTokens:   pointer.Ptr(int32(ClaudeMaxOutputTokens(modelID))),
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
