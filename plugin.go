package vivid

type PluginLoader interface {
	Name() string
	Version() string
	Func() func(string) interface{}
	Configure(v interface{}) error
	Cleanup() error
}

type PluginInstance interface {
	Destroy() error
}
