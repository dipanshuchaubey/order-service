package publisher

import (
	"github.com/google/wire"
)

// ProviderSet is publisher providers.
var ProviderSet = wire.NewSet(NewPublisher)

// Publisher interface
type PublisherInterface interface {
	PublishOrderEvents(eventName, eventData, orderId string) error
}
