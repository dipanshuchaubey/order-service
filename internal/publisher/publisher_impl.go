package publisher

import (
	"encoding/json"
	"order-service/internal/conf"
	"order-service/internal/constants"
	"order-service/internal/utils"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/go-kratos/kratos/v2/log"
)

type PublisherHandler struct {
	log  *log.Helper
	conf *conf.Publisher
	sns  *sns.SNS
}

func NewPublisher(logger log.Logger, conf *conf.Publisher) PublisherInterface {
	sess := utils.CreateAWSSession(conf.OrderPublisher.Region)
	svc := sns.New(sess)

	return PublisherHandler{log.NewHelper(logger), conf, svc}
}

func (p PublisherHandler) PublishOrderEvents(eventName, eventData, orderId string) error {
	message, msgErr := p.createMessage(eventName, eventData)
	if msgErr != nil {
		return msgErr
	}

	_, err := p.sns.Publish(&sns.PublishInput{
		Message:        &message,
		TopicArn:       &p.conf.OrderPublisher.TopicArn,
		MessageGroupId: &orderId,
	})
	if err != nil {
		p.log.Errorf("Error publishing message to SNS: %s", err)
		return err
	}

	p.log.Infof("Message published to SNS: %s", message)
	return nil
}

func (p PublisherHandler) createMessage(eventName, eventData string) (string, error) {
	message := struct {
		Event string `json:"event"`
		Data  string `json:"data"`
	}{
		Event: eventName,
		Data:  eventData,
	}

	byteMsg, err := json.Marshal(message)
	if err != nil {
		p.log.Errorf(constants.ErrorMarshalling, err)
		return constants.EmptyString, err
	}

	return string(byteMsg), nil
}
