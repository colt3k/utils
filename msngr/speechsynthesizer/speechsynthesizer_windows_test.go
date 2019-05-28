package speechsynthesizer

import (
	"testing"

	"github.com/colt3k/utils/config"
)

func TestNotification_Send(t *testing.T) {
	c := config.NewConfig()
	c.Load("../.env")

	n := Notification{
		Text:  "",
		Rate:  0,
		Voice: c.Util.GetString("voice"),
	}
	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
