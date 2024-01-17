package test

import (
	"context"
	"fmt"
	"gowhatsapp/whatsapp"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waproto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	godotenv.Load("../.env")
}

func TestCheckWaClientsExists(t *testing.T) {
	assert.Equal(t, whatsapp.WAClients["public"], nil)
}

func TestGetAllDevices(t *testing.T) {
	client := "dummy"

	dbName := fmt.Sprintf("%s_wa.db?_pragma=foreign_keys(1)", client)
	container, err := sqlstore.New("sqlite", dbName, nil)
	if err != nil {
		panic("Gagal set container :" + err.Error())
	}

	devices, err := container.GetAllDevices()
	if err != nil {
		panic("Gagal get all devices : " + err.Error())
	}

	for i, d := range devices {
		fmt.Println(i, ". This is user : ", d.ID.User)
		fmt.Println(i, ". This is pushname : ", d.PushName)
	}
}

func TestAuthSecondClient(t *testing.T) {
	client := "dummy"

	dbName := fmt.Sprintf("%s_wa.db?_pragma=foreign_keys(1)", client)
	container, err := sqlstore.New("sqlite", dbName, nil)
	if err != nil {
		panic("Gagal set container :" + err.Error())
	}

	devices := container.NewDevice()

	dummyClient := whatsmeow.NewClient(devices, nil)

	devices.Save()

	if dummyClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := dummyClient.GetQRChannel(context.Background())
		err = dummyClient.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {

				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stderr)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = dummyClient.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	time.Sleep(3 * time.Second)

	dummyClient.Disconnect()
}

func TestMultipleClientSendMessage(t *testing.T) {
	client := "dummy"

	dbName := fmt.Sprintf("%s_wa.db?_pragma=foreign_keys(1)", client)
	container, err := sqlstore.New("sqlite", dbName, nil)
	if err != nil {
		panic("Gagal set container :" + err.Error())
	}

	devices, err := container.GetAllDevices()
	if err != nil {
		panic("Gagal get all devices" + err.Error())
	}

	for _, d := range devices {
		wc := whatsmeow.NewClient(d, nil)

		err := wc.Connect()
		if err != nil {
			panic("Gagal start wa : " + err.Error())
		}

		time.Sleep(3 * time.Second)
		err = wc.SendPresence(types.PresenceAvailable)
		if err != nil {
			panic("Gagal start wa : " + err.Error())
		}

		jid, err := whatsapp.ParseJid(os.Getenv("DEV_WA_NUMBER"))
		if err != nil {
			panic("Gagal parsing : " + err.Error())
		}

		err = wc.SendChatPresence(jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		if err != nil {
			panic("Gagal cpc wa : " + err.Error())
		}

		time.Sleep(3 * time.Second)
		err = wc.SendChatPresence(jid, types.ChatPresencePaused, types.ChatPresenceMediaText)
		if err != nil {
			panic("Gagal cpp wa : " + err.Error())
		}

		time.Sleep(4 * time.Second)
		_, err = wc.SendMessage(context.Background(), jid, &waproto.Message{
			Conversation: proto.String("Test"),
		})

		if err != nil {
			panic("Gagal send message : " + err.Error())
		}

		time.Sleep(2 * time.Second)
		wc.Disconnect()
	}
	fmt.Println("done")
	// c := make(chan os.Signal)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// <-c
}

func TestConvertJID(t *testing.T) {
	nomorHp := "088294376130"
	parsed, err := whatsapp.ParseJid(nomorHp)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed)
}

func TestStartSessionWARestApi(t *testing.T) {
	// Happy path
	client := "testClient"
	whatsapp.StartSessionWARestApi(client)

}
