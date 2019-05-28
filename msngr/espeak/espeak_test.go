package espeak

import (
	"fmt"
	"testing"
)

func TestNotification_Send(t *testing.T) {

	title := "title"
	message := "body message"
	n := &Notification{
		Text:      fmt.Sprintf("%s %s", title, message),
		VoiceName: "english-us",
	}

	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
