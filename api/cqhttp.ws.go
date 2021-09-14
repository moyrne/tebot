package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CQWSServer interface {
	Event(ctx context.Context, message QMessage) (Reply, error)
	Ping()
}

const (
	readBufSize  = 1024 * 1024
	writeBufSize = 1024 * 1024
)

func RegisterCQWSServer(e *gin.Engine, server CQWSServer) {
	// TODO websocket: bad handshake
	e.POST("/", func(c *gin.Context) {
		u := websocket.Upgrader{ReadBufferSize: readBufSize, WriteBufferSize: writeBufSize}
		u.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			// don't return errors to maintain backwards compatibility
		}

		conn, err := u.Upgrade(c.Writer, c.Request, c.Request.Header)
		if err != nil {
			logrus.Error(errors.WithStack(err))
			return
		}
		defer conn.Close()

		conn.SetPingHandler(func(message string) error {
			fmt.Println("ping")
			server.Ping()
			err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second*10))
			if err == websocket.ErrCloseSent {
				return nil
			} else if e, ok := err.(net.Error); ok && e.Temporary() {
				return nil
			}
			return err
		})

		conn.SetPongHandler(func(appData string) error {
			fmt.Println("pong")
			server.Ping()
			return nil
		})

		go func() {
			for {
				var message QMessage
				if err := conn.ReadJSON(&message); err != nil {
					logrus.Error(errors.WithStack(err))
					return
				}
				fmt.Println("message", message)
				reply, err := server.Event(c.Request.Context(), message)
				if err != nil {
					_ = conn.WriteJSON(gin.H{"error": err.Error()})
					return
				}
				if err := conn.WriteJSON(reply); err != nil {
					logrus.Error(errors.WithStack(err))
				}
			}
		}()
	})
}
