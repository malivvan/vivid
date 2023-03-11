package vivid

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	AppName        = "vivid"
	AppDescription = "Violets inherently versatile Interpreter and Daemon"
	AppBinary      = func() string {
		name := AppName + "_" + runtime.GOOS + "_" + runtime.GOARCH
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		return name
	}()
	AppVersion = "undefined"
	AppCommit  = "undefined"
	AppBuild   = "undefined"
)

func init() {
	if !strings.HasPrefix(os.Args[0], filepath.Join(os.TempDir(), "go-build")) {
		if AppVersion == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application version")
		}
		if AppCommit == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application commit")
		}
		if AppBuild == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application build")
		}
		if AppVersion == "undefined" || AppCommit == "undefined" || AppBuild == "undefined" {
			fmt.Fprintln(os.Stderr, "please use the build script to build the application")
			os.Exit(1)
		}
	}
}
