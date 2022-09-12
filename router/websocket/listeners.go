package websocket

import (
	"github.com/luminous-gsm/fusion/event"
	eventModel "github.com/vmware/transport-go/model"
)

func (websock *WebsocketController) EventDockerPodCreateListener(m *eventModel.Message) {
	event := m.Payload.(event.FusionEvent[event.FusionDockerEventData])
	websock.SendJSON(event)
}
