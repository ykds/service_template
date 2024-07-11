package pusher

import (
	"service_template/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var _ Pusher = (*Websocket)(nil)

type Websocket struct {
	logger *logger.Logger
	opt    Option
}

func NewWebsocket(opt Option) *Websocket {
	ws := &Websocket{
		logger: opt.Logger,
		opt:    opt,
	}
	return ws
}

func (ws *Websocket) handleWebsocket(ctx *gin.Context) {
	_, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
}

func (ws *Websocket) Register(engine *gin.Engine, middleware []gin.HandlerFunc) {
	middleware = append(middleware, ws.handleWebsocket)
	engine.GET("/ws", middleware...)
}

func (ws *Websocket) PushAll(msg Message) error {
	return nil
}

func (ws *Websocket) Push(users []int, msg Message) error {
	return nil
}

func (ws *Websocket) Close() error {
	return nil
}
