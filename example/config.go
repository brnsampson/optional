package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/brnsampson/optional"
	"github.com/brnsampson/optional/file"
	"github.com/caarlos0/env"
)

const (
	DEFAULT_CONFIG_FILE = "conf.toml"
	DEFAULT_DEV_MODE    = false
	DEFAULT_LOG_LEVEL   = "info"
	DEFAULT_HOST        = "localhost"
	DEFAULT_IP          = "127.0.0.1"
	DEFAULT_PORT        = 1443
)

var (
	// Loader flags
	dev      bool
	confFile file.File = file.SomeFile(DEFAULT_CONFIG_FILE) // Setting a default value for an optional
	host     optional.Str
	logLevel optional.Str
	// subLoader flags
	ip   optional.Str
	port optional.Int
)

func setupFlags() {
	// Set up all of our command line flags here to make sure we don't have config scattered across
	// the whole world.

	// Loader
	flag.BoolVar(&dev, "dev", false, "Set default values appropriate for local development")
	flag.Var(&logLevel, "log", "set logging level [debug, info, warn, err]. Defaults to info.")
	flag.Var(&confFile, "config", "path to optional config file. Set to `none` to disable file loading entirely.")
	flag.Var(&host, "host", "hostname for a server or whatever")

	// subLoader
	flag.Var(&ip, "ip", "ip address for a server or whatever")
	flag.Var(&port, "port", "port for a server or whatever")
}

// Static and reloadable configs
type staticSubConfig struct {
	IP   string
	Port int
}

// If you has a module in your code that only needed this fragment of the config, you could add convenience
// methods to it like this:
func (c staticSubConfig) GetAddr() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

type staticConfig struct {
	ConfigFile string
	DevMode    bool
	LogLevel   slog.Level
	Host       string
	SubConfig  staticSubConfig
}

type subLoader struct {
	IP   optional.Str `env:"IP"`
	Port optional.Int `env:"PORT"`
}

func newSubLoader() subLoader {
	// These vars are defined at the top of the file globally and added to the flags in an init() funciton.
	// ip := optional.NoStr()
	// port := optional.NoInt()

	// fset := flag.NewFlagSet("subLoader", flag.ContinueOnError)
	// fset.Var(&ip, "ip", "ip address for a server or whatever")
	// fset.Var(&port, "port", "port for a server or whatever")

	// err := fset.Parse()
	return subLoader{IP: ip, Port: port}
}

func (l subLoader) Name() string {
	return "subLoader"
}

func (l *subLoader) ToStatic() (staticSubConfig, error) {
	return staticSubConfig{
		IP:   optional.GetOr(l.IP, DEFAULT_IP),
		Port: optional.GetOr(l.Port, DEFAULT_PORT),
	}, nil
}

type Loader struct {
	ConfigFile file.File    `env:"CONFIG_FILE"`
	DevMode    bool         `env:"DEV_MODE"`
	LogLevel   optional.Str `env:"LOG_LEVEL"`
	Host       optional.Str `env:"HOST"`
	SubLoader  subLoader
	current    staticConfig
}

func NewLoader() Loader {
	subLoader := newSubLoader()

	// fset := flag.NewFlagSet("loader", flag.ContinueOnError)
	// confPath := optional.SomeStr(DEFAULT_CONFIG_PATH)
	// host := optional.NoStr()
	// fset.Var(&confPath, "optional", "path to optional file. Set to `none` to disable loading from optional.")
	// fset.Var(&host, "host", "hostname for a server or whatever")
	// flag.Parse()

	return Loader{
		ConfigFile: confFile,
		DevMode:    dev,
		LogLevel:   logLevel,
		Host:       host,
		SubLoader:  subLoader,
		current:    staticConfig{},
	}
}

func (l Loader) Name() string {
	return "Loader"
}

func (l Loader) Current() staticConfig {
	return l.current
}

func (l *Loader) Update(configFile string) error {
	if configFile != "" {
		l.ConfigFile.Set(configFile)
	}

	if err := env.Parse(l); err != nil {
		return err
	}

	pretty, err := json.Marshal(l)
	if err != nil {
		slog.Error("Failed to print current state of loader... Does it contain a non-marshallable type?")
	} else {
		slog.Debug("Loaded env vars", "loader", string(pretty))
	}

	if l.ConfigFile.IsSome() {
		abs, err := l.ConfigFile.Abs()
		if err != nil {
			slog.Error(
				"Could not retrieve config file absolute file path from loader",
				slog.Any("path", abs),
				slog.Any("state", l),
				slog.Any("error", err),
			)
			return err
		}

		path, err := abs.Get()
		if err != nil {
			slog.Error(
				"Could not retrieve absolute file path from loader",
				slog.String("path", path),
				slog.Any("state", l),
				slog.Any("error", err),
			)
			return err
		}

		_, err = toml.DecodeFile(path, l)
		if err != nil {
			slog.Error(
				"Could not decode file into FullLoader",
				slog.String("path", path),
				slog.Any("state", l),
				slog.Any("error", err),
			)
			return err
		}
	}

	l.current, err = l.ToStatic()
	return err
}

func (l Loader) ToStatic() (staticConfig, error) {
	subconfig, err := l.SubLoader.ToStatic()
	if err != nil {
		slog.Error("Failed loading sub-component", "component", l.SubLoader.Name())
	}

	// Sometimes you have to do some logic on a value. Another great example would be reading TLS
	// certificates or signing keys.
	ll := l.LogLevel
	var level slog.Level
	if ll.Match("debug") || ll.Match("Debug") {
		level = slog.LevelDebug
	} else if ll.Match("warn") || ll.Match("Warn") || ll.Match("warning") || ll.Match("Warning") {
		level = slog.LevelWarn
	} else if ll.Match("err") || ll.Match("Err") || ll.Match("error") || ll.Match("Error") {
		level = slog.LevelError
	} else {
		level = slog.LevelInfo
	}

	// DevMode overrides
	if l.DevMode {
		level = slog.LevelDebug
	}

	return staticConfig{
		ConfigFile: optional.GetOr(l.ConfigFile, DEFAULT_CONFIG_FILE),
		DevMode:    l.DevMode,
		LogLevel:   level,
		Host:       optional.GetOr(l.Host, DEFAULT_HOST),
		SubConfig:  subconfig,
	}, nil
}
