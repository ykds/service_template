package pusher

import (
	"fmt"
	"net"
	"service_template/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

var _ Pusher = (*SSE)(nil)

type SSE struct {
	logger *logger.Logger
	opt    Option
}

func NewSSE(opt Option) *SSE {
	return &SSE{
		logger: opt.Logger,
		opt:    opt,
	}
}

func (sse *SSE) Register(engine *gin.Engine, middleware []gin.HandlerFunc) {
	middleware = append(middleware, sse.handleSSE)
	engine.GET("/sse", middleware...)
}

func (sse *SSE) handleSSE(ctx *gin.Context) {
	conn, _, err := ctx.Writer.Hijack()
	if err != nil {
		return
	}
	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/event-stream\r\n"))
	if err != nil {
		conn.Close()
		return
	}
	_, err = conn.Write([]byte(fmt.Sprintf("Access-Control-Allow-Origin: %s\r\n", ctx.Writer.Header().Get("Access-Control-Allow-Origin"))))
	if err != nil {
		conn.Close()
		return
	}
	_, err = conn.Write([]byte(fmt.Sprintf("Access-Control-Allow-Credentials: %s\r\n", ctx.Writer.Header().Get("Access-Control-Allow-Credential"))))
	if err != nil {
		conn.Close()
		return
	}
	_, err = conn.Write([]byte(fmt.Sprintf("Access-Control-Allow-Methods: %s\r\n", ctx.Writer.Header().Get("Access-Control-Allow-Methods"))))
	if err != nil {
		conn.Close()
		return
	}
	_, err = conn.Write([]byte(fmt.Sprintf("Access-Control-Allow-Headers: %s\r\n", ctx.Writer.Header().Get("Access-Control-Allow-Headers"))))
	if err != nil {
		conn.Close()
		return
	}
}

func (sse *SSE) hearbeat(userId int, conn net.Conn) {
	time.AfterFunc(10*time.Second, func() {
		_, err := conn.Write([]byte("event: heartbeat\ndata: \n\n"))
		if err != nil {
			logger.Errorf("write heartbeat to %d failed: %v, close connection", userId, err)
			conn.Close()
			return
		}
		sse.hearbeat(userId, conn)
	})
}

func (sse *SSE) PushAll(msg Message) error {
	return nil
}

func (sse *SSE) Push(users []int, msg Message) error {
	return nil
}

func (sse *SSE) Close() error {
	return nil
}
