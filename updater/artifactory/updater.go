package artifactory

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	iout "github.com/colt3k/utils/io"
	"github.com/colt3k/utils/netut/hc"
	"github.com/colt3k/utils/osut"
	"github.com/colt3k/utils/ques"
	"github.com/colt3k/utils/updater"
)

var ac *updater.AppConfig

func CheckUpdate(appName string, hosts []updater.Connection, version updater.Version) (*updater.AppConfig, bool, bool) {
	test := false
	testHosts(hosts)
	log.Logf(log.DEBUG, "- CheckUpdate")
	for _, d := range hosts {
		var autoUpdate bool

		// if all checks failed skip
		if !d.Available() && !d.HostPfx() && !d.HostSuffix() {
			continue
		}
		log.Logf(log.DEBUG, "-- attempt host %v", d.OnAvailable)
		var base bytes.Buffer
		var upUrl bytes.Buffer
		var autoUrl bytes.Buffer

		user := []byte(d.User)
		pass := []byte(d.PassOrToken)

		base.WriteString(d.URLPrefix)
		base.WriteString(d.Repository)
		base.WriteString(d.Path)

		upUrl.WriteString(base.String())
		if !test {
			upUrl.WriteString(appName + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".update")
		} else {
			upUrl.WriteString(appName + "-linux-" + runtime.GOARCH + ".update")
		}
		log.Logf(log.DEBUG, "-- Update File URL: %v", upUrl.String())

		var updateAvailable bool
		var compressdSuffix string
		url := upUrl.String()
		if !test {
			compressdSuffix = "-" + runtime.GOOS + "-" + runtime.GOARCH + ".tgz"
		} else {
			compressdSuffix = "-linux-" + runtime.GOARCH + ".tgz"
		}

		auth := &hc.Auth{Username: user, Password: pass}
		data, err := pullURLToString(url, auth, d.DisableValidateCert)
		if err != nil {
			log.Logf(log.WARN, "--- %v", err.Error())
			//return nil, updateAvailable
			continue
		}

		ac = new(updater.AppConfig)
		ac.BaseURL = base.String()
		ac.User = user
		ac.Pass = pass
		ac.DisableVerifyCert = d.DisableValidateCert

		dec := json.NewDecoder(ioutil.NopCloser(strings.NewReader(data)))
		if err = dec.Decode(&ac); bserr.Err(err, "error decoding") {
			return ac, updateAvailable, autoUpdate
		}

		// Check for Auto file and value
		autoUrl.WriteString(base.String())
		autoUrl.WriteString(appName + ".auto")
		log.Logf(log.DEBUG, "-- Auto File URL: %v", autoUrl.String())
		autoURI := autoUrl.String()
		autoDat, err := pullURLToString(autoURI, auth, d.DisableValidateCert)
		if err != nil {
			log.Logf(log.WARN, "--- %v", err.Error())
		}
		if len(autoDat) > 0 {
			var errAU error
			autoUpdate, errAU = strconv.ParseBool(autoDat)
			if errAU != nil {
				log.Logf(log.ERROR, "--- issue parsing auto update file %v", errAU)
			}
		}

		curVer, err := semver.Make(strings.TrimPrefix(version.Version, "v"))
		xVer, err := semver.Make(strings.TrimPrefix(ac.Version, "v"))

		log.Logf(log.DEBUG, "-- Current Version: %s, Remote Version: %s", curVer.String(), xVer.String())
		remoteTime := time.Unix(ac.Timestamp, 0)
		unx, _ := strconv.ParseInt(version.BuildDate, 10, 64)
		localTime := time.Unix(unx, 0)
		log.Logf(log.DEBUG, "-- LOCAL  App Name: %s, OS: %s/%s, Version: %s, App Time: %v, Converted: %v", appName, runtime.GOOS, runtime.GOARCH, version.Version, version.BuildDate, localTime)
		log.Logf(log.DEBUG, "-- REMOTE App Name: %s, OS: %s/%s, Version: %s, App Time: %v, Converted: %v", appName, ac.OS, ac.Arch, ac.Version, ac.Timestamp, remoteTime)

		if xVer.GT(curVer) {
			log.Logf(log.DEBUG, "-- remote version is newer %v > %v", xVer, curVer)
			base.WriteString(appName + compressdSuffix)
			ac.URL = base.String()
			ac.ArchiveName = appName + compressdSuffix
			updateAvailable = true
			return ac, updateAvailable, autoUpdate

		} else if localTime.Before(remoteTime) { // if current app is older than remote pull, could be a roll back
			//Check build time instead
			log.Logln(log.DEBUG, "-- remote time is newer than local %v > %v", remoteTime, localTime)
			base.WriteString(appName + compressdSuffix)
			ac.URL = base.String()
			ac.ArchiveName = appName + compressdSuffix
			if ac.Name == appName && ac.OS == runtime.GOOS {
				log.Logln(log.DEBUG, "-- update available!")
				updateAvailable = true
				return ac, updateAvailable, autoUpdate
			}
		} else {
			log.Logln(log.DEBUG, "-- local version is newer")
		}
		return ac, updateAvailable, autoUpdate
	}
	log.Logln(log.DEBUG, "- CheckUpdate END")
	return nil, false, false
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
	log.Logln(log.DEBUG, "- Perform Update")
	if ac, found, autoUpdate := CheckUpdate(appName, hosts, version); found {
		s := UpdateAvailableMsg()
		fmt.Println(s)
		if autoUpdate {
			downloadUpdate(ac)
		} else if !autoUpdate && question && ques.Confirm("\nPerform Update ? ") {
			downloadUpdate(ac)
		} else {
			return found
		}
	}
	log.Logln(log.DEBUG, "- Perform Update END")
	return false
}

func downloadUpdate(ac *updater.AppConfig) {
	log.Logln(log.DEBUG, "- Download Update")
	//Download
	if download(ac) {
		// success
		log.DisableTimestamp()
		log.Println("\n** successful download, exiting so you can restart the application **")
		log.EnableTimestamp()
		os.Exit(0)
	} else {
		// failed
		log.DisableTimestamp()
		log.Printf("\nupdate failed %v", ac.Issue)
		log.EnableTimestamp()
	}
	log.Logln(log.DEBUG, "- Download Update END")
}

func pullChangeLogAndDisplay(ac *updater.AppConfig) string {
	if ac != nil {
		log.Logln(log.DEBUG, "- Pull Change Log and Display")
		log.Logf(log.DEBUG,"-- change log: %v", ac.Changelog)
		log.Logf(log.DEBUG,"-- url used: %v", ac.BaseURL)
		if len(ac.Changelog) > 0 {
			httpClient := hc.NewClient(hc.HttpClientRequestTimeout(20), hc.DisableVerifyClientCert(ac.DisableVerifyCert))
			var err error
			url := ac.BaseURL + "/" + ac.Changelog
			auth := &hc.Auth{Username: ac.User, Password: ac.Pass}
			resp, err := httpClient.Fetch("GET", url, auth, nil, nil)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				log.Logf(log.WARN, "--- %v", err.Error())
			}

			// Read body to buffer
			body, err := ioutil.ReadAll(resp.Body)
			if bserr.Err(err, "Error reading body") {
				log.Logln(log.WARN, "--- error reading body on changelog")
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
	log.Logln(log.DEBUG, "- Pull Change Log and Display END")
	return ""
}

func download(ac *updater.AppConfig) bool {
	log.Logln(log.DEBUG, "- Download")
	var success bool
	httpClient := hc.NewClient(hc.HttpClientRequestTimeout(20), hc.DisableVerifyClientCert(ac.DisableVerifyCert))

	auth := &hc.Auth{Username: ac.User, Password: ac.Pass}
	resp, err := httpClient.Fetch("GET", ac.URL, auth, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Logf(log.ERROR, "--- %v", err.Error())
		ac.Issue = "update site unreachable"
		return success
	}

	// Read body to buffer
	body, err := ioutil.ReadAll(resp.Body)
	if bserr.Err(err, "--- error reading body") {
		ac.Issue = "unable to read response"
		return success
	}

	// write out
	log.Logf(log.DEBUG, "-- writing out update to %v", ac.ArchiveName)
	_, err = iout.WriteOut(body, ac.ArchiveName)
	if err != nil {
		log.Logf(log.ERROR, "--- issue writing out %v, %v", ac.ArchiveName, err)
		ac.Issue = "unable to write out archive"
		return success
	}

	// Pull executable from archive 'tgz'
	log.Logf(log.DEBUG, "-- extracting executable %v", ac.Name)
	cmd := "tar xvf " + ac.ArchiveName + " " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "/" + ac.Name
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		ac.Issue = "failed to extract"
		return success
	}

	// Validate HASH, pull file(s)
	archivePathDir := strings.TrimSuffix(ac.ArchiveName, ".tgz")
	sha256HashFileName := ac.Name + ".sha256"
	validHash := validateHash(sha256HashFileName, archivePathDir, ac.Name)
	if !validHash {
		log.Logln(log.INFO, "--- sha256 hash invalid")
		ac.Issue = "invalid sha256 hash"
		return success
	}

	// Move into location
	// get executable path and replace original
	s, err := os.Executable()
	if err != nil {
		return success
	}
	if strings.HasSuffix(s, "main") {
		log.Logln(log.WARN, "--- not a packaged executable")
		ac.Issue = "not a packaged executable"
		return success
	}
	cmd = "mv " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "/" + ac.Name + " " + s
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		ac.Issue = "failed to move/replace application"
		return success
	}

	// Clean up
	cmd = "rm -rf " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "*"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		ac.Issue = "failed to clean archive"
		return success
	}

	success = true
	log.Logln(log.DEBUG, "- Download END")
	return success
}

func testHosts(hosts []updater.Connection) {
	log.Logln(log.DEBUG, "- test hosts")
	host := osut.Hostname()
	log.Logf(log.DEBUG, "-- local host %v", host)
	for i, d := range hosts {
		// if neither is set skip entry
		if len(d.OnAvailable) > 0 {
			// test available
			var avail bool
			var err error
			if d.OnAvailableViaHTTP {
				log.Logf(log.DEBUG, "-- onAvailable check %v", d.OnAvailable)
				avail, err = hc.Reachable(d.OnAvailable, d.Name, 2, d.DisableValidateCert)
			} else {
				log.Logf(log.DEBUG, "-- ping host check %v", d.OnAvailable)
				avail, err = netut.Ping(d.OnAvailable)
			}
			if err != nil {
				log.Logf(log.DEBUG, "--- %v",err.Error())
			}

			log.Logf(log.DEBUG, "--- host available? %v: %s", avail, d.OnAvailable)
			hosts[i].SetAvailable(avail)
			//if avail {
			//	break
			//}
		}
		// if OnAvailable is set and OnHostname is NOT set but host is not resolvable skip
		if len(d.OnHostNamePrefix) > 0 && strings.HasPrefix(strings.ToLower(host), d.OnHostNamePrefix) {
			log.Logf(log.DEBUG, "--- on host starting with: %s hostname %s", d.OnHostNamePrefix, host)
			hosts[i].SetHostPfx(true)
		}
		if len(d.OnHostNameSuffix) > 0 && strings.HasSuffix(strings.ToLower(host), d.OnHostNameSuffix) {
			log.Logf(log.DEBUG, "--- on host ending with: %s hostname %s", d.OnHostNameSuffix, host)
			hosts[i].SetHostSfx(true)
		}
	}
}

func pullURLToString(url_ string, auth *hc.Auth, disableVerifyCert bool) (string, error) {
	log.Logln(log.DEBUG, "--- Pull URL Content To String")
	httpClient := hc.NewClient(hc.HttpClientRequestTimeout(20), hc.DisableVerifyClientCert(disableVerifyCert))

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

func hash_file_md5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}
func validateHash(hashFileName, archivePathDir, appName string) bool {
	log.Logln(log.DEBUG, "- Validate Hash")
	log.Logf(log.DEBUG, "-- extracting hash file %v", hashFileName)
	cmd := "tar xvf " + ac.ArchiveName + " " + archivePathDir + "/" + hashFileName
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		return false
	}

	// Read file contents
	hashContent, err := ioutil.ReadFile(archivePathDir + "/" + hashFileName)
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to read hash: %s", err.Error()))
		return false
	}
	// create hash from downloaded application
	var hash string
	if strings.HasSuffix(hashFileName, ".md5") {
		hash, err = hash_file_md5(archivePathDir + "/" + appName)
	} else if strings.HasSuffix(hashFileName, ".sha256") {
		hash, err = hash_file_sha256(archivePathDir + "/" + appName)
	}
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to create hash from application: %s", err.Error()))
		return false
	}

	if string(hashContent) == hash {
		log.Logln(log.INFO, "-- VALID hash")
		return true
	} else {
		log.Logf(log.INFO, "-- INVALID hash %s", hash)
	}
	log.Logln(log.DEBUG, "- Validate Hash END")
	return false
}
func hash_file_sha256(filePath string) (string, error) {
	var returnSHA256String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA256String, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA256String, err
	}
	hashInBytes := hash.Sum(nil)[:32]
	returnSHA256String = hex.EncodeToString(hashInBytes)
	return returnSHA256String, nil

}
