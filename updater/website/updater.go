package website

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/version"

	"github.com/colt3k/utils/updater"
)

const (
	defaultServerURI = "http://localhost:1234/apps/"
)

/*

{
	"os":"darwin",
    "name": "application_name",
    "timestamp": "1234567",
	"hash":"2345678ytgfghjk"
	"version":"v0.0.1"
}
*/

var ac *updater.AppConfig

func CheckUpdate(appName string) bool {

	var updateAvailable bool

	log.Logln(log.DEBUG, "checking for update..")
	var compressdSuffix string
	var resp *http.Response
	var err error
	url := defaultServerURI + appName + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".update"
	compressdSuffix = "-" + runtime.GOOS + "-" + runtime.GOARCH + ".tgz"
	log.Logln(log.DEBUG, "Update Check URL", url)

	resp, err = http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if bserr.WarnErr(err, "update site unreachable") {
		return updateAvailable
	}

	// Read body to buffer
	body, err := ioutil.ReadAll(resp.Body)
	if bserr.Err(err, "Error reading body") {
		return updateAvailable
	}

	// Because in go lang if you read the body then any subsequent calls
	// are unable to read the body again....
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	//log.Logln(log.INFO, "response body: ", resp.Body)
	ac = new(updater.AppConfig)

	dec := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
	if err := dec.Decode(&ac); bserr.Err(err, "error decoding") {
		return updateAvailable
	}

	curVer, err := semver.Make(strings.TrimPrefix(version.VERSION, "v"))
	xVer, err := semver.Make(strings.TrimPrefix(ac.Version, "v"))

	log.Logf(log.DEBUG, "Current Version: %s, Remote Version: %s", curVer.String(), xVer.String())
	remoteTime := time.Unix(ac.Timestamp, 0)
	unx, _ := strconv.ParseInt(version.BUILDDATE, 10, 64)
	localTime := time.Unix(unx, 0)
	log.Logf(log.DEBUG, "LOCAL  App Name: %s, OS: %s/%s, Version: %s, App Time: %v, Converted: %v", appName, runtime.GOOS, runtime.GOARCH, version.VERSION, version.BUILDDATE, localTime)
	log.Logf(log.DEBUG, "REMOTE App Name: %s, OS: %s/%s, Version: %s, App Time: %v, Converted: %v", appName, ac.OS, ac.Arch, ac.Version, ac.Timestamp, remoteTime)

	if xVer.GT(curVer) {
		ac.URL = defaultServerURI + appName + compressdSuffix
		updateAvailable = true
		return updateAvailable

	} else if localTime.Before(remoteTime) { // if current app is older than remote pull, could be a roll back
		//Check build time instead
		ac.URL = defaultServerURI + appName + compressdSuffix
		if ac.Name == appName && ac.OS == runtime.GOOS {
			log.Logln(log.DEBUG, "update available!")
			updateAvailable = true
			return updateAvailable
		}
	}
	return updateAvailable
}

func UpdateAvailableMsg() string {
	var buf bytes.Buffer
	tm := time.Unix(ac.Timestamp, 0)
	buf.WriteString("  ******************************************************************************************************\n\n")
	buf.WriteString(fmt.Sprintf("\tNEW! Update available, Version: %s, %v \n", ac.Version, tm))
	buf.WriteString(fmt.Sprintf("\tDownload Here: %s\n", ac.URL))
	buf.WriteString("\n  ******************************************************************************************************\n")

	return buf.String()
}
