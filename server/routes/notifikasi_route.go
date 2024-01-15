package routes

import (
	db "gowhatsapp/database"
	data "gowhatsapp/database/data"

	"github.com/gofiber/fiber/v2"
)

type NotifikasiPostReq struct {
	Id    int16  `json:"id" validate:"required"`
	Pesan string `json:"pesan" validate:"required"`
}

func RegisterNotifikasiRoute(app *fiber.App) {
	notifikasiRoute := app.Group("/notifikasi")

	notifikasiRoute.Get("/all", func(c *fiber.Ctx) error {
		dblocal := db.StartDBLocal()

		queryNotifikasi, err := dblocal.Query("SELECT * FROM notifikasi")
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Gagal mengambil data",
				"data":    nil,
			})
		}

		var notifikasis []data.Notifikasi
		for queryNotifikasi.Next() {
			notifikasi := new(data.Notifikasi)
			err := queryNotifikasi.Scan(&notifikasi.Id, &notifikasi.JenisNotifikasiId, &notifikasi.TujuanId, &notifikasi.Pesan, &notifikasi.Filename)
			if err != nil {
				return c.Status(400).JSON(&fiber.Map{
					"message": "Gagal mengambil data" + err.Error(),
					"data":    nil,
				})
			}

			notifikasis = append(notifikasis, *notifikasi)
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil mengambil data",
			"data":    notifikasis,
		})
	})

	notifikasiRoute.Patch("/set_pesan", func(c *fiber.Ctx) error {
		var body NotifikasiPostReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if body.Id <= 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Id tidak valid",
			})
		}

		db := db.StartDBLocal()

		result, err := db.Exec("UPDATE notifikasi SET pesan = ? WHERE id = ?", body.Pesan, body.Id)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"message": "Gagal memperbarui data",
			})
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Gagal memperbarui data",
			})
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil memperbarui data",
		})
	})
}
