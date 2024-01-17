package whatsapp

import (
	"fmt"
	"gowhatsapp/log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SendMessageUsingWARestApi(client string, number string, message string) error {
	ua := fiber.Post("http://localhost:3000/api/sendText")
	ua.CloseIdleConnections()
	return nil

}

func StartSessionWARestApi(client string) {
	ua := fiber.Post("http://localhost:3000/session/start/" + client)
	ua.Add("x-api-key", "your_global_api_key_here")

	statusCode, body, errs := ua.Bytes()
	if len(errs) > 0 {
		log.Log.Panic(errs)
	}
	if statusCode != 200 {
		log.Log.Panic(fmt.Sprintf("Failed to start session, status code: %d, body: %s", statusCode, string(body)))
	}
	log.Log.Info("Session started")
	log.Log.Info(string(body))

}

func InitializeThirdParty() {

	// Pengecekan apakah Node.js sudah terinstal
	if !isNodeInstalled() {
		log.Log.Info("Node.js belum terinstal. Memulai proses instalasi...")

		// Menjalankan perintah untuk mengunduh dan menginstal Node.js
		cmd := exec.Command("powershell", "Invoke-WebRequest", "-Uri", "https://nodejs.org/dist/latest/win-x64/node.exe", "-OutFile", "node_installer.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			log.Log.Panic("Gagal mengunduh atau menginstal Node.js:", err)
			return
		}

		// Menjalankan perintah instalasi Node.js
		cmd = exec.Command("node_installer.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			log.Log.Panic("Gagal menginstal Node.js:", err)
			return
		}

		log.Log.Info("Node.js berhasil diinstal!")
	} else {
		log.Log.Info("Node.js sudah terinstal. Tidak perlu melakukan instalasi.")
	}

	// Mengeksekusi file index.js setelah instalasi Node.js
	log.Log.Info("Menjalankan aplikasi...")
	runNodeApp()
}

func isNodeInstalled() bool {
	// Perintah untuk mencari lokasi instalasi Node.js dan npm
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "node", "npm")
	} else {
		cmd = exec.Command("which", "node", "npm")
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Jika output tidak mengandung "not found", berarti Node.js sudah terinstal
	return !strings.Contains(string(output), "not found")
}

// Fungsi untuk mengeksekusi file index.js
func runNodeApp() {
	cmd := exec.Command("node", "../whatsapp-api/server.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Log.Panicln("Gagal mengeksekusi file index.js:", err)
		return
	}

	log.Log.Info("Aplikasi berhasil dijalankan.")
}
