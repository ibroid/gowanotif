package test

import (
	"gowhatsapp/events"
	"testing"
)

func TestScheduleSipp(t *testing.T) {
	events.StartCronEvent()
}
