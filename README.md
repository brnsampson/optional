# optional
Go library for working with optional values

## What?

Have you ever needed to represent "something or nothing"? It's common in go to use a pointer for this, but in some
situations that can add complexity or is otherwise undesirable.

This package aims to provide an alternative, which is the "Optional". An Optional puts another type in a box which can
either be "Some" or "None". That, and a small handful of methods to interact with the inner value, is pretty much it.

The goal of this isn't to be particularly performant. This ain't no rust style "zero-cost abstraction". Don't use it in
the hot paths of your code - that's where the pointer trick goes.

The goal _is_ to provide the functionality that works in an intuitive way without any of the gotchas and questions that
using a pointer would bring so that you can implement what you need quickly without any thought. If you find performance
is an issue, then you can always replace these with pointers later.

Really, this is a generic solution to a common problem: How do you represent _NOT_ a thing? Often we can find a particular
"magic value" in the domain of our type which is not a valid value for our application and use that. Have a field for
email address? An empty string can be used as your None value. Port number? Pick anything above 65,535 to represent
"None selected". Sometimes, however, there is no obvious invalid value of your type in your use case. Moreover, you
sometimes need to differentiate between "default" and "None", and in those cases you actually need two values which are
both invalid for your application!

Disclamer: I'm sure you could find exciting new ways to shoot yourself in the foot if you tried to create an optional
of a pointer type (or a struct which contains pointer fields). You can definitely do it, I just don't advise doing so
unless you know what you are doing.

## Where?

Where Optionals really shine, are situations which are not performance sensitive in the first place. Configuration
is probably the best example, particularly if you are drawing configs from multiple sources and merging them together.
I actually first wrote this in a fit of irritation after tracking down a nil-pointer deference in a configuration
package for longer than it took to write the package.

## How?

See the example/ directory for a setup for brain-dead config parsing. Sure, it's verbose for two real parameters and there
are is a lot of boilderplate for what it does, but nothing is going to go wrong and it is fully extendable to a real
world project without needing any complicated additional libraries involved. Don't get me wrong, Cobra and Viper are
super powerful and well maintained, but 98% of the time I only really want one command per executable and every time I
touch Cobra or Viper I spend at least half an hour reading through documentation. Why bother for the vast majority of my
work what just involves [stupid things](https://github.com/brnsampson/go-partyparrot) like a slack-bot to render text
as [party parrots](https://cultofthepartyparrot.com/)?

Note that while this looks like a lot of code for what it does, it does have a good set of functionality for a small project:
 * Clear precedence of config sources
 * Only basic `flag` library used
 * Flag for debug output
 * Flag to choose config file (including option to skip file loading without needing a separate flag!)
 * Annotations for env var mapping and file loading of the config are defined by the loader struct itself
 * Default values kept directly above config loader structs for easy comparison
 * No super ugly long parameter sets to pass from flag parsing to initialize structs (builder pattern preferred)
 * Reloadable, so if you create an init ConfigLoader from flags you can then do hot reloads in response to e.g. a SIGHUP (see example/main.go)
 * No magic hidden in a library you need to look up

If you want to try it, it was written so that the precedence is flags > file > env. The default values in the code are
{ Host: "localhost", Port: 1443 }.

Try running it a few different ways and seeing what happens!
```bash
go run ./example
PORT=3000 go run ./example
PORT=3000 go run ./example --port 5000
PORT=3000 go run ./example --port 5000 --host "example.com"
PORT=3000 HOST=host.from.env go run ./example
PORT=3000 HOST=host.from.env go run ./example --config alt.toml
HOST=host.from.env go run ./example --config none
PORT=3000 HOST=host.from.env go run ./example --config alt.toml --port 5000 --host host.from.flag
PORT=3000 HOST=host.from.env go run ./example --config alt.toml --port 5000 --host host.from.flag
PORT=3000 HOST=host.from.env go run ./example --debug --config alt.toml --port 5000 --host host.from.flag
```

## State of the art

Currently, some other packages do provide a similar experience. There are a few issues that I have which are not addressed
by the ones I am aware of, however.

### [markphelips/optional](https://github.com/markphelps/optional)

#### Pros
 - should be somewhat performant due to just using a pointer under the hood
 - ergonomics pretty good

#### Cons
 - one optional type per underlying type
 - Requires a go generate script to support any non-implemented type
 - Doesn't support some useful features like applying transforms to the value inside the option
 - You could be surprised if you store one of these types in a struct since they are values wrapping pointers, and that would be a nasty debugging session.

### [leighmcculloch go-optional](https://github.com/leighmcculloch/go-optional)

#### Pros
 - Generic so it can be applied to any type
 - short and sweet
 - supports JSON and XML marshaling/unmarshaling

#### Cons
 - Generic over `any`, so things like equality can't be supported
 - I don't know if using an internal array is any more efficient than a boolean to represent Some/None
 - The Marshaling/Unmarshaling uses the zero values to represent None instead of `null`, removing much of the benefit there
 - Methods like `String()` uses Sprintf, which uses reflection. There is poor performance, then there is reflection performance.

### Magic Values
This is just when you use a specific value to represent your `None`. Sometimes this makes sense, such as when you have
a string field where an empty string would be meaningless.

#### Pros
 - Easy to use by just defining a const in your package
 - Possibly the most efficient way to do this

#### Cons
 - You don't always an invalid value to use as your magic value, so it's just impossible sometimes.
 - Ergonomics can get messy; different libraries may return different magic values which you have to translate between

### Pointers

#### Pros
 - Effective and fast. Just a nil check tells you if a value is set or not.

#### Cons
 - Nil checks everywhere. It's on you to check for nil before _every_ use, and the consequences of forgetting is a nil
pointer dereference panic.
 - Makes for more difficult to read and reason about code at a surprisingly low level of complexity
  - If you every pass a struct by value your invariants can be broken by methods modifying some fields only in the copied struct
    and others in the copied and original struct. You need to be _very_ careful if you do that.

### This package

#### Pros
 - Ergonomics pretty good
 - No surprises
 - Generic implementation means all types can be used in the same way
 - Methods for performing transforms on data without extracting it first
 - Specialized optional types can be built on top to provide any needed functionality for specific use cases.

#### Cons
 - Inefficient in terms of space and performance
 - Core `Option` type is limited in implementing convenient stdlib interfaces due to the use of generics.
 - To make the thing more useful, only `comparable` values can be wrapped currently. This isn't usually too big of a
deal for most _values_, but does mean that you cannot create an array option for example
(but why would you do that?!? Just check for zero len!)

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
will run a dev server or production node. It is very clunky to use anything else to do this.

Env vars: Where your code runs is important. What datacenter are you in? What region? Is this production or dev? You
_could_ have some system merging the environment configuration with the application configuration, but why? I have
seen more than enough puppet or chef configurations trying to handle application config to last a lifetime. But everyone
uses orchestration systems like k8s now right? I'm tired of seeing a pod with a handful of side-cars that all have the same
environment distilled into their own little files. Humans developed environment variables to solve a real problem, and
they are really good at solving those problems!

The wild west: Use YAML, TOML, JSON, BSON, JSON5, Consul, Vault, Zookeeper, etcd, whatever bad distributed key/value
store your company implemented on top of redis, or even DynoDB if you are a psycopath. I'm not gunna sweat it, because
if you use flags and env vars properly the only things that cares about this stuff is your application and I, as an
infrastructure engineer, am never going to have to interact with this. As long as you don't break your own config then
ask me to fix it. You can do that, I suppose, but know that some poor person that is taking time out of keeping that
old elasticsearch cluster alive for one more day is judging you hard.

Actually, on second thought just use TOML or whatever the company told you to use. You're probably just using it for
boolean feature flags anyways.

## FAQ's

### Why do you even have Unwrap()? It seems like a less convenient version of Get().

Yep. On the one hand, I like the idea of options being values. They could be passed around and referenced easily without
worrying about consequences.

In practice, I find I actually need pointers to options most of the time. They are things meant to be consumed rather than
passed around a lot, and as such when the thing that actually needs them uses them, I don't really want them to stick
around going forward.

As an example, consider a situation where you have a single config loader generating a nested config which includes
pointers. Maybe one field of your config is a sub-config meant for one component or domain of your application. If that
struct has pointer fields, then even if the struct is passed into your component by value there will be another pointer
floating around and accessible in your other code. That _should_ be okay, right? Surely nothing else will call a method
on the main config struct that might modify its fields right?

Well, maybe we should just make that a pointer so that we can nil it out after the component consumes its config. That
works, but now we have a nil pointer floating around and something else can still call methods on the main struct to modify
fields. Did you check that pointer _everywhere_ to make sure you are not dereferencing it when it's nil? Are you _sure_?

Well hey, just make the field an Optional and when the component consumes its part it can call Unwrap() to get the inner
value and the rest of the struct just has a None Optional left. It has all the same methods and code has to handle the
errors returned just like normal, but now it won't cause an unhandled panic if you lose focus for more than 10 seconds!


### Why didn't you just wrap a pointer then do the right thing? Isn't copying things around by value all the time expensive?

1. There is a whole world of hurt in golang around structs with pointer fields. If someone is so foolish as to blindly
pass such a thing around by value, bad things can happen quickly.
2. As such, I initially tried to only return pointers to optionals, but given the methods I wanted to provide this didn't
always work well.
3. You can always make an option with a pointer inner type if you want that. It probably isn't totally safe in all cases
and there is a very real chance that it won't do what you want. YMMV!

## Generating the keys and certs for testing

This is mostly a reminder for myself, given that the certs only have a lifetime of one year.

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
