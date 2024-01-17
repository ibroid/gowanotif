package events

import (
	"gowhatsapp/whatsapp"
)

func SendNotifMessage(client string, number string, message string) error {

	whatsapp.SendMessageViaWAHA("default", number, message)

	return nil
}
