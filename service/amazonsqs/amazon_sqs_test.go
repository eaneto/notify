package amazonsqs

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SQSSendMessageMock struct {
	mock.Mock
}

// SendMessage provides a mock for sending a message to SQS.
func (m *SQSSendMessageMock) SendMessage(ctx context.Context,
	params *sqs.SendMessageInput,
	optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

func TestAddReceivers(t *testing.T) {
	amazonSQS, error := New("", "", "")

	if error != nil {
		t.Error(error)
	}

	amazonSQS.AddReceivers("One receiver")
}

func TestSendMessageWithNoQueuesConfigured(t *testing.T) {
	mockSqs := new(SQSSendMessageMock)

	amazonSQS := AmazonSQS{
		sendMessageClient: mockSqs,
	}

	err := amazonSQS.Send(context.Background(), "i'm the subject", "i'm the messsage")

	assert.Nil(t, err)

	mockSqs.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendMessageWithSucessAndOneQueueConfigured(t *testing.T) {
	mockSqs := new(SQSSendMessageMock)
	output := sqs.SendMessageOutput{}
	mockSqs.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, nil)

	amazonSQS := AmazonSQS{
		sendMessageClient: mockSqs,
	}

	amazonSQS.AddReceivers("http://aws-sqs-my-queue-name.aws.sqs/queue-name")
	err := amazonSQS.Send(context.Background(), "i'm the subject", "i'm the messsage")

	assert.Nil(t, err)

	mockSqs.AssertExpectations(t)
	assert.Equal(t, 1, len(mockSqs.Calls))
}

func TestSendMessageWithSucessAndTwoQueuesConfigured(t *testing.T) {
	mockSqs := new(SQSSendMessageMock)
	output := sqs.SendMessageOutput{}
	mockSqs.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, nil)

	amazonSQS := AmazonSQS{
		sendMessageClient: mockSqs,
	}

	amazonSQS.AddReceivers(
		"http://aws-sqs-my-queue-name.aws.sqs/queue-name",
		"http://aws-sqs-my-queue-name.aws.sqs/queue-name",
	)
	err := amazonSQS.Send(context.Background(), "i'm the subject", "i'm the messsage")

	assert.Nil(t, err)

	mockSqs.AssertExpectations(t)
	assert.Equal(t, 2, len(mockSqs.Calls))
}

func TestSendMessageWithErrorAndOneQueueConfiguredShouldReturnError(t *testing.T) {
	mockSqs := new(SQSSendMessageMock)
	output := sqs.SendMessageOutput{}
	mockSqs.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, errors.New("Error on SQS"))

	amazonSQS := AmazonSQS{
		sendMessageClient: mockSqs,
	}

	amazonSQS.AddReceivers("http://aws-sqs-my-queue-name.aws.sqs/queue-name")
	err := amazonSQS.Send(context.Background(), "i'm the subject", "i'm the messsage")

	assert.NotNil(t, err)

	mockSqs.AssertExpectations(t)
	assert.Equal(t, 1, len(mockSqs.Calls))
}

func TestSendMessageWithErrorAndTwoQueuesConfiguredShouldReturnErrorOnFirst(t *testing.T) {
	mockSqs := new(SQSSendMessageMock)
	output := sqs.SendMessageOutput{}
	mockSqs.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, errors.New("Error on SQS"))

	amazonSQS := AmazonSQS{
		sendMessageClient: mockSqs,
	}

	amazonSQS.AddReceivers(
		"http://aws-sqs-my-queue-name.aws.sqs/queue-name",
		"http://aws-sqs-my-queue-name.aws.sqs/queue-name",
	)
	err := amazonSQS.Send(context.Background(), "i'm the subject", "i'm the messsage")

	assert.NotNil(t, err)

	mockSqs.AssertExpectations(t)
	assert.Equal(t, 1, len(mockSqs.Calls))
}
