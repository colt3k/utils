package nsuser

import (
	"testing"
)

func TestNotification_Send(t *testing.T) {
	//c := config.NewConfig()
	//c.Load("../.env")

	n := &Notification{
		Title:           "title",
		InformativeText: "message",
		//SoundName:       c.Util.GetString("soundName"),
	}

	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
