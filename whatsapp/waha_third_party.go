package whatsapp

import (
	"gowhatsapp/log"

	"github.com/gofiber/fiber/v2"
)

func SendMessageViaWAHA(session string, number string, message string) {
	ua := fiber.Post("http://192.168.0.202:3030/api/sendText")
	ua.Set("Content-Type", "application/json")
	ua.JSON(fiber.Map{
		"chatId":  number + "@c.us",
		"text":    message,
		"session": session,
	})

	statusCode, body, errs := ua.Bytes()
	if len(errs) > 0 {
		log.Log.Error("Gagal send message via WA HA:", errs)
	}

	if statusCode != 200 {
		log.Log.Error("Gagal send message via WA HA, status code:", statusCode)
	}

	log.Log.Info("Message sent via WA HA . Response :", string(body))
}
