package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func CreateTopic(ctx context.Context,
	params *sns.CreateTopicInput,
	optFns ...func(*sns.Options)) (*sns.CreateTopicOutput, error) {
	return nil, nil
}
