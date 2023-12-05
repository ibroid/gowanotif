package data

import (
	"context"
	"database/sql"
	"log"
)

type TemplateReply struct {
	Id      int16  `json:"id"`
	Pesan   string `json:"pesan"`
	Trigger string `json:"trigger"`
}

func InitDataTemplateReply(db *sql.DB) {

	defer log.Println("Insert table template success")

	data := []TemplateReply{
		{
			Id:      1,
			Pesan:   "Selamat datang di layanan Chat dan Notifikasi.\nSilahkan ketik : PROSES PERKARA lalu kirim ke nomor ini untuk melihat proses perkara anda.\nUntuk bertanya ke Admin/Petugas, silahkan ketik : TANYA PETUGAS, lalu kirim ke nomor ini. Setelah itu tunggu petugas untuk membalas pesan anda.\n\n_Dikirim secara otomatis. Mohon untuk *Tidak menelefon ke nomor ini*. Terima kasih_",
			Trigger: "INFO",
		},
		{
			Id:      2,
			Pesan:   "Silahkan Tunggu. Anda akan dihubungi jika admin/petugas sudah siap. Untuk permintaan pada hari libur akan di balas pada hari kerja.",
			Trigger: "TANYA PETUGAS",
		},
		{
			Id:      3,
			Pesan:   "Assalamualaikum. Terima kasih telah menghubungi.\nSilahkan ketik *info* lalu kirim ke nomor ini.\n\n_Pesan ini dikirim secara otomatis_",
			Trigger: "DEFAULT",
		},
		{
			Id:      4,
			Pesan:   "*PROSES PERKARA*.\n[PENDAFTARAN]\n{info_perkara}\n\n[PERSIDANGAN]\n{daftar_persidangan}\n\n[TRANSAKSI]\n{daftar_transaski}\n{total_transaksi}\n\n[SIDANG IKRAR]\n{tanggal_ikrar}\n\n[AKTA CERAI]\n{info_akta_cerai}",
			Trigger: "PROSES PERKARA",
		},
	}

	for _, q := range data {
		_, err := db.ExecContext(context.TODO(), "INSERT INTO template_reply (id, pesan, trigger) VALUES (?, ?, ?)", q.Id, q.Pesan, q.Trigger)
		if err != nil {
			log.Fatal("gagal exec insert tabel template_reply")
		}
	}
}
