package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/colt3k/utils/netut"

	"github.com/blang/semver"
	"github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/io"
	"github.com/colt3k/utils/netut/hc"
	"github.com/colt3k/utils/osut"
	"github.com/colt3k/utils/ques"
	"github.com/colt3k/utils/updater"
)

var ac *updater.AppConfig

func CheckUpdate(appName string, hosts []updater.Connection, version updater.Version) (*updater.AppConfig, bool) {

	testHosts(hosts)
	for _, d := range hosts {

		// if all checks failed skip
		if !d.Available() && !d.HostPfx() && !d.HostSuffix() {
			continue
		}
		var base bytes.Buffer
		var upUrl bytes.Buffer

		user := []byte(d.User)
		pass := []byte(d.PassOrToken)

		base.WriteString(d.URLPrefix)
		base.WriteString(d.Repository)
		base.WriteString(d.Path)

		upUrl.WriteString(base.String())
		upUrl.WriteString(appName + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".update")
		log.Logln(log.DEBUG, "URL:", upUrl.String())

		var updateAvailable bool
		var compressdSuffix string
		url := upUrl.String()
		compressdSuffix = "-" + runtime.GOOS + "-" + runtime.GOARCH + ".tgz"

		auth := &hc.Auth{Username: user, Password: pass}
		data, err := pullURLToString(url, auth)
		if bserr.WarnErr(err) {
			log.Logf(log.WARN, "update site unreachable %v", err.Error())
			//return nil, updateAvailable
			continue
		}

		ac = new(updater.AppConfig)
		ac.BaseURL = base.String()
		ac.User = user
		ac.Pass = pass

		dec := json.NewDecoder(ioutil.NopCloser(strings.NewReader(data)))
		if err := dec.Decode(&ac); bserr.Err(err, "error decoding") {
			return ac, updateAvailable
		}

		curVer, err := semver.Make(strings.TrimPrefix(version.Version, "v"))
		xVer, err := semver.Make(strings.TrimPrefix(ac.Version, "v"))

		log.Logf(log.DEBUG, "Current Version: %s, Remote Version: %s", curVer.String(), xVer.String())
		remoteTime := time.Unix(ac.Timestamp, 0)
		unx, _ := strconv.ParseInt(version.BuildDate, 10, 64)
		localTime := time.Unix(unx, 0)
		log.Logf(log.DEBUG, "LOCAL  App Name: %s, OS: %s/%s, Version: %s, App Time: %v, Converted: %v", appName, runtime.GOOS, runtime.GOARCH, version.Version, version.BuildDate, localTime)
		log.Logf(log.DEBUG, "REMOTE App Name: %s, OS: %s/%s, Version: %s, App Time: %v, Converted: %v", appName, ac.OS, ac.Arch, ac.Version, ac.Timestamp, remoteTime)

		if xVer.GT(curVer) {
			log.Logln(log.DEBUG, "remote version is newer")
			base.WriteString(appName + compressdSuffix)
			ac.URL = base.String()
			ac.ArchiveName = appName + compressdSuffix
			updateAvailable = true
			return ac, updateAvailable

		} else if localTime.Before(remoteTime) { // if current app is older than remote pull, could be a roll back
			//Check build time instead
			log.Logln(log.DEBUG, "remote time is newer than local")
			base.WriteString(appName + compressdSuffix)
			ac.URL = base.String()
			ac.ArchiveName = appName + compressdSuffix
			if ac.Name == appName && ac.OS == runtime.GOOS {
				log.Logln(log.DEBUG, "update available!")
				updateAvailable = true
				return ac, updateAvailable
			}
		} else {
			log.Logln(log.DEBUG, "local version is newer")
		}
		return ac, updateAvailable
	}

	return nil, false
}

func UpdateAvailableMsg() string {
	var buf bytes.Buffer
	tm := time.Unix(ac.Timestamp, 0)
	buf.WriteString("  ******************************************************************************************************\n\n")
	buf.WriteString(fmt.Sprintf("\tNEW! Update available, Version: %s, %v \n", ac.Version, tm))
	buf.WriteString(fmt.Sprintf("\tDownload Here: %s\n", ac.URL))
	chgs := pullChangeLogAndDisplay(ac)
	if len(chgs) > 0 {
		buf.WriteString(fmt.Sprintf("\tChanges:\n"))
		buf.WriteString(fmt.Sprintf("%s\n", chgs))
	}

	buf.WriteString("\n  ******************************************************************************************************\n")

	return buf.String()
}

func PerformUpdate(appName string, hosts []updater.Connection, version updater.Version, question bool) bool {

	/*
		1. Pull file from archive
			myappname-darwin-amd64/myappname
		2. Place in proper location (same as normal install)
		3. exit application and notify to restart - due to update or make option via input
	*/
	log.Logln(log.DEBUG, "check for update")
	if ac, found := CheckUpdate(appName, hosts, version); found {
		s := UpdateAvailableMsg()
		fmt.Println(s)
		if question && ques.Confirm("\nPerform Update ? ") {
			//Download
			if download(ac) {
				// success
				log.Println("\n** successful download, exiting so you can restart the application **")
				os.Exit(0)
			} else {
				// failed
				log.DisableTimestamp()
				log.Println("\nupdate failed")
				log.EnableTimestamp()
			}
		} else {
			return found
		}
	}
	return false
}

func pullChangeLogAndDisplay(ac *updater.AppConfig) string {
	if ac != nil {
		log.Println("Chg: ", ac.Changelog)
		log.Println("URL:", ac.BaseURL)
		if len(ac.Changelog) > 0 {
			httpClient := hc.NewClient(hc.HttpClientRequestTimeout(20))
			var err error
			url := ac.BaseURL + "/" + ac.Changelog
			auth := &hc.Auth{Username: ac.User, Password: ac.Pass}
			resp, err := httpClient.Fetch("GET", url, auth, nil, nil)
			if resp != nil {
				defer resp.Body.Close()
			}
			if bserr.WarnErr(err) {
				log.Logf(log.WARN, "update site unreachable %v", err.Error())
			}

			// Read body to buffer
			body, err := ioutil.ReadAll(resp.Body)
			if bserr.Err(err, "Error reading body") {
				log.Logln(log.WARN, "error reading body on changelog")
			}
			var byt bytes.Buffer
			lines := strings.SplitAfter(string(body), "\n")
			for _, j := range lines {
				byt.WriteString("\t\t")
				byt.WriteString(j)
			}
			return byt.String()
		}
	}
	return ""
}

func download(ac *updater.AppConfig) bool {

	var success bool
	httpClient := hc.NewClient(hc.HttpClientRequestTimeout(20))

	auth := &hc.Auth{Username: ac.User, Password: ac.Pass}
	resp, err := httpClient.Fetch("GET", ac.URL, auth, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if bserr.WarnErr(err) {
		log.Logf(log.WARN, "update site unreachable %v", err.Error())
		return success
	}

	// Read body to buffer
	body, err := ioutil.ReadAll(resp.Body)
	if bserr.Err(err, "Error reading body") {
		return success
	}

	// write out
	log.Logln(log.DEBUG, "Writing out update to ", ac.ArchiveName)
	io.WriteOut(body, ac.ArchiveName)

	// Pull executable from archive 'tgz'

	cmd := "tar xvf " + ac.ArchiveName + " " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "/" + ac.Name
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("Failed to execute command: %s %s", cmd, err.Error()))
		return success
	}

	// Move into location
	// get executable path and replace original
	s, err := os.Executable()
	if err != nil {
		return success
	}
	if strings.HasSuffix(s, "main") {
		log.Logln(log.WARN, "not a packaged executable")
		return success
	}
	cmd = "mv " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "/" + ac.Name + " " + s
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("Failed to execute command: %s %s", cmd, err.Error()))
		return success
	}

	// Clean up
	cmd = "rm -rf " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "*"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("Failed to execute command: %s %s", cmd, err.Error()))
		return success
	}

	success = true
	return success
}

func testHosts(hosts []updater.Connection) {
	log.Logln(log.DEBUG, "test hosts")
	host := osut.Hostname()
	log.Logln(log.DEBUG, "local host", host)
	for i, d := range hosts {
		// if neither is set skip entry
		if len(d.OnAvailable) > 0 {
			// test available
			var avail bool
			var err error
			if d.OnAvailableViaHTTP {
				avail, err = hc.Reachable(d.OnAvailable, d.Name, 2, d.DisableValidateCert)
			} else {
				avail, err = netut.Ping(d.OnAvailable)
			}
			if err != nil {
				log.Logf(log.DEBUG, err.Error())
			}

			log.Logf(log.DEBUG, "host available? %v: %s", avail, d.OnAvailable)
			hosts[i].SetAvailable(avail)
			//if avail {
			//	break
			//}
		}
		// if OnAvailable is set and OnHostname is NOT set but host is not resolvable skip
		if len(d.OnHostNamePrefix) > 0 && strings.HasPrefix(strings.ToLower(host), d.OnHostNamePrefix) {
			log.Logf(log.DEBUG, "on host starting with: %s hostname %s", d.OnHostNamePrefix, host)
			hosts[i].SetHostPfx(true)
		}
		if len(d.OnHostNameSuffix) > 0 && strings.HasSuffix(strings.ToLower(host), d.OnHostNameSuffix) {
			log.Logf(log.DEBUG, "on host ending with: %s hostname %s", d.OnHostNameSuffix, host)
			hosts[i].SetHostSfx(true)
		}
	}
}

func pullURLToString(url_ string, auth *hc.Auth) (string, error) {
	log.Logln(log.DEBUG, "obtain httpclient")
	httpClient := hc.NewClient(hc.HttpClientRequestTimeout(20))

	resp, err := httpClient.Fetch("GET", url_, auth, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url_, resp.Status)
	}
	slurp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url_, err)
	}
	return string(slurp), nil
}
