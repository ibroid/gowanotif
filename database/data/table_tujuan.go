package data

import (
	"context"
	"database/sql"
	"log"
)

type Tujuan struct {
	Id   int16  `json:"id"`
	Nama string `json:"nama"`
}

func InitDataTujuan(db *sql.DB) {

	defer log.Println("Insert table tujuan success")

	data := []Tujuan{
		{
			Id:   1,
			Nama: "Hakim",
		},
		{
			Id:   2,
			Nama: "Panitera",
		},
		{
			Id:   3,
			Nama: "Jurusita",
		},
		{
			Id:   4,
			Nama: "Admin",
		},
		{
			Id:   5,
			Nama: "Pihak",
		},
	}

	for _, q := range data {
		_, err := db.ExecContext(context.TODO(), "INSERT INTO tujuan (id, nama) VALUES (?, ?)", q.Id, q.Nama)
		if err != nil {
			log.Fatal("gagal exec insert tabel tujuan")
		}
	}
}
