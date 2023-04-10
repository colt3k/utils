package ss

import (
	"github.com/gorilla/mux"
	"net/http"
	"testing"
)

var (
	// AUTHTOKEN variable to hold build time authorization token prefix
	AUTHTOKEN = "myauthtokenplaceholder"
)

func AUTHVAL() string {
	return AUTHTOKEN
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// do whatever to in your main handler
}

func TestBasicAuth(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", BasicAuth(indexHandler, AUTHVAL(), "session_token", "johndoe", 3600, DynamicAuth()))
}
