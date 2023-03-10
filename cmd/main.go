package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/malivvan/vlang/internal/conf"
	"github.com/malivvan/vlang/internal/engine"
	"github.com/malivvan/vlang/internal/plugin"
	"github.com/malivvan/vlang/internal/repl"
	"github.com/malivvan/vlang/internal/vm"

	"github.com/dop251/goja"
	"github.com/integrii/flaggy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kardianos/service"
)

var version = "0.0.1"
var config = &service.Config{
	Name:        "Vortex",
	DisplayName: "Vortex",
	Description: "Vortex is a scripting engine for IoT devices.",
}

var (
	cmdRun       *flaggy.Subcommand
	cmdInstall   *flaggy.Subcommand
	cmdUninstall *flaggy.Subcommand
	cmdStart     *flaggy.Subcommand
	cmdStop      *flaggy.Subcommand
	cmdRestart   *flaggy.Subcommand
	cmdDevelop   *flaggy.Subcommand
)

var (
	varScript  string
	varUser    string
	varWorkdir = "./"
)

func init() {
	flaggy.SetName("Test Program")
	flaggy.SetDescription("A little example program")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	cmdRun = flaggy.NewSubcommand("run")
	cmdRun.Description = "run the service"
	cmdRun.AddPositionalValue(&varWorkdir, "workdir", 1, true, "path to the working directory")
	flaggy.AttachSubcommand(cmdRun, 1)

	cmdInstall = flaggy.NewSubcommand("install")
	cmdInstall.Description = "install the service"
	cmdInstall.AddPositionalValue(&varWorkdir, "workdir", 1, true, "path to the working directory")
	cmdInstall.AddPositionalValue(&varUser, "user", 2, true, "user to run the service as")
	flaggy.AttachSubcommand(cmdInstall, 1)

	cmdUninstall = flaggy.NewSubcommand("uninstall")
	cmdUninstall.Description = "uninstall the service"
	flaggy.AttachSubcommand(cmdUninstall, 1)

	cmdStart = flaggy.NewSubcommand("start")
	cmdStart.Description = "start the service"
	flaggy.AttachSubcommand(cmdStart, 1)

	cmdStop = flaggy.NewSubcommand("stop")
	cmdStop.Description = "stop the service"
	flaggy.AttachSubcommand(cmdStop, 1)

	cmdRestart = flaggy.NewSubcommand("restart")
	cmdRestart.Description = "restart the service"
	flaggy.AttachSubcommand(cmdRestart, 1)

	cmdDevelop = flaggy.NewSubcommand("develop")
	cmdDevelop.Description = "develop a script with live reload"
	cmdDevelop.AddPositionalValue(&varWorkdir, "workdir", 1, true, "path to the working directory")
	cmdDevelop.AddPositionalValue(&varScript, "script", 2, true, "name of the script to test")
	flaggy.AttachSubcommand(cmdDevelop, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {
	if cmdRun.Used {

		// Read the config file.
		var cfg conf.Config
		err := conf.LoadFile(filepath.Join(varWorkdir, "config.toml"), &cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Configure the plugin system.
		log.Info().Msg("Creating Plugin Factory")
		factory, err := plugin.NewFactory(cfg)
		if err != nil {
			panic(err)
		}

		// Create and run the service.
		config.Arguments = []string{"run", os.Args[2]}
		s, err := service.New(engine.New(factory), config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating service:", err)
			os.Exit(1)
		}
		err = s.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		os.Exit(0)
	} else if cmdDevelop.Used {
		Develop()
	} else {
		repl.Run(filepath.Join(varWorkdir, "config.toml"))
	}
}

func Develop() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var rt *goja.Runtime
	var wg sync.WaitGroup

	// Read the config file.
	var cfg conf.Config
	err := conf.LoadFile(filepath.Join(varWorkdir, "config.toml"), &cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Configure the plugin system.
	log.Info().Msg("Creating Plugin Factory")
	factory, err := plugin.NewFactory(cfg)
	if err != nil {
		panic(err)
	}
	provider := factory.NewProvider()

	// Schedule cleanup and shutdown on ctrl-c.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println()
		if rt != nil {
			rt.Interrupt("shutdown")
			wg.Wait()
		}
		provider.Destroy()
		os.Exit(0)
	}()

	// sanitize script var and get the path to the script.
	if !strings.HasSuffix(varScript, ".js") {
		varScript += ".js"
	}
	scriptPath := filepath.Join(varWorkdir, "scripts", varScript)

	lastInfo, err := os.Stat(scriptPath)
	if err != nil {
		panic(err)
	}
	for {

		// Run script in a goroutine.
		script, err := ioutil.ReadFile(scriptPath)
		if err != nil {
			panic(err)
		}
		rt = vm.New(provider)
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info().Str("script", varScript).Msg("Starting Script")
			_, err = rt.RunString(string(script))
			if err != nil {
				switch err.(type) {
				default:
					log.Err(err).Str("script", varScript).Err(err).Msg("Script Error")
				case *goja.InterruptedError:
					if err.(*goja.InterruptedError).Value() == "reload" {
						log.Info().Str("script", varScript).Str("reason", "reload").Msg("Stopping Script")
					} else if err.(*goja.InterruptedError).Value() == "shutdown" {
						log.Info().Str("script", varScript).Str("reason", "shutdown").Msg("Stopping Script")
					} else {
						panic(err) // unexpected interrupt
					}
				}
			}
		}()

		// On change, interrupt the script and wait for it to finish.
		for {
			time.Sleep(500 * time.Millisecond)
			info, err := os.Stat(scriptPath)
			if err != nil {
				panic(err)
			}
			if info.ModTime().After(lastInfo.ModTime()) {
				rt.Interrupt("reload")
				wg.Wait()
				lastInfo = info
				break
			}
		}
	}
}
