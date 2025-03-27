# optional

Go library for working with optional values

## How to use

```
import "github.com/brnsampson/optional"

// Create an Optional Int with no initial value
// The zero value of an optional is None
var i optional.Int

// This also works fine
i = optional.Int{}

// However, I normally use this functional form for symmetry with creating Options with values, optional.SomeInt(<my int>)
i = optional.NoInt()

// Check if i is None (empty)
if i.IsNone() {
  fmt.Println("i is empty!")
}

// Set the value to 42
i = optional.SomeInt(42)
// or
err := i.Replace(42)
if err != nil {
  fmt.Println("Hit an error while replacing the value of i!")
}

// Check if i has a value
if i.IsSome() {
  fmt.Println("i has a value!")
}

// Get i's value along with an 'ok' boolean telling us if the read is valid
val, ok := i.Get()
if ok {
  fmt.Println("Got i's value:")
  fmt.Println(val)
}

// Get i's value or a default value if i is None
tmp := optional.GetOr(i, 123)
fmt.Println("Got i's value or 123:")
fmt.Println(tmp)

// Get i's value or a default value AND set i to the default value if it is used
tmp, err = optional.GetOrInsert(i, 42)
if err != nil {
  fmt.Println("Error while replacing i's value with default")
} else {
  fmt.Println("Got i's value which should DEFINITELY be 42:")
  fmt.Println(tmp)
}

// Get i's value, but just panic if i is None
val = i.MustGet()

// For functions that automatically convert types into their string representation, the Option can be used directly:
fmt.Println("Printing i directly:")
fmt.Println(i)

// Modify i's value without having to unwrap it, do your thing, then re-wrap it. None's are not modified.
transform := func(x int) (int, error) { return x + 5, nil }
err = i.Transform(transform)
if err != nil {
  fmt.Println("Transformer function returned an error!")
}

// There is a helper function for if you want to replace all None values with a default value (the transformation
// is not applied to the default value)
err := optional.TransformOr(&o, transform, def)
if err != nil {
  fmt.Println("Transformer function returned an error OR a None value could not be updated to the default!")
}

// We can check to make sure i is now 42 + 5 + 5 = 52
if i.Match(52) {
  fmt.Println("i was indeed 52!")
} else {
  fmt.Println("uh-oh, what happened?")
}

// I won't show it here, but Json marshaling/unmarshaling works exactly as you would expect: null maps to None and
// any other value 'x' maps to Some(x)

// Other primative types work in the same manner, but the types of all method/function parameters are changed
// as you would expect.
```

## How to use the file package

```
import "github.com/brnsampson/optional"
import "github.com/brnsampson/optional/file"

// However we got it, we either have or do not have a path. For our example, let's assume we loaded this from a
// flag so we end up with a *string which could be nil

f := file.NoFile()
if path != nil {
  f = file.SomeFile(*path)
}

// Read back the path
p, ok := f.Get()
if ok {
  fmt.Println("Got path:")
  fmt.Println(p)
} else {
  fmt.Println("No path given!")
  os.Exit(1)
}

// Check if the given path is the same as some other path, matching all equivalent absolute and relative paths.
// In this case, check if the given path is equivalent to our working directory.
if f.Match(".") {
  fmt.Println("We are operating on our working directory. Be careful!")
}

// Get a new optional with any relative path converted to absolute path (also ensuring it is a valid path)
abs, err := f.Abs()
if err != nil {
  fmt.Println("Could not convert path into absolute path. Is it a valid path?")
  os.Exit(1)
}

// Stat the file, or just check if it exists if you don't care about other file info
// info, err := abs.Stat() // I don't care about the info
var opened *os.File
if abs.Exists() {
  opened, err = abs.Open()
  if err != nil {
    fmt.Println("Failed to open file:")
    fmt.Println(err)
    os.Exit(1)
  }
} else {
  opened, err = abs.Create()
  if err != nil {
    fmt.Println("Failed to create new file:")
    fmt.Println(err)
    os.Exit(1)
  }
}
defer opened.Close()

// Check that the file has permissions 700 and modify it if it does not
valid, err := abs.FilePermsValid(0700)
if err != nil {
  fmt.Println("Could not read file permissions!")
  os.Exit(1)
}

if !valid {
  err = abs.SetFilePerms(0700)
  if err != nil {
    fmt.Println("Failed to set file perms to 700")
    os.Exit(1)
  }
}
```

## How to load certificates and keys

```
import "github.com/brnsampson/optional/file"

// Similarly to the File type, the Cert and PrivateKey types make loading and using optional certificates
// easier and more intuitive. They both embed the Pem struct, which handles the loading of Pem format files.

// Create a Cert from a flag which requested the user to give the path to the certificate file.
certfile := file.NoCert()
if certPath != nil {
  certfile := file.SomeCert(*certPath)
}

// We can use all the same methods as the File type above, but it isn't necessary to go through all of the
// steps individually. The Cert type knows to check that the path is set, the file exists, and that the file permissions
// are correct as part of loading the certificates.
//
// certificates are returned as a []*x509.Certificate from the file now.
// Incidentally, we could write new certs to the file with certfile.WriteCerts(certs)
certs, err := certfile.ReadCerts()
if err != nil {
  fmt.Println("Error while reading certificates from file:")
  fmt.Println(err)
  os.Exit(1)
}

// Now we want to load a tls certificate. We typically need two files for this, the certificate(s) and private keyfile.
// Note: this specifically is for PEM format keys. There are other ways to store keys, but we have not yet implemented
// support for those. We do support most types of PEM encoded keyfiles though.

var privKeyFile file.PrivateKey  // Effectively the same as privKeyFile := file.NoPrivateKey()
if privKeyPath != nil {
  privKeyFile = file.SomePrivateKey(*privKeyPath)
}

// Again, we could manually do all the validity checks but those are also run as part of loading the TLS certificate.
// cert is of the type *tls.Certificate, not to be confused with *x509Certificate.
cert, err := privKey.ReadCert()
if err != nil {
  fmt.Println("Error while generating TLS certificate from PEM format key/cert files:")
  fmt.Println(err)
  os.Exit(1)
}

// Now we are ready to start up an TLS sever
tlsConf := tls.Config{
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

httpServ := &http.Server{
	Addr:         addr,
	TLSConfig:    tlsConf,
}

// The parameters ListenAndServeTLS takes are the cert file and keyfile, which may lead you to ask, "why did we bother
// with all of this then?" Essentially, we were able to do all of our validation and logic with our configuration
// loading and can put our http server somewhere that makes more sense without just getting panics in our server code
// when the user passes us an invalid path or something. We are also able to get more granular error messages than just
// "the server is panicing for some reason."
if e := httpServ.ListenAndServeTLS("", ""); e != nil {
  fmt.Println("TLS server exited")
  fmt.Println(err)
}

// In some situations you actually want to use a public/private keypair for signing instead.
// Here is how we would load those:
var privKeyFile file.PrivateKey  // Effectively the same as privKeyFile := file.NoPrivateKey()
if privKeyPath != nil {
  privKeyFile = file.SomePrivateKey(*privKeyPath)
}

var pubKeyFile file.PubKey  // Effectively the same as pubKeyFile := file.NoPubKey()
if pubKeyPath != nil {
  pubKeyFile = file.SomePubKey(*pubKeyPath)
}

// NOTE: as is usually the case with golang key loading, this returns pubKey as a []any and you have to kind of
// just know how to handle it yourself.
pubKeys, err := pubKeyFile.ReadPublicKeys()
if err != nil {
  fmt.Println("Error while reading public key(s) from file:")
  fmt.Println(err)
  os.Exit(1)
}

// While a public key file may have multiple public keys, private key files should only have a single key. This
// key is also returned as an any type which you will then need to sort out how to use just like any other key
// loading.
privKey, err := privKeyFile.ReadPrivateKey()
if err != nil {
  fmt.Println("Error while reading private key from file:")
  fmt.Println(err)
  os.Exit(1)
}
```

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
are is a lot of boilerplate for what it does, but nothing is going to go wrong and it is fully extendable to a real
world project without needing any complicated additional libraries involved. Don't get me wrong, Cobra and Viper are
super powerful and well maintained, but 98% of the time I only really want one command per executable and every time I
touch Cobra or Viper I spend at least half an hour reading through documentation. Why bother for the vast majority of my
public work just involves [stupid things](https://github.com/brnsampson/go-partyparrot) like a slack-bot to render text
as [party parrots](https://cultofthepartyparrot.com/)?

Note that while this looks like a lot of code for what it does, it does have a good set of functionality for a small project:

- Clear precedence of config sources
- Only basic `flag` library used
- Flag for debug output
- Flag to choose config file (including option to skip file loading without needing a separate flag!)
- Annotations for env var mapping and file loading of the config are defined by the loader struct itself
- Default values kept directly above config loader structs for easy comparison
- No super ugly long parameter sets to pass from flag parsing to initialize structs (builder pattern preferred)
- Reloadable, so if you create an init ConfigLoader from flags you can then do hot reloads in response to e.g. a SIGHUP (see example/main.go)
- No magic hidden in a library you need to look up

If you want to try it, it was written so that the precedence is flags > file > env. The default values in the code are
{ Host: "localhost", Port: 1443 }.

Try running it a few different ways and seeing what happens!

```bash
go run ./example
PORT=3000 go run ./example
PORT=3000 go run ./example --port 5000
PORT=3000 go run ./example --port 5000 --host "example.com"
PORT=3000 HOST=host.from.env go run ./example
PORT=3000 HOST=host.from.env go run ./example --config ./example/alt.toml
HOST=host.from.env go run ./example --config none
PORT=3000 HOST=host.from.env go run ./example --config ./example/alt.toml --port 5000 --host host.from.flag
PORT=3000 HOST=host.from.env go run ./example --config ./example/alt.toml --port 5000 --host host.from.flag
PORT=3000 HOST=host.from.env go run ./example --log debug --config ./example/alt.toml --port 5000 --host host.from.flag
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
- I don't know if using an internal array is any more efficient than a boolean or pointer to represent Some/None
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
- I hope you are documenting all of your magic numbers somewhere, because if anyone else looks at your code (including
  you in 6 months), they probably will have to go code spelunking to understand what is happening.

### Pointers

#### Pros

- Effective and fast. Just a nil check tells you if a value is set or not.

#### Cons

- Nil checks everywhere. It's on you to check for nil before _every_ use, and the consequences of forgetting is a nil
  pointer dereference panic.
- Makes for more difficult to read and reason about code at a surprisingly low level of complexity
- If you every pass a struct by value your invariants can be broken by methods modifying some fields only in the copied struct
  and others in the copied and original struct. You need to be _very_ careful if you do that.
- Encourages the "make everything a pointer" style, which encourages wishing you were using another language that
  doesn't require 10 nil checks in every function.

### This package

#### Pros

- Ergonomics pretty good. Built with merging values from multiple sources together, marshaling, and templating in mind
- No surprises
- Generic implementation means all derivative Option types can be used in the same way
- Methods for performing transforms on data without extracting it first, which is nice in loops
- Specialized optional types can be built on top to provide any needed functionality for specific use cases. Look at the
  file sub-package for an example.

#### Cons

- A bit inefficient in terms of space and performance
- Core `Option` type is unable to implement some convenient stdlib interfaces due to the use of generics.
- To make the thing more useful, only `comparable` values can be wrapped currently. This isn't usually too big of a
  deal for most _values_, but does mean that you cannot create an array option for example
  (but why would you do that?!? Just check for zero length!)

## Why? or: the BS configuration manifesto

There are many great things about golang, but being spoiled by choice is not typically one of them.

Not that there is anything wrong with having a good stdlib implementation that can actually be used in prod or there
is a lack of good libraries available; I mean there is literally no good way to represent "maybe a thing". If TypeScript
can have Optional values, why can't we?

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
store your company implemented on top of redis, or even DynamoDB if you are a psycopath. I'm not gunna sweat it, because
if you use flags and env vars properly the only things that cares about this stuff is your application and I, as an
dis-interested third party, am never going to have to interact with this. As long as you don't break your own config then
ask me to fix it. You can do that, I suppose, but know that some poor person that is taking time out of keeping that
old elasticsearch cluster alive for one more day is judging you hard.

Actually, on second thought just use TOML or whatever the company told you to use. You're probably just using it for
boolean feature flags, a listening address, and anyways.

## FAQ's

### What happens if I have a pointer to an Option?

It acts like a pointer to any other value. It is equally useful as a pointer to an int, though, so I wouldn't recommend
it.

In practice, I don't find that I actually need pointers to options very often. They are things meant to be calculated
once then immutably consumed rather than having a value mutated in a bunch of places. We have some source of truth which
we are trying to represent to different parts of our code, marshal and send over the wire, or cache/invalidate a local
state of that source. Values work just fine for that.

The one exception would be if you have a very large struct which you want to wrap in an Option. Accepting very large
things through a function call may not perform the way you want, in which case you could use a `*Option[MyStruct]` the
same way you might use a `*MyStruct`.

### What about a pointer to a struct which contains an Option?

It acts like any other value.

### What about an Option of a pointer?

That's tricky. It has all the same dangers as passing by value a struct with pointer fields. While it is technically
possible, I don't recommend it unless you have a very specific need and know your foot guns well.

Be aware that if you _do_ put a pointer inside an Option, `nil` _is a valid value and is NOT None_. This means that if
you call `myOption.IsNone` does not tell you if the inner value is a nil pointer.

Additionally, I bet there are complications with marshaling/unmarshaling. Would a json `null` unmarshal to a None-value
Option, or a Some-value Option with the value being a `nil` pointer? I'm not sure, and I'm not taking the time to think
about it.

You could make a specific Pointer derivative of the Option type that handles that sort of thing, but I just don't
know what the use case would be where you couldn't just use a pointer on it's own. I certainly don't plan on doing that.

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
