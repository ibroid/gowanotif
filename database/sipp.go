package database

import (
	"gowhatsapp/log"
	"os"

	"github.com/go-mysql-org/go-mysql/client"
)

type HakimPn struct {
}

func StartDBSipp() *client.Conn {

	conn, err := client.Connect(os.Getenv("DB_SIPP_HOST"), os.Getenv("DB_SIPP_USER"), os.Getenv("DB_SIPP_PASS"), os.Getenv("DB_SIPP_NAME"))
	if err != nil {
		log.Log.Errorln("Gagal connect db sipp : ", err)
	}

	if errp := conn.Ping(); errp != nil {
		log.Log.Errorln("Gagal ping db sipp : ", err)
	}

	return conn
}
