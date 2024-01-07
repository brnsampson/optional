package main

import (
	"github.com/BurntSushi/toml"
	"github.com/brnsampson/optional"
	"github.com/brnsampson/optional/config"
	"github.com/caarlos0/env"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

const (
	DEFAULT_HOST string = "localhost"
	DEFAULT_PORT int    = 1443
)

type SubConfigLoader struct {
	Port config.Int `env:"PORT"`
}

type ConfigLoader struct {
	ConfigPath config.Str `env:"CONFIG_FILE"`
	Host       config.Str `env:"HOST"`
	Nested     SubConfigLoader
}

type SubConfig struct {
	Port int
}

type Config struct {
	Host   string
	Nested *SubConfig
}

func NewSubConfigLoader() SubConfigLoader {
	return SubConfigLoader{config.NoInt()}
}

func SubConfigLoaderFromEnv() (loader SubConfigLoader, err error) {
	loader = NewSubConfigLoader()

	if err = env.Parse(&loader); err != nil {
		log.Error("Failed to load sub-config loader from env variables!")
		return
	}

	log.Debug("Loaded sub-config loader from env variables", "config", loader)
	return
}

func (l SubConfigLoader) WithPort(port config.Int) SubConfigLoader {
	l.Port = optional.Or(port, l.Port)
	return l
}

func (l SubConfigLoader) OrPort(port config.Int) SubConfigLoader {
	l.Port = optional.Or(l.Port, port)
	return l
}

func (l SubConfigLoader) Merged(other SubConfigLoader) SubConfigLoader {
	return l.OrPort(other.Port)
}

func (l SubConfigLoader) Finalize() (*SubConfig, error) {
	port := l.Port.UnwrapOr(DEFAULT_PORT)
	config := SubConfig{port}
	log.Info("Finalized subconfig", "config", config)
	return &config, nil
}

func (l *SubConfigLoader) Reload(init SubConfigLoader) error {
	env, err := SubConfigLoaderFromEnv()
	if err != nil {
		return err
	}

	*l = init.Merged(env)
	return nil
}

// Now for the main config loader
func NewConfigLoader() ConfigLoader {
	return ConfigLoader{config.NoStr(), config.NoStr(), NewSubConfigLoader()}
}

func ConfigLoaderFromEnv() (loader ConfigLoader, err error) {
	loader = NewConfigLoader()
	sc, err := SubConfigLoaderFromEnv()
	if err != nil {
		return
	}

	if err = env.Parse(&loader); err != nil {
		log.Error("Failed to load server config from env variables!")
		return
	}

	loader.Nested = sc
	log.Debug("Loaded server config loader from env variables", "config", loader)
	return
}

func ConfigLoaderFromFile(pathOpt config.Str) (loader ConfigLoader, err error) {
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

	loader = loader.WithConfigPath(config.SomeStr(abs))

	log.Debug("Loaded server config from file", "filename", abs, "config", loader)

	return
}

func LoadedConfigLoader(path config.Str, override ConfigLoader) (loader ConfigLoader, err error) {
	loader = override
	file, err := ConfigLoaderFromFile(path)
	if err != nil {
		return
	}
	loader = loader.Merged(file)

	env, err := ConfigLoaderFromEnv()
	if err != nil {
		return
	}
	loader = loader.Merged(env)
	log.Info("Loaded server config from all sources", "filename", path, "config", loader)

	return
}

func (l ConfigLoader) WithConfigPath(path config.Str) ConfigLoader {
	l.ConfigPath = optional.Or(path, l.ConfigPath)
	return l
}

func (l ConfigLoader) OrConfigPath(path config.Str) ConfigLoader {
	l.ConfigPath = optional.Or(l.ConfigPath, path)
	return l
}

func (l ConfigLoader) WithHost(host config.Str) ConfigLoader {
	l.Host = optional.Or(host, l.Host)
	return l
}

func (l ConfigLoader) OrHost(host config.Str) ConfigLoader {
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

func (l ConfigLoader) Finalize() (*Config, error) {
	host := l.Host.UnwrapOr(DEFAULT_HOST)
	nested, err := l.Nested.Finalize()
	if err != nil {
		log.Error("Failed to finalize config from config loader!")
		return nil, err
	}

	config := Config{host, nested}
	log.Info("Finalized server config", "config", config)
	return &config, nil
}

func (l *ConfigLoader) Reload(init ConfigLoader) error {
	newConfig, err := LoadedConfigLoader(l.ConfigPath, init)
	if err != nil {
		return err
	}

	*l = newConfig
	return nil
}
