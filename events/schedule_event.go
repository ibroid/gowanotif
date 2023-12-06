package events

import (
	"fmt"
	"gowhatsapp/database"
	"gowhatsapp/log"
	"os"

	"strings"
	"time"

	cron "github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

func StartCronEvent() {
	// jakartaTime := getLocation()
	scheduler = cron.New(cron.WithLocation(time.UTC))

	scheduleTime := os.Getenv("CRON_SCHEDULE_TIME")

	if scheduleTime == "" {
		log.Log.Warningln("Cron schedule time belum di tentukan. Menggunakan Default 0 8 * * *")
		scheduleTime = "0 8 * * *"
	}
	scheduler.AddFunc(scheduleTime, NotifPRSH)
	scheduler.AddFunc(scheduleTime, NotifPRSJ)
	scheduler.AddFunc(scheduleTime, NotifPRSP)
	go scheduler.Start()
}

func StopCronEvent() {
	scheduler.Stop()
}

func NotifPRSH() {
	dblocal := database.StartDBLocal()

	defer dblocal.Close()

	var (
		PesanNotifikasi string
		NamaNotifikasi  string
	)

	queryNotifikasi := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a
	JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id
	WHERE a.key = 'prs'
	AND b.tujuan_id = 1;`)

	if err := queryNotifikasi.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal parse prs : ", err)
		return
	}

	dbsipp := database.StartDBSipp()

	defer dbsipp.Close()

	activeHakim, err := dbsipp.Execute("SELECT * FROM hakim_pn WHERE aktif = 'Y'")
	if err != nil {
		log.Log.Errorln("Gagal select aktif hakim : ", err)
		return
	}

	for i := range activeHakim.Values {

		hakimId, _ := activeHakim.GetIntByName(i, "id")
		hakimName, _ := activeHakim.GetStringByName(i, "nama")

		queryRelaas, err := dbsipp.Execute(`SELECT nomor_perkara,d.hakim_nama FROM perkara AS a
		LEFT JOIN perkara_hakim_pn AS d USING(perkara_id) 
		WHERE a.perkara_id IN (
			SELECT b.perkara_id FROM perkara_jadwal_sidang AS b
			LEFT JOIN perkara_pelaksanaan_relaas AS c ON b.id = c.sidang_id
			WHERE b.tanggal_sidang = CURDATE()
			AND (b.agenda NOT LIKE '%mediasi%')
			AND b.urutan <= 2
			AND c.id IS NULL )
		AND d.jabatan_hakim_id = 1
		AND d.aktif = 'Y'
		AND d.hakim_id = ?`, hakimId)

		if err != nil {
			log.Log.Errorln("Gagal query no relaas di hakim ", hakimName, err)
			return
		}

		var daftarPerkara string

		for y := range queryRelaas.Values {
			nomorPerkara, _ := queryRelaas.GetStringByName(y, "nomor_perkara")

			daftarPerkara += fmt.Sprintf("%s \n", nomorPerkara)

		}

		nomorHakim, _ := activeHakim.GetStringByName(i, "keterangan")
		if len(daftarPerkara) > 0 {

			PesanNotifikasi := strings.Replace(PesanNotifikasi, "{daftar_relaas}", daftarPerkara, -1)
			log.Log.Info("Schedule Notifikasi PRS terkirim ke hakim : ", hakimName)

			// destHakim, _ := whatsapp.ParseJid(nomorHakim)
			// whatsapp.SendMessage("internal", destHakim, PesanNotifikasi)
			SendNotifMessage("internal", nomorHakim, PesanNotifikasi)
		}

		queryRelaas.Close()
	}

	defer activeHakim.Close()
}

func NotifPRSJ() {

	dblocal := database.StartDBLocal()

	defer dblocal.Close()

	var (
		PesanNotifikasi string
		NamaNotifikasi  string
	)

	queryNotifikasi := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a
	JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id
	WHERE a.key = 'prs'
	AND b.tujuan_id = 1;`)

	if err := queryNotifikasi.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal parse prs : ", err)
		return
	}

	dbsipp := database.StartDBSipp()

	defer dbsipp.Close()

	activeJs, err := dbsipp.Execute("SELECT * FROM jurusita WHERE aktif = 'Y'")

	if err != nil {
		log.Log.Errorln("Gagal select aktif hakim : ", err)
		return
	}

	defer activeJs.Close()

	for i := range activeJs.Values {

		jspId, _ := activeJs.GetStringByName(i, "id")
		jspName, _ := activeJs.GetStringByName(i, "nama")

		queryRelaas, err := dbsipp.Execute(`SELECT nomor_perkara,d.jurusita_nama FROM perkara AS a
		LEFT JOIN perkara_jurusita AS d USING(perkara_id) 
		WHERE a.perkara_id IN (
			SELECT b.perkara_id FROM perkara_jadwal_sidang AS b
			LEFT JOIN perkara_pelaksanaan_relaas AS c ON b.id = c.sidang_id
			WHERE b.tanggal_sidang = CURDATE()
			AND (b.agenda NOT LIKE '%mediasi%')
			AND b.urutan <= 2
			AND c.id IS NULL)
		AND d.aktif = 'Y'
		AND d.jurusita_id = ?`, jspId)

		if err != nil {
			log.Log.Errorln("Gagal Query Relaas JSP : ", err)
			return
		}

		var daftarPerkara string

		for y := range queryRelaas.Values {
			nomorPerkara, _ := queryRelaas.GetStringByName(y, "nomor_perkara")

			daftarPerkara += fmt.Sprintf("%s \n", nomorPerkara)

		}

		nomorJsp, _ := activeJs.GetStringByName(i, "keterangan")

		if len(daftarPerkara) > 0 {

			PesanNotifikasi := strings.Replace(PesanNotifikasi, "{daftar_relaas}", daftarPerkara, -1)
			log.Log.Infoln("Schedule Notifikasi PRS terkirim ke js : ", jspName)

			// destJsp, _ := whatsapp.ParseJid(nomorJsp)
			// whatsapp.SendMessage("internal", destJsp, PesanNotifikasi)
			SendNotifMessage("internal", nomorJsp, PesanNotifikasi)
		}

	}
}

func NotifPRSP() {
	dblocal := database.StartDBLocal()

	defer dblocal.Close()

	var (
		PesanNotifikasi string
		NamaNotifikasi  string
	)

	queryNotifikasi := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a
	JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id
	WHERE a.key = 'prs'
	AND b.tujuan_id = 1;`)

	if err := queryNotifikasi.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal parse prs : ", err)
		return
	}

	dbsipp := database.StartDBSipp()

	defer dbsipp.Close()

	queryRelaas, err := dbsipp.Execute(`SELECT nomor_perkara,d.jurusita_nama FROM perkara AS a
	LEFT JOIN perkara_jurusita AS d USING(perkara_id) 
	WHERE a.perkara_id IN (
		SELECT b.perkara_id FROM perkara_jadwal_sidang AS b
		LEFT JOIN perkara_pelaksanaan_relaas AS c ON b.id = c.sidang_id
		WHERE b.tanggal_sidang = CURDATE()
		AND (b.agenda NOT LIKE '%mediasi%')
		AND b.urutan <= 2
		AND c.id IS NULL)
	AND d.aktif = 'Y'`)

	if err != nil {
		log.Log.Errorln("Gagal query PRS to Panitera : ", err)
		return
	}

	var nomorPanitera string

	queryPengaturan := dblocal.QueryRow("SELECT value FROM pengaturan WHERE pengaturan.key = 'nomor_panitera'")

	errs := queryPengaturan.Scan(&nomorPanitera)
	if errs != nil {
		log.Log.Errorln("Gagal mengambil value dengan key nomor_panitera : ", errs)
	}

	var daftarPerkara string
	for y := range queryRelaas.Values {
		nomorPerkara, _ := queryRelaas.GetStringByName(y, "nomor_perkara")

		daftarPerkara += fmt.Sprintf("%s \n", nomorPerkara)

	}
	if nomorPanitera == "" {
		log.Log.Warn("Skip PRS ke Panitera. Nomor telepon belum ada")
		return
	}
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{daftar_relaas}", daftarPerkara, -1)
	log.Log.Info("Schedule Notifikasi PRS terkirim ke panitera")
	SendNotifMessage("internal", nomorPanitera, PesanNotifikasi)

}
