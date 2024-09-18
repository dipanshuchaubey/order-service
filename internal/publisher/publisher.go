package publisher

import (
	"context"

	"github.com/google/wire"
)

// ProviderSet is publisher providers.
var ProviderSet = wire.NewSet(NewPublisher)

// Publisher interface
type PublisherInterface interface {
	PublishOrderEvents(ctx context.Context, eventName, eventData, orderId string) error
}
