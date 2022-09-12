package websocket

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"github.com/luminous-gsm/fusion/event"
	"go.uber.org/zap"
)

type WebsocketController struct {
	Connection *ws.Conn
	Id         uuid.UUID
	Mu         sync.Mutex
}

func GetHandler(w http.ResponseWriter, r *http.Request) (*WebsocketController, error) {
	upgrader := ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	socketConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Errorw("websocket: could not upgrade connection to websocket", "error", err)
		return nil, err
	}

	connectionId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &WebsocketController{
		Connection: socketConnection,
		Id:         connectionId,
	}, nil

}

func (websock *WebsocketController) log() *zap.SugaredLogger {
	return zap.S().Named("websocket")
}

// Register/Subscribe to server events to send them them to the connected websocket client
func (websock *WebsocketController) RegisterServerEventListeners() error {
	websock.log().Debugw("websocket: register server events for websocket client", "id", websock.Id)

	handler, err := event.Instance().Bus().ListenRequestStream(event.EVENT_DOCKER_POD_CREATE)
	if err != nil {
		websock.log().Panicw("could not create handler for channel", "error", err, "channel", event.EVENT_DOCKER_POD_CREATE)
	}
	handler.Handle(
		websock.EventDockerPodCreateListener,
		event.Instance().DefaultErrorHandler,
	)
	return nil
}

func (websock *WebsocketController) UnregisterServerEventListeners() {
	zap.S().Debugw("websocket: unregister server events for websocket client", "id", websock.Id)

	// if err := event.Instance().Unsubscribe(event.EVENT_DOCKER_POD_CREATE, websock.EventDockerPodCreateListener); err != nil {
	// 	zap.S().Errorw("websocket: could not unregister event listener. No listener registered", "event", event.EVENT_DOCKER_POD_CREATE, "error", err)
	// }
}

func (controller *WebsocketController) SendJSON(v interface{}) error {
	controller.Mu.Lock()
	defer controller.Mu.Unlock()

	return controller.Connection.WriteJSON(v)
}

// Handle all events/messages coming from websocket
func (controller *WebsocketController) HandleEvent(eventMessage event.FusionEvent[map[string]interface{}]) error {
	zap.S().Debugw("websocket: new event", "event", eventMessage)
	if event.Instance().Bus().GetChannelManager().CheckChannelExists(eventMessage.Event) {
		event.FireEvent(
			eventMessage.Event,
			eventMessage,
		)

	} else {
		controller.log().Errorw("no channel registered on the event bus", "channel", eventMessage.Event)
	}

	return nil
}
