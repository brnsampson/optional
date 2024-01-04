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

type Cert struct {
	optional.Option[string]
}

func SomeCert(path string) Cert {
	return Cert{optional.Some(path)}
}

func NoCert() Cert {
	return Cert{optional.None[string]()}
}

// Override Str.Get() to update the behavior of all Get* and Unwrap* functions
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
		o.Set(tmp)
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

type PubKey struct {
	optional.Option[string]
}

func SomePubKey(path string) PubKey {
	return PubKey{optional.Some(path)}
}

func NoPubKey() PubKey {
	return PubKey{optional.None[string]()}
}

// Override Str.Get() to update the behavior of all Get* and Unwrap* functions
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
		o.Set(tmp)
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

// PrivateKey contains an optional path to a PEM format private key of some format. Methods are provided to read and
// parse the file into either the individual key or into a tls certificate.
type PrivateKey struct {
	optional.Option[string]
}

func SomePrivateKey(path string) PrivateKey {
	return PrivateKey{optional.Some(path)}
}

func NoPrivateKey() PrivateKey {
	return PrivateKey{optional.None[string]()}
}

// Override Str.Get() to update the behavior of all Get* and Unwrap* functions
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
		o.Set(tmp)
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

func (o PrivateKey) ReadCert(pub PubKeyOptional) (cert tls.Certificate, err error) {
	keyFile, err := o.Get()
	if err != nil {
		return
	}

	pubFile, err := o.Get()
	if err != nil {
		return
	}

	return tls.LoadX509KeyPair(pubFile, keyFile)
}

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
