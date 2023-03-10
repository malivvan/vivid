package vm

import (
	"math/big"
	"github.com/malivvan/vlang/internal/plugin"

	"github.com/dop251/goja"
)

func New(plugin *plugin.Provider) *goja.Runtime {
	vm := goja.New()
	vm.Set("plugin", plugin)
	vm.Set("encrypt", encrypt)
	vm.Set("decrypt", decrypt)
	vm.Set("sleep", sleep)
	vm.Set("print", print)
	vm.Set("pow", pow)
	return vm
}

func pow(i, e int64) *big.Int {
	var a, b = big.NewInt(i), big.NewInt(e)
	return a.Exp(a, b, nil)

}
