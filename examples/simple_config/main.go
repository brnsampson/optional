package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

const (
	DEFAULT_FLAG_CONFIG = "conf.toml"
)

var (
	help bool
	h    bool
)

func main() {
	setupFlags()

	flag.BoolVar(&help, "help", false, "Get usage message")
	flag.BoolVar(&h, "h", false, "Get usage message")
	flag.Parse()

	if help || h {
		flag.Usage()
		os.Exit(0)
	}

	loader := NewLoader()
	slog.Info("Initialized loader", "loader", loader)
	slog.Info("Initialized sub-loader", "loader", loader.SubLoader)

	err := loader.Update("")
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(2)
	} else {
		slog.Info("Loading completed")
	}
	config := loader.Current()

	// Set up logging
	programLevel := new(slog.LevelVar)
	programLevel.Set(config.LogLevel)
	// h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	logger := slog.New(h)
	slog.SetDefault(logger)

	slog.Debug("Loaded sub-config", "config", config.SubConfig)
	slog.Debug("Loaded config", "config", config)

	// We have a reloadable config! Easy, right?!?! ...right?
	fmt.Println("")
	fmt.Println("I loaded a config and my host is:")
	fmt.Println(config.Host)
	fmt.Println("")
	fmt.Println("And my address is:")
	fmt.Println(config.SubConfig.GetAddr())
}
