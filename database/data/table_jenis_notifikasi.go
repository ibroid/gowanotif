package data

import (
	"context"
	"database/sql"
	"log"
)

type JenisNotifikasi struct {
	Id             int16  `json:"id"`
	NamaNotifikasi string `json:"nama_notifikasi"`
	Key            string `json:"key"`
}

func InitDataJenisNotifikasi(db *sql.DB) {
	defer log.Println("Insert table jenis notifikasi success")

	data := []JenisNotifikasi{
		{
			Id:             1,
			NamaNotifikasi: "Penetapan Majelis Hakim",
			Key:            "pmh",
		},
		{
			Id:             2,
			NamaNotifikasi: "Penetapan Panitera Pengganti",
			Key:            "ppp",
		},
		{
			Id:             3,
			NamaNotifikasi: "Penetapan Jurusita",
			Key:            "pjs",
		},
		{
			Id:             4,
			NamaNotifikasi: "Penetapan Sidang Pertama",
			Key:            "phs",
		},
		{
			Id:             5,
			NamaNotifikasi: "Tundaan Sidang",
			Key:            "pts",
		},
		{
			Id:             6,
			NamaNotifikasi: "Sidang Putus Verstek",
			Key:            "spv",
		},
		{
			Id:             7,
			NamaNotifikasi: "Peringatan Relaas Sebelum Sidang",
			Key:            "prs",
		},
		{
			Id:             8,
			NamaNotifikasi: "Penerbitan Akta Cerai",
			Key:            "pac",
		},
		{
			Id:             9,
			NamaNotifikasi: "Pemberitahuan Upload Relaas",
			Key:            "pur",
		},
		{
			Id:             10,
			NamaNotifikasi: "Pemberitahuan Pendaftaran Perkara",
			Key:            "ppd",
		},
	}

	for _, q := range data {
		_, err := db.ExecContext(context.Background(), "INSERT INTO jenis_notifikasi (id, nama_notifikasi, key) values (?, ?, ?)", q.Id, q.NamaNotifikasi, q.Key)
		if err != nil {
			log.Fatal("gagal exec insert table jenis notifikasi : ", err)
		}
	}
}
