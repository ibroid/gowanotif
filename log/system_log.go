package log

import (
	"io"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type WebSocketWriter struct {
	conn *websocket.Conn
}

func NewWebSocketWriter(conn *websocket.Conn) *WebSocketWriter {
	return &WebSocketWriter{conn: conn}
}

func (ww *WebSocketWriter) Write(p []byte) (int, error) {
	err := ww.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (ww *WebSocketWriter) Error() error {
	return nil
}

func LogInit() {

	// log.SetOutput()

	file, err := os.OpenFile("gowhatsapp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Errorln("Fail load gowhatsapp.log")
	}

	mw := io.MultiWriter(os.Stderr, file)

	Log.SetOutput(mw)

	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
	})

}

func LogInitWithWs(c *websocket.Conn) {
	ww := NewWebSocketWriter(c)

	// Log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("gowhatsapp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Errorln("Fail load gowhatsapp.log")
	}

	mw := io.MultiWriter(os.Stderr, file, ww)

	Log.SetOutput(mw)
}
