package main

import (
	"context"
	"flag"
	"fmt"
	"gowhatsapp/database"
	"gowhatsapp/events"
	"gowhatsapp/log"
	"gowhatsapp/server"
	"gowhatsapp/whatsapp"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mdp/qrterminal/v3"
)

var Do string
var WAC string

func init() {
	godotenv.Load()

	flag.StringVar(&Do, "do", "", "Menjabarkan kegiatan yang akan dilakukan")
	flag.StringVar(&WAC, "wa_client", "", "Nama WA Client yang akan digunakan")
}

func main() {
	log.LogInit()

	flag.Parse()

	switch Do {
	case "event_only":
		RunSippEventOnly()
	case "api_only":
		RunApiOnly()
	case "wa_only":
		RunWAOnly()
	case "auth":
		RunWaAuth()
	case "schedule":
		events.StartCronEvent()
		return
	default:
		RunMain()
	}
}

func RunMain() {
	database.InitDBLocal()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err := whatsapp.RunEngine()
		if err != nil {
			log.Log.Fatalln("Gagal running whatsapp engine : ", err)
		}

		defer wg.Done()
	}()

	wg.Wait()

	go server.StartServer()
	events.StartCronEvent()
	events.StartSippEvent()

	if os.Getenv("SINGLE_WA") == "0" {
		whatsapp.StopWa("default")
	} else {
		whatsapp.StopWa("inernal")
		whatsapp.StopWa("public")
	}
}

func RunSippEventOnly() {
	events.StartSippEvent()
}

func RunApiOnly() {
	server.StartServer()
}

func RunWAOnly() {
	wr := make(chan *whatsapp.WAResponse, 1)
	log.Log.Warn("WA Only hanya berlaku 1 client saja")

	if WAC == "" {
		log.Log.Fatalln("Gagal run WA Only : Parameter wa_client tidak ada")
	}

	if whatsapp.WAClients[WAC] == nil {
		whatsapp.AppendDevice(WAC, "", true)
		go whatsapp.StartWa(WAC, wr)
	} else {
		go whatsapp.StartWa(WAC, wr)
	}

	<-wr

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	<-c
	whatsapp.StopWa(WAC)
}

func RunWaAuth() {

	if WAC == "" {
		log.Log.Fatalln("Gagal run WA Only : Parameter wa_client tidak ada. Gunakan --wa_client=default untuk single wa. Gunakan --wa_client=public lalu setelah selesai gunakan --wa_client=internal untuk double device")
	}

	handler := true
	if WAC == "internal" {
		handler = false
	}

	if os.Getenv("SINGLE_WA") == "0" {
		fmt.Println("Perhatian. Siapkan 2 Device untuk WA Notif Internal dan Wa Notif Public. Jika anda ingin menggunakan 1 Device untuk kedua nya. Silahkan set env SINGLE_WA ke 1")
		time.Sleep(3 * time.Second)
	}

	fmt.Println("Siapkan Device Anda yang akan digunakan untuk Wa Notif ", strings.ToUpper(WAC))

	var waNumber string

	switch WAC {
	case "internal":
		waNumber = os.Getenv("INTERNAL_WA_NUMBER")
	case "public":
		waNumber = os.Getenv("PUBLIC_WA_NUMBER")
	default:
		waNumber = os.Getenv("DEFAULT_WA_NUMBER")
	}

	whatsapp.AppendDevice(WAC, waNumber, handler)

	if whatsapp.WAClients[WAC].Store.ID != nil {
		log.Log.Fatalln("Default client sudah diotentikasi")
	}

	qrChan, _ := whatsapp.WAClients[WAC].GetQRChannel(context.Background())

	err := whatsapp.WAClients[WAC].Connect()
	if err != nil {
		log.Log.Errorln("Whatsapp gagal connect :", err)
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stderr)
				log.Log.Infoln("Whatsapp qr code")

			} else {
				log.Log.Infoln("Whatsapp Login event:", evt.Event)
				wg.Done()
			}
		}
	}()

	wg.Wait()
	time.Sleep(3 * time.Second)

	err = whatsapp.WAClients[WAC].Store.Save()
	if err != nil {
		log.Log.Warningln("Gagal save device internal")
	}

	time.Sleep(5 * time.Second)
	whatsapp.StopWa(WAC)
	time.Sleep(5 * time.Second)

	log.Log.Info("Otentikasi WA Client selesai. Aplikasi akan di close. Silahkan nyalakan kembali")

	time.Sleep(5 * time.Second)
}
