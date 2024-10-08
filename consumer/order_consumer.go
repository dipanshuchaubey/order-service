package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/consumer/handler"
	"order-service/internal/conf"
	"order-service/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kratos/kratos/v2/log"
)

type OrderConsumer struct {
	log     *log.Helper
	handler handler.SyncOrderInterface
	conf    *conf.Consumer
}

func NewOrderConsumer(conf *conf.Consumer, handler handler.SyncOrderInterface, logger log.Logger) (*OrderConsumer, error) {
	return &OrderConsumer{
		log:     log.NewHelper(logger),
		handler: handler,
		conf:    conf,
	}, nil
}

func (c *OrderConsumer) Consume() error {
	ctx := context.TODO()
	c.log.WithContext(ctx).Info("OrderConsumer:: Starting to consume messages from queue...")

	// Create SQS session =--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=
	must := utils.CreateAWSSession(c.conf.OrderConsumer.Region)
	svc := sqs.New(must)
	c.log.WithContext(ctx).Info("OrderConsumer:: Session created successfully...")

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
		c.log.WithContext(ctx).Infof("Messages read from queue: %d", len(msgResult.Messages))
		for _, message := range msgResult.Messages {
			var body handler.MessageBody
			decodeErr := json.Unmarshal([]byte(*message.Body), &body)
			if decodeErr != nil {
				errMsg := fmt.Sprintf("Error decoding message body: %s", decodeErr)
				c.log.WithContext(ctx).Errorf(errMsg)
				continue
			}

			var data handler.MessageData
			decodeErr = json.Unmarshal([]byte(body.Message), &data)
			if decodeErr != nil {
				errMsg := fmt.Sprintf("Error decoding message data: %s", decodeErr)
				c.log.WithContext(ctx).Errorf(errMsg)
				continue
			}

			// Process message =--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=
			err := c.handler.Handler(ctx, body.MessageId, data)
			if err != nil {
				c.log.WithContext(ctx).Errorf("Error processing message: %s", err)
				continue
			}

			ack = append(ack, &sqs.DeleteMessageBatchRequestEntry{Id: &body.MessageId, ReceiptHandle: message.ReceiptHandle})
		}

		if len(ack) > 0 {
			// Acknowledge messages =--=--=--=--=--=--=--=--=--=--=--=--=--=--=
			_, ackErr := c.acknowledgeQueueMessages(svc, ack)
			if ackErr != nil {
				errMsg := fmt.Sprintf("Error deleting messages from queue: %s", ackErr)
				c.log.WithContext(ctx).Errorf(errMsg)
				return ackErr
			}

			c.log.Infof("Acknowledge messages from queue: %d", len(ack))
		}
	}
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
		VisibilityTimeout:   aws.Int64(c.conf.OrderConsumer.WaitTime),
	})
}

func (c *OrderConsumer) acknowledgeQueueMessages(svc *sqs.SQS, ackMsg []*sqs.DeleteMessageBatchRequestEntry) (*sqs.DeleteMessageBatchOutput, error) {
	return svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(c.conf.OrderConsumer.QueueUrl),
		Entries:  ackMsg,
	})
}
