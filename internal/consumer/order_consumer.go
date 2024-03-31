package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/biz/order_sync"
	"order-service/internal/conf"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kratos/kratos/v2/log"
)

type OrderConsumer struct {
	log     *log.Helper
	handler order_sync.ISyncOrderHandler
	conf    *conf.Consumer
}

func NewOrderConsumer(conf *conf.Consumer, handler order_sync.ISyncOrderHandler, logger log.Logger) (*OrderConsumer, error) {
	return &OrderConsumer{
		log:     log.NewHelper(logger),
		handler: handler,
		conf:    conf,
	}, nil
}

func (c *OrderConsumer) Consume() error {
	c.log.Info("OrderConsumer:: Starting to consume messages from queue...")
	ctx := context.TODO()

	// Create SQS session =--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=
	must := c.createSQSSession()
	svc := sqs.New(must)
	c.log.Info("OrderConsumer:: Session created successfully...")

	for {
		// Read messages from queue =--=--=--=--=--=--=--=--=--=--=--=--=
		msgResult, msgErr := c.readMessagesFromQueue(svc)
		if msgErr != nil {
			fmt.Println("Error reading message from queue: ", msgErr)
			return msgErr
		}
		if len(msgResult.Messages) == 0 {
			continue
		}

		ack := make([]*sqs.DeleteMessageBatchRequestEntry, 0)
		c.log.Infof("Messages read from queue: %d", len(msgResult.Messages))
		for _, message := range msgResult.Messages {
			var body order_sync.MessageBody
			decodeErr := json.Unmarshal([]byte(*message.Body), &body)
			if decodeErr != nil {
				errMsg := fmt.Sprintf("Error decoding message body: %s", decodeErr)
				c.log.Errorf(errMsg)
				continue
			}

			var data order_sync.MessageData
			decodeErr = json.Unmarshal([]byte(body.Message), &data)
			if decodeErr != nil {
				errMsg := fmt.Sprintf("Error decoding message data: %s", decodeErr)
				c.log.Errorf(errMsg)
				continue
			}

			// Process message =--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=
			c.handler.Handler(ctx, body.MessageId, data)

			ack = append(ack, &sqs.DeleteMessageBatchRequestEntry{Id: &body.MessageId, ReceiptHandle: message.ReceiptHandle})
		}

		if len(ack) > 0 {
			// Acknowledge messages =--=--=--=--=--=--=--=--=--=--=--=--=--=--=
			_, ackErr := c.acknowledgeQueueMessages(svc, ack)
			if ackErr != nil {
				errMsg := fmt.Sprintf("Error deleting messages from queue: %s", ackErr)
				c.log.Errorf(errMsg)
				return ackErr
			}

			c.log.Infof("Acknowledge messages from queue: %d", len(ack))
		}
	}
}

func (c *OrderConsumer) createSQSSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(c.conf.OrderConsumer.Region),
			Credentials: credentials.NewEnvCredentials(),
		},
	}))
}

func (c *OrderConsumer) readMessagesFromQueue(svc *sqs.SQS) (*sqs.ReceiveMessageOutput, error) {
	return svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(c.conf.OrderConsumer.QueueUrl),
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(int64(c.conf.OrderConsumer.WaitTime)),
	})
}

func (c *OrderConsumer) acknowledgeQueueMessages(svc *sqs.SQS, ackMsg []*sqs.DeleteMessageBatchRequestEntry) (*sqs.DeleteMessageBatchOutput, error) {
	return svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(c.conf.OrderConsumer.QueueUrl),
		Entries:  ackMsg,
	})
}
