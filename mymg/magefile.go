package mymg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	toml "github.com/pelletier/go-toml"

	"github.com/colt3k/utils/crypt/genppk"
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
	dryRun    bool
	apps      applications
	arts      artifactories
	scpS      scps
	sftpS     sftps
	scpCustom scpcustoms
	timestamp = time.Now().Unix()
	baseDir   = ""
	buildDir  = ""
	prepDir   = ""

	versionPkg              = "github.com/colt3k/utils"
	versionFieldsTemplate   = `-X "%s/version.GITCOMMIT=%s" -X "%s/version.VERSION=%s" -X "%s/version.BUILDDATE=%s" -X "%s/version.GOVERSION=%s"`
	saltOverwriteValue      = ""
	goLDFlagsTemplate       = "-s -w %s"
	goLDFlags               string
	goLDFlagsStaticTemplate = "-s -w %s -extldflags -static"
	goLDFlagsStatic         string
	bump                    bool
	names                   []string
	prompt                  bool
	nostatic                bool

	buildTags     = ""
	crossBuildDir = "cross"

	toCleanFiles []string
	toCleanDirs  []string

	md5Exe    = "/bin/md5sum"
	sha1Exe   = "/bin/sha1sum"
	sha256Exe = "/bin/sha256sum"
	curlExe   = "/bin/curl"
	catExe    = "/bin/cat"
	gitExe    = "/bin/git"
	tarExe    = "/bin/tar"
	scpExe    = "/bin/scp"
	sftpExe   = "/bin/sftp"
	whichExe  = "/usr/bin/which"

)

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
	log.Println("Scps Obj:", scpS)
	return nil
}
func setupCustomScps(props map[string]interface{}) error {
	mapProps := props["scp-custom"]
	wrapper := make(map[string]interface{}, 1)
	wrapper["scp-custom"] = mapProps

	bytesWrapper, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesWrapper, &scpCustom)
	if err != nil {
		return err
	}
	log.Println("Custom Scps Obj:", scpCustom)

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
	log.Println("Sftps Obj:", sftpS)
	return nil
}
func setupArtifacts(props map[string]interface{}) error {
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
	log.Println("Artifacts Obj:", arts)
	return nil
}
func setupApps(props map[string]interface{}) error {
	appMap := props["application"]
	appWrapper := make(map[string]interface{}, 1)
	appWrapper["application"] = appMap

	bytesAppWrapper, err := json.MarshalIndent(appWrapper, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesAppWrapper, &apps)
	if err != nil {
		return err
	}

	sz := len(apps.Apps)
	if sz == 1 {
		prompt = false
	}
	fmt.Printf("Apps available %d\n\n", sz)
	// find absolute paths for files
	for i, d := range apps.Apps {
		if prompt && sz > 1 && ques.Confirm("process ("+d.Name+")? ") {
			apps.Apps[i].Enable = true
		} else if !prompt && sz == 1 {
			// if only one app enabled just run it, no confirmation needed
			apps.Apps[i].Enable = true
		}
		if !prompt && len(names) > 0 {
			for _, k := range names {
				if d.Name == k {
					apps.Apps[i].Enable = true
				}
			}
		}

		apps.Apps[i].VersionFile, err = filepath.Abs(d.VersionFile)
		if err != nil {
			return err
		}

		apps.Apps[i].ReadmeFile, err = filepath.Abs(d.ReadmeFile)
		if err != nil {
			return err
		}

		apps.Apps[i].ChangelogFile, err = filepath.Abs(d.ChangelogFile)
		if err != nil {
			return err
		}

		for j, k := range d.Files {
			apps.Apps[i].Files[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}

		for j, k := range d.OSDeployScripts {
			apps.Apps[i].OSDeployScripts[j], err = filepath.Abs(k)
			if err != nil {
				return err
			}
		}

		if len(d.SaltVariableOverwrite) > 0 {
			// prompt for value
			saltOverwriteValue = ques.Question("salt value ? ")
		}

		if len(d.YNPrompt) > 0 {
			prompt = ques.Confirm(d.YNPrompt)
			if !prompt {
				fmt.Println("Please pull first!!!, Exiting...")
				os.Exit(1)
			}
		}
	}
	log.Println("Apps Obj:", apps)
	return nil
}

func setupCleanFiles() {
	toCleanFiles = make([]string, 0)

	for _, d := range apps.Apps {
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
	path, err = filepath.Abs(path)

	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	props := tree.ToMap()

	err = setupApps(props)
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
	path, err = filepath.Abs(path)

	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	props := tree.ToMap()

	// APPS
	md5Exe = props["md5Exe"].(string)
	sha1Exe = props["sha1Exe"].(string)
	sha256Exe = props["sha256Exe"].(string)
	curlExe = props["curlExe"].(string)
	catExe = props["catExe"].(string)
	gitExe = props["gitExe"].(string)
	tarExe = props["tarExe"].(string)
	scpExe = props["scpExe"].(string)
	if props["sftpExe"] != nil {
		sftpExe = props["sftpExe"].(string)
	}
	err = setupScps(props)
	if err != nil {
		return err
	}
	err = setupCustomScps(props)
	if err != nil {
		return err
	}
	err = setupSftps(props)
	if err != nil {
		return err
	}
	err = setupArtifacts(props)
	if err != nil {
		return err
	}
	err = setupApps(props)
	if err != nil {
		return err
	}
	var processApp bool
	fmt.Println("") // clear line output
	for _, d := range apps.Apps {
		fmt.Println("App ", d.Name, "Enabled? ", d.Enable)
		if d.Enable {
			processApp = true
		}
	}
	if !processApp {
		fmt.Println("\nNo application selected to process.")
		os.Exit(-1)
	}
	setupCleanFiles()

	// Cleaning
	toCleanDirs = convertInterfaceArToStringAr(props["to_clean_dirs"].([]interface{}))

	// Building
	buildTags = props["build_tags"].(string)

	log.Println("toCleanFiles: ", toCleanFiles)
	log.Println("toCleanDirs: ", toCleanDirs)

	return nil
}

type scps struct {
	Instance []scp `json:"scp"`
}
type scp struct {
	Host     string `json:"host"`
	Path     string `json:"path"`
	SkipPing string `json:"skip_ping"`
}
type scpcustoms struct {
	Instance []scpcust `json:"scp-custom"`
}
type scpcust struct {
	Exec string `json:"exec"`
}
type sftps struct {
	Instance []sftp `json:"sftp"`
}
type sftp struct {
	Host     string `json:"host"`
	Path     string `json:"path"`
	SkipPing string `json:"skip_ping"`
}
type artifactories struct {
	Instance []artifactory `json:"artifactory"`
}
type artifactory struct {
	Host  string `json:"host"`
	Path  string `json:"path"`
	Creds string `json:"creds"`
}

type applications struct {
	Apps []application `json:"application"`
}

type application struct {
	Enable                bool     `json:"enable"`
	Name                  string   `json:"name"`
	OSTargets             []string `json:"ostargets"`
	OSDeployScripts       []string `json:"osdeployscripts"`
	Package               string   `json:"package"`
	ReadmeFile            string   `json:"readme"`
	VersionFile           string   `json:"version"`
	ChangelogFile         string   `json:"changelog"`
	Files                 []string `json:"files"`
	SaltVariableOverwrite string   `json:"salt_variable_overwrite"`
	YNPrompt              string   `json:"ynprompt"`
}

func Help() {
	fmt.Println()
	fmt.Println("Mage Usage:")
	fmt.Println("Global Parameters")
	fmt.Println("  dry=    Dry Run")
	fmt.Println("  bump=	Bump to next version")
	fmt.Println("  help=	Show this help")
	fmt.Println("  names=	Single or Comma separated list of targets to build")
	fmt.Println("  config=	Location of build.toml default is ./build.toml")
	fmt.Println("  nostatic= Do not set static flag for build")
	fmt.Println("Targets")
	fmt.Println("  help         show this help information")
	fmt.Println("  install      install the application to local machine")
	fmt.Println("  buildCross   create based on current local code and don't clean up")
	fmt.Println("  release      create and push release based on current local code")
	fmt.Println("  targets      show current project configured command's to build")
	fmt.Println("  auto         build true auto file and release")
	fmt.Println("  noauto       build false auto file and release")
	fmt.Println("Flags")
	fmt.Println("  -v		show verbose mage output")
	fmt.Println("  -d		show custom debug output")
	fmt.Println()

}

func Targets() {
	mg.SerialDeps(parseTargets)
	fmt.Println()
	fmt.Println("Targets")
	for _, d := range apps.Apps {
		fmt.Println("  ", d.Name)
	}

	fmt.Println()
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.SerialDeps(parseToml)

	gocmd := mg.GoCmd()
	fmt.Println("Building...")
	//$(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME) .
	var err error
	if _, ok := os.LookupEnv("CGO"); !ok {
		err = os.Setenv("CGO", "0")
		if err != nil {
			fmt.Println("issue setting CGO to 0:", err)
		}
	}

	for _, d := range apps.Apps {

		if !d.Enable {
			continue
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

		if !dryRun {
			err = sh.RunV(gocmd, "build", "-tags", buildTags, "-ldflags", goLDFlags, "-o", name, "./cmd/"+d.Name+"/.")
			if err != nil {
				return err
			}
		} else {
			var byt bytes.Buffer
			byt.WriteString(gocmd + " build -tags " + buildTags + " -ldflags " + goLDFlags + " -o " + name + " ./cmd/" + d.Name + "/.")
			fmt.Println("DRY_RUN: Building build", byt.String())
		}

	}

	return nil
}

func BumpVersion() error {
	mg.Deps(parseToml)
	gocmd := mg.GoCmd()

	if !dryRun {
		_, err := sh.Output(whichExe, "sembump")
		if err != nil {
			//update if not found
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

	for _, d := range apps.Apps {
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
				out2, err = sh.Output(catExe, d.ReadmeFile)
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
				err = sh.RunV(gitExe, "add", d.VersionFile, d.ReadmeFile)
				if err != nil {
					fmt.Println("issue on git add bump:", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + gitExe + " add " + d.VersionFile + " " + d.ReadmeFile)
			}

			if !dryRun {
				// Commit tag
				err = sh.RunV(gitExe, "commit", "-vsam", "Bump version to "+nVersion)
				if err != nil {
					fmt.Println("issue committing bump:", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + gitExe + " commit -vsam Bump version to " + nVersion)
			}

			// if there is more than one command for this project create a unique tag for it
			if strings.Index(d.VersionFile, "/cmd/") > -1 {
				nVersion = d.Name + "/" + nVersion
			}
			if !dryRun {
				// "Run make tag to create and push the tag for new version $(NEW_VERSION)"
				err = sh.RunV(gitExe, "tag", "-a", nVersion, "-m", nVersion)
				if err != nil {
					fmt.Println("issue tagging bump :", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + gitExe + " tag -a " + nVersion + " -m " + nVersion)
			}
			// Push Tag
			if !dryRun {
				err = sh.RunV(gitExe, "push", "origin", nVersion)
				if err != nil {
					fmt.Println("issue pushing tag :", err)
				}
			} else {
				fmt.Println("DRY_RUN: " + gitExe + "push origin " + nVersion)
			}

		}
	}
	return nil
}
func BuildCross() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range apps.Apps {
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

		cleaner(d.Name, true)

	}
	return nil
}

func Release() error {
	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range apps.Apps {
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
		err = scpCustomCopy(d.Name)
		if err != nil {
			fmt.Println("issue scp custom copy :", err)
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

	for _, d := range apps.Apps {
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
		err = scpCustomCopy(d.Name)
		if err != nil {
			fmt.Println("issue scp custom copy :", err)
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

func NoAuto() error {

	mg.SerialDeps(parseToml, BumpVersion)

	for _, d := range apps.Apps {
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
		err = scpCustomCopy(d.Name)
		if err != nil {
			fmt.Println("issue scp custom copy :", err)
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

// Build for all defined Architectures
func cross(app application) error {

	fmt.Println("CrossBuilding...")
	gocmd := mg.GoCmd()

	for i, d := range app.OSTargets {
		goosArch := strings.Split(d, "/")
		goos := goosArch[0]
		arch := goosArch[1]
		arm := ""
		if len(goosArch) == 3 {
			arm = goosArch[2]
		}
		err := os.Setenv("GOOS", goos)
		if err != nil {
			fmt.Println("issue setting GOOS :", goos, err)
		}
		err = os.Setenv("GOARCH", arch)
		if err != nil {
			fmt.Println("issue setting GOARCH :", arch, err)
		}
		err = os.Setenv("GOARM", arm)
		if err != nil {
			fmt.Println("issue setting GOARM :", arm, err)
		}
		if _, ok := os.LookupEnv("CGO"); !ok {
			err = os.Setenv("CGO", "0")
			if err != nil {
				fmt.Println("issue setting CGO to 0:", err)
			}
		}

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
				return fmt.Errorf("  failed to create %q: %v", path, err)
			}
		} else {
			fmt.Println("DRY_RUN: creating directory " + path)
		}
		fmt.Printf("  Packaging %s\n", goos)
		executableName := filepath.Join(path, name)
		if !dryRun {
			if nostatic {
				goLDFlagsStatic = goLDFlags
			}
			err = sh.RunV(gocmd, "build", "-tags", buildTags, "-ldflags", goLDFlagsStatic, "-o", executableName, "./cmd/"+app.Name+"/.")
			if err != nil {
				return err
			}
		} else {
			if nostatic {
				goLDFlagsStatic = goLDFlags
			}
			fmt.Println("DRY_RUN: " + gocmd + " build -tags " + buildTags + " -ldflags " + goLDFlagsStatic + " -o " + executableName + " ./cmd/" + app.Name + "/.")
		}
		// make release dir for this OS

		osarchDir := filepath.Join(baseDir, osarchName)
		if !dryRun {
			if err = os.MkdirAll(osarchDir, 0700); err != nil && !os.IsExist(err) {
				return fmt.Errorf("  failed to create %q: %v", path, err)
			}
		} else {
			fmt.Println("DRY_RUN: making directory " + osarchDir)
		}

		if !dryRun {
			// copy executable to release dir
			err = sh.Copy(filepath.Join(osarchDir, name), executableName)
			if err != nil {
				fmt.Println("issue copying :", executableName, err)
			}
		} else {
			fmt.Println("DRY_RUN: copying " + executableName + " to " + filepath.Join(osarchDir, name))
		}
		scriptContent := buildDeployScript(app.OSDeployScripts[i], app.Name)

		//${release_dir}/deploy_${goosarch[0]}.sh
		var scriptName string
		if goos != "windows" {
			scriptName = "deploy_" + goos + ".sh"
		} else {
			scriptName = "deploy_" + goos + ".txt"
		}
		if !dryRun {
			_, err = iout.WriteOut(scriptContent, filepath.Join(osarchDir, scriptName))
			if err != nil {
				return err
			}
		} else {
			fmt.Println("DRY_RUN: writing out deploy script to " + filepath.Join(osarchDir, scriptName))
		}
		// create sha files
		if !dryRun {
			md5sum, err := sh.Output(md5Exe, executableName)
			if err != nil {
				return err
			}
			md5sumParts := strings.Fields(md5sum)

			_, err = iout.WriteOut([]byte(md5sumParts[0]), executableName+".md5")
			if err != nil {
				return err
			}
		} else {
			fmt.Println("DRY_RUN: creating and writing md5 hash file")
		}

		if !dryRun {
			var shasum string
			shasum, err = sh.Output(sha256Exe, executableName)
			if err != nil {
				return err
			}
			sha256Parts := strings.Fields(shasum)
			_, err = iout.WriteOut([]byte(sha256Parts[0]), executableName+".sha256")
			if err != nil {
				return err
			}
			err = sh.Copy(filepath.Join(osarchDir, name+".md5"), executableName+".md5")
			if err != nil {
				fmt.Println("issue copying :", executableName+".md5", err)
			}
			err = sh.Copy(filepath.Join(osarchDir, name+".sha256"), executableName+".sha256")
			if err != nil {
				fmt.Println("issue copying :", executableName+".sha256", err)
			}
		} else {
			fmt.Println("DRY_RUN: creating and writing sha256 hash file")
		}

		if !dryRun {
			// Copy change log
			nm := filepath.Base(app.ChangelogFile)
			err = sh.Copy(filepath.Join(osarchDir, nm), app.ChangelogFile)
		} else {
			fmt.Println("DRY_RUN: copying Changelog file: ", app.ChangelogFile)
		}

		// Copy all prep files into osarchDir
		var files []os.FileInfo

		files, err = ioutil.ReadDir(prepDir)
		if err != nil {
			log.Fatal(err)
		}

		if dryRun && len(files) == 0 {
			fmt.Println("DRY_RUN: no files found in " + prepDir)
		}
		for _, f := range files {
			if goos == "windows" && f.Name() == app.Name+".bash" {
				continue
			}
			fmt.Println("  Found: ", f.Name())
			// Copy here
			fileTarget := filepath.Join(osarchDir, f.Name())
			log.Println("copy ", filepath.Join(prepDir, f.Name()), " to ", fileTarget)
			if !dryRun {
				err = sh.Copy(fileTarget, filepath.Join(prepDir, f.Name()))
				if err != nil {
					return err
				}
			} else {
				fmt.Println("DRY_RUN: copy ", filepath.Join(prepDir, f.Name()), " to ", fileTarget)
			}
		}

		if !dryRun {
			fmt.Println("  Creating archive ", osarchName)
			err = sh.RunV(tarExe, "cfz", osarchName+".tgz", osarchName)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("DRY_RUN:  Creating archive ", osarchName)
		}

		v := version(app.VersionFile)
		u := update{Os: goos, Arch: arch, Name: name, Timestamp: timestamp, Version: v, Changelog: app.Name + "-changes.txt"}
		updater, err := buildUpdateDir(u)
		if err != nil {
			return err
		}
		fmt.Println("  ", string(updater))
		if !dryRun {
			_, err = iout.WriteOut(updater, filepath.Join(baseDir, osarchName+".update"))
			if err != nil {
				return err
			}
		} else {
			fmt.Println("DRY_RUN: writing out update file " + filepath.Join(baseDir, osarchName+".update"))
		}
		// Copy changelog to release dir
		if !dryRun {
			err = sh.Copy(filepath.Join(baseDir, app.Name+"-changes.txt"), app.ChangelogFile)
		} else {
			nm := filepath.Base(app.ChangelogFile)
			fmt.Println("DRY_RUN: copying Changes file " + filepath.Join(baseDir, nm))
		}

		if !dryRun {
			// remove release dir
			err = os.RemoveAll(osarchDir)
			if err != nil {
				fmt.Println("issue removing all :", osarchDir, err)
			}
		} else {
			fmt.Println("DRY_RUN: removing all files in " + osarchDir)
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
				for _, d := range matches {
					fmt.Println("  scp'ing ", d)
					fmt.Println("    to ", k.Path)
					if !dryRun {
						out(scpExe, d, k.Path)
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

func scpCustomCopy(projectName string) error {

	for _, k := range scpCustom.Instance {
		fmt.Println("SCP Custom... ")

		matches := findFiles(projectName)
		if len(matches) == 0 {
			fmt.Println("  no files to transfer")
		}
		for _, d := range matches {
			f := filepath.Base(d)
			fmt.Printf("\tpassing \n\tparameter 1 %v,\n\tparameter 2 %v\n\tto %v\n", d, f, k.Exec)
			out(k.Exec, d, f)
		}
	}
	return nil
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
				for _, d := range matches {
					f := filepath.Base(d)
					exe := "echo put " + d + " " + k.Path + f + " | " + sftpExe + " " + k.Host
					fmt.Printf("Exe: |%v|\n", exe)

					var errorBuffer bytes.Buffer
					var errorBuffer2 bytes.Buffer
					c1 := exec.Command("echo", "put", d, k.Path+f)
					c2 := exec.Command(sftpExe, k.Host)
					c1.Stderr = &errorBuffer
					c2.Stderr = &errorBuffer2
					pr, pw := io.Pipe()
					c1.Stdout = pw
					c2.Stdin = pr

					var b2 bytes.Buffer
					c2.Stdout = &b2

					err := c1.Start()
					if err != nil {
						log.Printf("err: %v\n%v", err, string(errorBuffer.Bytes()))
					}
					err = c2.Start()
					if err != nil {
						log.Printf("err: %v\n%v", err, string(errorBuffer2.Bytes()))
					}
					err = c1.Wait()
					if err != nil {
						log.Printf("err: %v\n%v", err, string(errorBuffer.Bytes()))
					}
					err = pw.Close()
					if err != nil {
						log.Printf("err :%v\n", err)
					}
					err = c2.Wait()
					if err != nil {
						log.Printf("err :%v\n%v", err, string(errorBuffer2.Bytes()))
					}
					_, err = io.Copy(os.Stdout, &b2)
					if err != nil {
						log.Printf("err :%v\n%v", err, string(errorBuffer2.Bytes()))
					}
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

func artifactoryPush(projectName string) error {
	fmt.Println("Artifactory... ")
	for _, k := range arts.Instance {
		fmt.Println("  processing ", k.Host)
		if len(k.Host) > 0 {
			foundHost := ping(k.Host)
			fmt.Println("found? ", foundHost)
			if foundHost && len(k.Creds) > 0 {
				creds := loadArtifactoryCreds(k.Creds)
				var byt bytes.Buffer
				matches := findFiles(projectName)
				byt.WriteString("{")
				byt.WriteString(strings.Join(matches, ","))
				byt.WriteString("}")
				log.Println(byt.String())

				if byt.Len() > 2 {
					// build hash to upload
					for _, d := range matches {
						fmt.Println("  pushing via artifactory ", d)
						fmt.Println("    to ", k.Path)
						// hash each before uploading
						if !dryRun {

							md5sum, err := sh.Output(md5Exe, d)
							if err != nil {
								return err
							}
							md5parts := strings.Fields(md5sum)
							md5sum = md5parts[0]
							sha1sum, err := sh.Output(sha1Exe, d)
							if err != nil {
								return err
							}
							shaParts := strings.Fields(sha1sum)

							shasum256, err := sh.Output(sha256Exe, d)
							if err != nil {
								return err
							}
							sha256Parts := strings.Fields(shasum256)

							out(curlExe, "-u"+string(creds), "-sS", "-T", d, "-H", "X-Checksum-MD5:"+md5sum, "-H", "X-Checksum-Sha1:"+shaParts[0], "-H", "X-Checksum-Sha256:"+sha256Parts[0], k.Path)
						}
					}
					// Upload all at once without hashes
					//out(curlExe, "-u"+string(artifactoryCreds), "-T", byt.String(), artifactoryPath)
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
func buildDeployScript(scriptPath, name string) []byte {

	f, err := os.Open(scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	// replace $name with name of application
	s := strings.Replace(string(b), "$name", name, -1)
	return []byte(s)

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

	for _, d := range apps.Apps {
		if !d.Enable {
			continue
		}
		var err error
		if !dryRun {
			err = sh.RunV(gocmd, "install", "-a", "-tags", buildTags, "-ldflags", goLDFlags, "./cmd/"+d.Name+"/.")
			if err != nil {
				fmt.Println("!!!error: ", err)
			}
		} else {
			fmt.Println("DRY_RUN: " + gocmd + " install -a -tags " + buildTags + " -ldflags " + goLDFlags + " ./cmd/" + d.Name + "/.")
		}
		goPath := os.Getenv("GOPATH")
		binPath := filepath.Join(goPath, "/bin/", d.Name)
		err = os.Setenv("PROG", binPath)
		if err != nil {
			fmt.Println("issue setting PROG :", binPath, err)
		}

		if !dryRun {
			err = sh.Copy("/usr/local/etc/bash_completion.d/"+d.Name, "./pkgr/bash_autocomplete")
			if err != nil {
				fmt.Println("!!!error: ", err)
			}
		} else {
			fmt.Println("DRY_RUN copying ./pkgr/bash_autocomplete to /usr/local/etc/bash_completion.d/" + d.Name)
		}
		cleaner(d.Name, false)
	}

	return nil
}

func Clean() {
	mg.SerialDeps(parseToml)
	for _, d := range apps.Apps {
		if !d.Enable {
			continue
		}
		cleaner(d.Name, false)
	}
}

// Clean up after yourself
func cleaner(projectName string, dirsOnly bool) {
	fmt.Println("Cleaning...")
	for _, d := range toCleanDirs {
		fmt.Println("  Cleaning...", d)
		err := os.RemoveAll(d)
		if err != nil {
			log.Println(err)
		}
	}

	for _, d := range toCleanFiles {
		fmt.Println("  Cleaning...", d)
		matches := findFiles(projectName)
		if !dirsOnly {
			for _, k := range matches {
				fmt.Println("  Cleaning...", k)
				err := os.Remove(k)
				if err != nil {
					log.Println(err)
				}
			}
		}
		err := os.Remove(d)
		if err != nil {
			log.Println(err)
		}
	}
}

// *********** ONE OFF TASKS BELOW

// Setup preTask
func setup(app application) error {

	mg.SerialDeps(parseToml)
	fmt.Println("  retrieve version")
	ver := version(app.VersionFile)
	fmt.Println("  retrieve git commit hash")
	gitCommit := gitCommitHash()

	fmt.Println("  retrieve current path")
	cur := currentPath()
	baseDir = cur

	prepDir = filepath.Join(baseDir, "PREP")
	//app.Package += "/" + app.Name
	buildDir = filepath.Join(baseDir, crossBuildDir)
	log.Printf("** Project: %s \n** Project Pkg: %s\n** CurrentPath: %s", app.Name, app.Package, baseDir)

	fmt.Println("  setup ldflag version templates")
	versionFields := fmt.Sprintf(versionFieldsTemplate, versionPkg, gitCommit, versionPkg, ver, versionPkg, strconv.FormatInt(timestamp, 10), versionPkg, goVersion())
	fmt.Println("Version Fields: ", versionFields)
	if len(app.SaltVariableOverwrite) > 0 {
		versionFields += " " + fmt.Sprintf(app.SaltVariableOverwrite, saltOverwriteValue)
	}
	goLDFlags = fmt.Sprintf(goLDFlagsTemplate, versionFields)
	goLDFlagsStatic = fmt.Sprintf(goLDFlagsStaticTemplate, versionFields)

	fmt.Println("  load artifactory creds")

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
	log.Println("reading ", versionFile, "found", ver)
	return ver
}

func gitCommitHash() string {

	gitCommit := hash()

	gitCommit += gitStatus()

	log.Println("** GIT HASH:", gitCommit)

	return gitCommit
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hashResp, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hashResp
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
		4. on update verify sig with public key in application

		var publicKey = []byte(`
		-----BEGIN PUBLIC KEY-----
		MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEtrVmBxQvheRArXjg2vG1xIprWGuCyESx
		MMY8pjmjepSy2kuz+nl9aFLqmr+rDNdYvEBqQaZrYMc6k29gjvoQnQ==
		-----END PUBLIC KEY-----
		`)
	*/

	return nil
}
