package main

import (
	"os"
	//"errors"
	"path/filepath"
	"github.com/spf13/pflag"
	"github.com/caarlos0/env"
    "github.com/BurntSushi/toml"
    "github.com/charmbracelet/log"
    "github.com/brnsampson/optional/config"
)

const (
	DEFAULT_HOST string = "localhost"
	DEFAULT_PORT int = 1443
)

type SubConfigLoader struct {
	Port config.Int `env:"PORT"`
}

type ConfigLoader struct {
	Host config.Str `env:"HOST"`
	Nested *SubConfigLoader
}

type SubConfig struct {
	Port int
}

type Config struct {
    Host string
	Nested *SubConfig
}

// Fun fact, if your config has a sub-struct and you implement the same methods on it, it works pretty well with toml
// and the env / flag loading methods kind of just nest together like a russian doll.
func NewSubConfigLoader() *SubConfigLoader {
	return &SubConfigLoader{ config.NoInt() }
}

func (c *SubConfigLoader) LoadFromFlags(flags *pflag.FlagSet) error {
    tmp, err := flags.GetInt("port")
	if err != nil {
		return err
	} else if tmp == 0 {
		log.Debug("port not found in flags")
	} else {
        c.Port.Set(tmp)
    }

    return nil
}

func (c *SubConfigLoader) LoadFromEnv() error {
	if err := env.Parse(c); err != nil {
		log.Error("Failed to load nested config from env variables!")
		return err
	}

	log.Debug("Loaded server nested config from env variables", "config", c)
	return nil
}

func (c SubConfigLoader) Finalize() (*SubConfig, error) {
    port := c.Port.UnwrapOr(DEFAULT_PORT)
    return &SubConfig{port}, nil
}

// Now for the main config loader
func NewConfigLoader() *ConfigLoader {

    return &ConfigLoader{ config.NoStr(), NewSubConfigLoader() }
}

func (c *ConfigLoader) LoadFromFlags(flags *pflag.FlagSet) error {
    tmp, err := flags.GetString("host")
	if err != nil {
		return err
	} else if tmp == "" {
		log.Debug("Host not found in flags")
	} else {
        c.Host.Set(tmp)
    }

	return c.Nested.LoadFromFlags(flags)
}

func (c *ConfigLoader) LoadFromEnv() error {
	if err := c.Nested.LoadFromEnv(); err != nil {
		log.Error("Failed to load server config from env variables!")
		return err
	}

	if err := env.Parse(c); err != nil {
		log.Error("Failed to load server config from env variables!")
		return err
	}

	log.Debug("Loaded server config from env variables", "config", c)
	return nil
}

func (c *ConfigLoader) LoadFromFile(path string) error {
	// try to make it work if the user is in the project root or example directory
	tmp := path
	path = "./" + path
	if _, err := os.Stat(path); err != nil {
		path = "./example/" + tmp
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

    _, err = toml.DecodeFile(abs, &c)
    if err != nil {
		return err
    }

	log.Info("Loaded server config from file", "filename", abs, "config", c)

	return nil
}

func (c *ConfigLoader) LoadAll(flags *pflag.FlagSet, path string) error {
    if err := c.LoadFromEnv(); err != nil {
		return err
    }

    if err := c.LoadFromFile(path); err != nil {
		return err
    }

    if err := c.LoadFromFlags(flags); err != nil {
		return err
    }
	return nil
}

func (c ConfigLoader) Finalize() (*Config, error) {
    host := c.Host.UnwrapOr(DEFAULT_HOST)
	nested, err := c.Nested.Finalize()
	if err != nil {
		log.Error("Failed to finalize config from config loader!")
		return nil, err
	}

    return &Config{host, nested}, nil
}
