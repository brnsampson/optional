# optional
Go library for working with optional values

## What?

Have you ever needed to represent "something or nothing"? It's common in go to use a pointer for this.

This package aims to provide an alternative, which is the "Optional". An Optional puts another type in a box which can
either be "Some" or "None". That, and a small handful of methods to interact with the inner value, is pretty much it.

The goal of this isn't to be particularly performant. This ain't no rust style "zero-cost abstraction". Don't use it in
the hot paths of your code - that's where the pointer trick goes.

The goal _is_ to provide the functionality that works in an intuitive way without any of the gotchas and questions that
using a pointer would bring so that you can implement what you need quickly without any thought. If you find performance
is an issue, then you can always replace these with pointers later.

Really, this is a generic solution to a common problem: How do you represent _NOT_ a thing? Often we can find a particular
value in the domain of our type which is not a valid value for our application and use that. Have a field for email address?
An empty string can be used as your None value. Port number? Pick anything above 65,535 to represent "None selected". Sometimes,
however, there is no obvious invalid value of your type in your use case. Moreover, you sometimes need to differentiate
between "default" and "None", and in those cases you actually need two values which are both invalid for your application!

Disclamer: I'm sure you could find exciting new ways to shoot yourself in the foot if you tried to create an optional
of a pointer type (or a struct which contains pointer fields). You can definitely do it, I just don't advise doing so
unless you know what you are doing.

## State of the art

Currently, some other packages do provide a similar experience. There are a few issues that I have which are not addressed
by the ones I am aware of, however.

https://github.com/googleapis/google-cloud-go/tree/main/internal/optional

Pros: uh...
Cons: internal to a repo? Requires you to call a function that does type checking everytime you want to use a value? Meh.

https://github.com/markphelps/optional

Pros: should be somewhat performant due to just using a pointer under the hood, ergonomics pretty good
Cons: one optional type per underlying type. Requires a go generate script to support any non-implemented type. Doesn't
    support some useful features like applying transforms to the value inside the option. You could be surprised if you
    store one of these types in a struct since they are values wrapping pointers, and that would be a nasty debugging session.

https://github.com/leighmcculloch/go-optional

Pros: Generic so it can be applied to any type, short and sweet, supports JSON and XML marshaling/unmarshaling
Cons: Generic over `any`, so things like equality can't be supported. I don't know if using an internal array is any more
    efficient than a boolean to represent Some/None. The Marshaling/Unmarshaling uses the zero values to represent None
    instead of `null`, removing much of the benefit there. Methods like `String()` uses Sprintf, which uses reflection.
    There is poor performance, then there is reflection performance.

Magic Numbers:

Pros: Easy to use by just defining a const in your package. Possibly the most efficient way to do this.
Cons: You don't always an invalid value to use as your magic number, so it's just impossible sometimes.

Pointers:

Pros: Effective and fast. Just a nil check tells you if a value is set or not.
Cons: Nil checks everywhere. It's on you to check for nil before _every_ use, and the consequences of forgetting is a nil
    pointer dereference panic. Makes for more difficult to read and reason about code when things get complicated. If you
    every pass a struct by value your invariants can be broken by methods modifying some fields only in the copied struct
    and others in the copied and original struct. You need to be _very_ careful if you do that.

This package:

Pros: Ergonomics pretty good, no surprises, generic implementation means all types can be used in the same way. Methods
    for performing transforms on data without extracting it first. Specialized optional types can be built on top to
    provide any needed functionality for specific use cases.
Cons: Inefficient in terms of space and likely performance. Core Option type is limited in implementing convenient
    stdlib interfaces due to the use of generics.

## Why do you even have Unwrap()? It seems like a less convenient version of Get().

Yep. On the one hand, I like the idea of options being values. They could be passed around and referenced easily without
worrying about consequences.

In practice, I find I actually need pointers to options most of the time. They are things meant to be consumed rather than
passed around a lot, and as such when the thing that actually needs them uses them, I don't really want them to stick
around going forward.


## Why didn't you just wrap a pointer then do the right thing? Isn't copying things around by value all the time expensive?

1. There is a whole world of hurt in golang around structs with pointer fields. If someone is so foolish as to blindly
pass such a thing around by value, bad things can happen quickly.
2. As such, I initially tried to only return pointers to optionals, but given the methods I wanted to provide this didn't
always work well.
3. You can always make an option with a pointer inner type if you want that. It probably isn't totally safe in all cases tho.

## Where?

Where Optionals really shine, are situations which are not performance sensitive in the first place. Configuration
is probably the best example, particularly if you are drawing configs from multiple sources and merging them together.
I actually first wrote this in a fit of irritation after spending more time tracking down a nil-pointer deference in
a configuration package for longer than it took to write the package.

## How?

If you have a main package like this:

```go
package main

import (
    "os"
    "my/path/pkg/config"
    flag "github.com/spf13/pflag"
)

func main() {
    fs := flag.NewFlagSet("application", flag.PanicOnError)
    h := fs.String("host", "", "Host for whatever")
    fs.Parse(os.Args[1:])

    loader := config.NewConfigLoader()
    if err := loader.LoadFromFlags(fs); err != nil {
        panic()
    }
    if err := loader.LoadFromEnv(); err != nil {
        panic()
    }
    if err := loader.LoadFromFile("path/to/my/conf.toml"); err != nil {
        panic()
    }
    conf, err := loader.Finalize()
    if err != nil {
        panic()
    }

    // We have a static config! Easy, right?!?! ...right?
    fmt.Println(conf.Host)
}
```

You might create a configuration package like this:
```go
package config

import (
	"github.com/spf13/pflag"
	"github.com/caarlos0/env"
    "github.com/BurntSushi/toml"
    "github.com/charmbracelet/log"
    "github.com/brnsampson/optional/config"
)

const DEFAULT_HOST string = "localhost"

type ConfigLoader struct {
	Host config.StringOption `env:"HOST"`
}

type Config struct {
    Host string
}

// Fun fact, if your config has a sub-struct and you implement the same methods on it, it works pretty well with toml
// and the env / flag loading methods kind of just nest together like a russian doll.
func NewConfigLoader() *ConfigLoader {
    return ConfigLoager{ config.StringNone() }
}

func (c *ConfigLoader) LoadFromFlags(flags *pflag.FlagSet) error {
    tmp, err := flags.GetString("host")
	if err != nil || tmp == "" {
		log.Debug("Failed to load host from flags")
	} else {
        c.Host = option.NewStringOption(tmp)
    }

    return nil
}

func (c *ConfigLoader) LoadFromFile(filepath string) error {
	if _, err := os.Stat(filepath); errors.Is(os.ErrNotExist, err) {
		return err
	}

    _, err := toml.DecodeFile(filepath, &c)
    if err != nil {
		return err
    }

	log.Info("Loaded server config from file", "filename", filepath, "config", c)

	return nil
}

func (c *ConfigLoader) LoadFromEnv() error {
	if err := env.Parse(c); err != nil {
		log.Error("Failed to load server config from env variables!")
		return err
	}

	log.Debug("Loaded server config from env variables", "config", c)
	return nil
}

func (c ConfigLoader) Finalize() (*Config, error) {
    host := c.Host.UnwrapOrDefault(DEFAULT_HOST)
    return &Config{host}, nil
}

```

## Why?

There are many great things about golang, but being spoiled by choice is not typically one of them.

Not that there is anything wrong with having a good stdlib implementation that can actually be used in prod or there
is a lack of good libraries available; I mean there is literally no good way to represent "maybe a thing".

Sure, you can get by in most cases with a pointer. This is, in fact, the most performant way to do this. I do it all the time.

I kept hitting the same few situations, however, where this just caused me grief. The biggest was (is) configuration.

I get that many people feel strongly about configuration. Sometimes just having a quick TOML parser is all you need.
Sometimes you are deploying to something like nomad or kubernetes and get better functionality by using env vars.
Sometimes you are working at a company that has an in-house configuration system or has standardized on consul and you
have no choice. Everyone knows junior devs love nothing more than pre-planning a future outage by using YAML configs
and including a float field. And whatever you choose, you will probably also be accepting command-line flags. Flags are
just more configuration!

Personally, I don't really care. I find that I have a few goals with configuration, and accomplishing them in the
best way usually requires at least flags, env vars, and some kind of config file or networked configuration system.

Flags: Determines which entrypoint to use for code execution and also allows for a human override for any other configuration.
Think about your favorite single-binary system like consul, k8s, etc. The flag is what determines if a given execution
will run a server or worker node. It is very clunky to use anything else to do this.

Env vars: Where your code runs is important. What datacenter are you in? What region? Is this production or dev? You
_could_ have some system merging the environment configuration with the application configuration, but why? I have
seen more than enough puppet or chef configurations trying to handle application config to last a lifetime. But everyone
uses orchestration systems like k8s now right? I'm tired of seeing a pod with a handful of side-cars that all have the same
environment distilled into their own little files. Humans developed environment variables to solve a real problem!

The wild west: Use YAML, TOML, JSON, BSON, JSON5, Consul, Vault, Zookeeper, etcd, whatever bad distributed key/value
store your company implemented on top of redis, or even DynoDB if you are a psycopath. I'm not gunna sweat it, because
if you use flags and env vars properly the only things that cares about this stuff is your application and I, as an
infrastructure engineer, am never going to have to interact with this. As long as you don't break your own config then
ask me to fix it.

Actually, on second thought just use TOML or whatever the company told you to use. You're probably just using this for
boolean feature flags anyways.

## Generating the keys and certs for testing

### RSA

```bash
openssl genrsa -out tls/rsa/key.pem 4096
openssl rsa -in tls/rsa/key.pem -pubout -out tls/rsa/pubkey.pem
openssl req -new -key tls/rsa/key.pem -x509 -sha256 -nodes -subj "/C=US/ST=California/L=Who knows/O=BS Workshops/OU=optional/CN=www.whobe.us" -days 365 -out tls/rsa/cert.pem
```

### ECDSA

```bash
openssl ecparam -name secp521r1 -genkey -noout -out tls/ecdsa/key.pem
openssl ec -in tls/ecdsa/key.pem -pubout > tls/ecdsa/pub.pem
openssl req -new -key tls/ecdsa/key.pem -x509 -sha512 -nodes -subj "/C=US/ST=California/L=Who knows/O=BS Workshops/OU=optional/CN=www.whobe.us" -days 365 -out tls/ecdsa/cert.pem
```

### ED25519

```bash
openssl genpkey -algorithm ed25519 -out tls/ed25519/key.pem
openssl pkey -in tls/ed25519/key.pem -pubout -out tls/ed25519/pub.pem
openssl req -new -key tls/ed25519/key.pem -x509 -nodes -subj "/C=US/ST=California/L=Who knows/O=BS Workshops/OU=optional/CN=www.whobe.us" -days 365 -out tls/ed25519/cert.pem
```
