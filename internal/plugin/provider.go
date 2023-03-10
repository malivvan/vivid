package plugin

import (
	"github.com/malivvan/vlang/internal/conf"
	"github.com/malivvan/vlang/internal/plugin/influx"

	"github.com/rs/zerolog/log"
)

type Factory struct {
	influx map[string]func() *influx.Plugin
}

func (f *Factory) NewProvider() *Provider {
	return &Provider{
		factory: f,
		influx:  make(map[string]*influx.Plugin),
	}
}

type Provider struct {
	factory *Factory
	influx  map[string]*influx.Plugin
}

func NewFactory(config conf.Config) (*Factory, error) {

	// Create the plugin factory.
	factory := &Factory{
		influx: make(map[string]func() *influx.Plugin),
	}

	// Influx
	for name := range config.Influx {
		cfg := config.Influx[name]
		factory.influx[name] = func() *influx.Plugin {
			return influx.New(*cfg)
		}
		log.Info().Str("name", name).
			Str("host", cfg.Host).
			Str("org", cfg.Org).
			Str("bucket", cfg.Bucket).
			Msg("influx plugin configured")
	}

	return factory, nil
}

func (p *Provider) Destroy() {
	for name, plugin := range p.influx {
		plugin.Destroy()
		log.Info().Str("name", name).
			Msg("influx plugin destroyed")
	}
}

func (p *Provider) Influx(name string) *influx.Plugin {
	if _, ok := p.influx[name]; !ok {
		log.Info().Str("name", name).Msg("influx plugin created")
		p.influx[name] = p.factory.influx[name]()
	}
	return p.influx[name]
}
