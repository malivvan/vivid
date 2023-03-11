package vivid

import "github.com/kardianos/service"

type Service struct {
	runtime *Environment
}

func (e *Service) Start(s service.Service) error {
	println("Starting service...")

	return nil
}

func (e *Service) Stop(s service.Service) error {
	println("Stopping service...")
	return nil
}

func (e *Service) Routine() {

}
