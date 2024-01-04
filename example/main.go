package main

import (
    "os"
	"fmt"
    flag "github.com/spf13/pflag"
)

func main() {
    fs := flag.NewFlagSet("application", flag.PanicOnError)
    _ = fs.String("host", "", "Host for whatever")
    _ = fs.Int("port", 0, "Port for whatever")
    fs.Parse(os.Args[1:])

    loader := NewConfigLoader()
    if err := loader.LoadAll(fs, "conf.toml"); err != nil {
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
