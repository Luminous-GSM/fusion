package event

import (
	"encoding/json"
	"sync"

	"github.com/vmware/transport-go/bus"
	"go.uber.org/zap"
)

var (
	_conce       sync.Once
	eventService *EventService
)

type EventService struct{}

func InitEventBus() *EventService {
	_conce.Do(func() {
		eventService = &EventService{}
		zap.S().Info("eventbus: event bus configured")
	})

	return eventService
}

func (es EventService) InitEventChannels() {
	bus.GetBus().GetChannelManager().CreateChannel(EVENT_REQUEST_POD_CREATE)
	bus.GetBus().GetChannelManager().CreateChannel(EVENT_DOCKER_POD_CREATE)
	bus.GetBus().GetChannelManager().CreateChannel(EVENT_NODE_WARNING)

}

// Returns a Event Bus instance.
func Instance() *EventService {
	return eventService
}

func (es EventService) Bus() bus.EventBus {
	return bus.GetBus()
}

func UnmarshalUnknown(data map[string]interface{}, obj any) error {
	byteArray, err := json.Marshal(data)
	if err != nil {
		zap.S().Errorw("event: could not marshal map[string]interface{}", "error", err)
		return err
	}

	err = json.Unmarshal(byteArray, &obj)
	if err != nil {
		zap.S().Errorw("event: could not unmarshal byte array", "error", err)
		return err
	}

	return nil
}

func (es EventService) log() *zap.SugaredLogger {
	return zap.S().Named("event")
}

func (es EventService) DefaultErrorHandler(err error) {
	es.log().Errorw("error received on channel", "error", err)
}

func FireEvent[T FusionEventData](topic string, event FusionEvent[T]) {
	handler, err := bus.GetBus().RequestStream(
		topic,
		event,
	)
	if err != nil {
		zap.S().Named("event").Errorw("error requesting to stream", "error", err)
	}
	err = handler.Fire()
	if err != nil {
		zap.S().Named("event").Errorw("error requesting to stream", "error", err)
	}
}

func (es EventService) ListenWithPanic(channel string) bus.MessageHandler {
	handler, err := es.Bus().ListenRequestStream(channel)
	if err != nil {
		es.log().Panicw("could not create handler for channel", "error", err, "channel", channel)
	}
	return handler
}
