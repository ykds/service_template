package pusher

import (
	"service_template/pkg/logger"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var _ Pusher = (*SoektIO)(nil)

type SoektIO struct {
	*socketio.Server
	logger *logger.Logger
	opt    Option
}

type context struct {
	connId string
	userId int
}

func NewSocketIO(opt Option) *SoektIO {
	s := &SoektIO{
		logger: opt.Logger,
		opt:    opt,
	}
	s.Server = socketio.NewServer(nil)
	s.Server.OnConnect("/", func(c socketio.Conn) error {
		return nil
	})
	s.Server.OnError("/", func(c socketio.Conn, e error) {
	})
	s.Server.OnDisconnect("/", func(c socketio.Conn, reason string) {
	})
	go s.Server.Serve()
	return s
}

func (s *SoektIO) PushAll(msg Message) error {
	return nil
}

func (s *SoektIO) Push(users []int, msg Message) error {
	return nil
}

func (s *SoektIO) Register(engine *gin.Engine, middleware []gin.HandlerFunc) {
	engine.GET("/socket.io/*any", append(middleware, gin.WrapH(s))...)
	engine.POST("/socket.io/*any", append(middleware, gin.WrapH(s))...)
}

func (s *SoektIO) Close() error {
	return s.Server.Close()
}
