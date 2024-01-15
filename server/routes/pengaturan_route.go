package routes

import (
	"errors"
	db "gowhatsapp/database"
	data "gowhatsapp/database/data"
	"gowhatsapp/server/utils"

	"github.com/gofiber/fiber/v2"
)

type PengaturanPostReq struct {
	Value string `json:"value" validate:"required"`
	Key   string `json:"key" validate:"required"`
}

func RegisterPengaturanRoute(app *fiber.App) {
	pengaturanRoute := app.Group("/pengaturan")

	pengaturanRoute.Get("/all", func(c *fiber.Ctx) error {

		sql := db.StartDBLocal()
		queryPengaturan, err := sql.Query("SELECT * FROM pengaturan")
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Gagal mengambil data",
				"data":    nil,
			})
		}

		defer queryPengaturan.Close()

		var pengaturans []data.Pengaturan
		for queryPengaturan.Next() {
			pengaturan := new(data.Pengaturan)
			err := queryPengaturan.Scan(&pengaturan.Id, &pengaturan.Key, &pengaturan.Value, &pengaturan.Ket)
			if err != nil {
				return c.Status(400).JSON(&fiber.Map{
					"message": "Gagal mengambil data",
					"data":    nil,
				})
			}

			pengaturans = append(pengaturans, *pengaturan)
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil mengambil data",
			"data":    pengaturans,
		})
	})

	pengaturanRoute.Patch("/set", func(c *fiber.Ctx) error {
		postReq := &PengaturanPostReq{}

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

		sql := db.StartDBLocal()
		result, erre := sql.Exec("UPDATE pengaturan SET value = ? WHERE key = ?", postReq.Value, postReq.Key)

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 || erre != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Gagal memperbarui data",
			})
		}

		return c.JSON(&fiber.Map{
			"message": "Berhasil mengupdate data",
			"data":    nil,
		})
	})
}
