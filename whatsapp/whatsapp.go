package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"gowhatsapp/log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/gofiber/contrib/websocket"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var WAClients = make(map[string]*whatsmeow.Client)
var WALogs = make(map[string]waLog.Logger)

var IsRunning bool = false

type WAResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Qr      string `json:"qr"`
}

func RunEngine() error {
	if os.Getenv("SINGLE_WA") == "0" {
		if os.Getenv("INTERNAL_WA_NUMBER") == "" || os.Getenv("PUBLIC_WA_NUMBER") == "" {
			return errors.New("wa number public dan internal tidak ditemukan")
		}

		container, err := sqlstore.New("sqlite", "production_wa.db?_pragma=foreign_keys(1)", nil)
		if err != nil {
			return err
		}

		var devices []*store.Device
		if devices, err = container.GetAllDevices(); err != nil {
			return err
		}

		*store.DeviceProps.Os = *proto.String("Windows")
		*store.DeviceProps.PlatformType = *waProto.DeviceProps_CHROME.Enum()

		if len(devices) == 2 {
			for _, d := range devices {

				if d.ID.User == os.Getenv("PUBLIC_WA_NUMBER") {
					wr := make(chan *WAResponse, 1)
					WAClients["public"] = whatsmeow.NewClient(d, nil)
					SetEvenHandler("public")
					StartWa("public", wr)
					<-wr
					log.Log.Infoln("Whatsapp JID : ", WAClients["public"].Store.ID.User)
					log.Log.Infoln("Whatsapp NAME : ", WAClients["public"].Store.PushName)
					log.Log.Infoln("Whatsapp Handler : True")
				}

				if d.ID.User == os.Getenv("INTERNAL_WA_NUMBER") {
					wr := make(chan *WAResponse, 1)
					WAClients["internal"] = whatsmeow.NewClient(d, nil)
					StartWa("internal", wr)
					<-wr
					log.Log.Infoln("Whatsapp JID : ", WAClients["internal"].Store.ID.User)
					log.Log.Infoln("Whatsapp NAME : ", WAClients["internal"].Store.PushName)
					log.Log.Infoln("Whatsapp Handler : False")
				}
			}
		} else {
			return errors.New("tidak ditemukan device di dalam database. silahkan auth menggunakan --do=auth")
		}

	} else {
		wr := make(chan *WAResponse, 1)
		AppendDevice("default", os.Getenv("DEFAULT_WA_NUMBER"), false)
		StartWa("default", wr)
		<-wr

		log.Log.Infoln("Whatsapp JID : ", WAClients["default"].Store.ID.User)
		log.Log.Infoln("Whatsapp NAME : ", WAClients["default"].Store.PushName)
		log.Log.Infoln("Whatsapp Handler : True")
	}

	return nil
}

func setDevice(client string, number string) *store.Device {
	container := setContainer(client)

	if client == "dummy" || client == "default" {
		device, err := container.GetFirstDevice()

		if err != nil {
			log.Log.Fatalln("Gagal get first device :", err)
		}

		return device
	}

	jid, _ := ParseJid(number)

	oldDevice, err := container.GetDevice(jid)
	if err != nil {
		log.Log.Fatalln("Gagal get device : ", err)
	}

	if oldDevice != nil {
		return oldDevice
	}

	device := container.NewDevice()

	return device
}

func setContainer(client string) *sqlstore.Container {
	var dbName string

	if client == "dummy" {
		dbName = fmt.Sprintf("%s_wa.db?_pragma=foreign_keys(1)", client)
	} else {
		dbName = fmt.Sprintf("%s_wa.db?_pragma=foreign_keys(1)", "production")
	}

	container, err := sqlstore.New("sqlite", dbName, nil)
	if err != nil {
		log.Log.Panic(client, " Gagal store : ", err)
	}

	return container
}

func SetEvenHandler(client string) {

	WAClients[client].AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			log.Log.Debugln("Event handler trigered")
			if !v.Info.MessageSource.IsFromMe && !v.Info.Sender.IsEmpty() && v.Message.GetConversation() != "" {
				log.Log.Println("Pesan Baru :", v.Message.GetConversation(), ". Dari : ", v.Info.Sender.User)

				var wg sync.WaitGroup

				wg.Add(1)
				go ReadMessage(client, v.Info, &wg)
				wg.Wait()

				// SendMessage(client, v.Info.Sender.User, "Test Reply")
				autoReply := &AutoReply{
					message: v.Message.GetConversation(),
					sender:  v.Info.Sender.User,
					client:  client,
				}

				autoReply.Init()
			}
		}
	})
}

func ReadMessage(client string, msg types.MessageInfo, wg *sync.WaitGroup) {
	if err := WAClients[client].MarkRead([]string{msg.ID}, time.Now(), msg.Chat, msg.Sender); err != nil {
		log.Log.Errorln("Gagal read message :", err)
	} else {
		log.Log.Infoln("Sukses read message : ", msg.ID)
	}
	defer wg.Done()
}

func AppendDevice(client string, number string, handler bool) {
	// WALogs[client] = waLog.Stdout("Client", "DEBUG", true)
	WAClients[client] = whatsmeow.NewClient(setDevice(client, number), nil)
	*store.DeviceProps.Os = *proto.String("Windows")
	*store.DeviceProps.PlatformType = *waProto.DeviceProps_CHROME.Enum()
	if handler {
		SetEvenHandler(client)
	}
}

func StartWa(client string, wr chan<- *WAResponse) {
	if WAClients[client].Store.ID == nil {
		qrChan, _ := WAClients[client].GetQRChannel(context.Background())
		err := WAClients[client].Connect()
		if err != nil {
			log.Log.Errorln("Whatsapp gagal connect :", err)

			wr <- &WAResponse{
				Message: "Whatsapp gagal connect : " + err.Error(),
				Status:  false,
			}
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stderr)
				log.Log.Infoln("Whatsapp qr code")

				wr <- &WAResponse{
					Message: "Whatsapp berhasil berjalan dan menampilkan qr",
					Qr:      evt.Code,
					Status:  true,
				}
			} else {
				log.Log.Infoln("Whatsapp Login event:", evt.Event)
				wr <- &WAResponse{
					Message: "Whatsapp berhasil Login",
					Status:  true,
				}

			}
		}
	} else {
		err := WAClients[client].Connect()
		if err != nil {
			log.Log.Errorln("Whatsapp gagal connect", err)
			wr <- &WAResponse{
				Message: "Whatsapp gagal connect ",
				Status:  false,
			}
		}
		log.Log.Infoln("Whatsapp Berhasil dijalankan")
	}

	wr <- &WAResponse{
		Message: "Whatsapp Berhasil  berjalan ",
		Status:  true,
	}

	presenceAfterConnect(client)
}

func presenceAfterConnect(client string) {
	time.Sleep(time.Second * 5)

	if err := WAClients[client].SendPresence(types.PresenceAvailable); err != nil {
		log.Log.Errorln("Whatsapp gagal presensi : ", err)
	} else {
		log.Log.Infof("Presensi Berhasil")
	}
}

func StopWa(client string) error {

	if WAClients[client].IsLoggedIn() {
		if err := WAClients[client].SendPresence(types.PresenceUnavailable); err != nil {
			log.Log.Errorln("gagal logout : ", err)
			return err
		}
		time.Sleep(time.Second * 2)
	}
	log.Log.Infoln("Whatsapp Disconnected")
	WAClients[client].Disconnect()
	return nil
}

func LogoutWa(client string) error {

	if WAClients[client].IsLoggedIn() {
		if err := WAClients[client].Logout(); err != nil {
			log.Log.Errorln("gagal logout : ", err)
			return err
		}
		time.Sleep(2 * time.Second)
		WAClients[client].Disconnect()
	}
	return nil
}

func ParseJid(arg string) (types.JID, error) {
	if arg[0] == '0' {
		arg = "62" + arg[1:]
	}

	if arg[0] == '+' {
		arg = arg[1:]
	}

	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), nil
	}

	recipent, err := types.ParseJID(arg)

	return recipent, err
}

func ValidJid(client string, arg string) bool {
	recipent, err := ParseJid(arg)
	if err == nil {
		resp, err := WAClients[client].IsOnWhatsApp([]string{recipent.User})
		if err != nil {
			log.Log.Errorln("Gagal validasi jid", err)
			return false
		}
		return resp[0].IsIn
	}
	return true
}

func SendMessage(client string, dest types.JID, text string) error {
	if WAClients[client] == nil {
		return errors.New(client + "tidak terinstall")
	}

	msg := &waProto.Message{
		Conversation: proto.String(text),
	}

	WAClients[client].SendChatPresence(dest, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	time.Sleep(time.Second * 8)
	WAClients[client].SendChatPresence(dest, types.ChatPresencePaused, types.ChatPresenceMediaText)

	_, err := WAClients[client].SendMessage(context.Background(), dest, msg)

	if err != nil {
		return errors.New("gagal kirim pesan ke " + dest.User + err.Error())
	}

	return nil
}

func StartWaDummy(wg *sync.WaitGroup) {
	wr := make(chan *WAResponse, 1)

	if WAClients["dummy"] == nil {
		AppendDevice("dummy", "", true)
		go StartWa("dummy", wr)
	} else {
		go StartWa("dummy", wr)
	}

	<-wr

	wg.Done()

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGINT)
	// <-c
}
