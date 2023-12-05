package log

import (
	"io"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/siddontang/go-log/log"
)

func HttpLogInit() func(*fiber.Ctx) error {

	file, err := os.OpenFile("http.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorln("Gagal open http log:", err)
	}
	// defer file.Close()

	mw := io.MultiWriter(logger.ConfigDefault.Output, file)

	return logger.New(logger.Config{
		Output:     mw,
		Format:     "${time}:${method}-${path}​ ${status}:${resBody}\n",
		TimeFormat: time.RFC3339,
	})
}
