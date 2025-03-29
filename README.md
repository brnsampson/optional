# optional

Go library for working with optional values

## How to use

<details>
  <summary>Importing</summary>

```golang
import (
	"github.com/brnsampson/optional"      // For primative types
	"github.com/brnsampson/optional/file" // For files and certificates
)
```

</details>

<details>
  <summary>Setting optional values</summary>

```golang
	// There are types for all the primatives we would normally expect
	// Bool
	// Int Int16 Int32 Int64
	// Uint Uint16 Uint32 Uint64
	// Float32 Float64
	// Str
	// Time
	// and the generic Option[T comparable]

	// Create an Optional Int with no initial value
	// The zero value of an optional is None
	var i optional.Int

	// This also works fine
	i = optional.Int{}

	// However, I normally use this functional form for symmetry with creating
	// Options with values, optional.SomeInt(<my int>)
	i = optional.NoInt()

	// Check if i is None (empty)
	if i.IsNone() {
		fmt.Println("i is empty!")
	}

	// Set the value to some default if it was previously unset
	i.Default(42)

	// Update the value and get the previous value back for any comparisons you might need to do
	previous := i.Replace(42)

	// Some methods like Option.Replace() return an Optional interface type. This erases the
	// concrete type and hides all of the methods which could mutate the value,
	// as the previous value is only provided as a reference. Unfortunatly, this also
	// hides some convenient things like the implemntations of TextMarshaler and Stringer
	if previous.IsSome() {
		fmt.Println("Replaced previous value:")
		fmt.Println(previous.MustGet())
	}

	// Overwrite the previous value without care
	i = optional.SomeInt(42)
```

</details>

<details>
  <summary>Inspecting the value of an Optional</summary>

```golang
	i := optional.SomeInt(42)

	// Check if i has a value
	if i.IsSome() {
		fmt.Println("i has a value!")
	}

	// We can check to make sure i is 42
	if i.Match(42) {
		fmt.Println("i was indeed 42!")
	} else {
		panic("wtf?")
	}

	// Get i's value along with an 'ok' boolean telling us if the read is valid
	val, ok := i.Get()
	if ok {
		fmt.Println("Got i's value:")
		fmt.Println(val)
	}

	// Get i's value, but just panic if i is None
	val = i.MustGet()

	// Get i's value or a default value if i is None
	tmp := optional.GetOr(i, 123)
	fmt.Println("Got i's value or 123:")
	fmt.Println(tmp)

	// Get i's value or a default value AND set i to the default value if it is used
	// Note that helper functions require a MutableOptional interface, which only Option
	// Pointers fulful. That should be a given, since it's just like passing an int;
	// you can't expect a function to modify an int, it can only return a new int.
	tmp, err := optional.GetOrInsert(&i, 42)
	if err != nil {
		fmt.Println("Error while replacing i's value with default")
	} else {
		fmt.Println("Got i's value which should DEFINITELY be 42:")
		fmt.Println(tmp)
	}

	// For functions that automatically convert types into their string representation,
	// the Option can be used directly:
	fmt.Println("Printing i directly:")
	fmt.Println(i)
```

</details>

<details>
  <summary>Marshaling values</summary>

```golang
	// Define our value and transformation first
	i := optional.SomeInt(42)
	f := optional.SomeFloat32(12.34)
	s := optional.SomeStr("Hello!")
	nope := optional.NoStr()

	// Let's create a text string first using Sprintf. We can't use more specific verbs like
	// %d or %f because we have no way to represent None.
	newString := fmt.Sprintf("i: %s, f: %s, s: %s, nothing: %s", i, f, s, nope)
	fmt.Println(newString)

	// Options do have TextMarshaler and String methods implemented though, so we can equally well use %v
	newString = fmt.Sprintf("i: %v, f: %v, s: %v, nothing: %v", i, f, s, nope)
	fmt.Println(newString)

	// Now let's marshal a json string
	type MyStruct struct {
		Int        optional.Int
		Float      optional.Float32
		GoodString optional.Str
		BadString  optional.Str
	}

	myStruct := MyStruct{i, f, s, nope}
	jsonEncoded, err := json.Marshal(myStruct)
	if err != nil {
		fmt.Println("Failed to marshal json from struct!")
	} else {
		fmt.Println(string(jsonEncoded))
	}
```

</details>

<details>
  <summary>Transforming Values</summary>

```golang
	// Define our value and transformation first
	i := optional.SomeInt(42)
	transform := func(x int) (int, error) { return x + 5, nil }

	// Modify the value in an Option without unpacking it
	err := i.Transform(transform)
	if err != nil {
		fmt.Println("The transform function returned an error!")
	}

	// Apply our transform to a slice of options, while modifying None values to be their index in the slice.
	// Remember, the zero value is None
	opts := make([]optional.Int, 10)
	for i, opt := range opts {
		// Functions which modify options in place should accept the MutableOptional interface which
		// is implemented by Option pointer types, such as this helper function. Try to use optional.TransformOr
		// with opt instead of &opt. It doesn't work, just in the same way that passing an int into a function
		// and expecting the integer to be changed in place doesn't work.
		err = optional.TransformOr(&opt, transform, i)
		if err != nil {
			fmt.Println("The transform function returned an error!")
		}
	}
```

</details>

## How to use the file package

<details>
  <summary>Loading and reading from files</summary>

```golang
	// However we got it, we either have or do not have a path. For our example, let's assume we loaded this from a
	// flag so we end up with a *string which could be nil

	f := file.NoFile()
	if path != nil {
		f = file.SomeFile(*path)
	}

	// Just read the contents of a file. Acts like os.ReadFile(path)
	data, err := f.ReadFile()
	if err != nil {
		fmt.Println("Failed to read from file:")
		fmt.Println(err)
	} else {
		fmt.Println("Got file contents:")
		fmt.Println(string(data))
	}

	// Open a file for reading. File.Open() works just like os.Open(path),
	// so the file is opend in ReadOnly mode.
	var opened *os.File
	opened, err = f.Open()
	if err != nil {
		fmt.Println("Failed to open file for reading:")
		fmt.Println(err)
		return
	}
	defer opened.Close()

	// Now use the file handle exactly as you would if you called os.Open()
```

</details>

<details>
  <summary>Writing and deleting files</summary>

```golang
	f := file.NoFile()
	if path != nil {
		f = file.SomeFile(*path)
	}

	// Delete a file. Works like os.Remove, but also returns an error if the path is still None
	err := f.Remove()
	if err != nil {
		fmt.Println("Failed to remove file:")
		fmt.Println(err)
	}

	// Write the contents of a file. Acts like os.WriteFile(path)
	data := []byte("Hello, World!")
	err = f.WriteFile(data, 0644)
	if err != nil {
		fmt.Println("Failed to write file:")
		fmt.Println(err)
	}

	// Open a file for read/write. File.Create() works like like os.Create(path), which means
	// calling this will either create a file or truncate an existing file. If you want to
	// append to a file, you must use File.OpenFile(os.O_RDWR|os.O_CREATE, 0644) in the same way
	// that would need to when calling os.OpenFile. See https://pkg.go.dev/os#OpenFile for details.
	var opened *os.File
	opened, err = f.Create()
	if err != nil {
		fmt.Println("Failed to open/create file:")
		fmt.Println(err)
		return
	}
	defer opened.Close()

	// Now use the file handle exactly as you would if you called os.Create(path)
	opened.Write(data)
```

</details>

<details>
  <summary>Other file tools</summary>

```golang
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
	} else {
		fmt.Println("We are not in our working directory. Go nuts!")
	}

	// Get a new optional with any relative path converted to absolute path (also ensuring it is a valid path)
	abs, err := f.Abs()
	if err != nil {
		fmt.Println("Could not convert path into absolute path. Is it a valid path?")
		os.Exit(1)
	}

	// Stat the file, or just check if it exists if you don't care about other file info
	if abs.Exists() {
		fmt.Println("The file exists!")
	}

	info, err := abs.Stat() // I don't care about the info
	if err != nil {
		fmt.Println("Could not stat the file")
	} else {
		fmt.Println("Got file info:")
		fmt.Println(info)
	}

	// Check that the file has permissions 700 and modify it if it does not
	valid, err := abs.FilePermsValid(0644)
	if err != nil {
		fmt.Println("Could not read file permissions!")
		os.Exit(1)
	}

	if !valid {
		err = abs.SetFilePerms(0644)
		if err != nil {
			fmt.Println("Failed to set file perms to 700")
			os.Exit(1)
		}
	}
```

</details>

## How to load certificates and keys

<details>
  <summary>Loading a TLS certificate</summary>

```golang
	// Similarly to the File type, the Cert and PrivateKey types make loading and using optional certificates
	// easier and more intuitive. They both embed the Pem struct, which handles the loading of Pem format files.

	// Create a Cert from a flag which requested the user to give the path to the certificate file.
	// Certs and Key Options also return an error if the path cannot be resolved to an
	// absolute path or the file permissions are not correct for a certificate or key file.
	certFile := file.NoCert()
	var err error
	if certPath != nil {
		certFile, err = file.SomeCert(*certPath)
		if err != nil {
			fmt.Println("Failed to initialize cert Option")
			fmt.Println(err)
			return err
		}
	}

	// We can use all the same methods as the File type above, but it isn't necessary to go through all of the
	// steps individually. The Cert type knows to check that the path is set, the file exists, and that the file permissions
	// are correct as part of loading the certificates.
	//
	// certificates are returned as a []*x509.Certificate from the file now.
	// Incidentally, we could write new certs to the file with certfile.WriteCerts(certs)
	certs, err := certFile.ReadCerts()
	if err != nil {
		fmt.Println("Error while reading certificates from file:")
		fmt.Println(err)
		return err
	} else {
		fmt.Println("Found this many certs:")
		fmt.Println(len(certs))
	}

	// Now we want to load a tls certificate. We typically need two files for this, the certificate(s) and private keyfile.
	// Note: this specifically is for PEM format keys. There are other ways to store keys, but we have not yet implemented
	// support for those. We do support most types of PEM encoded keyfiles though.

	// Certs and Key Options also return an error if the path cannot be resolved to an
	// absolute path or the file permissions are not correct for a certificate or key file.
	var keyFile file.PrivateKey // Effectively the same as privKeyFile := file.NoPrivateKey()
	if keyPath != nil {
		keyFile, err = file.SomePrivateKey(*keyPath)
		if err != nil {
			fmt.Println("Failed to initialize private key Option")
			fmt.Println(err)
			return err
		}
	}

	// Again, we could manually do all the validity checks but those are also run as part of loading the TLS certificate.
	// cert is of the type *tls.Certificate, not to be confused with *x509Certificate.
	cert, err := keyFile.ReadCert(certFile)
	if err != nil {
		fmt.Println("Error while generating TLS certificate from PEM format key/cert files:")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Full *tls.Certificate loaded")

	// Now we are ready to start up an TLS sever
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

	httpServ := &http.Server{
		Addr:      "127.0.0.1:3000",
		TLSConfig: tlsConf,
	}

	// The parameters ListenAndServeTLS takes are the cert file and keyfile, which may lead you to ask, "why did we bother
	// with all of this then?" Essentially, we were able to do all of our validation and logic with our configuration
	// loading and can put our http server somewhere that makes more sense without just getting panics in our server code
	// when the user passes us an invalid path or something. We are also able to get more granular error messages than just
	// "the server is panicing for some reason."

	fmt.Println("Deferring https server halting for 1 second...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		haltctx, haltcancel := context.WithTimeout(context.Background(), time.Second)
		defer haltcancel()
		if err := httpServ.Shutdown(haltctx); err != nil {
			fmt.Println("Error haling http server")
			fmt.Println(err)
		}
	}()

	fmt.Println("Starting to listen on https...")
	if err = httpServ.ListenAndServeTLS("", ""); err != nil {
		fmt.Println("TLS server exited")
		fmt.Println(err)
	}

	return nil
```

</details>

<details>
  <summary>Loading Private and Public keys</summary>

```golang
	// In some situations you want to use a public/private keypair for signing instead.
	// Here is how we would load those:
	var privFile file.PrivateKey // Effectively the same as privKeyFile := file.NoPrivateKey()
	var err error
	if privPath != nil {
		privFile, err = file.SomePrivateKey(*privPath)
		if err != nil {
			fmt.Println("Failed to initialize private key Option")
			fmt.Println(err)
			return err
		}
	}

	var pubFile file.PubKey // Effectively the same as pubKeyFile := file.NoPubKey()
	if pubPath != nil {
		pubFile, err = file.SomePubKey(*pubPath)
		if err != nil {
			fmt.Println("Failed to initialize private key Option")
			fmt.Println(err)
			return err
		}
	}

	// NOTE: as is usually the case with golang key loading, this returns pubKey as a []any and you have to kind of
	// just know how to handle it yourself.
	pubKeys, err := pubFile.ReadPublicKeys()
	if err != nil {
		fmt.Println("Error while reading public key(s) from file:")
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Found this many public keys:")
		fmt.Println(len(pubKeys))
		return err
	}

	// While a public key file may have multiple public keys, private key files should only have a single key. This
	// key is also returned as an any type which you will then need to sort out how to use just like any other key
	// loading.
	privKey, err := privFile.ReadPrivateKey()
	if err != nil {
		fmt.Println("Error while reading private key from file:")
		fmt.Println(err)
		return err
	}

	fmt.Println("Loaded a private key from file")
	switch key := privKey.(type) {
	case *rsa.PrivateKey:
		fmt.Println("key is of type RSA:", key)
	case *dsa.PrivateKey:
		fmt.Println("key is of type DSA:", key)
	case *ecdsa.PrivateKey:
		fmt.Println("key is of type ECDSA:", key)
	case ed25519.PrivateKey:
		fmt.Println("key is of type Ed25519:", key)
	default:
		panic("unknown type of private key")
	}

	return nil
```

</details>

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

See the examples/simple_config/ directory for a setup for brain-dead config parsing. Sure, it's verbose for two real parameters and there
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
- Reloadable, so if you create an init ConfigLoader from flags you can then do hot reloads in response to e.g. a SIGHUP (see examples/simple_config/main.go)
- No magic hidden in a library you need to look up

If you want to try it, it was written so that the precedence is flags > file > env. The default values in the code are
{ Host: "localhost", Port: 1443 }.

See examples/simple_config/README.md for how to test it out.

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
