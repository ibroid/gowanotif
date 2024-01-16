package data

import (
	"context"
	"database/sql"
	"log"
)

func CreateAllTable(db *sql.DB) {
	defer log.Println("Success generate base table")

	queryList := []string{
		// CREATE TABLE ADMIN
		`CREATE TABLE IF NOT EXISTS admins (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, identifier TEXT, password TEXT, phone TEXT, status INTEGER ); `,

		// CREATE TABLE PENGATURAN
		`CREATE TABLE pengaturan (id INTEGER PRIMARY KEY AUTOINCREMENT, key TEXT, value TEXT, ket TEXT) `,

		// CREATE TABLE JENIS_NOTIFIKASI
		`CREATE TABLE IF NOT EXISTS jenis_notifikasi (id INTEGER PRIMARY KEY AUTOINCREMENT, nama_notifikasi TEXT, key TEXT);`,

		// CREATE TABLE TUJUAN
		`CREATE TABLE IF NOT EXISTS tujuan (id INTEGER PRIMARY KEY AUTOINCREMENT, nama TEXT );`,

		// CREATE TABLE NOTIFIKASI
		`CREATE TABLE IF NOT EXISTS notifikasi (id INTEGER PRIMARY KEY AUTOINCREMENT, jenis_notifikasi_id INTEGER, tujuan_id INTEGER, pesan TEXT, filename TEXT NULL);`,

		// CREATE NOREPLY
		`CREATE TABLE IF NOT EXISTS no_reply (id INTEGER PRIMARY KEY AUTOINCREMENT, number TEXT, valid DATETIME);`,

		// CREATE NOREPLY
		`CREATE TABLE IF NOT EXISTS template_reply (id INTEGER PRIMARY KEY AUTOINCREMENT, pesan TEXT, trigger TEXT);`,

		// CREATE NOREPLY
		`CREATE TABLE IF NOT EXISTS clients (id INTEGER PRIMARY KEY AUTOINCREMENT, client_name TEXT, jid TEXT, handler INTEGER, status INTEGER, service TEXT);`,
	}

	for no, query := range queryList {

		func(query string, no int) {

			_, err := db.ExecContext(context.Background(), query)
			if err != nil {
				log.Fatalln("Error exec table : ", err)
			}

		}(query, no)
	}

}
