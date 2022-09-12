package websocket

import (
	"encoding/json"

	"github.com/luminous-gsm/fusion/event"
	eventModel "github.com/vmware/transport-go/model"
)

func (websock *WebsocketController) EventDockerPodCreateListener(m *eventModel.Message) {
	event := m.Payload.(event.FusionEvent[event.FusionDockerEventData])
	data, err := json.Marshal(event)
	if err != nil {
		websock.log().Errorw("could not marshal event", "error", err)
	}
	websock.SendJSON(string(data))
}

// func (websock *WebsocketController) EventDockerPodCreateListener(event event.FusionEvent[event.FusionDockerEventData]) {
// 	zap.S().Debugw("websocket: received event", "eventData", event)
// 	data, err := json.Marshal(event)
// 	if err != nil {
// 		zap.S().Errorw("websocket: could not marshal event object to JSON", "error", err)
// 	}
// 	websock.SendJSON(string(data))
// }
