package vivid

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/dop251/goja"
	"github.com/fatih/color"
	"github.com/kardianos/service"
)

func CLI(config interface{}, loaders ...PluginLoader) {
	//	svcConfig := &service.Config{
	//		Name:        name,
	//		DisplayName: name,
	////	}

	args := os.Args[1:] // remove first argument
	if len(args) == 0 {

		// INTERPRETER: REPL in current directory.
		absPath, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
		cliInterpreter(absPath, true, false, config, loaders...)

	} else if len(args) == 2 && args[0] == "dev" {

		// INTERPRETER: Develop specified script.
		absPath, err := filepath.Abs(args[1])
		if err != nil {
			panic(err)
		}
		cliInterpreter(absPath, false, true, config, loaders...)

	} else if len(args) == 2 && args[0] == "run" {

		// DAEMON: Run.
		workdir := args[1]
		fmt.Println("run daemon in " + workdir)

	} else if len(args) == 3 && args[0] == "install" {

		// DAEMON: Install.
		workdir := args[1]
		user := args[2]
		fmt.Println("install daemon in " + workdir + " as " + user)

	} else if len(args) == 1 && args[0] == "uninstall" {

		// DAEMON: Uninstall.
		fmt.Println("uninstall daemon")

	} else if len(args) == 1 && args[0] == "start" {

		// DAEMON: Start.
		fmt.Println("start daemon")

	} else if len(args) == 1 && args[0] == "stop" {

		// DAEMON: Stop.
		fmt.Println("stop daemon")

	} else if len(args) == 1 && args[0] == "restart" {

		// DAEMON: Restart.
		fmt.Println("restart daemon")

	} else if len(args) >= 1 && args[0] == "doc" {

		// INFO: Show documentation.
		plugin := ""
		if len(args) >= 2 {
			plugin = args[1]
		}
		fmt.Println("show docs " + plugin)

	} else if len(args) >= 1 && args[0] == "update" {

		// INFO: Update to specified version.
		version := args[1]
		if len(args) >= 2 {
			version = args[1]
		}
		fmt.Println("update to version " + version)

	} else if len(args) == 1 && args[0] == "version" {

		fmt.Println("binary:  " + AppBinary)
		fmt.Println("version: " + AppVersion)
		fmt.Println("commit:  " + AppCommit)
		fmt.Println("build:   " + AppBuild)

	} else if len(args) == 1 && args[0] == "help" {

		// INFO: Show help.
		cliHelp(false)

	} else if len(args) == 1 {

		// Get absolute path of script or folder.
		absPath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println("script or folder at '" + args[0] + "' does not exist")
			os.Exit(1)
		}
		info, err := os.Stat(absPath)
		if os.IsNotExist(err) {
			fmt.Println("script or folder at '" + absPath + "' does not exist")
			os.Exit(1)
		}

		// Check if path is a directory or a file.
		if info.IsDir() {

			// INTERPRETER: REPL in specified directory.
			cliInterpreter(absPath, true, false, config, loaders...)
		} else {

			// INTERPRETER: Execute specified script.
			cliInterpreter(absPath, false, false, config, loaders...)
		}

	} else {
		fmt.Println("Invalid command line arguments!")
		fmt.Println("Please type '" + AppName + " help' for more information.")
		fmt.Println()
		os.Exit(1)
	}
}

func cliHelp(nocolor bool) {
	color.NoColor = nocolor
	ubw := color.New(color.FgHiWhite, color.Underline, color.Bold) // Underlined Bold White
	rrg := color.New(color.FgWhite)                                // Regular Regular Grey
	rbg := color.New(color.FgWhite, color.Bold)                    // Regular Bold Grey
	rbw := color.New(color.FgHiWhite, color.Bold)                  // Regular Bold White
	ubw.Println(strings.ToUpper(AppName) + " - " + AppDescription)
	fmt.Println()
	rbw.Println("USAGE")
	rbg.Println("  " + AppName + " ($script|$folder|command) (arguments)")
	fmt.Println()
	rbw.Println("INTERPRETER")
	rrg.Println("  In interpreter mode a single thread executes a single script. The")
	rrg.Println("  interpreter uses $home/." + AppName + " as app directory. GUI functions")
	rrg.Println("  are available in interpreter mode only.")
	fmt.Println()
	rbg.Println("  " + AppName + "                             launch REPL in current dir as workdir")
	rbg.Println("  " + AppName + " [$workdir]                  launch REPL in specified $workdir")
	rbg.Println("  " + AppName + " [$script]                   run a single $script from path")
	rbg.Println("  " + AppName + " dev [$script]               develop a $script with live reload")
	fmt.Println()
	rbw.Println("DAEMON")
	rrg.Println("  In daemon mode multiple threads execute the " + AppName + " scripts inside the")
	rrg.Println("  given application directory. GUI functions are NOT available in daemon mode.")
	fmt.Println()
	rbg.Println("  " + AppName + " run [$appdir]               run daemon in $appdir")
	rbg.Println("  " + AppName + " install [$appdir] [$user]   install daemon in $appdir as $user")
	rbg.Println("  " + AppName + " uninstall                   uninstall daemon")
	rbg.Println("  " + AppName + " start                       start deamon")
	rbg.Println("  " + AppName + " stop                        stop daemon")
	rbg.Println("  " + AppName + " restart                     restart daemon")
	fmt.Println()
	rbw.Println("INFO")
	rbg.Println("  " + AppName + " doc ($plugin)               print general and $plugin docs")
	rbg.Println("  " + AppName + " update ($version)           update to latest or specified $version")
	rbg.Println("  " + AppName + " version                     print version info")
	rbg.Println("  " + AppName + " help                        print this help")
	os.Exit(0)
}

func cliRunDaemon(svcConfig *service.Config, path string) {

}
func cliInstallDaemon(svcConfig *service.Config, path string, user string) {

}
func cliUninstallDaemon(svcConfig *service.Config) {

}
func cliControlDaemon(svcConfig *service.Config, action int) {

}

func cliInterpreter(path string, isDir bool, dev bool, config interface{}, loaders ...PluginLoader) {

	// 1. Determine and create home/.$AppName as interpreter application directory.
	appdir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error getting home directory:", err)
		os.Exit(1)
	}
	appdir = filepath.Join(appdir, "."+AppName)
	if err := os.MkdirAll(appdir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "error creating application directory:", err)
		os.Exit(1)
	}

	// 3.3. Create environment and spawn interpreter machine.
	env, err := New(appdir, config, loaders...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating runtime:", err)
		os.Exit(1)
	}

	// 2. Determine absolute working directory and change to it.
	workdir := path
	if !isDir {
		workdir = filepath.Dir(path)
	}
	if err := os.Chdir(workdir); err != nil {
		fmt.Fprintln(os.Stderr, "error changing to working directory:", err)
		os.Exit(1)
	}

	// 3. Create interpreter machine.
	machine, err := env.New("interpreter")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating machine:", err)
		os.Exit(1)
	}

	// 3. If path is a directory, launch REPL.
	if isDir {

		env.logger.Info().Str("machine", machine.Name()).Str("workdir", workdir).Msg("lauching repl")
		prompt.New(func(in string) {
			err := machine.Start(in, func(v goja.Value, err error) bool {
				if err != nil {
					color.HiRed("Error: %s\n", err)
				} else if v != nil && v.ExportType() != nil {
					fmt.Println(v.Export())
				}
				return false
			})
			if err != nil {
				color.HiRed("Error: %s\n", err)
			}
			machine.Wait()
		}, func(in prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		}, prompt.OptionLivePrefix(func() (prefix string, useLivePrefix bool) {
			return workdir + "> ", true
		}), prompt.OptionTitle(AppName+" "+AppVersion),
			prompt.OptionPrefixTextColor(prompt.Yellow),
			prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.ControlC,
				Fn: func(b *prompt.Buffer) {
					var cmd = exec.Command("/bin/stty", "sane")
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Run()
					env.Cleanup()
					os.Exit(0)
				},
			})).Run()

		// 3. If path is a file, run it.
	} else {

		// 3.4. Define run function.
		runfunc := func() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error reading script:", err)
				os.Exit(1)
			}
			machine.Start(string(data), func(v goja.Value, err error) bool {
				if err != nil {
					color.HiRed("Error: %s\n", err)
				} else if v != nil && v.ExportType() != nil {
					fmt.Println(v.Export())
				}
				return true
			})
			machine.Wait()
		}

		// 3.5. Run script once.
		runfunc()

		// 3.6. If dev mode is enabled, watch script for changes and run it again.
		if dev {
			newWatcher(path, 500, runfunc)
			select {}
		}
	}
}
