package consumer

import (
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

type MessageData struct {
	Event      string `json:"event"`
	OrderID    string `json:"order_id"`
	CartID     string `json:"cart_id"`
	PaymentRef string `json:"payment_ref"`
	OrderTime  int32  `json:"order_time"`
}

type MessageBody struct {
	Type           string `json:"Type"`
	MessageId      string `json:"MessageId"`
	SequenceNumber string `json:"SequenceNumber"`
	TopicArn       string `json:"TopicArn"`
	Subject        string `json:"Subject"`
	Message        string `json:"Message"`
	Timestamp      string `json:"Timestamp"`
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
	must := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(c.conf.OrderConsumer.Region),
			Credentials: credentials.NewStaticCredentials(
				"", "", "",
			),
		},
	}))

	svc := sqs.New(must)
	c.log.Info("OrderConsumer:: Session created successfully...")

	for {
		msgResult, msgErr := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
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
			var body MessageBody
			decodeErr := json.Unmarshal([]byte(*message.Body), &body)
			if decodeErr != nil {
				errMsg := fmt.Sprintf("Error decoding message body: %s", decodeErr)
				c.log.Errorf(errMsg)
				continue
			}

			var data MessageData
			decodeErr = json.Unmarshal([]byte(body.Message), &data)
			if decodeErr != nil {
				errMsg := fmt.Sprintf("Error decoding message data: %s", decodeErr)
				c.log.Errorf(errMsg)
				continue
			}

			ack = append(ack, &sqs.DeleteMessageBatchRequestEntry{Id: &body.MessageId, ReceiptHandle: message.ReceiptHandle})
		}

		if len(ack) > 0 {
			_, ackErr := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
				QueueUrl: aws.String(c.conf.OrderConsumer.QueueUrl),
				Entries:  ack,
			})
			if ackErr != nil {
				errMsg := fmt.Sprintf("Error deleting messages from queue: %s", ackErr)
				c.log.Errorf(errMsg)
				return ackErr
			}

			c.log.Infof("Acknowledge messages from queue: %d", len(ack))
		}
	}
}
