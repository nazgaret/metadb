package notifier

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/metadb-project/metadb/cmd/metadb/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

const matchSNSRegion = `^(?:[^:]+:){3}([^:]+).*`

type Notifier interface {
	Send(ctx context.Context, message string) error
}

func NewNoop() *noopNotifier {
	return &noopNotifier{}
}

type noopNotifier struct {
}

func (n *noopNotifier) Send(ctx context.Context, message string) error {
	return nil
}

type snsNotifier struct {
	client *sns.Client
	topic  string
}

func NewSNS(topic string) (*snsNotifier, error) {
	if len(topic) == 0 {
		return nil, errors.New("no topic provided for SNS")
	}

	// Get AWS Region from topic ARN
	re := regexp.MustCompile(matchSNSRegion)
	match := re.FindStringSubmatch(topic)
	if len(match) < 2 {
		return nil, fmt.Errorf("unable to get region from topic arn: %s", topic)
	}

	// Load credentials: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(match[1]))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	snsClient := sns.NewFromConfig(cfg)
	return &snsNotifier{snsClient, topic}, nil
}

func (n *snsNotifier) Send(ctx context.Context, message string) error {
	if n == nil || n.client == nil {
		return nil
	}

	res, err := n.client.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(n.topic),
	})
	if err != nil {
		return nil
	}

	log.Debug("message sent to SNS: %s", *res.MessageId)

	return nil
}

//
// Example of creating new implementation of Notifier:
//
// func NewFoo() *fooNotifier {
//	// init foo client
// 	return &fooNotifier{}
// }
//
// type fooNotifier struct {
//	// foo client
// }
//
// func (n *fooNotifier) Send(ctx context.Context, message string) error {
//  // sending notification using foo client
// 	return nil
// }
