package config

import (
	"os"
	"io"
	"io/fs"
	"fmt"
	"path/filepath"
	"encoding/pem"
	"crypto/tls"
	"crypto/x509"
	"github.com/brnsampson/optional"
)

// Verifying and setting file permissions for public/private keys and certificates use the following file mode masks.
// The *Perms modes are the desired permissions, while the *PermsMask consts are such that perms && mask should always
// be 0. The mask is only needed because _technically_ I suppose you could make a public key mode 600 or something if you
// really wanted.
const (
	KeyFilePerms		fs.FileMode = 0600
	PubKeyFilePerms		fs.FileMode = 0644
	KeyFilePermsMask    fs.FileMode = 0177
	PubKeyFilePermsMask fs.FileMode = 0133
)

func filePermsValid(path string, notValidPerms fs.FileMode) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	mode := stat.Mode()
	if (mode & notValidPerms) == 0 {
		// mode does not include one of the flags --x-wx-wx
		return true, nil
	}

	return false, nil
}

func setFilePerms(path string, perms fs.FileMode) error {
	err := os.Chmod(path, perms)
	if err != nil {
		return err
	}
	return nil
}

func readBlocks(path string, perms, notPerms fs.FileMode) (blocks []*pem.Block, err error) {
	valid, err := filePermsValid(path, notPerms)
	if err != nil {
		return
	}
	if !valid {
		err = fmt.Errorf("Cannot read blocks from %s: File permissions should be %o", path, PubKeyFilePerms)
		return
	}

	reader, err := os.Open(path)
	if err != nil {
		return
	}
	defer reader.Close()

	encoded, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	var block *pem.Block
	for {
		block, encoded = pem.Decode(encoded)
		if block == nil {
			break
		}

		blocks = append(blocks, block)
	}
	return
}

func writeBlocks(path string, perms, notPerms fs.FileMode, blocks []*pem.Block) error {
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	valid, err := filePermsValid(path, notPerms)
	if err != nil {
		return err
	}
	if !valid {
		err = setFilePerms(path, perms)
		if err != nil {
			return err
		}
	}

	for _, block := range blocks {
		err = pem.Encode(writer, block)
		if err != nil {
			return err
		}
	}

	return nil
}

// Cert wraps an optional path string and provides extra methods for reading, decoding, and writing pem files containing
// CERTIFICATE blocks.
type Cert struct {
	optional.Option[string]
}

func SomeCert(path string) Cert {
	return Cert{optional.Some(path)}
}

func NoCert() Cert {
	return Cert{optional.None[string]()}
}

// Overrides Option.Match to account for relative paths potentially being different strings but representing the same file.
func (o Cert) Match(probe string) bool {
	if o.IsNone() {
		return false
	} else {
		path, err := o.Get()
		if err != nil {
			// How did we get here...
			return false
		}

		abs, err := filepath.Abs(probe)
		if err != nil {
			// Invalid paths can never be equal!
			return false
		}
		return path == abs
	}
}

// Overrides Option.Get() to update the behavior of all Get* and Unwrap* functions in order to always return the absolute
// path of the desired file.
func (o Cert) Get() (string, error) {
	inner, err := o.Option.Get()
	if err != nil {
		return inner, err
	}

	abs, err := filepath.Abs(inner)
	if err != nil {
		return inner, err
	}

	return abs, nil
}

func (o Cert) Type() string {
	return "Cert"
}

func (o *Cert) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Cert) String() string {
	if o.IsNone() {
		return "None[Cert]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Cert]"
		}
		return tmp
	}
}

func (o Cert) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(tmp), err
	}
}

func (o *Cert) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		o.SetVal(tmp)
	}

	return nil
}

func (o Cert) FilePermsValid() (bool, error) {
	tmp, err := o.Get()
	if err != nil {
		return false, err
	}

	return filePermsValid(tmp, PubKeyFilePermsMask)
}

func (o Cert) SetFilePerms() error {
	tmp, err := o.Get()
	if err != nil {
		return err
	}

	return setFilePerms(tmp, PubKeyFilePerms)
}

func (o Cert) ReadCerts() (certs []*x509.Certificate, err error) {
	tmp, err := o.Get()
	if err != nil {
		return
	}

	blocks, err := readBlocks(tmp, PubKeyFilePerms, PubKeyFilePermsMask)
	if err != nil {
		return
	}

	var c []*x509.Certificate
	for _, block := range blocks {
		if block.Type == "CERTIFICATE" {
			c, err = x509.ParseCertificates(block.Bytes)
			if err != nil {
				return
			}
			certs = append(certs, c...)
		}
	}
	return
}

func (o Cert) WriteCerts(certs []*x509.Certificate) error {
	path, err := o.Get()
	if err != nil {
		return err
	}

	blocks := make([]*pem.Block, 0)
	for _, cert := range certs {
		blocks = append(blocks, &pem.Block{ Type: "CERTIFICATE", Bytes: cert.Raw })
	}

	return writeBlocks(path, PubKeyFilePerms, PubKeyFilePermsMask, blocks)
}

// PubKey wraps an optional path string and provides extra methods for reading, decoding, and writing pem files containing
// "* PUBLIC KEY" blocks.
type PubKey struct {
	optional.Option[string]
}

func SomePubKey(path string) PubKey {
	return PubKey{optional.Some(path)}
}

func NoPubKey() PubKey {
	return PubKey{optional.None[string]()}
}

// Overrides Option.Match to account for relative paths potentially being different strings but representing the same file.
func (o PubKey) Match(probe string) bool {
	if o.IsNone() {
		return false
	} else {
		path, err := o.Get()
		if err != nil {
			// How did we get here...
			return false
		}

		abs, err := filepath.Abs(probe)
		if err != nil {
			// Invalid paths can never be equal!
			return false
		}
		return path == abs
	}
}

// Overrides Option.Get() to update the behavior of all Get* and Unwrap* functions in order to always return the absolute
// path of the desired file.
func (o PubKey) Get() (string, error) {
	inner, err := o.Option.Get()
	if err != nil {
		return inner, err
	}

	abs, err := filepath.Abs(inner)
	if err != nil {
		return inner, err
	}

	return abs, nil
}

func (o PubKey) Type() string {
	return "PubKey"
}

func (o *PubKey) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o PubKey) String() string {
	if o.IsNone() {
		return "None[PubKey]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[PubKey]"
		}
		return tmp
	}
}

func (o PubKey) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(tmp), err
	}
}

func (o *PubKey) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		o.SetVal(tmp)
	}

	return nil
}

func (o PubKey) FilePermsValid() (bool, error) {
	tmp, err := o.Get()
	if err != nil {
		return false, err
	}

	return filePermsValid(tmp, PubKeyFilePermsMask)
}

func (o PubKey) SetFilePerms() error {
	tmp, err := o.Get()
	if err != nil {
		return err
	}

	return setFilePerms(tmp, PubKeyFilePerms)
}

// ReadPublicKeys will return all public keys found in the given filepath or error. The keys may be of type *rsa.PublicKey,
// *ecdsa.PublicKey, ed25519.PublicKey (Note: that is not a pointer), or *ecdh.PublicKey, depending on the contents of
// the file.
func (o PubKey) ReadPublicKeys() (pub []any, err error) {
	tmp, err := o.Get()
	if err != nil {
		return
	}

	blocks, err := readBlocks(tmp, PubKeyFilePerms, PubKeyFilePermsMask)
	if err != nil {
		return
	}

	var pubKey any
	for _, block := range blocks {
		switch block.Type {
		case "PUBLIC KEY":
			pubKey, err = x509.ParsePKIXPublicKey(block.Bytes)
		case "RSA PUBLIC KEY":
			pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		}
		if err == nil {
			pub = append(pub, pubKey)
		}
	}

	if len(pub) > 0 {
		// There were some public keys in the pem file along with *something* else
		err = nil
	}
	return
}

// WritePublicKey will accept any of an *rsa.PublicKey, *dsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey (Note:
// a pointer), or *ecdh.PublicKey. The key will be encoded and written to the path the PubKey option is set to
// with file permissions set appropriately.
func (o PubKey) WritePublicKeys(pubs []any) error {
	path, err := o.Get()
	if err != nil {
		return err
	}

	blocks := make([]*pem.Block, 0)
	for _, pub := range pubs {
		der, err := x509.MarshalPKIXPublicKey(pub)
		if err != nil {
			return err
		}
		block := pem.Block{ Type: "PUBLIC KEY", Bytes: der }
		blocks = append(blocks, &block)
	}

	return writeBlocks(path, PubKeyFilePerms, PubKeyFilePermsMask, blocks)
}

// PubKey wraps an optional path string and provides extra methods for reading, decoding, and writing pem files containing
// "* PRIVATE KEY" blocks.
type PrivateKey struct {
	optional.Option[string]
}

func SomePrivateKey(path string) PrivateKey {
	return PrivateKey{optional.Some(path)}
}

func NoPrivateKey() PrivateKey {
	return PrivateKey{optional.None[string]()}
}

// Overrides Option.Match to account for relative paths potentially being different strings but representing the same file.
func (o PrivateKey) Match(probe string) bool {
	if o.IsNone() {
		return false
	} else {
		path, err := o.Get()
		if err != nil {
			// How did we get here...
			return false
		}

		abs, err := filepath.Abs(probe)
		if err != nil {
			// Invalid paths can never be equal!
			return false
		}
		return path == abs
	}
}

// Overrides Option.Get() to update the behavior of all Get* and Unwrap* functions in order to always return the absolute
// path of the desired file.
func (o PrivateKey) Get() (string, error) {
	inner, err := o.Option.Get()
	if err != nil {
		return inner, err
	}

	abs, err := filepath.Abs(inner)
	if err != nil {
		return inner, err
	}

	return abs, nil
}

func (o PrivateKey) Type() string {
	return "PrivateKey"
}

func (o *PrivateKey) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o PrivateKey) String() string {
	if o.IsNone() {
		return "None[PrivateKey]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[PrivateKey]"
		}
		return tmp
	}
}

func (o PrivateKey) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(tmp), err
	}
}

func (o *PrivateKey) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		o.SetVal(tmp)
	}

	return nil
}

func (o PrivateKey) FilePermsValid() (bool, error) {
	tmp, err := o.Get()
	if err != nil {
		return false, err
	}

	stat, err := os.Stat(tmp)
	mode := stat.Mode()
	if (mode & KeyFilePermsMask) == 0 {
		// mode excludes all of the flags --xrwxrwx. That is to say, the permissions are 600
		return true, nil
	}

	return false, nil
}

func (o PrivateKey) SetFilePerms() error {
	tmp, err := o.Get()
	if err != nil {
		return err
	}

	err = os.Chmod(tmp, KeyFilePerms)
	if err != nil {
		return err
	}
	return nil
}

// ReadPrivateKey will return the first private key found in the given filepath or error. This may return an *rsa.PrivateKey,
// *ecdsa.PrivateKey, ed25519.PrivateKey (Note: that is not a pointer), or *ecdh.PrivateKey, depending on the contents of
// the file.
func (o PrivateKey) ReadPrivateKey() (key any, err error) {
	path, err := o.Get()
	if err != nil {
		return
	}
	blocks, err := readBlocks(path, KeyFilePerms, KeyFilePermsMask)
	if err != nil {
		return
	}

	var tmp any
	for _, block := range blocks {
		switch block.Type {
		case "PRIVATE KEY":
			tmp, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		case "RSA PRIVATE KEY":
			tmp, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		case "EC PRIVATE KEY":
			tmp, err = x509.ParseECPrivateKey(block.Bytes)
		}
		if err == nil {
			key = tmp
			break
		}
	}
	// if we only encountered errors, err will not be nil. If one of the blocks worked, then error will be nil for that
	// loop and we would have broken at that point.
	return
}

// ReadCert accepts a Cert struct and returns a tls.Certificate for the keypair if both Optionals are Some. This
// is going to be the most used case for anyone loading
func (o PrivateKey) ReadCert(certIn Cert) (cert tls.Certificate, err error) {
	keyFile, err := o.Get()
	if err != nil {
		return
	}

	certFile, err := certIn.Get()
	if err != nil {
		return
	}

	return tls.LoadX509KeyPair(certFile, keyFile)
}

// WritePrivateKey will accept any of an *rsa.PrivateKey, *dsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey (Note:
// a pointer), or *ecdh.PrivateKey. The key will be encoded and written to the path the PrivateKey option is set to
// with file permissions set appropriately.
func (o PrivateKey) WritePrivateKey(key any) error {
	path, err := o.Get()
	if err != nil {
		return err
	}

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	block := pem.Block{ Type: "PRIVATE KEY", Bytes: der }
	blocks := []*pem.Block{&block}

	return writeBlocks(path, KeyFilePerms, KeyFilePermsMask, blocks)
}
