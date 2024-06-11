package artifactory

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/colt3k/utils/debug"
	"github.com/colt3k/utils/file"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/colt3k/utils/netut"

	"github.com/blang/semver"
	log "github.com/colt3k/nglog/ng"
	iout "github.com/colt3k/utils/io"
	"github.com/colt3k/utils/netut/hc"
	"github.com/colt3k/utils/osut"
	"github.com/colt3k/utils/ques"
	"github.com/colt3k/utils/updater"
)

func CheckUpdate(appName string, hosts []updater.Connection, version updater.Version, checkOnly bool) (*updater.AppConfig, bool, bool) {
	var ac *updater.AppConfig
	test := false
	testHosts(hosts, checkOnly)
	if len(hosts) == 0 {
		log.Logln(log.DEBUG, "-- no configured hosts")
	}
	disabledHostCount := 0
	for _, d := range hosts {
		if !d.Available() && !d.HostPfx() && !d.HostSuffix() {
			disabledHostCount++
			continue
		}
		//if checkOnly {
		//	break
		//}
	}
	if disabledHostCount == len(hosts) {
		log.Logln(log.DEBUG, "-- no reachable hosts")
	}

	log.Logln(log.DEBUG, "")
	log.Logln(log.DEBUG, "- START CheckUpdate process")

	for _, d := range hosts {
		var autoUpdate bool

		// if all checks failed skip
		if !d.Available() && !d.HostPfx() && !d.HostSuffix() {
			continue
		}
		log.Logf(log.DEBUG, "-- attempt host %v", d.OnAvailable)
		var base bytes.Buffer
		var upURL bytes.Buffer
		var autoURL bytes.Buffer

		user := []byte(d.User)
		pass := []byte(d.PassOrToken)

		base.WriteString(d.URLPrefix)
		base.WriteString(d.Repository)
		base.WriteString(d.Path)

		upURL.WriteString(base.String())
		if !test {
			upURL.WriteString(appName + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".update")
		} else {
			upURL.WriteString(appName + "-linux-amd64.update")
		}
		log.Logf(log.DEBUG, "-- Update File URL: %v", upURL.String())

		var updateAvailable bool
		var compressdSuffix string
		url := upURL.String()
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

		dec := json.NewDecoder(io.NopCloser(strings.NewReader(data)))
		if err = dec.Decode(&ac); err != nil {
			log.Printf("error decoding %v\n", err)
			debug.PrintStack()
			return ac, updateAvailable, autoUpdate
		}

		// Check for Auto file and value, THIS IS OPTIONAL, ignore if not found
		autoURL.WriteString(base.String())
		autoURL.WriteString(appName + ".auto")
		log.Logf(log.DEBUG, "-- Auto File URL: %v", autoURL.String())
		autoURI := autoURL.String()
		autoDat, err := pullURLToString(autoURI, auth, d.DisableValidateCert)
		if err != nil {
			log.Logf(log.WARN, "--- %v", err.Error())
		}
		if len(autoDat) > 0 {
			var errAU error
			autoUpdate, errAU = strconv.ParseBool(autoDat)
			if errAU != nil {
				log.Logf(log.ERROR, "--- issue parsing auto update file %v", errAU)
				continue
			}
		}

		curVer, err := semver.Make(strings.TrimPrefix(version.Version, "v"))
		if err != nil {
			log.Logf(log.ERROR, "issue parsing version %v", err)
			continue
		}
		xVer, err := semver.Make(strings.TrimPrefix(ac.Version, "v"))
		if err != nil {
			log.Logf(log.ERROR, "issue parsing ac version %v", err)
			continue
		}

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
			log.Logf(log.DEBUG, "-- remote time is newer than local %v > %v", remoteTime, localTime)
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
	log.Logln(log.DEBUG, "- END CheckUpdate process")
	return nil, false, false
}

func UpdateAvailableMsg(ac *updater.AppConfig) string {
	var buf bytes.Buffer
	tm := time.Unix(ac.Timestamp, 0)
	buf.WriteString("  ******************************************************************************************************\n\n")
	buf.WriteString(fmt.Sprintf("\tNEW! Update available, Version: %s, %v \n", ac.Version, tm))
	buf.WriteString(fmt.Sprintf("\tDownload Here: %s\n", ac.URL))
	chgs := pullChangeLogAndDisplay(ac)
	if len(chgs) > 0 {
		buf.WriteString("\tChanges:\n")
		buf.WriteString(fmt.Sprintf("%s\n", chgs))
	}

	buf.WriteString("\n  ******************************************************************************************************\n")

	return buf.String()
}

// PerformUpdate return true(found update) if check only, return true when successfully updated
func PerformUpdate(appName string, hosts []updater.Connection, version updater.Version, question, checkOnly bool) bool {
	/*
		1. Pull file from archive
			myappname-darwin-amd64/myappname
		2. Place in proper location (same as normal install)
		3. exit application and notify to restart - due to update or make option via input
	*/
	log.Logln(log.DEBUG, "")
	log.Logln(log.DEBUG, "**** START Update process ****")
	if checkOnly {
		if _, found, _ := CheckUpdate(appName, hosts, version, checkOnly); found {
			log.Logf(log.WARN, "Update Available, use the application update call to obtain the latest version.")
			return found
		}
		return false
	}
	if appConfig, found, autoUpdate := CheckUpdate(appName, hosts, version, checkOnly); found {
		s := UpdateAvailableMsg(appConfig)
		fmt.Println(s)
		if autoUpdate {
			return downloadUpdate(appConfig)
		} else if !autoUpdate && question && ques.Confirm("\nPerform Update ? ") {
			return downloadUpdate(appConfig)
		} else {
			return found
		}
	}
	log.Logln(log.INFO, "**** END Update process ****")
	return false
}

// return true on success
func downloadUpdate(ac *updater.AppConfig) bool {
	log.Logln(log.DEBUG, "- Download Update")
	//Download
	if download(ac) {
		// success
		log.DisableTimestamp()
		log.Println("\n** successful download, exiting so you can restart the application **")
		log.EnableTimestamp()
		switch runtime.GOOS {
		case "windows":
			return true
		default: //Mac & Linux
			//os.Exit(0)
			return true
		}
	}
	// failed
	log.DisableTimestamp()
	log.Printf("\nupdate failed %v", ac.Issue)
	log.EnableTimestamp()

	log.Logln(log.DEBUG, "- Download Update END")
	return false
}

func pullChangeLogAndDisplay(ac *updater.AppConfig) string {
	if ac != nil {
		log.Logln(log.DEBUG, "")
		log.Logln(log.DEBUG, "- Pull Change Log and Display")
		log.Logf(log.DEBUG, "-- change log: %v", ac.Changelog)
		log.Logf(log.DEBUG, "-- url used: %v", ac.BaseURL)
		if len(ac.Changelog) > 0 {
			httpClient := hc.NewClient(hc.HTTPClientRequestTimeout(30), hc.DisableVerifyClientCert(ac.DisableVerifyCert))
			var err error
			url := ac.BaseURL
			if ac.BaseURL[len(ac.BaseURL)-1:] != "/" {
				url += "/"
			}
			url += ac.Changelog
			auth := &hc.Auth{Username: ac.User, Password: ac.Pass}
			resp, err := httpClient.Fetch("GET", url, auth, nil, nil)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				log.Logf(log.WARN, "--- %v", err.Error())
			}

			// Read body to buffer
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				debug.PrintStack()
				log.Logf(log.ERROR, "Error reading body %v", err)
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
	httpClient := hc.NewClient(hc.HTTPClientRequestTimeout(120), hc.DisableVerifyClientCert(ac.DisableVerifyCert))

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		debug.PrintStack()
		log.Logf(log.ERROR, "--- error reading body %v", err)
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
	switch runtime.GOOS {
	case "windows":
		_, err = exec.Command("cmd", "/C", cmd).Output()
	default: //Mac & Linux
		_, err = exec.Command("sh", "-c", cmd).Output()
	}
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		ac.Issue = "failed to extract"
		return success
	}

	// Validate HASH, pull file(s)
	archivePathDir := strings.TrimSuffix(ac.ArchiveName, ".tgz")
	sha256HashFileName := ac.Name + ".sha256"
	validHash := validateHash(sha256HashFileName, archivePathDir, ac)
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
	var output []byte
	cmd = "mv " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "/" + ac.Name + " " + s
	switch runtime.GOOS {
	case "windows":
		// move, copy robocopy; can't replace file while running on windows it's locked; provide notice here instead
		cmd = "move /y " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + file.PathSeparator() + ac.Name + " " + filepath.Dir(s) + file.PathSeparator() + ac.Name + ".new"
		output, err = exec.Command("cmd", "/C", cmd).CombinedOutput()
		if err == nil {
			log.Logf(log.WARN, "MANUAL: Due to Windows locking the new version is placed here %v, remove the original file and rename this by removing .new", filepath.Dir(s)+file.PathSeparator()+ac.Name+".new")
		}
	default: //Mac & Linux
		_, err = exec.Command("sh", "-c", cmd).Output()
	}
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s\nOutput: |%s|", cmd, err.Error(), string(output)))
		ac.Issue = "failed to move/replace application"
		return success
	}

	// Clean up
	cmd = "rm -rf " + strings.TrimSuffix(ac.ArchiveName, ".tgz") + "*"
	switch runtime.GOOS {
	case "windows":
		cmd = "rmdir /s /q " + strings.TrimSuffix(ac.ArchiveName, ".tgz")
		_, err = exec.Command("cmd", "/C", cmd).Output()
	default: //Mac & Linux
		_, err = exec.Command("sh", "-c", cmd).Output()
	}
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		ac.Issue = "failed to clean archive"
		return success
	}

	success = true
	log.Logln(log.DEBUG, "- Download END")
	return success
}

func testHosts(hosts []updater.Connection, checkOnly bool) {
	log.Logln(log.DEBUG, "- START testHosts process")
	host := osut.Hostname()
	log.Logf(log.DEBUG, "-- local host %v", host)
	for i, d := range hosts {
		log.Logln(log.DEBUG, "")
		var avail bool
		// if neither is set skip entry
		if len(d.OnAvailable) > 0 {
			// test available
			var err error
			if d.OnAvailableViaHTTP {
				if d.OnAvailableTimeout == 0 {
					d.OnAvailableTimeout = 10
				}
				log.Logf(log.DEBUG, "-- onAvailable check '%v' at (%v), timeout (%v seconds)", d.Name, d.OnAvailable, d.OnAvailableTimeout)
				avail, err = hc.Reachable(d.OnAvailable, d.Name, d.OnAvailableTimeout, d.DisableValidateCert)
			} else {
				log.Logf(log.DEBUG, "-- ping host check %v", d.OnAvailable)
				avail, err = netut.Ping(d.OnAvailable)
			}
			if err != nil {
				log.Logf(log.DEBUG, "--- %v", err.Error())
			}

			if avail {
				log.Logf(log.DEBUG, "--- host available? %v '%v' at (%s)", log.Green("%v", avail), d.Name, d.OnAvailable)
			} else {
				log.Logf(log.DEBUG, "--- host available? %v '%v' at (%s)", log.Red("%v", avail), d.Name, d.OnAvailable)
			}
			hosts[i].SetAvailable(avail)
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
		if checkOnly && avail {
			break
		}
	}
	log.Logln(log.DEBUG, "")
	log.Logln(log.DEBUG, "-- Available Hosts Found --")
	for _, d := range hosts {
		if !d.Available() && !d.HostPfx() && !d.HostSuffix() {
			continue
		}
		log.Logf(log.DEBUG, "    %v at (%s)", d.Name, d.OnAvailable)
	}
	log.Logln(log.DEBUG, "")
	log.Logln(log.DEBUG, "- END testHosts process")
}

func pullURLToString(url string, auth *hc.Auth, disableVerifyCert bool) (string, error) {
	log.Logln(log.DEBUG, "--- Pull URL Content To String")
	httpClient := hc.NewClient(hc.HTTPClientRequestTimeout(30), hc.DisableVerifyClientCert(disableVerifyCert))

	resp, err := httpClient.Fetch("GET", url, auth, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("%s: %v", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url, resp.Status)
	}
	slurp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url, err)
	}
	return string(slurp), nil
}

func hashFileMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New() //nolint:gosec
	if _, err = io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
func validateHash(hashFileName, archivePathDir string, ac *updater.AppConfig) bool {
	log.Logln(log.DEBUG, "- Validate Hash")
	log.Logf(log.DEBUG, "-- extracting hash file %v", hashFileName)
	cmd := "tar xvf " + ac.ArchiveName + " " + archivePathDir + "/" + hashFileName
	var err error
	switch runtime.GOOS {
	case "windows":
		_, err = exec.Command("cmd", "/C", cmd).Output()
	default: //Mac & Linux
		_, err = exec.Command("sh", "-c", cmd).Output()
	}
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to execute command: %s %s", cmd, err.Error()))
		return false
	}

	// Read file contents
	hashContent, err := os.ReadFile(archivePathDir + "/" + hashFileName)
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to read hash: %s", err.Error()))
		return false
	}
	// create hash from downloaded application
	var hash string
	if strings.HasSuffix(hashFileName, ".md5") {
		hash, err = hashFileMD5(archivePathDir + "/" + ac.Name)
	} else if strings.HasSuffix(hashFileName, ".sha256") {
		hash, err = hashFileSHA256(archivePathDir + "/" + ac.Name)
	}
	if err != nil {
		log.Logln(log.WARN, fmt.Sprintf("--- failed to create hash from application: %s", err.Error()))
		return false
	}

	if string(hashContent) == hash {
		log.Logln(log.INFO, "-- VALID hash")
		return true
	}
	log.Logf(log.INFO, "-- INVALID hash %s", hash)
	log.Logln(log.DEBUG, "- Validate Hash END")
	return false
}
func hashFileSHA256(filePath string) (string, error) {
	var returnSHA256String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA256String, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return returnSHA256String, err
	}
	hashInBytes := hash.Sum(nil)[:32]
	returnSHA256String = hex.EncodeToString(hashInBytes)
	return returnSHA256String, nil
}
