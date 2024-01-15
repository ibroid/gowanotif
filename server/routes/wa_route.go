package routes

import (
	"errors"
	"gowhatsapp/log"
	"gowhatsapp/server/utils"
	wa "gowhatsapp/whatsapp"

	"github.com/gofiber/fiber/v2"
)

func RegisterWaRoute(app *fiber.App) {

	waRouter := app.Group("/wa")

	waRouter.Post("/start", func(c *fiber.Ctx) error {
		postReq := &wa.WAStartRequest{}

		if err := c.BodyParser(postReq); err != nil {
			return errors.New("Gagal parsing :" + err.Error())
		}

		errs := utils.Validate(*postReq)
		if len(errs) != 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Validation Error",
				"field":   errs,
			})
		}

		wa.StartAuthedWClient(*postReq)

		return c.JSON(&fiber.Map{
			"status":  true,
			"message": "Wa client berhasil dijalankan",
			"data": &fiber.Map{
				"JID":  wa.WAClients[postReq.ClientName].Store.ID.User,
				"Name": wa.WAClients[postReq.ClientName].Store.PushName,
			},
		})
	})

	waRouter.Post("/auth", func(c *fiber.Ctx) error {
		postReq := &wa.WAAuthRequest{}

		if err := c.BodyParser(postReq); err != nil {
			return errors.New("Gagal parsing :" + err.Error())
		}

		errs := utils.Validate(*postReq)
		if len(errs) != 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Validation Error",
				"field":   errs,
			})
		}

		if postReq.Handler {
			wa.SetEvenHandler(postReq.ClientName)
		}

		wa.ClientAuthentication(*postReq, nil)

		return c.JSON(&fiber.Map{
			"status":  true,
			"message": "Wa client berhasil dijalankan",
			"data": &fiber.Map{
				"JID":  wa.WAClients[postReq.ClientName].Store.ID.User,
				"Name": wa.WAClients[postReq.ClientName].Store.PushName,
			},
		})
	})

	waRouter.Post("/send_message", func(c *fiber.Ctx) error {

		postRequest := &utils.MessagePostBody{}

		if err := c.BodyParser(postRequest); err != nil {
			log.Log.Panicln("Gagal body parsing : ", err)
		}

		errors := utils.Validate(postRequest)
		if len(errors) != 0 {
			log.Log.Panicln("Gagal validasi : ", errors)
		}

		if wa.CheckClient(postRequest.ClientName) {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Wa Client tidak berjalan",
				"status":  false,
			})

		} else {
			dest, err := wa.ParseJid(postRequest.Number)
			if err != nil {
				log.Log.Panicln("Gagal parsing JID : ", errors)
			}

			err = wa.SendMessage(postRequest.ClientName, dest, postRequest.Message)
			if err != nil {
				log.Log.Panicln("Gagal send message : ", errors)
			}

		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Request berhasil. Pesan terkirim",
		})
	})

	waRouter.Post("/stop", func(c *fiber.Ctx) error {
		postRequest := &utils.WaClient{}

		if err := c.BodyParser(postRequest); err != nil {
			log.Log.Panicln("Gagal body parsing : ", err)
		}

		errors := utils.Validate(postRequest)
		if errors != nil {
			log.Log.Panicln(errors)
		}

		if wa.CheckClient(postRequest.ClientName) {
			wa.ClientDisconnect(postRequest.ClientName)

		} else {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Wa Client tidak berjalan",
				"status":  false,
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  "success",
			"message": "Whatsapp Berhasil Dihentikan",
		})
	})
}
