package amazonsqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
)

// SQSSendMessageAPI Basic interface to send messages through SQS.
type SQSSendMessageAPI interface {
	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// SQSSendMessageClient Client specific for SQS using aws sdk v2.
type SQSSendMessageClient struct {
	client *sqs.Client
}

// SendMessage Send a message to SQS.
func (s SQSSendMessageClient) SendMessage(ctx context.Context,
	params *sqs.SendMessageInput,
	optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return s.client.SendMessage(ctx, params)
}

// AmazonSQS Basic structure with SQS information
type AmazonSQS struct {
	client            *sqs.Client
	sendMessageClient SQSSendMessageAPI
	queueNames        []string
}

// New creates a new AmazonSQS
func New(accessKeyID, secretKey, region string) (*AmazonSQS, error) {
	credProvider := credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, "")

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(credProvider),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)
	return &AmazonSQS{
		client:            client,
		sendMessageClient: SQSSendMessageClient{client: client},
	}, nil
}

// AddReceivers takes queue urls and adds them to the internal queues
// list. The Send method will send a given message to all those
// queues.
func (s *AmazonSQS) AddReceivers(queues ...string) {
	s.queueNames = append(s.queueNames, queues...)
}

// Send takes a message subject and a message body and sends them to
// all previously set queues.  This method is not atomic, so if one of
// the messages fails to be sent the other one may already have been
// published.
func (s AmazonSQS) Send(ctx context.Context, subject, message string) error {
	// Appends the subject and the message separated by a line break.
	sqsMessage := subject + "\n" + message
	// Loop through all configured queues to publish the same message.
	for _, queue := range s.queueNames {
		// Creates a SQS input with the queue data.
		input := &sqs.SendMessageInput{
			QueueUrl:    aws.String(queue),
			MessageBody: aws.String(sqsMessage),
		}
		_, err := s.sendMessageClient.SendMessage(ctx, input)

		if err != nil {
			return errors.Wrapf(err, "failed to send message using Amazon SQS to queue '%s'", queue)
		}
	}

	return nil
}
