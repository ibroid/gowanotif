package data

import (
	"context"
	"database/sql"
	"log"
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
			Key:   "nama_satker",
			Value: "Pengadilan Agama Jakarta Utara",
			Ket: sql.NullString{
				String: "Nama Satuan Kerja",
				Valid:  true,
			},
		},
		{
			Key:   "alamat_satker",
			Value: "Jl. Raya Plumpang Semper No.5, Kel. Tugu Selatan, Kec. Koja, Kota Jakarta Utara",
			Ket: sql.NullString{
				String: "Alamat Satuan Kerja",
				Valid:  true,
			},
		},
		{
			Key:   "email_satker",
			Value: "pengadilanagama.jakut@gmail.com",
			Ket: sql.NullString{
				String: "Email Satuan Kerja",
				Valid:  true,
			},
		},
		{
			Key:   "telepon_satker",
			Value: " 02143934701",
			Ket: sql.NullString{
				String: "Nomor Telepon Satuan Kerja",
				Valid:  true,
			},
		},
		// DATA KETUA
		{
			Key:   "nama_ketua",
			Value: "Dr. Drs. M. Fauzi Ardi, S.H., M.H.",
			Ket: sql.NullString{
				String: "Nama Ketua Satuan Kerja",
				Valid:  true,
			},
		},
		{
			Key:   "nomor_ketua",
			Value: "",
			Ket: sql.NullString{
				String: "Nomor Telepon (WA) Ketua",
				Valid:  true,
			},
		},
		// DATA WAKIL
		{
			Key:   "nama_wakil",
			Value: "Ruslan S.Ag., S.H., M.H.",
			Ket: sql.NullString{
				String: "Nama Wakil Ketua Satuan Kerja",
				Valid:  true,
			},
		},
		{
			Key:   "nomor_wakil",
			Value: "",
			Ket: sql.NullString{
				String: "Nomor Telepon (WA) Wakil Ketua",
				Valid:  true,
			},
		},
		// DATA PANITERA
		{
			Key:   "nama_panitera",
			Value: "Aday, S.Ag., M.H.",
			Ket: sql.NullString{
				String: "Nama Panitera",
				Valid:  true,
			},
		},
		{
			Key:   "nomor_panitera",
			Value: "",
			Ket: sql.NullString{
				String: "Nomor Telepon (WA) Panitera",
				Valid:  true,
			},
		},
		// DATA SEKRETARIS
		{
			Key:   "nama_sekretaris",
			Value: "Drs. SAFE`I AGUSTIAN",
			Ket: sql.NullString{
				String: "Nama Sekretaris",
				Valid:  true,
			},
		},
		{
			Key:   "nomor_sekretaris",
			Value: "",
			Ket: sql.NullString{
				String: "Nomor Telepon (WA) Sekretaris",
				Valid:  true,
			},
		},
	}

	for _, q := range data {

		func(q Pengaturan) {
			_, err := db.ExecContext(context.Background(), "INSERT INTO PENGATURAN ( key, value, key) VALUES (?, ?, ?)", q.Key, q.Value, q.Ket)

			if err != nil {
				log.Fatal("gagal exec insert table pengaturan : ", err)
			}
		}(q)
	}
}
