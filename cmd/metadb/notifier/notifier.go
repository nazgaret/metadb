package notifier

import (
	"context"
	"fmt"
	"regexp"

	"github.com/metadb-project/metadb/cmd/metadb/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

const matchRegion = `^(?:[^:]+:){3}([^:]+).*`

type Notifier struct {
	client *sns.Client
	topic  string
}

func Init(topic string) (*Notifier, error) {
	if len(topic) == 0 {
		return nil, nil
	}

	re := regexp.MustCompile(matchRegion)
	match := re.FindStringSubmatch(topic)
	if len(match) < 2 {
		return nil, fmt.Errorf("unable to get region from topic arn: %s", topic)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(match[1]))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)
	return &Notifier{snsClient, topic}, nil
}

func (n *Notifier) Send(ctx context.Context, message string) error {
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
