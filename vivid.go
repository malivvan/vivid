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
	//////////////////////////////////////////////////////////////////////
	AppBinary = func() string {
		name := AppName + "_" + runtime.GOOS + "_" + runtime.GOARCH
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		return name
	}()
	AppRepo    = "undefined"
	AppVersion = "undefined"
	AppCommit  = "undefined"
	AppBuild   = "undefined"
	AppPubkey  = "undefined"
)

func init() {
	if !strings.HasPrefix(os.Args[0], filepath.Join(os.TempDir(), "go-build")) {
		if AppRepo == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application repository")
		}
		if AppVersion == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application version")
		}
		if AppCommit == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application commit")
		}
		if AppBuild == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application build")
		}
		if AppPubkey == "undefined" {
			fmt.Fprintln(os.Stderr, "invalid application pubkey")
		}
		if AppRepo == "undefined" || AppVersion == "undefined" || AppCommit == "undefined" || AppBuild == "undefined" || AppPubkey == "undefined" {
			fmt.Fprintln(os.Stderr, "please use the build script to build the application")
			os.Exit(1)
		}
	}
}
