package mymg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/colt3k/utils/crypt/genppk"
	"github.com/colt3k/utils/stringut"
	"github.com/pelletier/go-toml/v2"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/colt3k/utils/ques"

	iout "github.com/colt3k/utils/io"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// EXAMPLE:
// config=./build.toml mage -v -d tunler/ install
// config=./build.toml mage -v -d tunler/ release
// mage -v install or release
// bump= mage -v install or release

var (
	dryRun         bool
	displayOnly    bool
	configMessages bytes.Buffer
	config         Config
	build          BuildData
	postclean      PostClean
	apps           Apps
	prjkts         Projects
	arts           Artifactories
	scpS           SCPs
	sftpS          SFTPs
	pushCustom     PushCustoms
	pullCustom     PullCustoms
	timestamp      = time.Now().Unix()
	baseDir        = ""
	buildDir       = ""
	prepDir        = ""

	versionPkg            = "github.com/colt3k/utils"
	versionFieldsTemplate = `-X "%s/version.GITCOMMIT=%s" -X "%s/version.GITBRANCH=%s" -X "%s/version.VERSION=%s" -X "%s/version.BUILDDATE=%s" -X "%s/version.GOVERSION=%s"`
	overwriteValues       []string
	goLDFlagsTemplate     = "-s -w %s"
	goLDFlags             string
	// extldflags	Set space-separated flags to pass to the external linker.
	// -static 		means do not link against shared libraries
	goLDFlagsStaticTemplate = "-s -w %s -extldflags -static"
	goLDFlagsStatic         string
	bump                    bool
	names                   []string
	prompt                  bool
	nostatic                bool

	buildTags     = ""
	crossBuildDir = "cross"

	toCleanFiles []string
)

func setupBuild(props map[string]interface{}) error {
	mapProps := props["build"]
	if mapProps != nil {
		mp := mapProps.(map[string]interface{})
		build.Tags = propRtv(mp, "tags")
		build.UseAltApps = propRtv(mp, "useAltApps")
	} else {
		configMessages.WriteString("WARN: [build] section not declared\n")
	}

	config.Build = build
	return nil
}
func setupPostClean(props map[string]interface{}) error {
	mapProps := props["postclean"]
	if mapProps != nil {
		mp := mapProps.(map[string]interface{})
		dirs := make([]string, 0)
		if val, ok := mp["dirs"]; ok {
			for _, v := range val.([]interface{}) {
				dirs = append(dirs, v.(string))
			}
		}
		files := make([]string, 0)
		if val, ok := mp["files"]; ok {
			for _, v := range val.([]interface{}) {
				files = append(files, v.(string))
			}
		}
		postclean.Dirs = dirs
		postclean.Files = files
	} else {
		configMessages.WriteString("WARN: [postclean] section not declared\n")
	}

	config.PostClean = postclean
	return nil
}
func propRtv(p map[string]interface{}, key string) string {
	if val, ok := p[key]; ok {
		// verify this is found on local for apps
		if strings.HasSuffix(key, "Exe") {
			log.Printf("find value for key '%v'", key)

			if !fileExistsAndIsNotADir(val.(string)) {
				if strings.ToLower(build.UseAltApps) == "yes" || strings.ToLower(build.UseAltApps) == "y" ||
					strings.ToLower(build.UseAltApps) == "1" || strings.ToLower(build.UseAltApps) == "t" ||
					strings.ToLower(build.UseAltApps) == "true" {
					log.Printf("- !!! '%v' not found for key '%v' !!!\n", val.(string), key)
					path, err := findExec(filepath.Base(val.(string)))
					if err != nil {
						log.Fatalf("- nor on path %v\n", err)
					}
					log.Printf("- Using Found Alternative: %v\n", path)
					return path
				}
				log.Fatalf("- !!! '%v' not found for key '%v', useAltApps is off !!!\n", val.(string), key)
			}
		}
		return val.(string)
	} else if runtime.GOOS == "linux" {
		switch key {
		case "md5Exe":
			return "/bin/md5sum"
		case "sha1Exe":
			return "/bin/sha1sum"
		case "sha256Exe":
			return "/bin/sha256sum"
		case "curlExe":
			return "/bin/curl"
		case "catExe":
			return "/bin/cat"
		case "gitExe":
			return "/bin/git"
		case "tarExe":
			return "/bin/tar"
		case "scpExe":
			return "/bin/scp"
		case "sftpExe":
			return "/bin/sftp"
		case "whichExe":
			return "/usr/bin/which"
		}
	}
	return ""
}
func setupApps(props map[string]interface{}) error {
	mapProps := props["apps"]
	if mapProps != nil {
		mp := mapProps.(map[string]interface{})
		apps.MD5Exe = propRtv(mp, "md5Exe")
		apps.SHA1Exe = propRtv(mp, "sha1Exe")
		apps.SHA256Exe = propRtv(mp, "sha256Exe")
		apps.CurlExe = propRtv(mp, "curlExe")
		apps.CatExe = propRtv(mp, "catExe")
		apps.GitExe = propRtv(mp, "gitExe")
		apps.TarExe = propRtv(mp, "tarExe")
		apps.ScpExe = propRtv(mp, "scpExe")
		apps.SftpExe = propRtv(mp, "sftpExe")
		apps.UPXExe = propRtv(mp, "upxExe")
		apps.WhichExe = propRtv(mp, "whichExe")
	} else {
		configMessages.WriteString("WARN: [apps] section not declared\n")
	}

	config.Apps = apps
	return nil
}
func setupScps(props map[string]interface{}) error {
	mapProps := props["scp"]
	wrapper := make(map[string]interface{}, 1)
	wrapper["scp"] = mapProps

	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesWrapper, &scpS)
	if err != nil {
		return err
	}

	config.SCP = scpS
	return nil
}

func setupCustomPushes(props map[string]interface{}) error {
	mapProps := props["push-custom"]
	wrapper := make(map[string]interface{}, 1)
	wrapper["push-custom"] = mapProps

	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesWrapper, &pushCustom)
	if err != nil {
		return err
	}

	config.PushCustom = pushCustom
	return nil
}

func setupCustomPull(props map[string]interface{}) error {
	mapProps := props["pull-custom"]
	wrapper := make(map[string]interface{}, 1)
	wrapper["pull-custom"] = mapProps

	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesWrapper, &pullCustom)
	if err != nil {
		return err
	}

	config.PullCustom = pullCustom
	return nil
}

func setupSftps(props map[string]interface{}) error {
	mapProps := props["sftp"]
	wrapper := make(map[string]interface{}, 1)
	wrapper["sftp"] = mapProps

	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesWrapper, &sftpS)
	if err != nil {
		return err
	}

	config.SFTP = sftpS
	return nil
}
func setupArtifactory(props map[string]interface{}) error {
	artMap := props["artifactory"]
	artWrapper := make(map[string]interface{}, 1)
	artWrapper["artifactory"] = artMap

	bytesAppWrapper, err := json.MarshalIndent(artWrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesAppWrapper, &arts)
	if err != nil {
		return err
	}

	config.Artifactory = arts
	return nil
}

func setupProjects(props map[string]interface{}) error {
	appMap := props["project"]
	if appMap == nil {
		configMessages.WriteString("ERROR: [project] section not declared\n")
	}
	appWrapper := make(map[string]interface{}, 1)
	appWrapper["project"] = appMap

	bytesAppWrapper, err := json.MarshalIndent(appWrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesAppWrapper, &prjkts)
	if err != nil {
		return err
	}

	sz := len(prjkts.Projects)
	if sz == 1 {
		prompt = false
	}
	if sz > 0 {
		fmt.Printf("Projects available %d\n", sz)
	}
	// find absolute paths for files
	for i, d := range prjkts.Projects {
		if prompt && sz > 1 && ques.Confirm("process ("+d.Name+")? ") {
			prjkts.Projects[i].Enable = true
		} else if !prompt && sz == 1 {
			// if only one app enabled just run it, no confirmation needed
			prjkts.Projects[i].Enable = true
		}
		if !prompt && len(names) > 0 {
			for _, k := range names {
				if d.Name == k {
					prjkts.Projects[i].Enable = true
				}
			}
		}

		fmt.Printf("d.VersionFile: %v\n", d.VersionFile)
		prjkts.Projects[i].VersionFile, err = filepath.Abs(d.VersionFile)
		if err != nil {
			return err
		}

		prjkts.Projects[i].ReadmeFile, err = filepath.Abs(d.ReadmeFile)
		if err != nil {
			return err
		}

		fmt.Printf("d.ChangelogFile: %v\n", d.ChangelogFile)
		if len(strings.TrimSpace(d.ChangelogFile)) > 0 {
			prjkts.Projects[i].ChangelogFile, err = filepath.Abs(d.ChangelogFile)
			if err != nil {
				return err
			}
		}
		fmt.Printf("d.ChangelogFullFile: %v\n", d.ChangelogFullFile)
		if len(strings.TrimSpace(d.ChangelogFullFile)) > 0 {
			prjkts.Projects[i].ChangelogFullFile, err = filepath.Abs(d.ChangelogFullFile)
			if err != nil {
				return err
			}
		}

		for j, k := range d.Files {
			prjkts.Projects[i].Files[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}
		for j, k := range d.WindowsFiles {
			prjkts.Projects[i].WindowsFiles[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}
		for j, k := range d.LinuxFiles {
			prjkts.Projects[i].LinuxFiles[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}
		for j, k := range d.MacFiles {
			prjkts.Projects[i].MacFiles[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}

		for j, k := range d.OSDeployScripts {
			prjkts.Projects[i].OSDeployScripts[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}

		if len(d.OverrideVarFiles) > 0 {
			// find files; if one doesn't exist throw an error
			for j, k := range d.OverrideVarFiles {
				prjkts.Projects[i].OverrideVarFiles[j], err = filepath.Abs(k)
				if err != nil {
					return err
				}
				if !fileExistsAndIsNotADir(prjkts.Projects[i].OverrideVarFiles[j]) {
					fmt.Printf("\n!! One or more override_var_files NOT Found\n\t%v\n", prjkts.Projects[i].OverrideVarFiles[j])
					os.Exit(1)
				} else {
					fmt.Printf("Found\n\t%v\n", prjkts.Projects[i].OverrideVarFiles[j])
				}
			}
		}
		if len(d.OverrideVariables) > 0 && !displayOnly {
			// Split on semi-colon, create an array to store answers
			overrides := strings.Split(d.OverrideVariables, ";")
			overwriteValues = make([]string, 0)
			if len(d.OverrideVarFiles) > 0 && len(d.OverrideVarFiles) == len(overrides) {
				// set answers and skip asking for each; output on cmd line to show what is being used for each
				for _, q := range d.OverrideVarFiles {
					overwriteValues = append(overwriteValues, readAndOutput(q))
				}
			} else {
				for _, k := range overrides {
					ans := ques.Question(k + " value ? ")
					overwriteValues = append(overwriteValues, ans)
				}
			}
			fmt.Printf("Answers: %v\n", overwriteValues)
		}

		if len(d.YNPrompt) > 0 && !displayOnly {
			prompt = ques.Confirm(d.YNPrompt)
			if !prompt {
				fmt.Println("Please pull first!!!, Exiting...")
				os.Exit(1)
			}
		}
	}

	config.Project = prjkts
	return nil
}

func setupCleanFiles() {
	toCleanFiles = make([]string, 0)

	for _, d := range prjkts.Projects {
		toCleanFiles = append(toCleanFiles, d.Name)
	}
}

func parseTargets() error {
	var err error
	path, ok := os.LookupEnv("config")
	if !ok {
		// not defined set default
		path = "./build.toml"
	}
	path, _ = filepath.Abs(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var props map[string]interface{}
	err = toml.Unmarshal(data, &props)
	if err != nil {
		panic(err)
	}
	//tree, err := toml.LoadFile(path)
	//if err != nil {
	//	return err
	//}
	//props := tree.ToMap()

	err = setupProjects(props)
	if err != nil {
		return err
	}

	return nil
}

func parseToml() error {
	var err error

	_, ok := os.LookupEnv("nostatic")
	if ok {
		nostatic = true
	}
	_, ok = os.LookupEnv("bump")
	if ok {
		bump = true
	}
	_, ok = os.LookupEnv("dry")
	if ok {
		dryRun = true
	}

	namesIn, namesFound := os.LookupEnv("names")
	if !namesFound {
		prompt = true
	} else {
		names = strings.Split(namesIn, ",")
	}

	path, ok := os.LookupEnv("config")
	if !ok {
		// not defined set default
		path = "./build.toml"
	}
	path, _ = filepath.Abs(path)

	if !fileExistsAndIsNotADir(path) {
		return fmt.Errorf("toml file not found at %v", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var props map[string]interface{}
	err = toml.Unmarshal(data, &props)
	if err != nil {
		panic(err)
	}

	// APPS
	err = setupBuild(props)
	if err != nil {
		return err
	}
	err = setupPostClean(props)
	if err != nil {
		return err
	}
	for _, b := range postclean.Dirs {
		if !exists(b) && b != "PREP/" && b != "cross" {
			return fmt.Errorf("postclean dir not found %v", b)
		}
	}
	for _, b := range postclean.Files {
		if !fileExistsAndIsNotADir(b) {
			return fmt.Errorf("postclean file not found %v", b)
		}
	}
	err = setupApps(props)
	if err != nil {
		return err
	}
	err = setupScps(props)
	if err != nil {
		return err
	}
	//for _, b := range scpS.Instance {
	//	if !fileExistsAndIsNotADir(b.Path) {
	//		return fmt.Errorf("scp path not found %v", b.Path)
	//	}
	//}
	err = setupCustomPushes(props)
	if err != nil {
		return err
	}
	err = setupCustomPull(props)
	if err != nil {
		return err
	}
	//for _, b := range scpCustom.Instance {
	//	if !fileExistsAndIsNotADir(b.Exec) {
	//		return fmt.Errorf("scp-custom exec path not found %v", b.Exec)
	//	}
	//}

	err = setupSftps(props)
	if err != nil {
		return err
	}
	err = setupArtifactory(props)
	if err != nil {
		return err
	}
	// validate credential paths
	for _, b := range arts.Instance {
		//if !fileExistsAndIsNotADir(b.Creds) && !fileExistsAndIsNotADir(b.CredsPath) {
		//	return fmt.Errorf("af credentials not found %v", b.Creds)
		//}
		if len(b.CredsPath) > 0 {
			if !fileExistsAndIsNotADir(b.CredsPath) {
				fmt.Printf("af credential path not found %v\n", b.CredsPath)
				os.Exit(1)
			}
			credsPath := readAndOutput(b.CredsPath)
			if !fileExistsAndIsNotADir(credsPath) {
				fmt.Printf("af creds path declared but credentials not found %v\n", b.Creds)
				os.Exit(1)
			}
			fmt.Printf("creds path found %v\n", b.CredsPath)
		} else if len(b.Creds) > 0 && !fileExistsAndIsNotADir(b.Creds) {
			fmt.Printf("af credentials not found %v\n", b.Creds)
			os.Exit(1)
		}
	}
	err = setupProjects(props)
	if err != nil {
		return err
	}
	// validate paths
	for _, b := range prjkts.Projects {
		for _, q := range b.OSDeployScripts {
			if !fileExistsAndIsNotADir(q) {
				return fmt.Errorf("os deploy script not found %v for %v", q, b.Name)
			}
		}
		//if !exists(b.Package) {
		//	return fmt.Errorf("package path not found %v", b.Package)
		//}
		if !fileExistsAndIsNotADir(b.VersionFile) {
			return fmt.Errorf("version file not found %v", b.VersionFile)
		}
		if !fileExistsAndIsNotADir(b.ReadmeFile) {
			return fmt.Errorf("readme file not found %v", b.ReadmeFile)
		}
		for _, q := range b.Files {
			if !fileExistsAndIsNotADir(q) {
				return fmt.Errorf("file not found %v for %v", q, b.Name)
			}
		}
		for _, q := range b.WindowsFiles {
			if !fileExistsAndIsNotADir(q) {
				return fmt.Errorf("windows file not found %v for %v", q, b.Name)
			}
		}
		for _, q := range b.LinuxFiles {
			if !fileExistsAndIsNotADir(q) {
				return fmt.Errorf("linux file not found %v for %v", q, b.Name)
			}
		}
		for _, q := range b.MacFiles {
			if !fileExistsAndIsNotADir(q) {
				return fmt.Errorf("mac file not found %v for %v", q, b.Name)
			}
		}
	}

	var processApp bool
	//fmt.Println("") // clear line output
	if len(prjkts.Projects) == 0 {
		log.Println()
		log.Println("!!! No projects defined !!!")
		log.Fatalf("")
	}
	for _, d := range prjkts.Projects {
		fmt.Println("Project", d.Name, "Enabled? ", d.Enable)
		if d.Enable {
			processApp = true
		}
	}
	if !processApp {
		fmt.Println("\nNo project selected to process.")
		os.Exit(-1)
	}
	setupCleanFiles()

	log.Println("toCleanFiles: ", toCleanFiles)
	log.Println("toCleanPostDirs: ", postclean.Dirs)
	log.Println("toCleanPostFiles: ", postclean.Files)

	return nil
}

//type configuration struct {
//	EnvVars       map[string]string `json:"env_vars"`
//	BuildTags     string            `json:"build_tags"`
//	CleanDirs     []string          `json:"clean_dirs"`
//	Apps          map[string]string `json:"apps"`
//	Scps          scps              `json:"scps"`
//	ScpsCustom    scpcustoms        `json:"scps_custom"`
//	Sftp          sftps             `json:"sftp"`
//	Artifactories artifactories     `json:"artifactories"`
//	Applications  projects          `json:"projects"`
//}
//type scps struct {
//	Instance []scp `json:"scp"`
//}
//type scp struct {
//	Host     string `json:"host"`
//	Path     string `json:"path"`
//	SkipPing string `json:"skip_ping"`
//}

//type scpcustoms struct {
//	Instance []scpcust `json:"scp-custom"`
//}
//type scpcust struct {
//	Exec string `json:"exec"`
//}
//type sftps struct {
//	Instance []sftp `json:"sftp"`
//}
//type sftp struct {
//	Host     string `json:"host"`
//	Path     string `json:"path"`
//	SkipPing string `json:"skip_ping"`
//}
//type artifactories struct {
//	Instance []artifactory `json:"artifactory"`
//}
//type artifactory struct {
//	Host  string `json:"host"`
//	Path  string `json:"path"`
//	Creds string `json:"creds"`
//}
//
//type projects struct {
//	Projects []Project `json:"project"`
//}

type applications struct {
	Apps []application `json:"application"`
}

type application struct {
	Enable            bool     `json:"enable"`
	Name              string   `json:"name"`
	OSTargets         []string `json:"ostargets"`
	OSDeployScripts   []string `json:"osdeployscripts"`
	Package           string   `json:"package"`
	ReadmeFile        string   `json:"readme"`
	VersionFile       string   `json:"version"`
	ChangelogFile     string   `json:"changelog"`
	Files             []string `json:"files"`
	WindowsFiles      []string `json:"windowsfiles"`
	LinuxFiles        []string `json:"linuxfiles"`
	MacFiles          []string `json:"macfiles"`
	OverrideVariables string   `json:"override_variables"`
	OverrideVarFiles  []string `json:"override_var_files"`
	YNPrompt          string   `json:"ynprompt"`
}

func Help() {
	fmt.Println()
	fmt.Println("Mage Usage:")
	fmt.Println("Global Parameters")
	fmt.Println("  dry=		Dry Run")
	fmt.Println("  bump=		Bump to next version")
	fmt.Println("  help=		Show this help")
	fmt.Println("  names=	Single or Comma separated list of targets to build")
	fmt.Println("  config=	Location of build.toml default is ./build.toml")
	fmt.Println("  nostatic=	Do not set static flag for build")
	fmt.Println("Targets")
	fmt.Println("  help         show this help information")
	fmt.Println("  install      install the project to local machine")
	fmt.Println("  buildCross   create based on current local code and don't clean up")
	fmt.Println("  release      create and push release based on current local code")
	fmt.Println("  targets      show current project configured command's to build")
	fmt.Println("  display      show information after reading configuration")
	fmt.Println("  genconf      create an empty configuration build.toml")
	fmt.Println("  convert      convert build.toml to new format as migrated.toml")
	fmt.Println("  clean        clean any artifacts or build directories")
	fmt.Println("  auto         build true auto file and release")
	fmt.Println("  noauto       build false auto file and release")
	fmt.Println("  autostatus   show status of the auto file if it exists")
	fmt.Println("Flags")
	fmt.Println("  -v		show verbose mage output")
	fmt.Println("  -d		show custom debug output")
	fmt.Println()

}

type Config struct {
	Build       BuildData     `json:"build"`
	PostClean   PostClean     `json:"postclean"`
	Apps        Apps          `json:"apps"`
	SCP         SCPs          `json:"scp"`
	SFTP        SFTPs         `json:"sftp"`
	PushCustom  PushCustoms   `json:"push-custom"`
	PullCustom  PullCustoms   `json:"pull-custom"`
	Artifactory Artifactories `json:"artifactory"`
	Project     Projects      `json:"projects"`
}
type GenConfig struct {
	Build       BuildData         `json:"build" toml:"build" comment:"Build Options"`
	PostClean   PostClean         `json:"postclean" toml:"postclean" comment:"Directories to clean when complete"`
	Apps        Apps              `json:"apps" toml:"apps" comment:"Application paths"`
	SCP         []ScpData         `json:"scp" toml:"scp,omitempty" comment:"Array of Secure Copy Configurations\n- host must be available via ping check, requires ppk setup"`
	SFTP        []SftpData        `json:"sftp" toml:"sftp,omitempty" comment:"Array of SFTP Configurations"`
	PushCustom  []PushCustom      `json:"push-custom" toml:"push-custom,omitempty" comment:"Array of Custom Pushes, provide executable, parameters if needed and \n\twhich values (1,passes only full path) and (2, passes both full path and name)"`
	PullCustom  []PullCustom      `json:"pull-custom" toml:"pull-custom,omitempty" comment:"Array of Custom Pulls, provide executable, parameters if needed and \n\twhich values (1,passes only full path) and (2, passes both full path and name)"`
	Artifactory []ArtifactoryData `json:"artifactory" toml:"artifactory,omitempty" comment:"Array of Artifactory instances\n- host must be available via http check"`
	Project     []Project         `json:"project" toml:"project" comment:"Array of projects to build\n- duplicate this section for each project to build"`
}
type Projects struct {
	Projects []Project `json:"project" toml:"project"`
}
type SCPs struct {
	Instance []ScpData `json:"scp" toml:"scp"`
}
type PushCustoms struct {
	Instance []PushCustom `json:"push-custom" toml:"push-custom"`
}
type PullCustoms struct {
	Instance []PullCustom `json:"pull-custom" toml:"pull-custom"`
}
type SFTPs struct {
	Instance []SftpData `json:"sftp" toml:"sftp"`
}
type Artifactories struct {
	Instance []ArtifactoryData `json:"artifactory" toml:"artifactory"`
}
type BuildData struct {
	Tags       string `json:"tags" toml:"tags"`
	UseAltApps string `json:"useAltApps" toml:"useAltApps"`
}
type PostClean struct {
	Dirs  []string `json:"dirs" toml:"dirs"`
	Files []string `json:"files" toml:"files"`
}
type ScpData struct {
	Host     string `json:"host" toml:"host"`
	Path     string `json:"path" toml:"path"`
	SkipPing string `json:"skip_ping" toml:"skip_ping"`
	Backup   bool   `json:"backup" toml:"backup"`
}

func (s *ScpData) UnmarshalJSON(data []byte) error {
	type scpDataAlias ScpData
	scpData := &scpDataAlias{
		SkipPing: "false",
	}
	err := json.Unmarshal(data, scpData)
	if err != nil {
		return err
	}
	*s = ScpData(*scpData)
	return nil
}

type SftpData struct {
	Host     string `json:"host" toml:"host"`
	Path     string `json:"path" toml:"path"`
	SkipPing string `json:"skip_ping" toml:"skip_ping"`
	Backup   bool   `json:"backup" toml:"backup"`
}
type PushCustom struct {
	Exec       string `json:"exec" toml:"exec"`
	Parameters string `json:"parameters" toml:"parameters"`
	PassValues int64  `json:"values" toml:"values"`
}
type PullCustom struct {
	Exec       string `json:"exec" toml:"exec"`
	Parameters string `json:"parameters" toml:"parameters"`
	PassValues int64  `json:"values" toml:"values"`
	PrefixPath string `json:"prefix_path" toml:"prefix_path"`
}
type ArtifactoryData struct {
	Host      string `json:"host" toml:"host"`
	Path      string `json:"path" toml:"path"`
	Creds     string `json:"creds" toml:"creds"`
	CredsPath string `json:"creds_path" toml:"creds_path"`
	Backup    bool   `json:"backup" toml:"backup"`
}
type Project struct {
	Enable            bool     `json:"-" toml:"-"`
	Name              string   `json:"name" toml:"name"`
	OSTargets         []string `json:"ostargets" toml:"ostargets"`
	OSEnvFlags        []string `json:"os_env_flags" toml:"os_env_flags"`
	OSDeployScripts   []string `json:"osdeployscripts" toml:"osdeployscripts"`
	Package           string   `json:"package" toml:"package"`
	VersionFile       string   `json:"version" toml:"version"`
	ReadmeFile        string   `json:"readme" toml:"readme"`
	ChangelogFile     string   `json:"changelog" toml:"changelog"`
	ChangelogFullFile string   `json:"changelogfull" toml:"changelogfull"`
	Files             []string `json:"files" toml:"files"`
	WindowsFiles      []string `json:"windowsfiles" toml:"windowsfiles"`
	LinuxFiles        []string `json:"linuxfiles" toml:"linuxfiles"`
	MacFiles          []string `json:"macfiles" toml:"macfiles"`
	OverrideVariables string   `json:"override_variables" toml:"override_variables,omitempty" comment:"needs to be full repository path to variable defined semi-colon separated"`
	OverrideVarFiles  []string `json:"override_var_files" toml:"override_var_files,omitempty" comment:"local files used to fill override variables in order"`
	YNPrompt          string   `json:"ynprompt" toml:"ynprompt,omitempty"`
}
type Apps struct {
	MD5Exe    string `json:"md5Exe" toml:"md5Exe"`
	SHA1Exe   string `json:"sha1Exe" toml:"sha1Exe"`
	SHA256Exe string `json:"sha256Exe" toml:"sha256Exe"`
	CurlExe   string `json:"curlExe" toml:"curlExe"`
	CatExe    string `json:"catExe" toml:"catExe"`
	GitExe    string `json:"gitExe" toml:"gitExe"`
	TarExe    string `json:"tarExe" toml:"tarExe"`
	ScpExe    string `json:"scpExe" toml:"scpExe"`
	SftpExe   string `json:"sftpExe" toml:"sftpExe"`
	UPXExe    string `json:"upxExe" toml:"upxExe"`
	WhichExe  string `json:"whichExe" toml:"whichExe"`
}

func findApp(app string) (string, bool) {
	if !fileExistsAndIsNotADir(app) {
		path, err := findExec(filepath.Base(app))
		if err != nil {
			log.Printf("- Not Found on Path: %v\n", app)
			return app, false
		}
		log.Printf("- Not Found: %v -> Using: %v\n", app, path)
		return path, true
	}
	return app, true
}

func GenConf() {
	fmt.Println("- building config")
	c := &GenConfig{}
	c.Build.Tags = ""
	c.Build.UseAltApps = "yes"
	c.PostClean.Dirs = []string{"PREP/", "cross"}
	c.PostClean.Files = []string{}

	// Build c.Apps from found in environment and add fillers as needed
	c.Apps.MD5Exe, _ = findApp("/sbin/md5sum")
	c.Apps.SHA1Exe, _ = findApp("/usr/local/bin/sha1sum")
	c.Apps.SHA256Exe, _ = findApp("/usr/local/bin/sha256sum")
	c.Apps.CurlExe, _ = findApp("/usr/bin/curl")
	c.Apps.CatExe, _ = findApp("/bin/cat")
	c.Apps.GitExe, _ = findApp("/usr/local/bin/git")
	c.Apps.TarExe, _ = findApp("/usr/bin/tar")
	c.Apps.ScpExe, _ = findApp("/usr/bin/scp")
	c.Apps.SftpExe, _ = findApp("/usr/bin/sftp")
	c.Apps.UPXExe, _ = findApp("/usr/local/bin/upx")
	c.Apps.WhichExe, _ = findApp("/usr/bin/which")
	c.SCP = []ScpData{{Host: "main.domain.com", Path: "main:/root/apps", SkipPing: "false"}}
	c.SFTP = []SftpData{{Host: "main.domain.com", Path: "/apps/", SkipPing: "true"}}
	c.PushCustom = []PushCustom{{Exec: "path to executable", Parameters: "-d abc -f xzy", PassValues: 2}}
	c.Artifactory = []ArtifactoryData{{Host: "main.domain.com", Path: "http://main.domain.com:8081/artifactory/artifactoryreponame/appname/",
		Creds:     "/Users/username/keys/auths/.myartifactorycreds",
		CredsPath: "./pkgr/creds.txt"}}
	c.Project = []Project{{Name: "appname", OSTargets: []string{"darwin/amd64"}, OSDeployScripts: []string{"./pkgr/deploy_darwin.sh"},
		Package: "go.domain.com/colt3k/appname", VersionFile: "cmd/appname/VERSION.txt", ReadmeFile: "cmd/appname/README.md",
		ChangelogFile: "cmd/appname/CHANGES.txt", Files: []string{"./pkgr/bash_autocomplete", "cmd/appname/README.md"},
		WindowsFiles: []string{}, LinuxFiles: []string{}, MacFiles: []string{},
		YNPrompt: "Did you pull the latest? (y/n), will exit on 'n'", OverrideVariables: "-X \"domain.com/projectx/util.SEED=%s\";-X \"domain.com/projectx/util.KEY=%s\""}}

	b, err := toml.Marshal(c)
	if err != nil {
		log.Fatalf("issue marshalling config %v", err)
	}
	_, err = iout.WriteOut(b, "demo.toml")
	if err != nil {
		log.Fatalf("issue writing demo.toml %v", err)
	}
	fmt.Println("- build complete")
	fmt.Println()
	fmt.Println("*****************************************************")
	fmt.Println("update the demo.toml and rename to build.toml for use")
	fmt.Println("*****************************************************")
	fmt.Println()
}

func Convert() {
	// read old format and output new format
	fmt.Println("- Converting")

	// Find and load file
	path, ok := os.LookupEnv("config")
	if !ok {
		// not defined set default
		path = "./build.toml"
	}
	path, _ = filepath.Abs(path)
	if !fileExistsAndIsNotADir(path) {
		log.Fatalf("toml file not found at %v", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("issue reading file %v", err)
	}
	var props map[string]interface{}
	err = toml.Unmarshal(data, &props)
	if err != nil {
		panic(err)
	}
	c := &GenConfig{}
	fmt.Println("- Converting - build_tags")
	c.Build.Tags = props["build_tags"].(string)
	c.Build.UseAltApps = "yes"
	fmt.Println("- Converting - to_clean_dirs")
	c.PostClean.Dirs = convertInterfaceArToStringAr(props["to_clean_dirs"].([]interface{}))
	fmt.Println("- Converting - app paths")
	if props["md5Exe"] != nil {
		c.Apps.MD5Exe = props["md5Exe"].(string)
	} else {
		c.Apps.MD5Exe, _ = findApp("/sbin/md5sum")
	}
	if props["sha1Exe"] != nil {
		c.Apps.SHA1Exe = props["sha1Exe"].(string)
	} else {
		c.Apps.SHA1Exe, _ = findApp("/usr/local/bin/sha1sum")
	}
	if props["sha256Exe"] != nil {
		c.Apps.SHA256Exe = props["sha256Exe"].(string)
	} else {
		c.Apps.SHA256Exe, _ = findApp("/usr/local/bin/sha256sum")
	}
	if props["curlExe"] != nil {
		c.Apps.CurlExe = props["curlExe"].(string)
	} else {
		c.Apps.CurlExe, _ = findApp("/usr/bin/curl")
	}
	if props["catExe"] != nil {
		c.Apps.CatExe = props["catExe"].(string)
	} else {
		c.Apps.CatExe, _ = findApp("/bin/cat")
	}
	if props["gitExe"] != nil {
		c.Apps.GitExe = props["gitExe"].(string)
	} else {
		c.Apps.GitExe, _ = findApp("/usr/local/bin/git")
	}
	if props["tarExe"] != nil {
		c.Apps.TarExe = props["tarExe"].(string)
	} else {
		c.Apps.TarExe, _ = findApp("/usr/bin/tar")
	}
	if props["scpExe"] != nil {
		c.Apps.ScpExe = props["scpExe"].(string)
	} else {
		c.Apps.ScpExe, _ = findApp("/usr/bin/scp")
	}
	if props["sftpExe"] != nil {
		c.Apps.SftpExe = props["sftpExe"].(string)
	} else {
		c.Apps.SftpExe, _ = findApp("/usr/bin/sftp")
	}
	if props["upxExe"] != nil {
		c.Apps.UPXExe = props["upxExe"].(string)
	} else {
		c.Apps.UPXExe, _ = findApp("/usr/local/bin/upx")
	}
	if props["whichExe"] != nil {
		c.Apps.WhichExe = props["whichExe"].(string)
	} else {
		c.Apps.WhichExe, _ = findApp("/usr/bin/which")
	}
	fmt.Println("- Converting - SCP")
	scps := convertOldToGenConf(props, "scp").SCP
	if len(scps) > 0 {
		c.SCP = scps
	}
	fmt.Println("- Converting - SCP Custom")
	scpsCust := convertOldToGenConf(props, "scp-custom").PushCustom
	if len(scpsCust) > 0 {
		c.PushCustom = scpsCust
	}
	fmt.Println("- Converting - SFTP")
	sftps := convertOldToGenConf(props, "sftp").SFTP
	if len(sftps) > 0 {
		c.SFTP = convertOldToGenConf(props, "sftp").SFTP
	}
	fmt.Println("- Converting - Artifactory")
	afs := convertOldToGenConf(props, "artifactory").Artifactory
	if len(afs) > 0 {
		c.Artifactory = convertOldToGenConf(props, "artifactory").Artifactory
	}
	fmt.Println("- Converting - Projects/Applications")
	var tmp applications
	//log.Printf("scp : %v\n", props["scp"])
	mapProps := props["application"]
	wrapper := make(map[string]interface{}, 1)
	wrapper["application"] = mapProps
	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		log.Fatalf("issue parsing application %v", err)
	}
	err = json.Unmarshal(bytesWrapper, &tmp)
	if err != nil {
		log.Fatalf("issue marshalling application %v", err)
	}
	// loop through applications and convert to project type
	projects := make([]Project, 0)
	for _, m := range tmp.Apps {
		p := Project{}
		p.Enable = m.Enable
		p.Name = m.Name
		p.OSTargets = m.OSTargets
		p.OSDeployScripts = m.OSDeployScripts
		p.Package = m.Package
		p.ReadmeFile = m.ReadmeFile
		p.VersionFile = m.VersionFile
		p.ChangelogFile = m.ChangelogFile
		p.Files = m.Files
		if len(m.OverrideVariables) > 0 {
			p.OverrideVariables = m.OverrideVariables
		}
		if len(m.YNPrompt) > 0 {
			p.YNPrompt = m.YNPrompt
		}
		projects = append(projects, p)
	}
	c.Project = projects

	b, err := toml.Marshal(c)
	if err != nil {
		log.Fatalf("issue marshalling config %v", err)
	}
	tmpOut := strings.TrimSpace(string(b))
	_, err = iout.WriteOut([]byte(tmpOut), "migrated.toml")
	if err != nil {
		log.Fatalf("issue writing migrated.toml %v", err)
	}

	fmt.Println("- conversion complete")
	fmt.Println()
	fmt.Println("*************************************************************************")
	fmt.Println("compare the migrated.toml to your configuration build.toml and replace it")
	fmt.Println("*************************************************************************")
	fmt.Println()
}
func Display() {
	displayOnly = true
	mg.SerialDeps(parseToml)

	s, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println("- Configuration:", string(s))
	var byt bytes.Buffer
	for _, m := range config.Project.Projects {
		targs := len(m.OSTargets)
		ostargs := len(m.OSEnvFlags)
		depScripts := len(m.OSDeployScripts)
		if targs != depScripts {
			byt.WriteString(fmt.Sprintf("-- In project %v ostargets and osdeployscripts count should match one for one.\n", m.Name))
		}
		if targs != ostargs && ostargs > 0 {
			byt.WriteString(fmt.Sprintf("-- In project %v ostargets and os_env_flags count should match one for one.\n", m.Name))
		}
	}
	if byt.Len() > 0 || configMessages.Len() > 0 {
		fmt.Println("Warning/Errors Found")
		if configMessages.Len() > 0 {
			fmt.Println(configMessages.String())
		}
		if byt.Len() > 0 {
			fmt.Println(byt.String())
		}
	}
}

func Comment() {
	mg.SerialDeps(parseToml)
	fmt.Println()
	lastTag, _ := sh.Output(apps.GitExe, "describe", "--tags", "--abbrev=0")
	//fmt.Printf("last Tag : %v\n", lastTag)
	comments, _ := sh.Output(apps.GitExe, "log", lastTag+"..HEAD", "--oneline")
	fmt.Printf("Comments : %v\n", comments)
	lines := strings.Split(comments, "\n")
	var byt bytes.Buffer
	for _, l := range lines {
		// all after first space
		var bytLn bytes.Buffer
		flds := strings.Fields(l)
		for i, k := range flds {
			if i > 0 {
				bytLn.WriteString(k)
				bytLn.WriteString(" ")
			}
		}
		byt.WriteString(bytLn.String())
		byt.WriteString("\n")
	}
	fmt.Printf("formatted comments: %v\n", byt.String())
}
func Targets() {
	mg.SerialDeps(parseTargets)
	fmt.Println()
	fmt.Println("Targets")
	for _, d := range prjkts.Projects {
		fmt.Println("  ", d.Name)
	}

	fmt.Println()
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.SerialDeps(parseToml)
	skipUPX := false
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	fmt.Println("goos: " + goos + " goarch: " + goarch)
	if goos == "darwin" && (goarch == "arm64" || goarch == "amd64") {
		skipUPX = true
	}
	gocmd := mg.GoCmd()
	fmt.Println("Building...")
	//$(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME) .
	var err error
	if cgoval, ok := os.LookupEnv("CGO_ENABLED"); !ok {
		err = os.Setenv("CGO_ENABLED", "0")
		if err != nil {
			fmt.Println("issue setting CGO_ENABLED to 0:", err)
		}
	} else {
		fmt.Printf("CGO_ENABLED set to %v\n", cgoval)
	}

	for _, d := range prjkts.Projects {
		if exists(d.Name) {
			log.Fatalf("\n!!! ERROR: file or directory already exists with application name '%s', exiting !!!\n\n", d.Name)
		}
		if !d.Enable {
			continue
		}

		var flags string
		if len(d.OSEnvFlags) > 0 {
			for i, q := range d.OSTargets {
				if q == goos+"/"+goarch {
					flags = d.OSEnvFlags[i]
					fmt.Printf("flags found %v", flags)
					err = setEnvFlags(flags, dryRun)
					if err != nil {
						return fmt.Errorf("  !!! failed to set env vars %v", err)
					}
				}
			}
		}
		fmt.Println("Building ", d.Name)
		err = setup(d)
		if err != nil {
			fmt.Println("issue with setup :", err)
		}
		cleaner(d.Name, false)
		err = Format()
		if err != nil {
			fmt.Println("issue formatting :", err)
		}

		err = Lint()
		if err != nil {
			fmt.Println("issue linting :", err)
		}
		err = Test()
		if err != nil {
			fmt.Println("issue testing :", err)
		}
		err = Vet()
		if err != nil {
			fmt.Println("issue vetting :", err)
		}
		name := d.Name
		if runtime.GOOS == "windows" {
			name += ".exe"
		}

		projectMainDir := "./cmd/" + d.Name + "/"
		if !exists(projectMainDir) {
			fmt.Println("Path doesn't exist: " + projectMainDir + " using local dir '.' instead")
			projectMainDir = "."
		} else {
			projectMainDir = projectMainDir + "."
		}

		if !dryRun {
			err = sh.RunV(gocmd, "build", "-trimpath", "-tags", buildTags, "-ldflags", goLDFlags, "-o", name, projectMainDir)
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
			if exists(apps.UPXExe) && !skipUPX {
				fmt.Printf("\n*** START UPX binary compression on  %v ***\n", name)
				fi, _ := os.Stat(name)
				fmt.Printf("\n")
				err = sh.RunV(apps.UPXExe, "-q", "-q", "-q", name)
				if err != nil {
					clearEnvFlags(flags, dryRun)
					return err
				}
				fmt.Printf("\n")
				fmt.Printf("    - prior to compression: %v\n", stringut.HRByteCount(fi.Size(), false))
				fi, _ = os.Stat(name)
				fmt.Printf("    - post compression: %v\n", stringut.HRByteCount(fi.Size(), false))
				fmt.Printf("\n*** END UPX binary compression on  %v ***\n", name)
			} else if skipUPX {
				fmt.Println("- upx not applicable for darwin/arm64 - ")
			} else {
				fmt.Println("- no upx available for binary compression - ")
			}
			clearEnvFlags(flags, dryRun)
		} else {
			clearEnvFlags(flags, dryRun)
			var byt bytes.Buffer
			byt.WriteString(gocmd + " build -trimpath -tags " + buildTags + " -ldflags " + goLDFlags + " -o " + name + " " + projectMainDir)
			if exists(apps.UPXExe) {
				byt.WriteString("\n")
				byt.WriteString(apps.UPXExe + " -q -q -q " + name)
			}
			fmt.Println("DRY_RUN: Building build", byt.String())
		}
	}

	return nil
}

func BumpVersion() error {
	mg.Deps(parseToml)
	gocmd := mg.GoCmd()

	if !dryRun {
		_, err := sh.Output(apps.WhichExe, "sembump")
		if err != nil {
			// update if not found
			fmt.Println("Updating sembump")
			err = sh.RunV(gocmd, "get", "-u", "github.com/colt3k/utils/sembump@latest")
			if err != nil {
				log.Println(err)
				return err
			}
		}
	} else {
		fmt.Println("DRY_RUN: " + gocmd + " get -u github.com/colt3k/utils/sembump@latest")
	}

	for _, d := range prjkts.Projects {
		if d.Enable && bump {
			fmt.Println("Bumping Version...")
			ver := version(d.VersionFile)
			nVersion, err := sh.Output("sembump", "--kind", "patch", ver)
			if err != nil {
				return err
			}
			fmt.Printf("  Bumping VERSION.txt from %s to %s\n", ver, nVersion)
			if !dryRun {
				_, err = iout.WriteOut([]byte(nVersion), d.VersionFile)
				if err != nil {
					return err
				}
			}

			fmt.Printf("  Updating links to download binaries in README.md\n")
			// read in modify content and write out instead
			if !dryRun {
				var out2 string
				out2, err = sh.Output(apps.CatExe, d.ReadmeFile)
				if err != nil {
					return err
				}
				readmeContent := strings.Replace(out2, ver, nVersion, -1)
				_, err = iout.WriteOut([]byte(readmeContent), d.ReadmeFile)
				if err != nil {
					return err
				}
			}

			fmt.Printf("  updated %s", d.ReadmeFile)

			if !dryRun {
				// Add Version and Readme file to git prep
				err = sh.RunV(apps.GitExe, "add", d.VersionFile, d.ReadmeFile)
				if err != nil {
					fmt.Println("issue on git add bump:", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + apps.GitExe + " add " + d.VersionFile + " " + d.ReadmeFile)
			}

			if !dryRun {
				// Commit tag
				err = sh.RunV(apps.GitExe, "commit", "-vsam", "Bump version to "+nVersion)
				if err != nil {
					fmt.Println("issue committing bump:", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + apps.GitExe + " commit -vsam Bump version to " + nVersion)
			}

			// if there is more than one command for this project create a unique tag for it
			if strings.Contains(d.VersionFile, "/cmd/") {
				nVersion = d.Name + "/" + nVersion
			}
			if !dryRun {
				// "Run make tag to create and push the tag for new version $(NEW_VERSION)"
				err = sh.RunV(apps.GitExe, "tag", "-a", nVersion, "-m", nVersion)
				if err != nil {
					fmt.Println("issue tagging bump :", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + apps.GitExe + " tag -a " + nVersion + " -m " + nVersion)
			}
			// Push Tag
			if !dryRun {
				err = sh.RunV(apps.GitExe, "push", "origin", nVersion)
				if err != nil {
					fmt.Println("issue pushing tag :", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + apps.GitExe + "push origin " + nVersion)
			}

			if len(strings.TrimSpace(d.ChangelogFullFile)) > 0 && fileExistsAndIsNotADir(d.ChangelogFullFile) {
				fmt.Printf("ChangelogFullFile %v\n", d.ChangelogFullFile)
				fmt.Println()
				fmt.Println("Rebuilding CHANGELOG.md")
				// Rebuild CHANGELOG.md d.ChangelogFullFile

				var chglogPath string
				var found bool
				if chglogPath, found = findApp("chglog"); found {
					output, err := sh.Output(chglogPath, "init", "-file", d.ChangelogFullFile)
					if err != nil {
						return err
					}
					fmt.Printf("%v", output)
				} else {
					fmt.Printf("chglog application not found at %v\n", chglogPath)
				}
			} else {
				fmt.Println("ChangelogFullFile Not Set")
			}
			if len(strings.TrimSpace(d.ChangelogFile)) > 0 && fileExistsAndIsNotADir(d.ChangelogFile) {
				fmt.Printf("ChangelogFile %v\n", d.ChangelogFile)
				fmt.Println()
				fmt.Println("Updating CHANGES.txt")
				// Update CHANGES.txt d.ChangelogFile
				var chglogPath string
				var found bool
				if chglogPath, found = findApp("chglog"); found {
					output, err := sh.Output(chglogPath, "changes", "-file", d.ChangelogFile)
					if err != nil {
						return err
					}
					fmt.Printf("%v", output)
				} else {
					fmt.Printf("chglog application not found at %v\n", chglogPath)
				}
			} else {
				fmt.Println("ChangelogFile Not Set")
			}
		}
	}
	return nil
}
func BuildCross() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		fmt.Println("\n** Setup")
		err := setup(d)
		if err != nil {
			fmt.Println("issue setup :", err)
		}
		fmt.Println("\n** Release Prep")
		// make dir PREP
		fmt.Println("  - create PREP dir")
		if err = os.MkdirAll(prepDir, 0700); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create %q: %v", prepDir, err)
		}

		// loop through included files and place in PREP
		fmt.Println("  - add files to PREP")
		for _, k := range d.Files {
			fileName := filepath.Base(k)
			if fileName == "bash_autocomplete" {
				fileName = d.Name + ".bash"
			}

			fileTarget := filepath.Join(prepDir, fileName)
			var fullFilePath string
			fullFilePath, err = filepath.Abs(k)
			if err != nil {
				return err
			}
			log.Printf("    copying %s to %s\n", fullFilePath, fileTarget)
			err = sh.Copy(fileTarget, fullFilePath)
			if err != nil {
				return err
			}
		}
		// end of loop

		// cross building
		err = cross(d)
		if err != nil {
			log.Println(err)
			return err
		}

		cleaner(d.Name, true)
	}
	return nil
}

func Release() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		fmt.Println("\nSetup")
		err := setup(d)
		if err != nil {
			fmt.Println("issue setup :", err)
		}
		fmt.Println("Release Prep")
		// make dir PREP
		fmt.Println("  create PREP dir")
		if err = os.MkdirAll(prepDir, 0700); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create %q: %v", prepDir, err)
		}

		// loop through included files and place in PREP
		fmt.Println("  add files to PREP")
		for _, k := range d.Files {
			fileName := filepath.Base(k)
			if fileName == "bash_autocomplete" {
				fileName = d.Name + ".bash"
			}

			fileTarget := filepath.Join(prepDir, fileName)
			var fullFilePath string
			fullFilePath, err = filepath.Abs(k)
			if err != nil {
				return err
			}
			log.Printf("copying %s to %s\n", fullFilePath, fileTarget)
			err = sh.Copy(fileTarget, fullFilePath)
			if err != nil {
				return err
			}
		}
		// end of loop

		// cross building
		err = cross(d)
		if err != nil {
			log.Println(err)
			return err
		}

		err = scpCopy(d.Name)
		if err != nil {
			fmt.Println("issue scp copy :", err)
		}
		err = customCopy(d.Name)
		if err != nil {
			fmt.Println("issue push custom copy :", err)
		}
		err = sftpCopy(d.Name)
		if err != nil {
			fmt.Println("issue sftp copy :", err)
		}
		err = artifactoryPush(d.Name)
		if err != nil {
			fmt.Println("issue artifactory push :", err)
		}

		cleaner(d.Name, false)
	}
	return nil
}

func Auto() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		fmt.Println("\nSetup")
		err := setup(d)
		if err != nil {
			fmt.Println("issue setup :", err)
		}

		_, err = iout.WriteOut([]byte("true"), filepath.Join(baseDir, d.Name+".auto"))
		if err != nil {
			return err
		}

		//_, err = iout.WriteOut([]byte("false"), filepath.Join(baseDir, d.Name+".auto"))
		//if err != nil {
		//	return err
		//}

		err = scpCopy(d.Name)
		if err != nil {
			fmt.Println("issue scp copy :", err)
		}
		err = customCopy(d.Name)
		if err != nil {
			fmt.Println("issue push custom copy :", err)
		}
		err = sftpCopy(d.Name)
		if err != nil {
			fmt.Println("issue sftp copy :", err)
		}
		err = artifactoryPush(d.Name)
		if err != nil {
			fmt.Println("issue artifactory push :", err)
		}

		//cleaner(d.Name, false)
	}

	return nil
}
func AutoStatus() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		fmt.Println("\nSetup")
		err := setup(d)
		if err != nil {
			fmt.Println("issue setup :", err)
		}

		err = scpCopyAutoStatus(d.Name)
		if err != nil {
			fmt.Println("issue scp pull auto status :", err)
		}
		err = customPull(d.Name)
		if err != nil {
			fmt.Println("issue pull auto status :", err)
		}
		err = sftpCopyAutoStatus(d.Name)
		if err != nil {
			fmt.Println("issue sftp pull auto status :", err)
		}

		err = artifactoryPullAutoStatus(d.Name)
		if err != nil {
			fmt.Println("issue artifactory pull auto status :", err)
		}
	}
	return nil
}
func NoAuto() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		fmt.Println("\nSetup")
		err := setup(d)
		if err != nil {
			fmt.Println("issue setup :", err)
		}

		_, err = iout.WriteOut([]byte("false"), filepath.Join(baseDir, d.Name+".auto"))
		if err != nil {
			return err
		}

		err = scpCopy(d.Name)
		if err != nil {
			fmt.Println("issue scp copy :", err)
		}
		err = customCopy(d.Name)
		if err != nil {
			fmt.Println("issue push custom copy :", err)
		}
		err = sftpCopy(d.Name)
		if err != nil {
			fmt.Println("issue sftp copy :", err)
		}
		err = artifactoryPush(d.Name)
		if err != nil {
			fmt.Println("issue artifactory push :", err)
		}

		cleaner(d.Name, false)
	}

	return nil
}

func setEnvFlags(flags string, dryRun bool) error {
	fmt.Println("  - Set Environment Vars")
	if len(strings.TrimSpace(flags)) > 0 {
		for _, q := range strings.Split(flags, ",") {
			flgParts := strings.Split(q, "=")
			p1 := strings.TrimSpace(flgParts[0])
			p2 := strings.TrimSpace(flgParts[1])
			if strings.Contains(p2, "\"") {
				p2 = strings.ReplaceAll(p2, "\"", "")
			}
			fmt.Printf("    setting env var %v = %v\n", p1, p2)
			//if !dryRun {
			err := os.Setenv(p1, p2)
			if err != nil {
				fmt.Printf("    !!! issue setting %v to %v: %v\n", p1, p2, err)
				return err
			}
			//}
		}
	}
	return nil
}
func clearEnvFlags(flags string, dryRun bool) error {
	// clear any prior flags
	fmt.Println("  Remove Environment Vars")
	if len(strings.TrimSpace(flags)) > 0 {
		for _, q := range strings.Split(flags, ",") {
			flgParts := strings.Split(q, "=")
			p1 := strings.TrimSpace(flgParts[0])
			p2 := strings.TrimSpace(flgParts[1])
			if strings.Contains(p2, "\"") {
				p2 = strings.ReplaceAll(p2, "\"", "")
			}
			fmt.Printf("    clearing env var %v = %v\n", p1, p2)
			if !dryRun {
				err := os.Unsetenv(p1)
				if err != nil {
					fmt.Printf("    !!! issue setting %v to %v: %v\n", p1, p2, err)
					return err
				}
			}
		}
	}
	return nil
}

// Build for all defined Architectures
func cross(app Project) error {
	fmt.Println("\n** CrossBuilding...")
	gocmd := mg.GoCmd()
	for i, d := range app.OSTargets {
		fmt.Println("\n  Building...")
		goosArch := strings.Split(d, "/")
		goos := goosArch[0]
		arch := goosArch[1]
		arm := ""
		if len(goosArch) == 3 {
			arm = goosArch[2]
		}
		err := os.Setenv("GOOS", goos)
		if err != nil {
			fmt.Println("  !!! issue setting GOOS :", goos, err)
		} else {
			fmt.Printf("  GOOS set to %v\n", goos)
		}
		err = os.Setenv("GOARCH", arch)
		if err != nil {
			fmt.Println("  !!! issue setting GOARCH :", arch, err)
		} else {
			fmt.Printf("  GOARCH set to %v\n", arch)
		}
		err = os.Setenv("  GOARM", arm)
		if err != nil {
			fmt.Println("  !!! issue setting GOARM :", arm, err)
		} else {
			if len(arm) > 0 {
				fmt.Printf("  GOARM set to %v\n", arm)
			} else {
				fmt.Println("  GOARM NOT set")
			}
		}

		skipUPX := false
		if goos == "darwin" && (arch == "arm64" || arch == "amd64") {
			skipUPX = true
		}
		//goosAct := runtime.GOOS
		//goarchAct := runtime.GOARCH

		var flags string
		if len(app.OSEnvFlags) > 0 {
			flags = app.OSEnvFlags[i]
			fmt.Printf("flags found %v", flags)
			err = setEnvFlags(flags, dryRun)
			if err != nil {
				return fmt.Errorf("  !!! failed to set env vars %v", err)
			}
		}

		//fmt.Printf("actual %v/%v, building %v/%v\n", goosAct, goarchAct, goos, arch)
		//if goosAct == goos && goarchAct == arch {
		//	if cgoval, ok := os.LookupEnv("CGO_ENABLED"); ok {
		//		wasCGO = cgoval
		//		err = os.Setenv("CGO_ENABLED", "0")
		//		if err != nil {
		//			fmt.Println("issue setting CGO_ENABLED to 0:", err)
		//		}
		//	}
		//} else if wasCGO != "0" {
		//	err = os.Setenv("CGO_ENABLED", wasCGO)
		//	if err != nil {
		//		fmt.Println("issue setting CGO_ENABLED to 0:", err)
		//	}
		//}
		//if cgoval, ok := os.LookupEnv("CGO_ENABLED"); !ok {
		//	err = os.Setenv("CGO_ENABLED", "0")
		//	if err != nil {
		//		fmt.Println("issue setting CGO_ENABLED to 0:", err)
		//	}
		//} else {
		//	fmt.Printf("CGO_ENABLED set to %v was CGO Set %v\n", cgoval, wasCGO)
		//}

		name := app.Name
		if goos == "windows" {
			name += ".exe"
		}
		osarchName := app.Name + "-" + goos + "-" + arch
		if len(arm) > 0 {
			osarchName = app.Name + "-" + goos + "-" + arch + "-" + arm
		}
		path := filepath.Join(buildDir, osarchName)

		if !dryRun {
			if err = os.MkdirAll(path, 0700); err != nil && !os.IsExist(err) {
				clearEnvFlags(flags, dryRun)
				return fmt.Errorf("  !!! failed to create %q: %v", path, err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: creating directory\n\t%v\n", path)
		}
		fmt.Printf("  - Packaging %s/%s\n", goos, arch)
		executableName := filepath.Join(path, name)

		projectMainDir := "./cmd/" + app.Name + "/"
		if !exists(projectMainDir) {
			fmt.Println("    !!! Path doesn't exist: " + projectMainDir + " using local dir '.' instead")
			projectMainDir = "."
		} else {
			projectMainDir = projectMainDir + "."
		}

		if !dryRun {
			if nostatic {
				goLDFlagsStatic = goLDFlags
			}
			err = sh.RunV(gocmd, "build", "-tags", buildTags, "-ldflags", goLDFlagsStatic, "-o", executableName, projectMainDir)
			if err != nil {
				clearEnvFlags(flags, dryRun)
				fmt.Println("\n** If failed cross build on arm64 for amd64 with CGO_ENABLED try this on the command line instead. **\n")
				fmt.Printf("GOOS=%v GOARCH=%v CGO_ENABLED=1 go build -tags %v '-ldflags %v' -o %v %v\n\n", goos, arch, buildTags, goLDFlagsStatic, executableName, projectMainDir)
				return err
			}
			if exists(apps.UPXExe) && !skipUPX {
				fmt.Printf("\n  *** START UPX binary compression on  %v ***\n", executableName)
				fi, _ := os.Stat(executableName)
				fmt.Printf("\n")
				err = sh.RunV(apps.UPXExe, "-q", "-q", "-q", executableName)
				if err != nil {
					clearEnvFlags(flags, dryRun)
					return err
				}
				fmt.Printf("\n")
				fmt.Printf("    - prior to compression: %v\n", stringut.HRByteCount(fi.Size(), false))
				fi, _ = os.Stat(executableName)
				fmt.Printf("    - post compression: %v\n", stringut.HRByteCount(fi.Size(), false))
				fmt.Printf("\n  *** END UPX binary compression on  %v ***\n", executableName)
			} else if skipUPX {
				fmt.Println("- upx not applicable for darwin/arm64 - ")
			} else {
				fmt.Println("- no upx available for binary compression - ")
			}
		} else {
			if nostatic {
				goLDFlagsStatic = goLDFlags
			}
			fmt.Println("    - DRY_RUN: " + gocmd + " build -trimpath -tags " + buildTags + " '-ldflags " + goLDFlagsStatic + "' -o " + executableName + " " + projectMainDir)
		}
		// make release dir for this OS

		osarchDir := filepath.Join(baseDir, osarchName)
		if !dryRun {
			if err = os.MkdirAll(osarchDir, 0700); err != nil && !os.IsExist(err) {
				clearEnvFlags(flags, dryRun)
				return fmt.Errorf("  failed to create %q: %v", path, err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: making directory\n\t%v\n", osarchDir)
		}

		if !dryRun {
			// copy executable to release dir
			err = sh.Copy(filepath.Join(osarchDir, name), executableName)
			if err != nil {
				fmt.Println("  !!! issue copying :", executableName, err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: copying\n\t%v\n\tto  %v\n", executableName, filepath.Join(osarchDir, name))
		}
		scriptName := filepath.Base(app.OSDeployScripts[i])
		scriptContent := buildDeployScript(app.OSDeployScripts[i])
		fmt.Printf("  - Reading %v\n", scriptName)
		if !dryRun {
			_, err = iout.WriteOut(scriptContent, filepath.Join(osarchDir, scriptName))
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
		} else {
			fmt.Printf("    DRY_RUN: writing out deploy script to\n\t%v\n", filepath.Join(osarchDir, scriptName))
		}

		switch goos {
		case "windows":
			for _, x := range app.WindowsFiles {
				scriptName = filepath.Base(x)
				scriptContent = buildDeployScript(x)
				if !dryRun {
					_, err = iout.WriteOut(scriptContent, filepath.Join(osarchDir, scriptName))
					if err != nil {
						clearEnvFlags(flags, dryRun)
						return err
					}
				} else {
					fmt.Println("  - DRY_RUN: writing out windows deploy script to " + filepath.Join(osarchDir, scriptName))
				}
			}
		case "linux":
			for _, x := range app.LinuxFiles {
				scriptName = filepath.Base(x)
				if scriptName == "bash_autocomplete" {
					scriptName = app.Name + ".bash"
				}
				scriptContent = buildDeployScript(x)
				if !dryRun {
					_, err = iout.WriteOut(scriptContent, filepath.Join(osarchDir, scriptName))
					if err != nil {
						clearEnvFlags(flags, dryRun)
						return err
					}
				} else {
					fmt.Println("  - DRY_RUN: writing out linux deploy script to " + filepath.Join(osarchDir, scriptName))
				}
			}
		case "darwin":
			for _, x := range app.MacFiles {
				scriptName = filepath.Base(x)
				if scriptName == "bash_autocomplete" {
					scriptName = app.Name + ".bash"
				}
				scriptContent = buildDeployScript(x)
				if !dryRun {
					_, err = iout.WriteOut(scriptContent, filepath.Join(osarchDir, scriptName))
					if err != nil {
						clearEnvFlags(flags, dryRun)
						return err
					}
				} else {
					fmt.Println("  - DRY_RUN: writing out mac deploy script to " + filepath.Join(osarchDir, scriptName))
				}
			}
		}

		// create hash files
		if !dryRun {
			var md5sum string
			md5sum, err = sh.Output(apps.MD5Exe, executableName)
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
			md5sumParts := strings.Fields(md5sum)

			_, err = iout.WriteOut([]byte(md5sumParts[0]), executableName+".md5")
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
			err = sh.Copy(filepath.Join(osarchDir, name+".md5"), executableName+".md5")
			if err != nil {
				fmt.Println("  !!! issue copying :", executableName+".md5", err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: creating/writing and copying md5 hash file \n\t%v\n\t%v\n", executableName+".md5", filepath.Join(osarchDir, name+".md5"))
		}

		if !dryRun {
			var shasum string
			shasum, err = sh.Output(apps.SHA256Exe, executableName)
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
			sha256Parts := strings.Fields(shasum)
			_, err = iout.WriteOut([]byte(sha256Parts[0]), executableName+".sha256")
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
			err = sh.Copy(filepath.Join(osarchDir, name+".sha256"), executableName+".sha256")
			if err != nil {
				fmt.Println("  !!! issue copying :", executableName+".sha256", err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: creating/writing and copying sha256 hash file \n\t%v\n\t%v\n", executableName+".sha256", filepath.Join(osarchDir, name+".sha256"))
		}

		if !dryRun {
			// Copy change log
			nm := filepath.Base(app.ChangelogFile)
			err = sh.Copy(filepath.Join(osarchDir, nm), app.ChangelogFile)
			if err != nil {
				log.Printf("  !!! ERROR issue copying changelog %v", err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: copying Changelog file: \n\t%v\n", app.ChangelogFile)
		}

		// Copy all prep files into osarchDir
		var files []os.FileInfo

		entries, err := os.ReadDir(prepDir)
		if err != nil {
			clearEnvFlags(flags, dryRun)
			log.Fatal(err)
		}
		for _, entry := range entries {
			var info fs.FileInfo
			info, err = entry.Info()
			if err != nil {
				clearEnvFlags(flags, dryRun)
				log.Fatal(err)
			}
			files = append(files, info)
		}

		if dryRun && len(files) == 0 {
			fmt.Printf("  - DRY_RUN: no files found in \n\t%v\n", prepDir)
		} else {
			fmt.Printf("  - DRY_RUN: writing out prep files in \n\t%v\n", prepDir)
		}
		for _, f := range files {
			if goos == "windows" && f.Name() == app.Name+".bash" {
				continue
			}
			fmt.Println("    Found: ", f.Name())
			// Copy here
			fileTarget := filepath.Join(osarchDir, f.Name())
			if !dryRun {
				log.Printf("    copy \n\t%v\n\tto  %v\n", filepath.Join(prepDir, f.Name()), fileTarget)
				err = sh.Copy(fileTarget, filepath.Join(prepDir, f.Name()))
				if err != nil {
					clearEnvFlags(flags, dryRun)
					return err
				}
			} else {
				fmt.Printf("    DRY_RUN: copy \n\t%v\n\tto  %v\n", filepath.Join(prepDir, f.Name()), fileTarget)
			}
		}

		if !dryRun {
			fmt.Println("  Creating archive ", osarchName)
			err = sh.RunV(apps.TarExe, "cfz", osarchName+".tgz", osarchName)
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
		} else {
			fmt.Println("  - DRY_RUN:  Creating archive ", osarchName)
		}

		v := version(app.VersionFile)
		u := update{Os: goos, Arch: arch, Name: name, Timestamp: timestamp, Version: v, Changelog: app.Name + "-changes.txt"}
		updater, err := buildUpdateDir(u)
		if err != nil {
			clearEnvFlags(flags, dryRun)
			return err
		}
		fmt.Println("  ", string(updater))
		if !dryRun {
			_, err = iout.WriteOut(updater, filepath.Join(baseDir, osarchName+".update"))
			if err != nil {
				clearEnvFlags(flags, dryRun)
				return err
			}
		} else {
			fmt.Printf("  - DRY_RUN: writing out update file\n\t%v\n", filepath.Join(baseDir, osarchName+".update"))
		}
		// Copy changelog to release dir
		if !dryRun {
			err = sh.Copy(filepath.Join(baseDir, app.Name+"-changes.txt"), app.ChangelogFile)
			if err != nil {
				log.Printf("  !!! ERROR issue copying changelogfile %v", err)
			}
		} else {
			nm := filepath.Base(app.ChangelogFile)
			fmt.Printf("  - DRY_RUN: copying Changes file\n\t%v\n", filepath.Join(baseDir, nm))
		}

		if !dryRun {
			// remove release dir
			err = os.RemoveAll(osarchDir)
			if err != nil {
				fmt.Println("  !!! issue removing all :", osarchDir, err)
			}
		} else {
			fmt.Printf("  - DRY_RUN: removing all files in\n\t%v\n", osarchDir)
		}
		clearEnvFlags(flags, dryRun)
	}

	return nil
}

func scpCopyAutoStatus(projectName string) error {
	for _, k := range scpS.Instance {
		fmt.Println("\nSCP Pull... ")
		if len(k.Host) > 0 {
			foundHost := false
			if strings.ToLower(k.SkipPing) != "true" && strings.ToLower(k.SkipPing) != "y" {
				foundHost = ping(k.Host)
			} else {
				foundHost = true
			}

			if foundHost {
				// copy remote to local
				// k.Path+projectName+".auto", "-o", projectName+".auto"
				fmt.Println("  scp'ing ", k.Path+projectName+".auto")
				fmt.Println("    to ", projectName+".auto")
				out(apps.ScpExe, k.Path+projectName+".auto", projectName+".auto")
				fmt.Printf("Read And Output content: |%v|\n", projectName+".auto")
				readAndOutput(projectName + ".auto")
				err := os.Remove(projectName + ".auto")
				if err != nil {
					return err
				}
			} else if !foundHost {
				fmt.Println("  scp not configured")
			}
		} else {
			fmt.Println("  scp not configured")
		}
	}
	return nil
}
func scpCopy(projectName string) error {
	for _, k := range scpS.Instance {
		fmt.Println("SCP... ")
		if len(k.Host) > 0 {
			foundHost := false
			if strings.ToLower(k.SkipPing) != "true" && strings.ToLower(k.SkipPing) != "y" {
				foundHost = ping(k.Host)
			} else {
				foundHost = true
			}

			if foundHost {
				matches := findFiles(projectName)
				if len(matches) == 0 {
					fmt.Println("  no files to transfer")
				}
				t := time.Now()
				unxTime := t.Unix()
				for _, d := range matches {
					if k.Backup {
						fmt.Println("  BACKUP started... ")
						// pull backup here
						bkupDir := filepath.Join(filepath.Dir(d), "backup")
						if err := os.MkdirAll(bkupDir, 0700); err != nil && !os.IsExist(err) {
							fmt.Printf("failed to create %q: %v", bkupDir, err)
							os.Exit(1)
						}
						bkupSCPDir := filepath.Join(bkupDir, "scp")
						if err := os.MkdirAll(bkupSCPDir, 0700); err != nil && !os.IsExist(err) {
							fmt.Printf("failed to create %q: %v", bkupSCPDir, err)
							os.Exit(1)
						}
						iout.WriteOut([]byte(strconv.Itoa(int(unxTime))+"\n"+t.Format(time.RFC3339)+"\n"), filepath.Join(bkupSCPDir, "last_timestamp.txt"))
						fn := filepath.Base(d)
						fmt.Printf("**** BACKUP TO \n%v\nfrom\n%v\n", bkupSCPDir, k.Path+fn)
						out(apps.ScpExe, k.Path+fn, bkupSCPDir)
						fmt.Println("  BACKUP complete. ")
					}
					fmt.Println("  scp'ing ", d)
					fmt.Println("    to ", k.Path)
					if !dryRun {
						out(apps.ScpExe, d, k.Path)
					} else {
						fmt.Println("DRY_RUN: scp " + d + " to " + k.Path)
					}
				}
			} else if !foundHost {
				fmt.Println("  scp not configured")
			}
		} else {
			fmt.Println("  scp not configured")
		}
	}
	return nil
}

func customCopy(projectName string) error {
	for _, k := range pushCustom.Instance {
		fmt.Println("\nPush Custom... ")

		matches := findFiles(projectName)
		if len(matches) == 0 {
			fmt.Println("  no files to transfer")
		}
		for _, d := range matches {
			f := filepath.Base(d)
			var oldformat bool
			exe := "" + k.Exec
			fmt.Printf("Exec %v, Parameters: %v, PassValues: %v\n", k.Exec, k.Parameters, k.PassValues)
			if k.PassValues == 1 {
				fmt.Println("NEW FORMAT")
				exe += " " + k.Parameters + " " + d
			} else if k.PassValues == 2 {
				fmt.Println("NEW FORMAT")
				exe += " " + k.Parameters + " " + d + " " + f
			} else {
				fmt.Println("OLD FORMAT")
				oldformat = true
				// old format
				exe += " " + d + " " + f
			}
			fmt.Printf("Exe: |%v|\n", exe)

			params := make([]string, 0)
			if len(k.Parameters) > 0 {
				parts := strings.Split(k.Parameters, " ")
				for _, q := range parts {
					params = append(params, strings.TrimSpace(q))
				}
			}
			if !oldformat {
				if k.PassValues == 1 && len(k.Parameters) > 0 {
					params = append(params, d)
					out(k.Exec, params...)
				} else if k.PassValues == 1 && len(k.Parameters) == 0 {
					out(k.Exec, d)
				} else if k.PassValues == 2 && len(k.Parameters) > 0 {
					params = append(params, d)
					params = append(params, f)
					out(k.Exec, params...)
				} else if k.PassValues == 2 && len(k.Parameters) == 0 {
					params = append(params, d)
					params = append(params, f)
					out(k.Exec, d, f)
				}
			} else {
				out(k.Exec, d, f)
			}

		}
	}
	return nil
}

func customPull(projectName string) error {
	for _, k := range pullCustom.Instance {
		fmt.Println("\nPull Custom... ")

		var oldformat bool
		exe := "" + k.Exec
		fmt.Printf("Exec %v, Parameters: %v, PassValues: %v\n", k.Exec, k.Parameters, k.PassValues)
		if k.PassValues == 1 {
			fmt.Println("NEW FORMAT")
			exe += " " + k.Parameters + " " + k.PrefixPath + projectName + ".auto"
		} else if k.PassValues == 2 {
			fmt.Println("NEW FORMAT")
			exe += " " + k.Parameters + " " + k.PrefixPath + projectName + ".auto" + " ."
		} else {
			fmt.Println("OLD FORMAT")
			oldformat = true
			// old format
			exe += " " + k.PrefixPath + projectName + ".auto" + " " + "."
		}
		fmt.Printf("Exe: |%v|\n", exe)

		params := make([]string, 0)
		if len(k.Parameters) > 0 {
			parts := strings.Split(k.Parameters, " ")
			for _, q := range parts {
				params = append(params, strings.TrimSpace(q))
			}
		}
		if !oldformat {
			if k.PassValues == 1 && len(k.Parameters) > 0 {
				params = append(params, k.PrefixPath+projectName+".auto")
				out(k.Exec, params...)
			} else if k.PassValues == 1 && len(k.Parameters) == 0 {
				out(k.Exec, k.PrefixPath+projectName+".auto")
			} else if k.PassValues == 2 && len(k.Parameters) > 0 {
				params = append(params, k.PrefixPath+projectName+".auto")
				params = append(params, ".")
				out(k.Exec, params...)
			} else if k.PassValues == 2 && len(k.Parameters) == 0 {
				params = append(params, k.PrefixPath+projectName+".auto")
				params = append(params, ".")
				out(k.Exec, k.PrefixPath+projectName+".auto", ".")
			}
		} else {
			out(k.Exec, k.PrefixPath+projectName+".auto", ".")
		}

		fmt.Printf("Read And Output content: |%v|\n", projectName+".auto")
		readAndOutput(projectName + ".auto")
		err := os.Remove(projectName + ".auto")
		if err != nil {
			return err
		}
	}
	return nil
}
func sftpCopyAutoStatus(projectName string) error {
	for _, k := range sftpS.Instance {
		fmt.Println("\nSFTP Pull... ")
		if len(k.Host) > 0 {
			foundHost := false
			if strings.ToLower(k.SkipPing) != "true" && strings.ToLower(k.SkipPing) != "y" {
				foundHost = ping(k.Host)
			} else {
				foundHost = true
			}

			if foundHost {
				exe := "echo get " + k.Path + projectName + ".auto" + " " + projectName + ".auto" + " | " + apps.SftpExe + " " + k.Host
				fmt.Printf("Exe: |%v|\n", exe)

				var errorBuffer bytes.Buffer
				var errorBuffer2 bytes.Buffer
				c1 := exec.Command("echo", "get", k.Path+projectName+".auto", projectName+".auto")
				c2 := exec.Command(apps.SftpExe, k.Host)
				c1.Stderr = &errorBuffer
				c2.Stderr = &errorBuffer2
				pr, pw := io.Pipe()
				c1.Stdout = pw
				c2.Stdin = pr

				var b2 bytes.Buffer
				c2.Stdout = &b2

				err := c1.Start()
				if err != nil {
					log.Printf("err: %v\n%v", err, errorBuffer.String())
				}
				err = c2.Start()
				if err != nil {
					log.Printf("err: %v\n%v", err, errorBuffer2.String())
				}
				err = c1.Wait()
				if err != nil {
					log.Printf("err: %v\n%v", err, errorBuffer.String())
				}
				err = pw.Close()
				if err != nil {
					log.Printf("err :%v\n", err)
				}
				err = c2.Wait()
				if err != nil {
					log.Printf("err :%v\n%v", err, errorBuffer2.String())
				}
				_, err = io.Copy(os.Stdout, &b2)
				if err != nil {
					log.Printf("err :%v\n%v", err, errorBuffer2.String())
				}

				fmt.Printf("Read And Output content: |%v|\n", projectName+".auto")
				readAndOutput(projectName + ".auto")
				err = os.Remove(projectName + ".auto")
				if err != nil {
					return err
				}
			} else if !foundHost {
				fmt.Println("  sftp not configured")
			}
		} else {
			fmt.Println("  sftp not configured")
		}
	}
	return nil
}

func sftpCall(name, hostname string, args ...string) {
	var errorBuffer bytes.Buffer
	var errorBuffer2 bytes.Buffer
	c1 := exec.Command(name, args...)
	c2 := exec.Command(apps.SftpExe, hostname)
	c1.Stderr = &errorBuffer
	c2.Stderr = &errorBuffer2
	pr, pw := io.Pipe()
	c1.Stdout = pw
	c2.Stdin = pr

	var b2 bytes.Buffer
	c2.Stdout = &b2

	err := c1.Start()
	if err != nil {
		log.Printf("err: %v\n%v", err, errorBuffer.String())
	}
	err = c2.Start()
	if err != nil {
		log.Printf("err: %v\n%v", err, errorBuffer2.String())
	}
	err = c1.Wait()
	if err != nil {
		log.Printf("err: %v\n%v", err, errorBuffer.String())
	}
	err = pw.Close()
	if err != nil {
		log.Printf("err :%v\n", err)
	}
	err = c2.Wait()
	if err != nil {
		log.Printf("err :%v\n%v", err, errorBuffer2.String())
	}
	_, err = io.Copy(os.Stdout, &b2)
	if err != nil {
		log.Printf("err :%v\n%v", err, errorBuffer2.String())
	}
}
func sftpCopy(projectName string) error {
	for _, k := range sftpS.Instance {
		fmt.Println("SFTP... ")
		if len(k.Host) > 0 {
			foundHost := false
			if strings.ToLower(k.SkipPing) != "true" && strings.ToLower(k.SkipPing) != "y" {
				foundHost = ping(k.Host)
			} else {
				foundHost = true
			}

			if foundHost {
				matches := findFiles(projectName)
				if len(matches) == 0 {
					fmt.Println("  no files to transfer")
				}
				if k.Backup && len(matches) > 0 {
					fmt.Println("  BACKUP started... ")
					t := time.Now()
					unxTime := t.Unix()
					// pull backup here
					bkupDir := filepath.Join(filepath.Dir(matches[0]), "backup")
					if err := os.MkdirAll(bkupDir, 0700); err != nil && !os.IsExist(err) {
						fmt.Printf("failed to create %q: %v", bkupDir, err)
						os.Exit(1)
					}
					bkupSFTPDir := filepath.Join(bkupDir, "sftp")
					if err := os.MkdirAll(bkupSFTPDir, 0700); err != nil && !os.IsExist(err) {
						fmt.Printf("failed to create %q: %v", bkupSFTPDir, err)
						os.Exit(1)
					}
					iout.WriteOut([]byte(strconv.Itoa(int(unxTime))+"\n"+t.Format(time.RFC3339)+"\n"), filepath.Join(bkupSFTPDir, "last_timestamp.txt"))
					sftpCall("echo", k.Host, "get", k.Path+projectName+"-*", bkupSFTPDir)
					fmt.Println("  BACKUP complete. ")
				}
				for _, d := range matches {
					f := filepath.Base(d)
					exe := "echo put " + d + " " + k.Path + f + " | " + apps.SftpExe + " " + k.Host
					fmt.Printf("Exe: |%v|\n", exe)

					sftpCall("echo", k.Host, "put", d, k.Path+f)
				}
			} else if !foundHost {
				fmt.Println("  sftp not configured")
			}
		} else {
			fmt.Println("  sftp not configured")
		}
	}
	return nil
}

func artifactoryPullAutoStatus(projectName string) error {
	fmt.Println("\nArtifactory Pull... ")
	for _, k := range arts.Instance {
		fmt.Println("  processing ", k.Host)
		if len(k.Host) > 0 {
			foundHost := ping(k.Host)
			fmt.Println("found? ", foundHost)
			if foundHost && (len(k.Creds) > 0 || len(k.CredsPath) > 0) {
				var creds []byte
				if len(k.Creds) > 0 {
					creds = loadArtifactoryCreds(k.Creds)
				} else if len(k.CredsPath) > 0 {
					credsPath := readAndOutput(k.CredsPath)
					creds = loadArtifactoryCreds(credsPath)
				}
				//fmt.Printf("to pull: %v\n", k.Path+projectName+".auto")
				out(apps.CurlExe, "-u"+string(creds), "-sS", k.Path+projectName+".auto", "-o", projectName+".auto")
				fmt.Printf("Read And Output content: |%v|\n", projectName+".auto")
				readAndOutput(projectName + ".auto")
				err := os.Remove(projectName + ".auto")
				if err != nil {
					return err
				}
			} else if foundHost && len(k.Creds) == 0 {
				fmt.Println("  no artifactory credentials found")
			} else if !foundHost {
				fmt.Println("  artifactory not configured")
			}
		} else {
			fmt.Println("  artifactory not configured")
		}
	}
	return nil
}
func artifactoryPush(projectName string) error {
	fmt.Println("Artifactory... ")
	t := time.Now().Unix()
	for _, k := range arts.Instance {
		fmt.Println("  processing ", k.Host)
		if len(k.Host) > 0 {
			foundHost := ping(k.Host)
			fmt.Println("found? ", foundHost)
			if foundHost && (len(k.Creds) > 0 || len(k.CredsPath) > 0) {
				var creds []byte
				if len(k.Creds) > 0 {
					creds = loadArtifactoryCreds(k.Creds)
				} else if len(k.CredsPath) > 0 {
					credsPath := readAndOutput(k.CredsPath)
					creds = loadArtifactoryCreds(credsPath)
					//fmt.Printf("Creds found %v\n", string(creds))
				}
				var byt bytes.Buffer
				matches := findFiles(projectName)
				byt.WriteString("{")
				byt.WriteString(strings.Join(matches, ","))
				byt.WriteString("}")
				log.Println(byt.String())

				if byt.Len() > 2 {
					// build hash to upload
					for _, d := range matches {
						if k.Backup {
							fmt.Println("  BACKUP started... ")
							// pull backup here
							bkupDir := filepath.Join(filepath.Dir(d), "backup")
							if err := os.MkdirAll(bkupDir, 0700); err != nil && !os.IsExist(err) {
								fmt.Printf("failed to create %q: %v", bkupDir, err)
								os.Exit(1)
							}
							fn := filepath.Base(d)
							svFile := filepath.Join(bkupDir, fn) + ".af." + strconv.Itoa(int(t))
							fmt.Printf("**** BACKUP TO \n%v\nfrom\n%v\n", svFile, k.Path+fn)
							out(apps.CurlExe, "-u"+string(creds), "-sS", "-o", svFile, k.Path+fn)
							fmt.Println("  BACKUP complete. ")
						}
						fmt.Println("  pushing via artifactory ", d)
						fmt.Println("    to ", k.Path)
						// hash each before uploading
						if !dryRun {
							md5sum, err := sh.Output(apps.MD5Exe, d)
							if err != nil {
								return err
							}
							md5parts := strings.Fields(md5sum)
							md5sum = md5parts[0]
							sha1sum, err := sh.Output(apps.SHA1Exe, d)
							if err != nil {
								return err
							}
							shaParts := strings.Fields(sha1sum)

							shasum256, err := sh.Output(apps.SHA256Exe, d)
							if err != nil {
								return err
							}
							sha256Parts := strings.Fields(shasum256)

							out(apps.CurlExe, "-u"+string(creds), "-sS", "-T", d, "-H", "X-Checksum-MD5:"+md5sum, "-H", "X-Checksum-Sha1:"+shaParts[0], "-H", "X-Checksum-Sha256:"+sha256Parts[0], k.Path)
						}
					}
					// Upload all at once without hashes
					//out(curlExe, "-u"+string(artifactoryCreds), "-T", byt.String(), artifactoryPath)
				}
			} else if foundHost && len(k.Creds) == 0 && len(k.CredsPath) == 0 {
				fmt.Println("  no artifactory credentials found")
			} else if !foundHost {
				fmt.Println("  artifactory not configured")
			}
		} else {
			fmt.Println("  artifactory not configured")
		}
	}

	return nil
}

func findFiles(projectName string) []string {
	dir := filepath.Join(baseDir, projectName)
	matches, err := filepath.Glob(dir + "*")
	if err != nil {
		log.Println(err)
	}
	return matches
}
func out(cmd string, args ...string) {
	outResp, err := sh.Output(cmd, args...)
	if err != nil {
		log.Println("err:", err)
	}
	if len(outResp) > 0 {
		log.Println("Out:", outResp)
	}
}
func ping(host string) bool {
	err := sh.Run("ping", "-n", "-c", "1", host)
	if err != nil {
		log.Println("error:", err)
		return false
	}

	return true
}

type update struct {
	Os        string `json:"os"`
	Arch      string `json:"arch"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
	Changelog string `json:"changelog"`
}

func buildUpdateDir(u update) ([]byte, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func buildDeployScript(scriptPath string) []byte {
	f, err := os.Open(scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	// replace $name with name of project
	//s := strings.Replace(string(b), "$name", name, -1)
	return b
}

// Format your go code
func Format() error {
	fmt.Println("Formatting...")
	// format simplify and list files whose formatting differs, exclude '.pb.go:', exclude vendor
	err := sh.RunV("gofmt", "-s", "-l", ".")
	return err
}

// Perform Lint checks on your project
func Lint() error {
	fmt.Println("Lint Checks...")

	err := sh.RunV("golint", "./...")
	return err
}

// Run Tests on your project
func Test() error {
	fmt.Println("Testing...")

	gocmd := mg.GoCmd()
	err := sh.RunV(gocmd, "test", "-v", "-tags", buildTags+" cgo", "./...")

	return err
}

// Vet your code
func Vet() error {
	fmt.Println("Vet'ting...")

	gocmd := mg.GoCmd()
	err := sh.RunV(gocmd, "vet", "./...")
	if err != nil {
		log.Println("go vet found issues")
	}
	return nil
}

func Install() error {
	mg.SerialDeps(Build)
	fmt.Println("Installing...")
	gocmd := mg.GoCmd()

	skipUPX := false
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	fmt.Println("goos: " + goos + " goarch: " + goarch)
	if goos == "darwin" && (goarch == "arm64" || goarch == "amd64") {
		skipUPX = true
	}
	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		var err error
		projectMainDir := "./cmd/" + d.Name + "/"
		if !exists(projectMainDir) {
			fmt.Println("Path doesn't exist: " + projectMainDir + " using local dir '.' instead")
			projectMainDir = "."
		} else {
			projectMainDir = projectMainDir + "."
		}
		err = setup(d)
		if err != nil {
			fmt.Println("issue with setup :", err)
		}
		if !dryRun {
			err = sh.RunV(gocmd, "install", "-a", "-tags", buildTags, "-ldflags", goLDFlags, projectMainDir)
			if err != nil {
				fmt.Println("!!!error: ", err)
			}
		} else {
			fmt.Println("DRY_RUN: " + gocmd + " install -a -tags " + buildTags + " -ldflags " + goLDFlags + " " + projectMainDir)
		}
		goPath := os.Getenv("GOPATH")
		binPath := filepath.Join(goPath, "/bin/", d.Name)
		if exists(binPath) {
			if exists(apps.UPXExe) && !skipUPX {
				fmt.Printf("\n*** START UPX binary compression on  %v ***\n", binPath)
				fi, _ := os.Stat(binPath)
				fmt.Printf("\n")
				err = sh.RunV(apps.UPXExe, "-q", "-q", "-q", binPath)
				if err != nil {
					return err
				}
				fmt.Printf("\n")
				fmt.Printf("    - prior to compression: %v\n", stringut.HRByteCount(fi.Size(), false))
				fi, _ = os.Stat(binPath)
				fmt.Printf("    - post compression: %v\n", stringut.HRByteCount(fi.Size(), false))
				fmt.Printf("\n*** END UPX binary compression on  %v ***\n", binPath)
			} else if skipUPX {
				fmt.Println("- upx not applicable for darwin/arm64 - ")
			} else {
				fmt.Println("- no upx available for binary compression - ")
			}
		}
		err = os.Setenv("PROG", binPath)
		if err != nil {
			fmt.Println("issue setting PROG :", binPath, err)
		}
		bashCompletionPath := "/usr/local/etc/bash_completion.d/"
		if goos == "darwin" && goarch == "arm64" {
			bashCompletionPath = "/opt/homebrew/etc/bash_completion.d/"
		} else if goos == "linux" {
			bashCompletionPath = "/etc/bash_completion.d/"
		}
		if !dryRun && goos != "windows" {
			err = sh.Copy(bashCompletionPath+d.Name, "./pkgr/bash_autocomplete")
			if err != nil {
				fmt.Println("!!!error: ", err)
			}
		} else if goos == "windows" {
			fmt.Println("No bash completion on windows skipping setup.")
		} else {
			fmt.Println("DRY_RUN copying ./pkgr/bash_autocomplete to " + bashCompletionPath + d.Name)
		}
		cleaner(d.Name, false)
	}

	return nil
}

func Clean() {
	mg.SerialDeps(parseToml)
	for _, d := range prjkts.Projects {
		if !d.Enable {
			continue
		}
		cleaner(d.Name, false)
	}
}

// Clean up after yourself
func cleaner(projectName string, dirsOnly bool) {
	fmt.Println("Cleaning Post Clean Dirs...")
	for _, d := range postclean.Dirs {
		if exists(d) {
			fmt.Println("  Cleaning...", d)
			err := os.RemoveAll(d)
			if err != nil {
				log.Printf("while remove directories: %v\n", err)
			}
		}
	}

	for _, d := range toCleanFiles {
		fmt.Println("  Cleaning created files/folders...", d)
		matches := findFiles(projectName)
		if !dirsOnly {
			fmt.Println("finding files starting with project name prefix")
			for _, k := range matches {
				fmt.Println("  Cleaning...", k)
				err := os.Remove(k)
				if err != nil {
					log.Printf("while removing files (matches): %v\n", err)
				}
			}
		}
		if exists(d) {
			err := os.Remove(d)
			if err != nil {
				log.Printf("while removing files: %v\n", err)
			}
		}
	}
}

// *********** ONE OFF TASKS BELOW

// Setup preTask
func setup(app Project) error {
	mg.SerialDeps(parseToml)
	fmt.Println("  - retrieve version")
	ver := version(app.VersionFile)
	fmt.Printf("  - version found %s\n", ver)
	fmt.Println("  - retrieve git commit hash")
	gitCommit := gitCommitHash()
	gitBranch := gitbranch()

	fmt.Println("  - retrieve current path")
	cur := currentPath()
	baseDir = cur

	prepDir = filepath.Join(baseDir, "PREP")
	//app.Package += "/" + app.Name
	buildDir = filepath.Join(baseDir, crossBuildDir)
	log.Printf("  - Reading\n   - Project: %s \n     Project Pkg: %s\n     CurrentPath: %s", app.Name, app.Package, baseDir)

	fmt.Println("  - setup ldflag version templates")
	versionFields := fmt.Sprintf(versionFieldsTemplate, versionPkg, gitCommit, versionPkg, gitBranch, versionPkg, ver, versionPkg, strconv.FormatInt(timestamp, 10), versionPkg, goVersion())
	fmt.Println("    Version Fields: ", versionFields)
	if len(app.OverrideVariables) > 0 {
		overrides := strings.Split(app.OverrideVariables, ";")
		for i, k := range overrides {
			val := fmt.Sprintf(k, overwriteValues[i])
			fmt.Printf("IDX: %d. VARIABLE: %v VALUE: %v - together %v \n", i, k, overwriteValues[i], val)
			versionFields += " " + val
		}
		fmt.Printf("version fields: %v\n", versionFields)
	}
	goLDFlags = fmt.Sprintf(goLDFlagsTemplate, versionFields)
	goLDFlagsStatic = fmt.Sprintf(goLDFlagsStaticTemplate, versionFields)

	fmt.Println("  - load artifactory creds")

	return nil
}

func currentPath() string {
	cur, _ := filepath.Abs(".")
	return cur
}
func version(versionFile string) string {
	ver, err := sh.Output("cat", versionFile)
	if err != nil {
		log.Println(err)
	}
	log.Println("    reading ", versionFile, "found", ver)
	return ver
}

func readAndOutput(targetFile string) string {
	if !exists(targetFile) {
		fmt.Printf("file not found: %v\n", targetFile)
	}
	content, err := sh.Output("cat", targetFile)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("content: %v\n", content)
	return content
}

func gitCommitHash() string {
	gitCommit := hash()

	gitCommit += gitStatus()

	log.Println("    ** GIT HASH:", gitCommit)

	return gitCommit
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hashResp, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hashResp
}

// gitbranch returns the git branch for the current repo or "" if none.
func gitbranch() string {
	// git branch --format='%(refname:short)'
	branchResp, _ := sh.Output("git", "branch", "--format=%(refname:short)")
	return branchResp
}

func gitStatus() string {
	s, _ := sh.Output("git", "status", "--porcelain", "--untracked-files=no")
	if len(s) > 0 {
		return "-dirty"
	}
	return ""
}

func goVersion() string {
	resp, _ := sh.Output("go", "version")
	flds := strings.Fields(resp)
	return flds[2]
}
func loadArtifactoryCreds(path string) []byte {
	if len(path) > 0 {
		creds, _ := sh.Output("cat", path)
		return []byte(creds)
	}
	return nil
}

func convertInterfaceArToStringAr(data []interface{}) []string {
	tmp := make([]string, 0)
	for _, d := range data {
		tmp = append(tmp, d.(string))
	}
	return tmp
}

func convertOldToGenConf(props map[string]interface{}, key string) GenConfig {
	var tmp GenConfig
	mapProps := props[key]
	wrapper := make(map[string]interface{}, 1)
	wrapper[key] = mapProps
	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		log.Fatalf("issue parsing %v %v", key, err)
	}
	err = json.Unmarshal(bytesWrapper, &tmp)
	if err != nil {
		log.Fatalf("issue marshalling %v %v", key, err)
	}
	return tmp
}

// Used with Ping on SCP, host must be available via ping check, requires ppk setup
func PPK() error {
	fmt.Println("Building PPK...")

	ppk := genppk.PPK{PrivateFilename: "mypriv", PublicFilename: "mypub"}
	ppk.GenerateKeys(0)
	ppk.SavePrivateKeyAsPEM()
	ppk.SavePublicKeyAsPEM()
	fmt.Println("Finished Building PPK...")

	// Generate .go file with data in it
	/*
		1. check for key.go
		2. create if it doesn't exist for the project
		3. create signature with private key and place in update file
		4. on update verify sig with public key in project

		var publicKey = []byte(`
		-----BEGIN PUBLIC KEY-----
		MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEtrVmBxQvheRArXjg2vG1xIprWGuCyESx
		MMY8pjmjepSy2kuz+nl9aFLqmr+rDNdYvEBqQaZrYMc6k29gjvoQnQ==
		-----END PUBLIC KEY-----
		`)
	*/

	return nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
func fileExistsAndIsNotADir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func findExec(app string) (string, error) {
	path, err := exec.LookPath(app)
	if err != nil {
		return "", fmt.Errorf("not found on path %v", app)
	}
	return path, nil
}
