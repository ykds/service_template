package pusher

import (
	"service_template/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	WebsocketPusher = "websocket"
	SSEPusher       = "sse"
	SocketIOPusher  = "socketio"
)

type Message struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	CreatedTime int64  `json:"created_time"`
}

type Option struct {
	Type        string `json:"type" yaml:"type"`
	AllowOrigin string `json:"allow_origin" yaml:"allow_origin"`
	Logger      *logger.Logger
}

type Pusher interface {
	Register(engine *gin.Engine, middleware []gin.HandlerFunc)
	Push([]int, Message) error
	PushAll(Message) error
	Close() error
}

func NewPusher(opt Option) Pusher {
	switch opt.Type {
	case WebsocketPusher:
		return NewWebsocket(opt)
	case SSEPusher:
		return NewSSE(opt)
	case SocketIOPusher:
		return NewSocketIO(opt)
	default:
		panic("unsupported pusher type")
	}
}
