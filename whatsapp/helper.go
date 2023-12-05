package whatsapp

import (
	"gowhatsapp/database"
	"gowhatsapp/log"
	"strings"
	"time"
)

func ReadableDate(date string, erra error) string {
	if erra != nil {
		log.Log.Errorln("Gagal convert readable date dari atas: ", erra)
		return "null"
	}

	dateParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Log.Errorln("Gagal convert readable date : ", err)
		return "null"
	}

	dateReadable := dateParse.Format("2 January 2006")
	return dateReadable
}

func PerkaraIDFromJID(number string) (string, error) {
	number = strings.Replace(number, "62", "0", 1)

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	query, err := dbsipp.Execute("SELECT perkara_id FROM perkara_pihak1 AS c JOIN pihak AS d ON d.id = c.pihak_id WHERE d.telepon = ?", number)

	if err != nil {
		log.Log.Errorln("Gagal query get perkara_id from jid : ", err)
		return "", err
	}

	perkaraId, _ := query.GetStringByName(0, "perkara_id")

	return perkaraId, nil
}

func PerkaraIDFromPar3(nomorPerkara string) (string, error) {

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	query, err := dbsipp.Execute("SELECT perkara_id FROM perkara WHERE nomor_perkara = ?", nomorPerkara)

	if err != nil {
		log.Log.Errorln("Gagal query get perkara_id from jid : ", err)
		return "", err
	}

	perkaraId, _ := query.GetStringByName(0, "perkara_id")

	return perkaraId, nil
}
