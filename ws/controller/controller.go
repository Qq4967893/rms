package controller

import (
	"backend-rms/core/logger"
	"backend-rms/ws"
	"backend-rms/ws/response"
	"time"

	"nhooyr.io/websocket"
)

type TickerHandlerFunc func(c *ws.Context) (data any, err error)

func TickerHandle(c *ws.Context, d time.Duration, h TickerHandlerFunc) {
	ticker := time.NewTicker(d)
	stopCh := make(chan struct{})
	defer func() {
		ticker.Stop()
		close(stopCh)
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				data, err := h(c)
				if err != nil {
					response.Error(c, err)
					return
				}
				if data != nil {
					response.Write(c, response.WithData(data))
				}
			case <-stopCh:
				return
			}
		}
	}()

	ch := make(chan []byte, 16)

	for {
		select {
		default:
			typ, p, err := c.Read()
			if err != nil || typ != websocket.MessageText {
				return
			}
			if string(p) == "PING" {
				err = c.Conn.Write(c, typ, []byte("PONG"))
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					return
				}
				if err != nil {
					logger.Error(err.Error())
					return
				}
				continue
			}
			ch <- p
		case p := <-ch:
			ticker.Reset(d)
			c.Set("params", p)
			data, err := h(c)
			if err != nil {
				response.Error(c, err)
				return
			}
			if data != nil {
				response.Write(c,
					response.WithCode(response.StatusOK),
					response.WithData(data),
				)
			}
		}
	}
}
