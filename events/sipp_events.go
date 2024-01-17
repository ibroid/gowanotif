package events

import (
	"errors"
	"gowhatsapp/database"
	"gowhatsapp/log"
	"gowhatsapp/whatsapp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	if os.Getenv("DEV") == "1" {
		log.Log.Infof("%s ke table %v. data : %v. len : %v", e.Action, string(e.Table.Name), e.Rows[0], len(e.Rows[0]))
	}

	if os.Getenv("SINGLE_WA") == "1" {
		if os.Getenv("SERVICE_ONLY") == "internal" {
			if string(e.Table.Name) == "perkara_pihak1" || string(e.Table.Name) == "perkara_akta_cerai" {
				return nil
			}
		}

		if os.Getenv("SERVICE_ONLY") == "public" {
			if string(e.Table.Name) != "perkara_pihak1" && string(e.Table.Name) != "perkara_akta_cerai" {
				return nil
			}
		}
	}

	switch string(e.Table.Name) {
	case "perkara_hakim_pn":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_hakim_pn. ID :", e.Rows[0][0], " .Notif to Hakim")
			NotifPMH(e.Rows[0])
		}

	case "perkara_panitera_pn":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_panitera. ID :", e.Rows[0][0], " .Notif to PP")
			NotifPPP(e.Rows[0])
		}

	case "perkara_jurusita":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_jurusita. ID :", e.Rows[0][0], " .Notif to JS")
			NotifPJS(e.Rows[0])
		}

	case "perkara_pelaksanaan_relaas":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_pelaksanaan_relaas. ID :", e.Rows[0][0], " .Notif to Hakim")
			NotifPUR(e.Rows[0])
		}

	case "perkara_jadwal_sidang":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_jadwal_sidang. ID :", e.Rows[0][0], " .Notif to JS")
			NotifPHS(e.Rows[0])
		}

	case "perkara_putusan":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_putusan. ID :", e.Rows[0][0], " .Notif to JS")
			NotifSPV(e.Rows[0])
		}

	case "perkara_akta_cerai":
		if e.Action == "update" {
			log.Log.Infoln("Update event perkara_akta_cerai. ID :", e.Rows[0][0], " .Notif to Pihak")
			NotifPAC(e.Rows[1], e.Rows[0])
		}

	case "perkara_pihak1":
		if e.Action == "insert" {
			log.Log.Infoln("Insert event perkara_pihak1. ID :", e.Rows[0][0], " .Notif to Pihak")
		}

	default:
	}

	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func StartSippEvent() {

	cfg := canal.NewDefaultConfig()
	cfg.Addr = os.Getenv("DB_SIPP_HOST")
	cfg.User = os.Getenv("DB_SIPP_USER")
	cfg.Password = os.Getenv("DB_SIPP_PASS")
	cfg.ServerID = 1
	cfg.IncludeTableRegex = []string{
		"perkara_pihak1",
		"perkara_pihak2",
		"perkara_putusan",
		"perkara_akta_cerai",
		"perkara_jadwal_sidang",
		"perkara_hakim_pn",
		"perkara_panitera",
		"perkara_jurusita",
		"perkara_putusan_pemberitahuan_putusan",
		"perkara_pelaksanaan_relaas",
	}

	cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Log.Fatalln("Gagal new canal : ", err)
	}

	c.SetEventHandler(&MyEventHandler{})

	mp, _ := c.GetMasterPos()
	c.RunFrom(mp)

}

// NOTIF PENETAPAN MAJELIS HAKIM
func NotifPMH(rowData []interface{}) error {
	if rowData[2] == "18" {
		log.Log.Warningln("Skip Notif PMH. Note : PMH ini bersifat ikrar talak. ID : ", rowData[0])
		return nil
	}

	jabatanHakimId := rowData[6].(uint32)

	if jabatanHakimId != 1 {
		log.Log.Warningln("Skip Notif PMH. Note : PMH ini untuk anggota majelis. ID : ", rowData[0])
		return nil
	}

	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a
	JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id
	WHERE a.key = 'pmh'
	AND b.tujuan_id = 1;`)

	err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi)

	if err != nil {
		log.Log.Errorln("Gagal query notifikasi pmh hakim :", err, ". ID : ", rowData[0])
		return err
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataPmh, error := dbsipp.Execute("SELECT nomor_perkara, nama, c.keterangan FROM perkara AS a LEFT JOIN perkara_hakim_pn AS b USING(perkara_id) LEFT JOIN hakim_pn AS c ON b.hakim_id = c.id WHERE a.perkara_id = ? AND b.jabatan_hakim_id = 1 AND b.`aktif` = 'Y'", rowData[1])

	if error != nil {
		log.Log.Errorln("Gagal query data pmh : ", err, ". ID : ", rowData[0])
		return err
	}

	noHpHakim, err := queryDataPmh.GetStringByName(0, "keterangan")
	if err != nil {
		log.Log.Errorln("Gagal get uint by name di kolom keterangan : ", err)
		return err
	}

	nomorPerkara, err := queryDataPmh.GetStringByName(0, "nomor_perkara")
	if err != nil {
		log.Log.Errorln("Gagal get uint by name di kolom keterangan : ", err)
		return err
	}

	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)

	// whatsapp.SendMessage("internal", noHpHakim, PesanNotifikasi)
	if err := SendNotifMessage("internal", noHpHakim, PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal kirim notif PMH : ", err)
	}

	return nil

}

// NOTIF PENERAPAN PANITERA PENGGANTI
func NotifPPP(rowData []interface{}) error {

	if rowData[2] == "18" {
		log.Log.Warningln("Skip Notif Penetapan PP. Note : Penetapan ini bersifat ikrar talak. ID : ", rowData[0])
		return nil
	}

	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a
	JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id
	WHERE a.key = 'ppp'
	AND b.tujuan_id = 2;`)

	err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi)

	if err != nil {
		log.Log.Errorln("Gagal scan query notifikasi penetapan pp :", err, ". ID : ", rowData[0])
		return err
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataPpp, error := dbsipp.Execute("SELECT nomor_perkara, nama, c.keterangan FROM perkara AS a LEFT JOIN perkara_panitera_pn AS b USING(perkara_id) LEFT JOIN panitera_pn AS c ON b.panitera_id = c.id WHERE a.perkara_id = ? AND b.aktif = 'Y' AND b.`urutan` = 1", rowData[1])

	if error != nil {
		log.Log.Errorln("Gagal query data penetapan PP : ", err, ". ID : ", rowData[0])
		return err
	}

	noHpPP, err := queryDataPpp.GetStringByName(0, "keterangan")
	if err != nil {
		log.Log.Errorln("Gagal get uint by name di kolom keterangan : ", err)
		return err
	}

	nomorPerkara, err := queryDataPpp.GetStringByName(0, "nomor_perkara")
	if err != nil {
		log.Log.Errorln("Gagal get uint by name di kolom keterangan : ", err)
		return err
	}

	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)
	// whatsapp.SendMessage("internal", noHpPP, PesanNotifikasi)
	if noHpPP == "" {
		log.Log.Warningln("Skip Notif PPP. Nomor hp tidak ada")
		return errors.New("nomor hp tidak ada")
	}

	if err := SendNotifMessage("internal", noHpPP, PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal kirim notif PPP : ", err)
	}

	return nil
}

func NotifPJS(rowData []interface{}) error {
	if rowData[2] == "18" {
		log.Log.Warningln("Skip Notif Penetapan PP. Note : Penetapan ini bersifat ikrar talak. ID : ", rowData[0])
		return nil
	}

	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'pjs' AND b.tujuan_id = 3;`)

	err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi)

	if err != nil {
		log.Log.Errorln("Gagal scan query notifikasi penetapan jurusita :", err, ". ID : ", rowData[0])
		return err
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataPjs, error := dbsipp.Execute(`SELECT nomor_perkara, nama, c.keterangan 
	FROM perkara AS a LEFT JOIN perkara_jurusita AS b USING(perkara_id) 
	LEFT JOIN jurusita AS c ON b.jurusita_id = c.id 
	WHERE a.perkara_id = ? AND b.aktif = 'Y' `, rowData[1])

	if error != nil {
		log.Log.Errorln("Gagal query data penetapan jurusita :", err, ". ID : ", rowData[0])
		return err
	}

	for i := range queryDataPjs.Values {
		noHpJs, err := queryDataPjs.GetStringByName(i, "keterangan")
		if err != nil {
			log.Log.Errorln("Gagal get no hp js : ", err)
			return err
		}

		nomorPerkara, err := queryDataPjs.GetStringByName(i, "nomor_perkara")
		if err != nil {
			log.Log.Errorln("Gagal get nomor perkara notif pjs : ", err)
			return err
		}

		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)
		// whatsapp.SendMessage("internal", noHpJs, PesanNotifikasi)
		if err := SendNotifMessage("internal", noHpJs, PesanNotifikasi); err != nil {
			log.Log.Errorln("Gagal kirim notif PJS : ", err)
		}
	}

	return nil
}

func NotifPHS(rowData []interface{}) error {
	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	urutanSidang := rowData[5].(uint32)

	switch urutanSidang {
	case 1:
		queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'phs' AND b.tujuan_id = 3;`)

		if err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
			log.Log.Errorln("Gagal scan quuery notif phs : ", err)
			return err
		}
	case 2:
		queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'pts' AND b.tujuan_id = 3;`)

		if err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
			log.Log.Errorln("Gagal scan quuery notif phs : ", err)
			return err
		}
	default:
		log.Log.Warningln("Skip notifikasi phs to js : Urutan sidang melebihi 2")
		return nil
	}

	queryDataPjs, error := dbsipp.Execute(`SELECT nomor_perkara, nama, c.keterangan 
		FROM perkara AS a LEFT JOIN perkara_jurusita AS b USING(perkara_id) 
		LEFT JOIN jurusita AS c ON b.jurusita_id = c.id 
		WHERE a.perkara_id = ? AND b.aktif = 'Y' `, rowData[1])
	if error != nil {
		log.Log.Errorln("Gagal quuery data pjs notif phs : ", error)
		return error
	}

	for i := range queryDataPjs.Values {

		noHpJs, err := queryDataPjs.GetStringByName(i, "keterangan")
		if err != nil {
			log.Log.Errorln("Gagal get no hp js : ", err)
			return err
		}

		nomorPerkara, err := queryDataPjs.GetStringByName(i, "nomor_perkara")
		if err != nil {
			log.Log.Errorln("Gagal get nomor perkara notif pjs : ", err)
			return err
		}

		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)

		urutan, _ := rowData[5].(string)
		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{urutan}", urutan, -1)

		alasan, _ := rowData[16].(string)
		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{alasan_tunda}", alasan, -1)

		tanggal_sidang, _ := rowData[6].(string)
		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{tanggal_sidang}", whatsapp.ReadableDate(tanggal_sidang, nil), -1)

		kehadiran, _ := rowData[14].(string)

		kehadiranInt, _ := strconv.Atoi(kehadiran)
		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{status_pihak}", StatusPihak(kehadiranInt), -1)

		// whatsapp.SendMessage("internal", noHpJs, PesanNotifikasi)
		if err := SendNotifMessage("internal", noHpJs, PesanNotifikasi); err != nil {
			log.Log.Errorln("Gagal kirim notif PHS : ", err)
		}

	}

	return nil
}

func NotifSPV(rowData []interface{}) error {

	if rowData[2] != "Y" {
		log.Log.Warningln("Skip notif SPV. Note : Bukan perkara putus verstek")
		return nil
	}

	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'spv' AND b.tujuan_id = 3;`)

	if err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal scan query notif SPV : ", err)
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataPjs, error := dbsipp.Execute(`SELECT nomor_perkara, nama, c.keterangan 
	FROM perkara AS a LEFT JOIN perkara_jurusita AS b USING(perkara_id) 
	LEFT JOIN jurusita AS c ON b.jurusita_id = c.id 
	WHERE a.perkara_id = ? AND b.aktif = 'Y' `, rowData[0])

	if error != nil {
		log.Log.Errorln("Gagal quuery data pjs notif SPV : ", error)
		return error
	}

	for i := range queryDataPjs.Values {

		noHpJs, err := queryDataPjs.GetStringByName(i, "keterangan")
		if err != nil {
			log.Log.Errorln("Gagal get no hp js : ", err)
			return err
		}

		nomorPerkara, err := queryDataPjs.GetStringByName(i, "nomor_perkara")
		if err != nil {
			log.Log.Errorln("Gagal get nomor perkara notif pjs : ", err)
			return err
		}

		tanggalPutus, _ := rowData[1].(string)

		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)
		PesanNotifikasi = strings.Replace(PesanNotifikasi, "{tanggal_putus}", whatsapp.ReadableDate(tanggalPutus, nil), -1)

		// whatsapp.SendMessage("internal", noHpJs, PesanNotifikasi)
		if err := SendNotifMessage("internal", noHpJs, PesanNotifikasi); err != nil {
			log.Log.Errorln("Gagal kirim notif SPV : ", err)
		}
	}

	return nil
}

func NotifPAC(rowData []interface{}, rowOldData []interface{}) error {
	if rowData[3] == nil && rowData[4] == nil {
		log.Log.Warningln("Skip PAC. Data data masih nil")
		return errors.New("nomor dan tanggal ac masih nil. skip")
	}

	if rowOldData[3] != nil && rowOldData[4] == nil {
		log.Log.Warningln("Skip PAC. Duplikasi edit data")
		return errors.New("aksi terduplikasi")
	}

	timeTerbit, errParse := time.Parse("2006-01-02", rowData[4].(string))
	if errParse != nil {
		log.Log.Errorln("Gagal parsing tanggal terbit : ", errParse)
		return errParse
	}

	if !time.Now().Equal(timeTerbit) {
		log.Log.Warningln("Skip PAC : Tanggal terbit bukan hari ini")
		return errors.New("skip Notif PAC. Tanggal terbit tidak sesuai hari ini")
	}

	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'pac' AND b.tujuan_id = 5;`)

	if err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal scan query notif PAC : ", err)
		return err
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataPihak, err := dbsipp.Execute(`SELECT a.nama,telepon,nomor_perkara FROM pihak AS a JOIN perkara_pihak1 AS b ON a.id = b.pihak_id JOIN perkara AS r USING(perkara_id) WHERE b.perkara_id = ?`, rowData[0])

	if err != nil {
		log.Log.Errorln("Gagal query sipp notif PAC : ", err)
		return err
	}

	nomorTelepon, _ := queryDataPihak.GetStringByName(0, "telepon")
	nomor_perkara, _ := queryDataPihak.GetStringByName(0, "nomor_perkara")

	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomor_perkara, -1)
	nomorAC := rowData[3].(string)
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_ac}", nomorAC, -1)
	tanggalTerbit := rowData[4].(string)
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{tanggal_terbit}", whatsapp.ReadableDate(tanggalTerbit, nil), -1)

	// whatsapp.SendMessage("public", nomorTelepon, PesanNotifikasi)
	if nomorTelepon == "" {
		log.Log.Warningln("Skip Notif PAC. Nomor telepon kosong")
		return errors.New("nomor telepon kosong")
	}
	if err := SendNotifMessage("public", nomorTelepon, PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal kirim notif PAC : ", err)
	}

	return nil

}

func NotifPUR(rowData []interface{}) error {
	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'pur' AND b.tujuan_id = 1;`)

	if err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal scan query notif PUR : ", err)
		return err
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataHakim, err := dbsipp.Execute(`SELECT nomor_perkara, nama, c.keterangan FROM perkara AS a LEFT JOIN perkara_hakim_pn AS b USING(perkara_id) LEFT JOIN hakim_pn AS c ON b.hakim_id = c.id WHERE a.perkara_id = ? AND b.jabatan_hakim_id = 1 AND b.aktif = 'Y'`, rowData[1])

	if err != nil {
		log.Log.Errorln("Gagal query hakim notif PUR : ", err)
	}

	queryDataPihak, err := dbsipp.Execute(`SELECT nama FROM pihak WHERE id = ?`, rowData[4])

	if err != nil {
		log.Log.Errorln("Gagal query pihak notif PUR : ", err)
	}

	noHpHakim, err := queryDataHakim.GetStringByName(0, "keterangan")
	if err != nil {
		log.Log.Errorln("Gagal get uint by name di kolom keterangan : ", err)
		return err
	}

	nomorPerkara, err := queryDataHakim.GetStringByName(0, "nomor_perkara")
	if err != nil {
		log.Log.Errorln("Gagal get string by name di kolom nomor_perkara : ", err)
		return err
	}

	namaPihak, err := queryDataPihak.GetStringByName(0, "nama")
	if err != nil {
		log.Log.Errorln("Gagal get string by name di kolom nama : ", err)
		return err
	}

	var (
		tanggalRelaas    string
		statusRelaas     string
		keteranganRelaas []uint8
	)

	if rowData[6] != nil {
		tanggalRelaas = rowData[6].(string)
	}

	if rowData[8] != nil {
		statusRelaas = rowData[8].(string)
	}

	if rowData[9] != nil {
		keteranganRelaas = rowData[9].([]uint8)
	}

	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{tanggal_pelaksanaan}", whatsapp.ReadableDate(tanggalRelaas, nil), -1)
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{status_panggilan}", StatusRelaas(statusRelaas), -1)
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{detail}", string(keteranganRelaas), -1)
	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nama_pihak}", namaPihak, -1)

	// whatsapp.SendMessage("internal", noHpHakim, PesanNotifikasi)
	if err := SendNotifMessage("internal", noHpHakim, PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal kirim notif PUR : ", err)
	}
	return nil
}

func NotifPPD(rowData []interface{}) error {

	var (
		NamaNotifikasi  string
		PesanNotifikasi string
	)

	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	queryNotif := dblocal.QueryRow(`SELECT nama_notifikasi, pesan FROM jenis_notifikasi AS a JOIN notifikasi AS b ON a.id = b.jenis_notifikasi_id WHERE a.key = 'ppd' AND b.tujuan_id = 5;`)

	if err := queryNotif.Scan(&NamaNotifikasi, &PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal scan query notif PUR : ", err)
		return err
	}

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryData, err := dbsipp.Execute(`SELECT nomor_perkara,telepon FROM perkara_pihak1 JOIN perkara USING(perkara_id) JOIN pihak ON perkara_pihak1.pihak_id = pihak.id WHERE perkara.perkara_id = ? `, rowData[1])

	if err != nil {
		log.Log.Errorln("Gagal query data PPD : ", err)
		return err
	}

	nomorPerkara, err := queryData.GetStringByName(0, "nomor_perkara")
	if err != nil {
		log.Log.Errorln("Gagal get nomor_perkara di notif PPD : ", err)
		return err
	}

	teleponPihak, err := queryData.GetStringByName(0, "telepon")
	if err != nil {
		log.Log.Errorln("Gagal get telepon di notif PPD : ", err)
		return err
	}

	PesanNotifikasi = strings.Replace(PesanNotifikasi, "{nomor_perkara}", nomorPerkara, -1)

	// whatsapp.SendMessage("public", teleponPihak, PesanNotifikasi)
	if teleponPihak == "" {
		log.Log.Warningln("Skip Notif PPD. Nomor telepon kosong")
		return errors.New("nomor telepon kosong")
	}

	if err := SendNotifMessage("public", teleponPihak, PesanNotifikasi); err != nil {
		log.Log.Errorln("Gagal kirim notif PMH : ", err)
	}

	return nil
}

func StatusPihak(id int) string {
	if id == 1 {
		return "Semua Pihak"
	}
	if id == 2 {
		return "Penggugat Saja"
	}
	if id == 3 {
		return "Tergugat Saja"
	}
	if id == 10 {
		return "Sebagian Penggugat"
	}
	if id == 20 {
		return "Sebagian Tergugat"
	}
	return "Tidak dapat melihat status kehadiran"
}

func StatusRelaas(id string) string {
	if id == "Y" {
		return "Bertemu"
	}

	if id == "T" {
		return "Tidak Bertemu"
	}

	if id == "E" {
		return "Ecourt / E-Summon / Elektronik"
	}

	if id == "S" {
		return "Surat Tercatat / POS"
	}

	return "Tidak Dapat Dilaksanakan"
}
