package data

import (
	"context"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func InitDataAdmin(db *sql.DB) {

	defer log.Println("Insert table admin success")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("kuyabatok"), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal("Error saat bycrpt password")
	}

	_, err = db.ExecContext(context.Background(), `INSERT INTO admins (name, identifier, password, phone, status) VALUES (?, ?, ?, ?, ?)`, "Imal Malik", "admin", string(hashedPassword), "6289636811489", 1)

	if err != nil {
		log.Fatal("Error saat buat akun admin :", err)
	}
}
