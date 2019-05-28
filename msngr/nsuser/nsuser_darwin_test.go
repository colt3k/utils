package nsuser

import (
	"testing"

	"github.com/colt3k/utils/config"
)

func TestNotification_Send(t *testing.T) {
	c := config.NewConfig()
	c.Load("../.env")

	n := &Notification{
		Title:           "title",
		InformativeText: "message",
		SoundName:       c.Util.GetString("soundName"),
	}

	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
