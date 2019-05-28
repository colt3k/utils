package pushbullet

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/colt3k/utils/config"
)

func TestNotification_Send(t *testing.T) {

	n := Notification{
		Title:       "title",
		Body:        "mesg",
		AccessToken: "token",
		Client:      &http.Client{Timeout: 3 * time.Second},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.Header.Get("Access-Token") == "" {
			t.Error("missing access token")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("content type should be application/json")
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("missing payload")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	PB_API_URL = ts.URL
	mockResp.ErrorCode = "" // success
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.ErrorCode = "error" // failure
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}

func TestSend(t *testing.T) {
	c := config.NewConfig()
	c.Load("../.env")

	n := Notification{
		Type:        "note",
		Title:       "my title2",
		Body:        "my mesg2",
		AccessToken: c.Util.GetString("accesstoken"),
		DeviceIden:  c.Util.GetString("deviceid"),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 3 * time.Second,
		},
	}

	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
