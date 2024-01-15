package whatsapp

import (
	"context"
	"gowhatsapp/log"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
)

type WAAuthRequest struct {
	ClientName  string `json:"client_name"`
	PhoneNumber string `json:"phone_number"`
	Handler     bool   `json:"handler"`
}

type WAStartRequest struct {
	ClientName  string `json:"client_name"`
	PhoneNumber string `json:"phone_number"`
}

func StartAuthedWClient(req WAStartRequest) {
	if !CheckClient(req.ClientName) {
		panic("WAC not initialized")
	}

	device, err := setDevice(req.ClientName, req.PhoneNumber)
	if err != nil {
		panic("Error setting device: " + err.Error())
	}

	WAClients[req.ClientName] = whatsmeow.NewClient(device, nil)
	err = WAClients[req.ClientName].Connect()
	if err != nil {
		panic("Whatsapp client failed to connect: " + err.Error())
	}
}

func ClientAuthentication(req WAAuthRequest, wareschan chan<- WAResponse) {
	if CheckClient(req.ClientName) {
		panic("WAC already initialized")
	}

	device, err := setDevice(req.ClientName, req.PhoneNumber)
	if err != nil {
		panic("Error setting device: " + err.Error())
	}

	if !WAClients[req.ClientName].IsConnected() {
		WAClients[req.ClientName] = whatsmeow.NewClient(device, nil)
		err = WAClients[req.ClientName].Connect()
		if err != nil {
			panic("Whatsapp client failed to connect: " + err.Error())
		}
	}

	ctxQr, ctxQrCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer ctxQrCancel()

	qrChan, _ := WAClients[req.ClientName].GetQRChannel(ctxQr)

	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stderr)
			log.Log.Infoln("Whatsapp qr code")

		} else {
			log.Log.Infoln("Whatsapp Login event:", evt.Event)

			ClientDisconnect(req.ClientName)
		}
	}
}

func ClientDisconnect(client string) {
	time.Sleep(time.Second * 6)
	if CheckClient(client) {
		WAClients[client].Disconnect()
	}
}
