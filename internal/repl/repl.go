package repl

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"github.com/malivvan/vlang/internal"
	"github.com/malivvan/vlang/internal/conf"
	"github.com/malivvan/vlang/internal/plugin"
	"github.com/malivvan/vlang/internal/vm"

	"github.com/c-bata/go-prompt"
	"github.com/dop251/goja"
	"github.com/fatih/color"
)

type REPL struct {
	Runtime *goja.Runtime
}

func (repl *REPL) executor(in string) {
	ret, err := repl.Runtime.RunString(in)
	if err != nil {
		color.HiRed("Error: %s", err)
	} else if ret != nil && ret.ExportType() != nil {
		fmt.Println(ret.Export())
	}
}

func (repl *REPL) completer(in prompt.Document) []prompt.Suggest {

	return []prompt.Suggest{}
}
func (repl *REPL) exit(buf *prompt.Buffer) {
	var cmd = exec.Command("/bin/stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	os.Exit(0)
}

func Run(configPath string) {

	// Read the config file.
	var cfg conf.Config
	err := conf.LoadFile(configPath, &cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Configure the plugin system.
	factory, err := plugin.NewFactory(cfg)
	if err != nil {
		panic(err)
	}
	provider := factory.NewProvider()
	repl := REPL{Runtime: vm.New(provider)}

	color.New(color.FgHiWhite, color.Bold).Print(internal.Name + " " + internal.Version)
	color.White(" [" + runtime.Version() + "-" + runtime.GOOS + "_" + runtime.GOARCH + "]")
	prompt.New(repl.executor, repl.completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle(internal.Name+" "+internal.Version),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn:  repl.exit,
		})).Run()

}
