package server

import (
	"gowhatsapp/log"
	"gowhatsapp/server/routes"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var Http *fiber.App

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func StartServer() {
	Http = fiber.New(fiber.Config{
		AppName:      "GoWhatsapp",
		ServerHeader: "Fiber Go",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	Http.Use(recover.New())
	Http.Use(log.HttpLogInit())

	Http.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"app_name": "Go Whatsapp",
			"version":  "0.1",
		})
	})

	Http.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	Http.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.LogInitWithWs(c)
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Log.Debugln("read:", err)
				break
			}
			log.Log.Debugf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Log.Debugln("write:", err)
				break
			}
		}
	}))

	routes.RegisterWaRoute(Http)

	Http.Listen("0.0.0.0:8099")
}
