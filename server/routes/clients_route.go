package routes

import (
	db "gowhatsapp/database"
	"gowhatsapp/database/data"
	"gowhatsapp/server/utils"

	"github.com/gofiber/fiber/v2"
)

type PatchClientsServiceReq struct {
	ClientName string `json:"client_name" validate:"required"`
	Service    string `json:"service" validate:"required"`
}

func RegisterClientsRoute(app *fiber.App) {
	clientsRoute := app.Group("/clients")

	clientsRoute.Get("/", func(c *fiber.Ctx) error {
		dblocal := db.StartDBLocal()
		defer dblocal.Close()

		queryGetAllClients, err := dblocal.Query("SELECT * FROM clients")
		if err != nil {
			panic("Gagal get clients : " + err.Error())
		}
		defer queryGetAllClients.Close()

		var clients []data.Clients
		for queryGetAllClients.Next() {
			var client data.Clients
			err = queryGetAllClients.Scan(&client.Id, &client.ClientName, &client.Jid, &client.Handler, &client.Status, &client.Service)
			if err != nil {
				panic("Gagal get clients : " + err.Error())
			}

			clients = append(clients, client)
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil mengambil data clients",
			"data":    clients,
		})
	})

	clientsRoute.Patch("/", func(c *fiber.Ctx) error {
		var body data.Clients
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		errs := utils.Validate(body)
		if len(errs) != 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Validation Error",
				"field":   errs,
			})
		}

		dblocal := db.StartDBLocal()
		defer dblocal.Close()

		queryUpdateClient, err := dblocal.Exec("UPDATE clients SET client_name =?, jid =?, handler =?, status =?, service =? WHERE id =?", body.ClientName, body.Jid, body.Handler, body.Status, body.Service, body.Id)
		if err != nil {
			panic("Gagal update client : " + err.Error())
		}

		updatedClientId, _ := queryUpdateClient.LastInsertId()

		rowsAffected, err := queryUpdateClient.RowsAffected()
		if err != nil {
			panic("Gagal update client : " + err.Error())
		}

		if rowsAffected == 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Client tidak ditemukan",
			})
		}

		updatedClient := dblocal.QueryRow("SELECT * FROM clients WHERE id =?", updatedClientId)

		updatedClientData := new(data.Clients)
		err = updatedClient.Scan(&updatedClientData.Id, &updatedClientData.ClientName, &updatedClientData.Jid, &updatedClientData.Handler, &updatedClientData.Status, &updatedClientData.Service)

		if err != nil {
			panic("Gagal get last updated client : " + err.Error())
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil mengubah service client",
			"data":    updatedClientData,
		})

	})

	clientsRoute.Post("/", func(c *fiber.Ctx) error {
		var body data.Clients
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		errs := utils.Validate(body)
		if len(errs) != 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Validation Error",
				"field":   errs,
			})
		}

		dblocal := db.StartDBLocal()
		defer dblocal.Close()

		queryInsertClient, err := dblocal.Exec("INSERT INTO clients (client_name, jid, handler, status, service) VALUES (?,?,?,?,?)", body.ClientName, body.Jid, body.Handler, body.Status, body.Service)
		if err != nil {
			panic("Gagal insert client : " + err.Error())
		}

		insertedId, _ := queryInsertClient.LastInsertId()

		rowsAffected, err := queryInsertClient.RowsAffected()
		if err != nil {
			panic("Gagal insert client : " + err.Error())
		}

		if rowsAffected == 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Gagal insert client",
			})
		}

		lastInsertedClient := dblocal.QueryRow("SELECT * FROM clients WHERE id =?", insertedId)

		var lastInsertedClientData data.Clients

		err = lastInsertedClient.Scan(&lastInsertedClientData.Id, &lastInsertedClientData.ClientName, &lastInsertedClientData.Jid, &lastInsertedClientData.Handler, &lastInsertedClientData.Status, &lastInsertedClientData.Service)

		if err != nil {
			panic("Gagal get last inserted client : " + err.Error())
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil insert client",
			"data":    lastInsertedClientData,
		})
	})

	clientsRoute.Delete("/:id", func(c *fiber.Ctx) error {
		clientId := c.Params("id")
		if clientId == "" {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Id tidak valid",
			})
		}
		dblocal := db.StartDBLocal()
		defer dblocal.Close()

		queryDeleteClient, err := dblocal.Exec("DELETE FROM clients WHERE id =?", clientId)

		if err != nil {
			panic("Gagal delete client : " + err.Error())
		}
		rowsAffected, err := queryDeleteClient.RowsAffected()
		if err != nil {
			panic("Gagal delete client : " + err.Error())
		}
		if rowsAffected == 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Client tidak ditemukan",
			})
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil menghapus client",
		})
	})
}
