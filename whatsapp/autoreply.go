package whatsapp

import (
	"database/sql"
	"errors"
	"fmt"
	"gowhatsapp/database"
	"gowhatsapp/log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow/types"
)

type AutoReply struct {
	message string
	sender  string
	client  string
}

func (ar *AutoReply) messenger(number string, message string) error {
	dest, err := ParseJid(number)
	if err != nil {
		return err
	}

	if err := SendMessage(ar.client, dest, message); err != nil {
		return err
	}

	return nil
}

func (ar *AutoReply) Init() error {

	if os.Getenv("MAINTENANCE") == "1" {
		ar.messenger(ar.sender, "Layanan sedang dalam proses maintenance. Mohon tunggu beberapa saat lagi")
		return nil
	}

	if err := ar.CheckWaitList(); err != nil {
		log.Log.Warningln(err)
		return err
	}

	switch strings.ToUpper(ar.message) {
	case "INFO":
		ar.Info()

	case "TANYA PETUGAS":
		ar.TanyaPetugas()

	default:
		if strings.Contains(strings.ToUpper(ar.message), "PROSES PERKARA") {

			ar.ProsesPerkara()
		} else {

			ar.UnknownMessage()
		}
	}

	return nil
}

func (ar *AutoReply) CheckWaitList() error {
	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	today := time.Now()
	var Number string

	queryCekSender := dblocal.QueryRow("SELECT number FROM no_reply WHERE number = ? AND DATE(valid) = ?", ar.sender, today.Format("2006-01-02"))

	if err := queryCekSender.Scan(&Number); err != nil {
		if err == sql.ErrNoRows {
			// Handle if user has not in pending session
		} else {
			return err
		}
	}

	if Number != "" {
		return errors.New("skip auto reply. sender sedang pending")
	}

	return nil
}

func (ar *AutoReply) ProsesPerkara() error {
	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	splitMessage := strings.Split(ar.message, " ")
	var perkaraId string

	if len(splitMessage) == 2 {

		if err := ar.CheckIdentity(); err != nil {
			ar.messenger(ar.sender, "Maaf terjadi kesalahan. "+err.Error())
			return err
		}

		var err error
		perkaraId, err = PerkaraIDFromJID(ar.sender)

		if err != nil {
			ar.messenger(ar.sender, "Terjadi Kesalahan. Silahkan tunggu beberapa saat lagi.")
			return err
		}
	}

	if len(splitMessage) == 3 {
		var err error
		perkaraId, err = PerkaraIDFromPar3(splitMessage[2])

		if err != nil {
			ar.messenger(ar.sender, "Terjadi Kesalahan. Silahkan tunggu beberapa saat lagi.")
			log.Log.Errorln(err)
			return err
		}
	}

	recipent, _ := ParseJid(ar.sender)

	WAClients[ar.client].SendChatPresence(recipent, types.ChatPresenceComposing, types.ChatPresenceMediaText)

	var (
		PesanReply string
	)

	queryReply := dblocal.QueryRow("SELECT pesan FROM template_reply WHERE trigger = 'PROSES PERKARA'")
	queryReply.Scan(&PesanReply)

	var wg sync.WaitGroup
	wg.Add(5)

	// fetch proses pendaftaran
	go prosesFetchPendaftaran(&wg, &perkaraId, &PesanReply)

	// fetch proses persidangan
	go prosesFetchPersidangan(&wg, &perkaraId, &PesanReply)

	// fetch proses transaksi
	go prosesFetchTransaksi(&wg, &perkaraId, &PesanReply)

	// fetch proses ikrar
	go prosesFetchIkrar(&wg, &perkaraId, &PesanReply)

	// fetch proses akta
	go prosesFetchAktaCerai(&wg, &perkaraId, &PesanReply)

	wg.Wait()
	time.Sleep(15 * time.Second)

	WAClients[ar.client].SendChatPresence(recipent, types.ChatPresencePaused, types.ChatPresenceMediaText)

	ar.messenger(ar.sender, PesanReply)
	return nil
}

func (ar *AutoReply) UnknownMessage() {
	ar.messenger(ar.sender, "Mohon maaf kata kunci tidak dikenali. Silahkan Ketik : Info, lalu kirim ke nomor ini.")
}

func (ar *AutoReply) Info() error {
	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	var Pesan string

	queryReply := dblocal.QueryRow("SELECT pesan FROM template_reply WHERE trigger = 'INFO'")

	if err := queryReply.Scan(&Pesan); err != nil {
		log.Log.Errorln("Gagal scan pesan info : ", err)
		return err
	}

	ar.messenger(ar.sender, Pesan)
	return nil
}

func (ar *AutoReply) TanyaPetugas() {
	dblocal := database.StartDBLocal()
	defer dblocal.Close()

	_, err := dblocal.Exec("INSERT INTO no_reply (number, valid) VALUES (?, ?)", ar.sender, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Log.Println("Gagal insert no reply at tanya petugas : ", err)
	}
	queryReply := dblocal.QueryRow("SELECT pesan FROM template_reply WHERE trigger = 'TANYA PETUGAS'")

	var Pesan string
	queryReply.Scan(&Pesan)

	ar.messenger(ar.sender, Pesan)
}

func (ar *AutoReply) CheckIdentity() error {

	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	queryDataPihak, err := dbsipp.Execute("SELECT telepon,id FROM pihak WHERE telepon = ?", ar.sender)

	if err != nil {
		// log.Log.Errorln("Terjadi Kesalahan saat reply proses perkara query sipp : ", err)
		return err
	}

	idPihak, _ := queryDataPihak.GetStringByName(0, "id")

	if idPihak == "" {
		return errors.New("nomor ini tidak ter register di pengadilan agama jakarta utara. Coba tambahkan nomor perkara di akhir. Contoh : PROSES PERKARA " + os.Getenv("DEMO_PERKARA"))
	}

	return nil
}

func prosesFetchPendaftaran(wg *sync.WaitGroup, perkaraId *string, pesan *string) {
	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	defer wg.Done()

	queryDataPerkara, err := dbsipp.Execute("SELECT nomor_perkara, jenis_perkara_text, para_pihak, tanggal_pendaftaran FROM perkara WHERE perkara_id = ?", *perkaraId)
	if err != nil {
		log.Log.Errorln("Gagal query data perkara di info proses : ", err)
	}
	nomorPerkara, _ := queryDataPerkara.GetStringByName(0, "nomor_perkara")
	jenisPperkara, _ := queryDataPerkara.GetStringByName(0, "jenis_perkara_text")
	paraPihak, _ := queryDataPerkara.GetStringByName(0, "para_pihak")
	tanggalDaftar := ReadableDate(queryDataPerkara.GetStringByName(0, "tanggal_pendaftaran"))

	pesanReply := fmt.Sprintf("Nomor Perkara : %v\nJenis Perkara : %v\nTanggal Daftar : %v\nPara Pihak : %v", nomorPerkara, jenisPperkara, tanggalDaftar, strings.ReplaceAll(paraPihak, "<br />", "\n"))

	*pesan = strings.Replace(*pesan, "{info_perkara}", pesanReply, -1)
}

func prosesFetchPersidangan(wg *sync.WaitGroup, perkaraId *string, pesan *string) {
	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	defer wg.Done()

	queryPerkaraSidang, err := dbsipp.Execute("SELECT tanggal_sidang, agenda, alasan_ditunda, urutan, keterangan FROM  perkara_jadwal_sidang WHERE perkara_id = ?", *perkaraId)
	if err != nil {
		log.Log.Errorln("Gagal query data sidang perkara di info proses : ", err)
	}

	var infoSidangTemplate string

	if isnull, _ := queryPerkaraSidang.IsNull(0, 0); isnull {
		infoSidangTemplate = "Belum Ditentukan"

	} else {

		for v := range queryPerkaraSidang.Values {
			readableTanggalSidang, _ := queryPerkaraSidang.GetStringByName(v, "tanggal_sidang")
			agenda, _ := queryPerkaraSidang.GetStringByName(v, "agenda")
			alasan, _ := queryPerkaraSidang.GetStringByName(v, "alasan_ditunda")
			keterangan, _ := queryPerkaraSidang.GetStringByName(v, "keterangan")
			urutan, _ := queryPerkaraSidang.GetStringByName(v, "urutan")

			infoSidangTemplate += fmt.Sprintf("*Sidang ke %v* tanggal %v. Agenda : %v. Alasan ditunda : %v. %v\n", urutan, readableTanggalSidang, agenda, alasan, keterangan)
		}

	}

	*pesan = strings.Replace(*pesan, "{daftar_persidangan}", infoSidangTemplate, -1)
}

func prosesFetchTransaksi(wg *sync.WaitGroup, perkaraId *string, pesan *string) {
	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	defer wg.Done()

	queryPerkaraTransakski, err := dbsipp.Execute("SELECT jenis_transaksi, uraian, jumlah, tanggal_transaksi FROM perkara_biaya WHERE perkara_id = ?", *perkaraId)
	if err != nil {
		log.Log.Errorln("Gagal query data transaksi perkara di info proses : ", err)
	}

	var infoTransaskiTemplate string
	var totalTransaskiTemplate string
	var countTransaksiMasuk int
	var countTransaksiKeluar int
	for t := range queryPerkaraTransakski.Values {
		// readableTanggalTransaksi := ReadableDate(queryPerkaraTransakski.GetStringByName(t, "tanggal_transaksi"))
		readableTanggalTransaksi, _ := queryPerkaraTransakski.GetStringByName(t, "tanggal_transaksi")
		jenis_transaksi, _ := queryPerkaraTransakski.GetStringByName(t, "jenis_transaksi")
		uraian, _ := queryPerkaraTransakski.GetStringByName(t, "uraian")
		jumlah, _ := queryPerkaraTransakski.GetStringByName(t, "jumlah")

		var jenis string
		if jenis_transaksi == "-1" {
			jenis = "keluar"

			tr, _ := strconv.Atoi(strings.Replace(jumlah, ".00", "", -1))
			countTransaksiKeluar += tr
		} else {
			jenis = "masuk"
			tr, _ := strconv.Atoi(strings.Replace(jumlah, ".00", "", -1))
			countTransaksiMasuk += tr
		}

		infoTransaskiTemplate += fmt.Sprintf("*Transaksi %v* senilai %v pada tanggal %v. Keterangan : %v\n", jenis, strings.Replace(jumlah, ".00", "", -1), readableTanggalTransaksi, uraian)
	}

	totalTransaskiTemplate = strconv.Itoa(countTransaksiMasuk - countTransaksiKeluar)

	*pesan = strings.Replace(*pesan, "{daftar_transaski}", infoTransaskiTemplate, -1)
	*pesan = strings.Replace(*pesan, "{total_transaksi}", "Sisa Panjar : "+totalTransaskiTemplate, -1)
}

func prosesFetchIkrar(wg *sync.WaitGroup, perkaraId *string, pesan *string) {
	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	defer wg.Done()

	queryPerkaraSidang, err := dbsipp.Execute("SELECT tanggal_sidang, agenda, alasan_ditunda, urutan keterangan FROM  perkara_jadwal_sidang WHERE perkara_id = ? AND ikrar_talak = 'Y'", *perkaraId)
	if err != nil {
		log.Log.Errorln("Gagal query data sidang ikrar perkara di info proses : ", err)
	}

	infoSidangTemplate, _ := queryPerkaraSidang.GetStringByName(0, "tanggal_sidang")

	if infoSidangTemplate == "" {
		infoSidangTemplate = "TIDAK/BELUM ADA IKRAR"
	}

	*pesan = strings.Replace(*pesan, "{tanggal_ikrar}", infoSidangTemplate, -1)
}

func prosesFetchAktaCerai(wg *sync.WaitGroup, perkaraId *string, pesan *string) {
	dbsipp := database.StartDBSipp()
	defer dbsipp.Close()

	defer wg.Done()

	queryPerkaraAkta, err := dbsipp.Execute("SELECT nomor_akta_cerai,tgl_akta_cerai,no_seri_akta_cerai FROM perkara_akta_cerai WHERE perkara_id = ?", *perkaraId)
	if err != nil {
		log.Log.Errorln("Gagal query data akta cerai perkara di info proses : ", err)
	}

	var infoAktaTemplate string

	nomorAkta, _ := queryPerkaraAkta.GetStringByName(0, "nomor_akta_cerai")

	if nomorAkta != "" {
		tanggalAkta, _ := queryPerkaraAkta.GetStringByName(0, "tgl_akta_cerai")
		nomorSeri, _ := queryPerkaraAkta.GetStringByName(0, "no_seri_akta_cerai")
		infoAktaTemplate = fmt.Sprintf("Nomor Akta : %v\nTanggal Terbit : %v\nNomor Seri :%v", nomorAkta, tanggalAkta, nomorSeri)
	} else {
		infoAktaTemplate = "TIDAK/BELUM TERBIT AKTA CERAI"
	}

	*pesan = strings.Replace(*pesan, "{info_akta_cerai}", infoAktaTemplate, -1)
}

func MaskingDev(sender *string) {
	if os.Getenv("DEV") == "1" {
		*sender = os.Getenv("DEMO_WA_NUMBER")
	}
}
