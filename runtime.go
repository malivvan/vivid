package vivid

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/malivvan/vivid/internal/conf"
)

type Config struct {
	Log LogConfig `toml:"log"`
}

type Environment struct {
	mutex    sync.Mutex
	appdir   string
	config   interface{}
	logger   *Logger
	plugins  []PluginLoader
	machines []*Machine
	watcher  *Watcher
}

var DefaultConfig = Config{
	Log: LogConfig{
		Level: "info",
	},
}

func New(appdir string, config interface{}, loaders ...PluginLoader) (*Environment, error) {

	// Sanity check input
	if info, err := os.Stat(appdir); err != nil || !info.IsDir() {
		return nil, errors.New("vlang: workdir does not exist or is not a directory")
	}
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		return nil, errors.New("vlang: config defaults must be a pointer")
	}

	// Create env, load config file and configure initially.
	env := &Environment{
		appdir:  appdir,
		config:  config,
		plugins: loaders,
	}
	err := env.reconfigure()
	if err != nil {
		return nil, errors.New("vlang: failed to configure: " + err.Error())
	}
	//rt.logger.Info().Msg("vlang node configured")

	// Create config file watcher
	//	rt.watcher = newWatcher(filepath.Join(rt.workdir, "vivid.conf"), 500, func() {
	//	err = rt.reconfigure()
	//	if err != nil {
	//		rt.logger.Warn().Err(err).Msg("failed to reconfigure node")
	//		return
	//	}
	//	rt.logger.Info().Msg("vlang node reconfigured")
	//})

	return env, nil
}

func (env *Environment) Cleanup() error {
	env.mutex.Lock()
	defer env.mutex.Unlock()
	//////////////////////
	env.logger.Info().Msg("vlang node cleanup")

	// 1. Stop all VM's.
	if len(env.machines) > 0 {
		for _, vm := range env.machines {
			vm.Stop()
		}
		env.machines = []*Machine{}
	}

	// 2. Close the old logger.
	if env.logger != nil {
		err := env.logger.Close()
		if err != nil {
			return errors.New("vlang: failed to close logger: " + err.Error())
		}
	}

	return nil
}

func (env *Environment) reconfigure() error {
	env.mutex.Lock()
	defer env.mutex.Unlock()
	//////////////////////

	// 1. Stop all VM's.
	if len(env.machines) > 0 {
		for _, vm := range env.machines {
			vm.Stop()
		}
		env.machines = []*Machine{}
	}

	// 2. Close the old logger.
	if env.logger != nil {
		err := env.logger.Close()
		if err != nil {
			return errors.New("vlang: failed to close logger: " + err.Error())
		}
	}

	// 3. Reload the config file.
	err := conf.LoadFile(filepath.Join(env.appdir, AppName+".conf"), env.config)
	if err != nil {
		return errors.New("vlang: failed to load config file: " + err.Error())
	}

	// 4. Create the new logger.
	env.logger = env.newLogger(AppName)
	env.logger.Info().Str("version", AppVersion).Str("appdir", env.appdir).Msg("creating environment")

	// 5. Configure all plugins.
	for _, plugin := range env.plugins {

		// Cleanup the plugin.
		err := plugin.Cleanup()
		if err != nil {
			return errors.New("vlang: failed to cleanup plugin '" + plugin.Name() + " v" + plugin.Version() + "': " + err.Error())
		}

		// Validate and extract the plugin config.
		pluginConfigField := reflect.ValueOf(env.config).Elem().FieldByName(plugin.Name())
		if !(pluginConfigField.Kind() == reflect.Map || pluginConfigField.Type().Key().Kind() == reflect.String || pluginConfigField.Type().Elem().Kind() == reflect.Ptr) {
			panic("field '" + plugin.Name() + "' of type map[string]*" + strings.ToLower(plugin.Name()) + "." + "Config is missing in config")
		}
		pluginConfig := pluginConfigField.Interface()

		// Configure the plugin.
		err = plugin.Configure(pluginConfig)
		if err != nil {
			env.logger.Err(err).
				Str("name", plugin.Name()).
				Str("version", plugin.Version()).
				Msg("failed to configure plugin")
			return err
		}

		// Log the plugin config.
		configs := []string{}
		for _, key := range pluginConfigField.MapKeys() {
			configs = append(configs, key.String())
		}
		env.logger.Info().
			Str("name", plugin.Name()).
			Str("version", plugin.Version()).
			Strs("configs", configs).
			Msg("plugin configured")
	}

	env.logger.Info().Msg("environment created")
	return nil
}
