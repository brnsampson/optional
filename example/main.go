package main

import (
	"fmt"
	"github.com/brnsampson/optional/config"
    "github.com/charmbracelet/log"
	"flag"
)

const (
	DEFAULT_FLAG_CONFIG = "conf.toml"
)

var (
	debug bool
)

func main() {
	confPath := config.SomeStr(DEFAULT_FLAG_CONFIG)
	host := config.NoStr()
	port := config.NoInt()
	flag.BoolVar(&debug, "debug", false,"set logging level to debug")
	flag.Var(&confPath, "config", "path to config file. Set to `none` to disable loading from config.")
	flag.Var(&host, "host", "hostname for a server or whatever")
	flag.Var(&port, "port", "port for a server or whatever")
    flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	sl := NewSubConfigLoader().WithPort(port)
	flagConf := NewConfigLoader().WithHost(host).WithNested(sl)

	log.Debug("Loaded sub config loader from flags", "config", sl)
	log.Debug("Loaded config loader from flags", "config", flagConf)

    loader, err := LoadedConfigLoader(confPath, flagConf)
    if err != nil {
		fmt.Println(err)
        panic("Error while loading config!")
    }

    conf, err := loader.Finalize()
    if err != nil {
		fmt.Println(err)
        panic("Error while finalizing config!")
    }

    // We have a static config! Easy, right?!?! ...right?
	fmt.Println("")
	fmt.Println("I loaded a config and my name is:")
    fmt.Println(conf.Host)
	fmt.Println("")
	fmt.Println("And my port is:")
    fmt.Println(conf.Nested.Port)
}
