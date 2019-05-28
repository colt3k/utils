package say

import (
	"fmt"
	"testing"

	"github.com/colt3k/utils/config"
)

func TestSayNotification_Send(t *testing.T) {

	c := config.NewConfig()
	c.Load("../.env")

	n := &Notification{
		Voice: c.Util.GetString("voice"),
		Text:  fmt.Sprintf("%s %s", "title", "message"),
		Rate:  200,
	}
	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
