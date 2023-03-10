package engine

import (
	"sync"
	"github.com/malivvan/vlang/internal/plugin"

	"github.com/kardianos/service"
)

type Engine struct {
	wg      sync.WaitGroup
	factory *plugin.Factory
}

func New(factory *plugin.Factory) *Engine {
	return &Engine{
		factory: factory,
	}
}

func (e *Engine) Start(s service.Service) error {
	println("Starting service...")

	return nil
}

func (e *Engine) Stop(s service.Service) error {
	println("Stopping service...")
	return nil
}

func (e *Engine) Routine() {

}
