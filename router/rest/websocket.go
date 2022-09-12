package rest

import (
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	"github.com/luminous-gsm/fusion/event"
	"github.com/luminous-gsm/fusion/router/websocket"
	"go.uber.org/zap"
)

func RunWebSocket(c *gin.Context) {
	// s := middleware.GetServerManager(c)

	// ctx, cancel := context.WithCancel(c.Request.Context())
	// defer cancel()

	websockController, err := websocket.GetHandler(c.Writer, c.Request)
	if err != nil {
		NewError(err).Abort(c)
		return
	}
	defer websockController.Connection.Close()

	if err := websockController.RegisterServerEventListeners(); err != nil {
		zap.S().Error("websocket: could not register event listeners for websocket connection", "error", err)
		websockController.Connection.Close()
	}
	defer websockController.UnregisterServerEventListeners()

	for {

		var eventMessage event.FusionEvent[map[string]interface{}]
		//Read Message from client
		err := websockController.Connection.ReadJSON(&eventMessage)
		if err != nil {
			zap.S().Errorw("websocket: read JSON error", "error", err.Error())
			if ws.IsUnexpectedCloseError(err) {
				zap.S().Errorw("error handling websocket message for server", "error", err)
				break
			}
			continue
		}

		go func(eventMessage event.FusionEvent[map[string]interface{}]) {
			if err := websockController.HandleEvent(eventMessage); err != nil {
				zap.S().Errorw("websocket: error handling event", "error", err)
			}
		}(eventMessage)

	}
}
