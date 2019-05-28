package slack

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/colt3k/utils/config"
)

func TestNotification_Send(t *testing.T) {
	c := config.NewConfig()
	c.Load("../.env")

	n := Notification{
		Text:     "hi",
		Token:    c.Util.GetString("token"),
		Channel:  "#apps",
		Username: c.Util.GetString("username"),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 3 * time.Second,
		}}
	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
