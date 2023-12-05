package database

import (
	"database/sql"
	"gowhatsapp/database/data"
	"gowhatsapp/log"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

var DBS *sql.DB

func InitDBLocal() {

	if err := os.Remove("localstore.db"); err != nil {
		log.Log.Errorln("Gagal hapus db awal : ", err)
	}

	var conErr error
	DBS, conErr = sql.Open("sqlite", "localstore.db?_pragma=foreign_keys(1)")
	if conErr != nil {
		log.Log.Fatalln("Gagal connect sqlite local : ", conErr)
	}

	defer log.Log.Println("Inisialisasi sqlite berhasil")
	defer DBS.Close()

	data.CreateAllTable(DBS)
	data.InitDataAdmin(DBS)
	data.InitDataNotifikasi(DBS)
	data.InitDataJenisNotifikasi(DBS)
	data.InitDataPengaturan(DBS)
	data.InitDataTujuan(DBS)
	data.InitDataTemplateReply(DBS)
}

func StartDBLocal() *sql.DB {

	dbs, conErr := sql.Open("sqlite", "localstore.db?_pragma=foreign_keys(1)")
	if conErr != nil {
		log.Log.Errorln("Gagal open db local : ", conErr)
	}

	return dbs
}
