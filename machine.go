package vivid

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
	"github.com/malivvan/vivid/stdlib"
)

type Machine struct {
	env    *Environment
	name   string
	logger *Logger

	waitgroup sync.WaitGroup
	running   atomic.Value
	runtime   *goja.Runtime
}

func (machine *Machine) Name() string {
	return machine.name
}

func (env *Environment) New(name string) (*Machine, error) {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	for _, machine := range env.machines {
		if machine.name == name {
			return nil, errors.New("machine already exists")
		}
	}

	machine := &Machine{
		name:   name,
		logger: env.newLogger(name),
		env:    env,
	}
	machine.running.Store(false)
	env.machines = append(env.machines, machine)

	machine.logger.Info().Str("name", machine.name).Msg("machine created")
	return machine, nil
}

func (machine *Machine) Start(code string, callback func(goja.Value, error) bool) error {
	if machine.running.Swap(true).(bool) {
		return errors.New("machine is running")
	}
	machine.waitgroup.Add(1)
	go func() {
		defer func() {
			machine.running.Store(false)
			machine.waitgroup.Done()
		}()

		// Create a new runtime if not exist.
		if machine.runtime == nil {
			machine.runtime = goja.New()
			machine.runtime.Set("print", stdlib.Print)
			machine.runtime.Set("sleep", stdlib.Sleep)
			for _, loader := range machine.env.plugins {
				machine.runtime.Set(strings.ToLower(loader.Name()), loader.Func)
			}
		}

		// Run the code.
		value, err := machine.runtime.RunString(code)

		// Handle interrupt.
		_, interrupted := err.(*goja.InterruptedError)
		if interrupted {
			callback(nil, nil)    // callback without value and error on interrupt
			machine.runtime = nil // ensure runtime is reset after
			return
		}

		// Handle return value and error.
		resetRuntime := callback(value, err)
		if resetRuntime {
			machine.runtime = nil
		} else {
			machine.runtime.ClearInterrupt()
		}
	}()
	return nil
}

func (machine *Machine) Running() bool {
	return machine.running.Load().(bool)
}

func (machine *Machine) Wait() {
	if machine.running.Load().(bool) {
		machine.waitgroup.Wait()
	}
}

func (machine *Machine) Stop() {
	if machine.running.Load().(bool) {
		machine.runtime.Interrupt("INTERRUPT")
		machine.waitgroup.Wait()
	}
}
