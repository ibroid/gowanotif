package test

import (
	"context"
	"errors"
	"fmt"
	"gowhatsapp/events"
	"gowhatsapp/log"
	"gowhatsapp/whatsapp"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
	log.LogInit()
}

func TestStartSippEventCanal(t *testing.T) {
	// events.StartSippEventCanal()
}

func TestNoRelaasToHakimEvent(t *testing.T) {

}

// go test -run TestNotifPMH gowhatsapp/test -v
func TestNotifPMH(t *testing.T) {
	godotenv.Load("../.env")
	// fmt.Println(os.Getenv("DEV"))
	var wg sync.WaitGroup

	wg.Add(1)
	go whatsapp.StartWaDummy(&wg)
	wg.Wait()
	time.Sleep(5 * time.Second)

	events.NotifPMH([]interface{}{"101474", "26432", "10", "2023-09-25", "", "1", "1", "Hakim Ketua", "109", "C.4", "196712251994031005", "Drs. Sarnoto, MH.", nil, "Y", nil, "fauzi", "2023-09-25 09:02:54", nil, nil})
}

// go test -run TestNotifPpp gowhatsapp/test -v
func TestNotifPpp(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go whatsapp.StartWaDummy(&wg)
	wg.Wait()
	time.Sleep(5 * time.Second)

	events.NotifPPP([]interface{}{"29787", "26838", "10", "2023-11-09", nil, "1", "117", "yos", "198603082015031002", "Yosie Ahmad Diantoro,�S.H.", nil, "Y", nil, "", "fauzi", "2023-11-09 11:45:08", nil, nil})

	defer whatsapp.StopWa("dummy")
}

// go test -run TestNotifPjs gowhatsapp/test -v
func TestNotifPjs(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go whatsapp.StartWaDummy(&wg)
	wg.Wait()
	time.Sleep(5 * time.Second)

	events.NotifPJS([]interface{}{"29528", "26839", "10", "2023-11-09", "", "1", "12", "", "196801052003121001", "Syamsuddin", nil, "Y", nil, "", "fauzi", "2023-11-09 11:49:25", "", nil})

	defer whatsapp.StopWa("dummy")
}

func TestNotifPhs(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go whatsapp.StartWaDummy(&wg)
	wg.Wait()
	time.Sleep(5 * time.Second)

	events.NotifPHS([]interface{}{"76187", "26778", "T", "T", "T", "1", "2023-11-08", "09:00:00", nil, nil, "SIDANG PERTAMA", "1", "Umar Bin Khatab", "T", nil, "T", nil, 0, "Y", nil, nil, nil, nil, "C10", "2023-11-02 15:51:21", "C10", "2023-11-02 15:54:42"})

	defer whatsapp.StopWa("dummy")
}

func TestNotifPTS(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(1)
	go whatsapp.StartWaDummy(&wg)
	wg.Wait()
	time.Sleep(5 * time.Second)

	events.NotifPHS([]interface{}{
		"76117",
		"26697",
		"T",
		"T",
		"T",
		"2",
		"2023-11-07",
		"09:50:00",
		"11:00:00",
		nil,
		"Memanggil Tergugat",
		"2",
		"Abu Musa Al Asyari",
		"T",
		"2",
		"T",
		"0",
		"0",
		"Y",
		"Putus Verstek",
		nil,
		nil,
		nil,
		"ario",
		"2023-10-31 01:12:05",
		"ario",
		"2023-11-07 09:59:59"})

	defer whatsapp.StopWa("dummy")

}

func TestStatusKehadiran(t *testing.T) {
	fmt.Println(events.StatusPihak(1))
}

func TestNotifSPV(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go whatsapp.StartWaDummy(&wg)
	wg.Wait()
	time.Sleep(5 * time.Second)

	events.NotifSPV([]interface{}{
		"25936",
		"2023-09-11",
		"Y",
		"10,8,11,16",
		"62",
		nil,
		nil,
		nil,
		nil,
		nil,
		"resources/file/doc/2023/09/PA.JU_2023_Pdt.G_1915_putusan_akhir.docx",
		"resources/file/anonimisasi/2023/09/PA.JU_2023_Pdt.G_1915_putusan_anonimisasi.docx",
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		"2023-09-11",
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		[]interface{}{},
		nil,
		nil,
		nil,
		nil,
		"dudi",
		"2023-09-11 02:30:43",
		"panmudgugatan2020",
		"2023-11-09 01:09:09"})

	defer whatsapp.StopWa("dummy")
}

func TestNotifPRS(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wr := make(chan *whatsapp.WAResponse, 1)

		if whatsapp.WAClients["dummy"] == nil {
			whatsapp.AppendDevice("dummy", "", true)
			go whatsapp.StartWa("dummy", wr)
		} else {
			go whatsapp.StartWa("dummy", wr)
		}

		<-wr

		time.Sleep(5 * time.Second)

		events.NotifPRSH()
		wg.Done()
	}()

	wg.Wait()

	// time.Sleep(30 * time.Second)

	defer whatsapp.StopWa("dummy")
}

func TestNotifPAC(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wr := make(chan *whatsapp.WAResponse, 1)

		if whatsapp.WAClients["dummy"] == nil {
			whatsapp.AppendDevice("dummy", "", true)
			go whatsapp.StartWa("dummy", wr)
		} else {
			go whatsapp.StartWa("dummy", wr)
		}

		<-wr

		time.Sleep(5 * time.Second)

		events.NotifPAC([]interface{}{"25418", "2023", "1626", "1626/AC/2023/PA.JU", "2023-08-28", "J 023719", "Cerai Gugat", "9", 1, 1, 1, nil, nil, "resources/file/doc/2023/06/PAJU_2023_PdtG_1462_aktacerai_1693282515.pdf", 1, nil, nil, "C10", "2023-08-10 01:34:17", "meja3gugatan", "2023-08-28 10:37:14"})
		wg.Done()
	}()

	wg.Wait()

	// time.Sleep(30 * time.Second)

	defer whatsapp.StopWa("dummy")
}

func TestNotifPUR(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wr := make(chan *whatsapp.WAResponse, 1)

		if whatsapp.WAClients["dummy"] == nil {
			whatsapp.AppendDevice("dummy", "", true)
			go whatsapp.StartWa("dummy", wr)
		} else {
			go whatsapp.StartWa("dummy", wr)
		}

		<-wr

		time.Sleep(5 * time.Second)

		events.NotifPUR([]interface{}{"41079", "23645", "76159", "105", "65986", "", "2023-11-09", nil, "T", "", "resources/file/doc/2022/12/PAJU_2022_PdtG_3177_relaas_23645_76159_65986_1700002480.pdf", "0", "indahf", "2023-11-15 05:54:26", nil, nil})
		wg.Done()
	}()

	wg.Wait()

	// time.Sleep(30 * time.Second)

	defer whatsapp.StopWa("dummy")
}

func TestNotifPPD(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wr := make(chan *whatsapp.WAResponse, 1)

		if whatsapp.WAClients["dummy"] == nil {
			whatsapp.AppendDevice("dummy", "", true)
			go whatsapp.StartWa("dummy", wr)
		} else {
			go whatsapp.StartWa("dummy", wr)
		}

		<-wr

		time.Sleep(5 * time.Second)

		events.NotifPPD([]interface{}{"31294", "26837", "1", "77104", "1", "HERI RAHARJO BIN A KAMDJI", "Jalan Sungai Kendal No. 66 B, RT.004 RW.008, Kelurahan Rorotan, Kecamatan Cilincing, Kota Jakarta Utara", "", nil, nil, "fauzi 2023-11-09 11:41:32", nil, nil})
		wg.Done()
	}()

	wg.Wait()

	// time.Sleep(30 * time.Second)

	defer whatsapp.StopWa("dummy")
}

func TestNotifWithContext(t *testing.T) {

	err := errors.New("Sample error")
	ctx, cancel := context.WithCancel(context.TODO())

	defer cancel()
	if err != nil {

		fmt.Println("Was cancel")
		cancel()
	}

	fmt.Println("sampe sini")
	ctx.Done()
}
