package data

import (
	"database/sql"
	"log"
)

type Notifikasi struct {
	Id                int16          `json:"id"`
	JenisNotifikasiId int16          `json:"jenis_notifikasi_id"`
	TujuanId          int16          `json:"tujuan_id"`
	Pesan             string         `json:"pesan"`
	Filename          sql.NullString `json:"filename"`
}

func InitDataNotifikasi(db *sql.DB) {
	defer log.Println("Insert table notifikasi success")

	data := []Notifikasi{
		{
			Id:                1,
			JenisNotifikasiId: 1,
			TujuanId:          1,
			Pesan: `*NOTIFIKASI PMH* 
Anda ditunjuk sebagai Ketua Majelis Hakim pada perkara nomor {nomor_perkara}.
Silahkan buka SIPP untuk penetapan sidang pertama`,
			Filename: sql.NullString{
				String: "add_majelis.php",
				Valid:  true,
			},
		},
		{
			Id:                2,
			JenisNotifikasiId: 2,
			TujuanId:          2,
			Pesan: `*NOTIFIKASI PENETAPAN PANITERA PENGGANTI*
Anda ditunjuk sebagai Panitera Pengganti pada perkara nomor {nomor_perkara}.`,
			Filename: sql.NullString{
				String: "add_panitera_m.php",
				Valid:  true,
			},
		},
		{
			Id:                3,
			JenisNotifikasiId: 3,
			TujuanId:          3,
			Pesan: `*NOTIFIKASI PENUNJUKAN JURUSITA*
Anda ditunjuk sebagai jurusita pada perkara nomor {nomor_perkara}`,
			Filename: sql.NullString{
				String: "add_jurusita_m.php",
				Valid:  true,
			},
		},
		{
			Id:                4,
			JenisNotifikasiId: 4,
			TujuanId:          5,
			Pesan: `*NOTIFIKASI SIDANG PERTAMA*
Nomor Perkara : {nomor_perkara}
Para pihak : {para_pihak}
Tanggal Sidang : {tanggal_sidang}
Agenda : {agenda_sidang}`,
			Filename: sql.NullString{
				String: "add_sidang_pertama_m.php",
				Valid:  true,
			},
		},
		{
			Id:                5,
			JenisNotifikasiId: 4,
			TujuanId:          3,
			Pesan: `*NOTIFIKASI PHS*
Nomor perkara {nomor_perkara} telah ditetapkan hari sidang pertama pada tanggal {tanggal_sidang}`,
			Filename: sql.NullString{
				String: "add_sidang_pertama_m.php",
				Valid:  true,
			},
		},
		{
			Id:                6,
			JenisNotifikasiId: 5,
			TujuanId:          3,
			Pesan: `*NOTIFIKASI TUNDAAN SIDANG*
Sidang ke {urutan} nomor perkara {nomor_perkara} telah ditunda dengan alasan {alasan_tunda}, tunda ke tanggal {tanggal_sidang}. Kehadiran pihak {status_pihak}.

Silahkan cek di sipp untuk melihat pihak yang anda panggil hadir dalam persidangan.`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                7,
			JenisNotifikasiId: 5,
			TujuanId:          5,
			Pesan: `*NOTIFIKASI SIDANG SELANJUTNYA*
Nomor Perkara: {nomor_perkara}
Sidang Selanjutnya : {tanggal_sidang}
Agenda : {agenda_sidang}`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                8,
			JenisNotifikasiId: 6,
			TujuanId:          3,
			Pesan: `*NOTIFIKASI PUTUS VERSTEK*
Berikut adalah perkara dengan putusan verstek.
Nomor Perkara : {nomor_perkara}
Tanggal Putus : {tanggal_putus}
Silahkan periksa kehadiran pihak yang anda panggil`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                9,
			JenisNotifikasiId: 6,
			TujuanId:          5,
			Pesan: `*NOTIFIKASI SISA PANJAR*
Nomor Perkara : {nomor_perkara}
Sisa Panjar : {sisa_panjar}
{terbilang_sisa_panjar}
_Sisa panjar belum termasuk biaya Pemberitahuan isi putusan (apabila verstek)_`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                10,
			JenisNotifikasiId: 8,
			TujuanId:          5,
			Pesan: `*NOTIFIKASI AKTA CERAI*
Nomor Perkara : {nomor_perkara}
Nomor AC : {nomor_ac}
Tanggal terbit : {tanggal_terbit}
Silahkan ambil akta cerai anda sehari setelah pesan ini disampaikan.
Pastikan pengambilan akta cerai pada hari kerja.`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                11,
			JenisNotifikasiId: 7,
			TujuanId:          1,
			Pesan: `*NOTIFIKASI SEBELUM SIDANG*
Berikut nomor perkara yang sidang hari ini namun belum ada hasil relaas panggilan.

{daftar_relaas}`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                12,
			JenisNotifikasiId: 9,
			TujuanId:          1,
			Pesan: `*NOTIFIKASI HASIL RELAAS*
Berikut hasil relaas panggilan nomor {nomor_perkara}.
Tanggal Pelaksanaan : {tanggal_pelaksanaan}
Nama Pihak : {nama_pihak}
Status : {status_panggilan} ({detail})`,
			Filename: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
		{
			Id:                13,
			JenisNotifikasiId: 4,
			TujuanId:          2,
			Pesan: `*NOTIFIKASI PHS*
Nomor perkara {nomor_perkara} telah ditetapkan hari sidang pertama pada tanggal {tanggal_sidang}`,
			Filename: sql.NullString{
				String: "add_sidang_pertama_m.php",
				Valid:  true,
			},
		},
		{
			Id:                14,
			JenisNotifikasiId: 7,
			TujuanId:          3,
			Pesan: `*NOTIFIKASI PERINGATAN RELAAS*
Berikut daftar relaas yangg belum terupload untuk sidang hari ini
{daftar_relaas}`,
			Filename: sql.NullString{
				String: "add_sidang_pertama_m.php",
				Valid:  true,
			},
		},
		{
			Id:                15,
			JenisNotifikasiId: 10,
			TujuanId:          5,
			Pesan: `*NOTIFIKASI PENDAFTARAN PERKARA*
Perkara anda berhasil didaftarkan dengan nomor perkara : {nomor_perkara}. Nomor perkara ini sebagai acuan identitas perkara anda di Pengadilan Agama Jakarta Utara. Untuk melihat proses perkara anda silahkan ketik : PROSES PERKARA, kirim ke nomor ini.

_Dikirim secara otomatis. Mohon untuk *tidak membalas pesan ini* dan *tidak menelefon ke nomor ini*. Terima Kasih_`,
			Filename: sql.NullString{
				String: "add_sidang_pertama_m.php",
				Valid:  true,
			},
		},
	}

	for _, q := range data {
		_, err := db.Exec(`INSERT INTO notifikasi (id, jenis_notifikasi_id, tujuan_id, pesan, filename) VALUES (?, ?, ?, ?, ?)`, q.Id, q.JenisNotifikasiId, q.TujuanId, q.Pesan, q.Filename)
		if err != nil {
			log.Fatal("gagal exec insert table notifikasi : ", err)
		}
	}
}
