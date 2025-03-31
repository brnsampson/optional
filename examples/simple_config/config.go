package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log/slog"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/brnsampson/optional"
	"github.com/brnsampson/optional/file"

	"go-simpler.org/env"
)

const (
	DEFAULT_CONFIG_FILE = "conf.toml"
	DEFAULT_DEV_MODE    = false
	DEFAULT_LOG_LEVEL   = "info"
	DEFAULT_HOST        = "localhost"
	DEFAULT_IP          = "127.0.0.1"
	DEFAULT_PORT        = 1443
	DEFAULT_TLS_KEY     = "../../testing/rsa/key.pem"
	DEFAULT_TLS_CERT    = "../../testing/rsa/cert.pem"
)

var (
	// Loader flags
	dev      optional.Bool
	confFile file.File
	logLevel optional.Str
	// subLoader flags
	host optional.Str
	ip   optional.Str
	port optional.Int
	cert file.Cert
	key  file.PrivateKey
)

func setupFlags() {
	// Set up all of our command line flags here to make sure we don't have config scattered across
	// the whole world.

	// Loader
	flag.Var(&dev, "dev", "Set default values appropriate for local development")
	flag.Var(&logLevel, "log", "set logging level [debug, info, warn, err]. Defaults to info.")
	flag.Var(&confFile, "config", "path to optional config file. Set to `none` to disable file loading entirely.")
	flag.Var(&host, "host", "hostname for a server or whatever")

	// subLoader
	flag.Var(&ip, "ip", "ip address for a server or whatever")
	flag.Var(&port, "port", "port for a server or whatever")
	flag.Var(&cert, "cert", "TLS certificate to use")
	flag.Var(&port, "key", "TLS key to use")
}

// Static and reloadable configs
type SubConfig struct {
	Host    string
	IP      string
	Port    int
	TlsConf *tls.Config
}

// If you has a module in your code that only needed this fragment of the config, you could add convenience
// methods to it like this:
func (c SubConfig) GetAddr() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

type Config struct {
	ConfigFile file.File
	DevMode    bool
	LogLevel   slog.Level
	SubConfig  SubConfig
}

type subLoader struct {
	Host    optional.Str    `env:"HOST"`
	IP      optional.Str    `env:"IP"`
	Port    optional.Int    `env:"PORT"`
	TlsCert file.Cert       `env:"TLS_CERTIFICATE"`
	TlsKey  file.PrivateKey `env:"TLS_KEY"`
}

func newSubLoader() subLoader {
	// We only need to actually specify fields we want to initialize to some flag-controlled variable.
	return subLoader{
		Host:    host,
		IP:      ip,
		Port:    port,
		TlsCert: cert,
		TlsKey:  key,
	}
}

func (l *subLoader) LoadFlags() {
	l.Host = optional.Or(host, l.Host)
	l.IP = optional.Or(ip, l.IP)
	l.Port = optional.Or(port, l.Port)
	l.TlsCert = optional.Or(cert, l.TlsCert)
	l.TlsKey = optional.Or(key, l.TlsKey)
}

func (l *subLoader) ToStatic() (SubConfig, error) {
	// we want flags to override all other options while also applying defaults if
	// no source set them.
	l.LoadFlags()
	l.TlsCert.Default(DEFAULT_TLS_CERT)
	l.TlsKey.Default(DEFAULT_TLS_KEY)

	cert, err := l.TlsKey.ReadCert(l.TlsCert)
	if err != nil {
		return SubConfig{}, err
	}

	tlsConf := &tls.Config{
		Certificates:             []tls.Certificate{cert},
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	return SubConfig{
		Host:    optional.GetOr(l.Host, DEFAULT_HOST),
		IP:      optional.GetOr(l.IP, DEFAULT_IP),
		Port:    optional.GetOr(l.Port, DEFAULT_PORT),
		TlsConf: tlsConf,
	}, nil
}

type Loader struct {
	ConfigFile file.File     `env:"CONFIG_FILE"`
	DevMode    optional.Bool `env:"DEV_MODE"`
	LogLevel   optional.Str  `env:"LOG_LEVEL"`
	SubLoader  subLoader
	current    Config
}

func NewLoader() Loader {
	return Loader{}

}

func (l Loader) Current() Config {
	return l.current
}

func (l *Loader) Update() error {
	if err := env.Load(l, nil); err != nil {
		slog.Error("Failed to load env vars")
		return err
	}

	// Prefer the flag setting, if present.
	// Also, convert relative paths into absolute to ensure it is possible
	// and so that users can see when relative paths do not point where they
	// expect based on log messages.
	l.ConfigFile = optional.Or(confFile, l.ConfigFile)
	tmp, err := l.ConfigFile.Abs()
	if err != nil {
		slog.Error(
			"Could not determine config file absolute file path",
			slog.Any("path", l.ConfigFile),
			slog.Any("error", err),
		)
		return err
	}
	l.ConfigFile = tmp

	pretty, err := json.Marshal(l)
	if err != nil {
		slog.Error("Failed to print current state of loader... Does it contain a non-marshallable type?")
	} else {
		slog.Debug("Loaded env vars", "loader", string(pretty))
	}

	path, ok := l.ConfigFile.Get()
	if ok {
		_, err = toml.DecodeFile(path, l)
		if err != nil {
			slog.Error(
				"Could not decode toml file into Loader",
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

func (l *Loader) LoadFlags() {
	// ConfigFile is also set from the flag when a Loader is created, but if both
	// flag and env vars changed it we want to reset it to the flag value.
	l.ConfigFile = optional.Or(confFile, l.ConfigFile)
	l.LogLevel = optional.Or(logLevel, l.LogLevel)
	l.DevMode = optional.Or(dev, l.DevMode)
}

func (l Loader) ToStatic() (Config, error) {
	// We want flags to win no matter what
	l.LoadFlags()

	subconfig, err := l.SubLoader.ToStatic()
	if err != nil {
		slog.Error("Failed loading sub-loader", "error", err)
	}

	// Ideally there would be a better way to map log level strings to the levels. Maybe there is, and I'm just not familiar with it?
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

	// DevMode overrides log level to debug. Bool optionals have a special method True() for convenience.
	if l.DevMode.True() {
		level = slog.LevelDebug
	}

	return Config{
		ConfigFile: l.ConfigFile,
		DevMode:    optional.GetOr(l.DevMode, false),
		LogLevel:   level,
		SubConfig:  subconfig,
	}, nil
}
