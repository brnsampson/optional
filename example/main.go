package main

import (
	"flag"
	"fmt"
	"github.com/brnsampson/optional/confopt"
	"github.com/charmbracelet/log"
)

const (
	DEFAULT_FLAG_CONFIG = "conf.toml"
)

var (
	debug bool
	dev bool
	DEV_HOST = confopt.SomeStr("127.0.0.1")
	DEV_PORT = confopt.SomeInt(8088)
)

func main() {
	confPath := confopt.SomeStr(DEFAULT_FLAG_CONFIG)
	host := confopt.NoStr()
	port := confopt.NoInt()
	flag.BoolVar(&debug, "debug", false, "set logging level to debug")
	flag.BoolVar(&dev, "dev", false, "Set default values appropriate for local development")
	flag.Var(&confPath, "confopt", "path to confopt file. Set to `none` to disable loading from confopt.")
	flag.Var(&host, "host", "hostname for a server or whatever")
	flag.Var(&port, "port", "port for a server or whatever")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	sl := NewSubConfigLoader().WithPort(port)
	flagConf := NewConfigLoader().WithHost(host)

	if dev {
		sl = sl.WithPort(DEV_PORT)
		flagConf = flagConf.WithHost(DEV_HOST)
	}

	flagConf = flagConf.WithNested(sl)

	log.Debug("Loaded sub confopt loader from flags", "confopt", sl)
	log.Debug("Loaded confopt loader from flags", "confopt", flagConf)


	config, err := NewConfig(flagConf)
	if err != nil {
		fmt.Println(err)
		panic("Error while finalizing the config!")
	}

	// Reload with config.Reload()
	err = config.Reload()
	if err != nil {
		fmt.Println(err)
		panic("Error while reloading config!")
	}

	// We have a reloadable config! Easy, right?!?! ...right?
	fmt.Println("")
	fmt.Println("I loaded a config and my name is:")
	fmt.Println(config.Host)
	fmt.Println("")
	fmt.Println("And my port is:")
	fmt.Println(config.Nested.Port)
}
