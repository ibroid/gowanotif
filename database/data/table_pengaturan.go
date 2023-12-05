package data

import (
	"context"
	"database/sql"
	"log"
	"sync"
)

type Pengaturan struct {
	Id    int16          `json:"id"`
	Key   string         `json:"key"`
	Value string         `json:"value"`
	Ket   sql.NullString `json:"ket"`
}

func InitDataPengaturan(db *sql.DB) {

	defer log.Println("Insert table pengaturan success")

	data := []Pengaturan{
		{
			Id:    1,
			Key:   "nama_satker",
			Value: "Pengadilan Agama Jakarta Utara",
			Ket: sql.NullString{
				String: "Nama Satuan Kerja",
				Valid:  true,
			},
		},
		{
			Id:    2,
			Key:   "nama_panitera",
			Value: "",
			Ket: sql.NullString{
				String: "Nama Panitera",
				Valid:  true,
			},
		},
		{
			Id:    3,
			Key:   "nomor_panitera",
			Value: "",
			Ket: sql.NullString{
				String: "Nomor Telepon Panitera",
				Valid:  true,
			},
		},
	}

	var wg sync.WaitGroup

	for _, q := range data {
		wg.Add(1)

		go func(q Pengaturan) {
			_, err := db.ExecContext(context.Background(), "INSERT INTO PENGATURAN (id, key, value, key) VALUES (?, ?, ?, ?)", q.Id, q.Key, q.Value, q.Ket)

			if err != nil {
				log.Fatal("gagal exec insert table pengaturan")
			}
		}(q)
	}
}
