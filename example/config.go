package main

import (
	"github.com/BurntSushi/toml"
	"github.com/brnsampson/optional"
	"github.com/brnsampson/optional/confopt"
	"github.com/caarlos0/env"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

const (
	DEFAULT_HOST string = "localhost"
	DEFAULT_PORT int    = 1443
)

// Static and reloadable configs

type staticSubConfig struct {
	Port int
}

// This is just an example of how you would implement an interface that is then carried over to your SubConfig. More likely,
// you would have config fields for Host and Port and expose a method like GetAddr().
func (c staticSubConfig) GetPort() int {
	return c.Port
}

type staticConfig struct {
	Host   string
	Nested staticSubConfig
}

type SubConfig struct {
	staticSubConfig
	loader *ConfigLoader
}

func (c *SubConfig) Reload() error {
	l, err := c.loader.WithReload()
	if err != nil {
		return err
	}

	static, err := l.Finalize()
	if err != nil {
		return err
	}

	c.staticSubConfig = static.Nested
	return nil
}

type Config struct {
	staticConfig
	loader *ConfigLoader
}

func NewConfig(loader ConfigLoader) (conf Config, err error) {
	l, err := loader.WithReload()
	if err != nil {
		return
	}

	static, err := l.Finalize()
	if err != nil {
		return
	}

	return Config{static, &loader}, err
}

func (c *Config) Reload() error {
	l, err := c.loader.WithReload()
	if err != nil {
		return err
	}

	static, err := l.Finalize()
	if err != nil {
		return err
	}

	c.staticConfig = static
	return nil
}

// Loaders here.
type SubConfigLoader struct {
	Port confopt.Int `env:"PORT"`
}

type ConfigLoader struct {
	ConfigPath confopt.Str `env:"CONFIG_FILE"`
	Host       confopt.Str `env:"HOST"`
	Nested     SubConfigLoader
}

func NewSubConfigLoader() SubConfigLoader {
	return SubConfigLoader{confopt.NoInt()}
}

func NewConfigLoader() ConfigLoader {
	return ConfigLoader{confopt.NoStr(), confopt.NoStr(), NewSubConfigLoader()}
}

func configLoaderFromEnv() (loader ConfigLoader, err error) {
	loader = NewConfigLoader()

	if err = env.Parse(&loader); err != nil {
		log.Error("Failed to load server confopt from env variables!")
		return
	}

	if err = env.Parse(&loader.Nested); err != nil {
		log.Error("Failed to load sub-config loader from env variables!")
		return
	}

	log.Debug("Loaded server confopt loader from env variables", "confopt", loader)
	return
}

func configLoaderFromFile(pathOpt confopt.Str) (loader ConfigLoader, err error) {
	loader = NewConfigLoader()

	// Short-circuit and return default loader if no path was given.
	path, err := pathOpt.Get()
	if err != nil {
		log.Debug("config file path was None. Skipping...")
		return loader, nil
	}

	// try to make it work if the user is in the project root or example directory
	tmp := path
	path = "./" + path
	if _, err = os.Stat(path); err != nil {
		path = "./example/" + tmp
		if _, err = os.Stat(path); err != nil {
			return
		}
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}

	_, err = toml.DecodeFile(abs, &loader)
	if err != nil {
		return
	}

	loader = loader.WithConfigPath(confopt.SomeStr(abs))

	log.Debug("Loaded server confopt from file", "filename", abs, "confopt", loader)

	return
}

// SubConfigLoader methods
func (l SubConfigLoader) WithPort(port confopt.Int) SubConfigLoader {
	l.Port = optional.Or(port, l.Port)
	return l
}

func (l SubConfigLoader) OrPort(port confopt.Int) SubConfigLoader {
	l.Port = optional.Or(l.Port, port)
	return l
}

func (l SubConfigLoader) Merged(other SubConfigLoader) SubConfigLoader {
	return l.OrPort(other.Port)
}

func (l SubConfigLoader) Finalize() (staticSubConfig, error) {
	port := optional.GetOr(l.Port, DEFAULT_PORT)
	conf := staticSubConfig{port}
	log.Info("Finalized static subconfig from SubConfigLoader", "config", conf)
	return conf, nil
}

// Now for the main config loader
func loadedConfigLoader(override optional.Option[ConfigLoader]) (loader ConfigLoader, err error) {
	loader = optional.GetOrElse(override, NewConfigLoader)
	file, err := configLoaderFromFile(loader.ConfigPath)
	if err != nil {
		return
	}
	loader = loader.Merged(file)

	env, err := configLoaderFromEnv()
	if err != nil {
		return
	}
	loader = loader.Merged(env)
	log.Info("Loaded server config from all sources", "filename", loader.ConfigPath, "loader_state", loader)

	return
}

func (l ConfigLoader) AsOption() optional.Option[ConfigLoader] {
	return optional.Some(l)
}

func (l ConfigLoader) WithConfigPath(path confopt.Str) ConfigLoader {
	l.ConfigPath = optional.Or(path, l.ConfigPath)
	return l
}

func (l ConfigLoader) OrConfigPath(path confopt.Str) ConfigLoader {
	l.ConfigPath = optional.Or(l.ConfigPath, path)
	return l
}

func (l ConfigLoader) WithHost(host confopt.Str) ConfigLoader {
	l.Host = optional.Or(host, l.Host)
	return l
}

func (l ConfigLoader) OrHost(host confopt.Str) ConfigLoader {
	l.Host = optional.Or(l.Host, host)
	return l
}

func (l ConfigLoader) WithNested(nested SubConfigLoader) ConfigLoader {
	l.Nested = nested
	return l
}

func (l ConfigLoader) Merged(other ConfigLoader) ConfigLoader {
	return l.OrConfigPath(other.ConfigPath).OrHost(other.Host).WithNested(l.Nested.Merged(other.Nested))
}

func (l ConfigLoader) Finalize() (config staticConfig, err error) {
	host := optional.GetOr(l.Host, DEFAULT_HOST)
	nested, err := l.Nested.Finalize()
	if err != nil {
		log.Error("Failed to finalize config from config loader!")
		return
	}

	config = staticConfig{host, nested}
	log.Info("Finalized server config", "config", config)
	return config, nil
}

func (l ConfigLoader) WithReload() (loader ConfigLoader, err error) {
	return loadedConfigLoader(l.AsOption())
}
