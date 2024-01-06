package config_test

import (
	"os"
	"bytes"
	"path/filepath"
	"reflect"
	"testing"
	"crypto"
	"crypto/x509"
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
)

func TestCertType(t *testing.T) {
	o := config.SomeCert("/not/a/real/path")
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestCertString(t *testing.T) {
	path := "../tls/rsa/cert.pem"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	o := config.SomeCert(path)
	assert.Equal(t, abs, o.String())
}

func TestCertMarshalText(t *testing.T) {
	path := "../tls/rsa/cert.pem"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	o := config.SomeCert(path)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, abs, string(s))
}

func TestCertUnmarshalText(t *testing.T) {
	path := "../tls/rsa/cert.pem"
	nullCert := "null"
	intCert := "42"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)
	intAbs, err := filepath.Abs(intCert)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)


	// Text sucessful unmarshaling
	o := config.NoCert()
	err = o.UnmarshalText([]byte(path))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, abs, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullCert))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-path. This will work because it should interpret this as a file named "41", which you could
	// theoretically have.
	err = o.UnmarshalText([]byte(intCert))
	assert.NilError(t, err)

	ret, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, intAbs, ret)
}

func TestCertFilePermsValid(t *testing.T) {
	valid_path := "../tls/rsa/cert.pem"
	invalid_path := "../tls/rsa/cert_bad_perms.pem"
	v := config.SomeCert(valid_path)
	i := config.SomeCert(invalid_path)

	good, err := v.FilePermsValid()
	assert.NilError(t, err)

	bad, err := i.FilePermsValid()
	assert.NilError(t, err)

	assert.Assert(t, good)
	assert.Assert(t, !bad)
}

func TestCertSetFilePerms(t *testing.T) {
	f, err := os.CreateTemp("", "cert")
	assert.NilError(t, err)

	path := f.Name()
	defer os.Remove(path)
	f.Chmod(0666)

	o := config.SomeCert(path)
	o.SetFilePerms()
	s, err := f.Stat()
	assert.NilError(t, err)
	assert.Equal(t, config.PubKeyFilePerms, s.Mode())
}

func TestCertReadRSACerts(t *testing.T) {
	expectedCertIssuer := "CN=www.whobe.us,OU=optional,O=BS Workshops,L=Who knows,ST=California,C=US"
	certPath := "../tls/rsa/cert.pem"

	c := config.SomeCert(certPath)
	certs, err := c.ReadCerts()
	assert.NilError(t, err)
	assert.Assert(t, len(certs) == 1)

	cert := certs[0]
	assert.Equal(t, x509.SHA256WithRSA, cert.SignatureAlgorithm)
	assert.Equal(t, expectedCertIssuer, cert.Issuer.String())
}

func TestCertReadECDSACerts(t *testing.T) {
	expectedCertIssuer := "CN=www.whobe.us,OU=optional,O=BS Workshops,L=Who knows,ST=California,C=US"
	certPath := "../tls/ecdsa/cert.pem"

	c := config.SomeCert(certPath)
	certs, err := c.ReadCerts()
	assert.NilError(t, err)
	assert.Assert(t, len(certs) == 1)

	cert := certs[0]
	assert.Equal(t, x509.ECDSAWithSHA512, cert.SignatureAlgorithm)
	assert.Equal(t, expectedCertIssuer, cert.Issuer.String())
}

func TestCertReadED25519Certs(t *testing.T) {
	expectedCertIssuer := "CN=www.whobe.us,OU=optional,O=BS Workshops,L=Who knows,ST=California,C=US"
	certPath := "../tls/ed25519/cert.pem"

	c := config.SomeCert(certPath)
	certs, err := c.ReadCerts()
	assert.NilError(t, err)
	assert.Assert(t, len(certs) == 1)

	cert := certs[0]
	assert.Equal(t, x509.PureEd25519, cert.SignatureAlgorithm)
	assert.Equal(t, expectedCertIssuer, cert.Issuer.String())
}

func TestCertFileWriteCerts(t *testing.T) {
	// Read in valid certificates (tested above)
	certPath := "../tls/rsa/cert.pem"
	c := config.SomeCert(certPath)
	certs, err := c.ReadCerts()
	assert.NilError(t, err)

	f, err := os.CreateTemp("", "cert")
	assert.NilError(t, err)

	// Create temporary file (empty)
	path := f.Name()
	defer os.Remove(path)
	f.Chmod(0644)
	f.Close()

	tc := config.SomeCert(path)
	// Write certificates to new file
	err = tc.WriteCerts(certs)
	assert.NilError(t, err)

	// Check that reading the certs back gives us the expected values
	newCerts, err := tc.ReadCerts()
	assert.NilError(t, err)

	// validate that we are testing what we think we are testing
	assert.Assert(t, len(certs) == 1)
	assert.Assert(t, len(newCerts) == 1)
	assert.Assert(t, certs[0].Equal(newCerts[0]))
}

func TestPubKeyType(t *testing.T) {
	o := config.SomePubKey("/not/a/real/path")
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestPubKeyString(t *testing.T) {
	path := "../tls/rsa/pubkey.pem"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	o := config.SomePubKey(path)
	assert.Equal(t, abs, o.String())
}

func TestPubKeyMarshalText(t *testing.T) {
	path := "../tls/rsa/pubkey.pem"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	o := config.SomePubKey(path)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, abs, string(s))
}

func TestPubKeyUnmarshalText(t *testing.T) {
	path := "../tls/rsa/pubkey.pem"
	nullPubKey := "null"
	intPubKey := "42"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)
	intAbs, err := filepath.Abs(intPubKey)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)


	// Text sucessful unmarshaling
	o := config.NoPubKey()
	err = o.UnmarshalText([]byte(path))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, abs, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullPubKey))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-path. This will work because it should interpret this as a file named "41", which you could
	// theoretically have.
	err = o.UnmarshalText([]byte(intPubKey))
	assert.NilError(t, err)

	ret, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, intAbs, ret)
}

func TestPubKeyFilePermsValid(t *testing.T) {
	valid_path := "../tls/rsa/pubkey.pem"
	invalid_path := "../tls/rsa/pubkey_bad_perms.pem"
	v := config.SomePubKey(valid_path)
	i := config.SomePubKey(invalid_path)

	good, err := v.FilePermsValid()
	assert.NilError(t, err)

	bad, err := i.FilePermsValid()
	assert.NilError(t, err)

	assert.Assert(t, good)
	assert.Assert(t, !bad)
}

func TestPubKeySetFilePerms(t *testing.T) {
	f, err := os.CreateTemp("", "pubkey")
	assert.NilError(t, err)

	path := f.Name()
	defer os.Remove(path)
	f.Chmod(0666)

	o := config.SomePubKey(path)
	o.SetFilePerms()
	s, err := f.Stat()
	assert.NilError(t, err)
	assert.Equal(t, config.PubKeyFilePerms, s.Mode())
}

func TestPubKeyReadPublicKeysRSA(t *testing.T) {
	pubPath := "../tls/rsa/pubkey.pem"
	p := config.SomePubKey(pubPath)
	keys, err := p.ReadPublicKeys()
	assert.NilError(t, err)

	// There is only one cert in the file...
	assert.Assert(t, len(keys) == 1)
	key := keys[0]

	switch key.(type) {
	case *rsa.PublicKey:
	default:
		panic("rsa pem file was not read into rsa key!")
	}
}

func TestPubKeyReadPublicKeysECDSA(t *testing.T) {
	pubPath := "../tls/ecdsa/pub.pem"

	p := config.SomePubKey(pubPath)
	keys, err := p.ReadPublicKeys()
	assert.NilError(t, err)

	// There is only one cert in the file...
	assert.Assert(t, len(keys) == 1)
	key := keys[0]

	switch key.(type) {
	case *ecdsa.PublicKey:
	default:
		panic("rsa pem file was not read into rsa key!")
	}
}

func TestPubKeyReadPublicKeysED25519(t *testing.T) {
	pubPath := "../tls/ed25519/pub.pem"

	p := config.SomePubKey(pubPath)
	keys, err := p.ReadPublicKeys()
	assert.NilError(t, err)

	// There is only one cert in the file...
	assert.Assert(t, len(keys) == 1)
	key := keys[0]

	switch key.(type) {
	case ed25519.PublicKey:
	default:
		panic("rsa pem file was not read into rsa key!")
	}
}

func TestPubKeyWritePubKey(t *testing.T) {
	// Read in valid certificates (tested above)
	pubPath := "../tls/rsa/pubkey.pem"
	c := config.SomePubKey(pubPath)
	keys, err := c.ReadPublicKeys()
	assert.NilError(t, err)

	f, err := os.CreateTemp("", "pubkey")
	assert.NilError(t, err)

	// Create temporary file (empty)
	path := f.Name()
	defer os.Remove(path)
	f.Chmod(0666)
	f.Close()

	tc := config.SomePubKey(path)
	// Write certificates to new file
	err = tc.WritePublicKeys(keys)
	assert.NilError(t, err)

	// Check that reading the certs back gives us the expected values
	newKeys, err := tc.ReadPublicKeys()
	assert.NilError(t, err)

	// validate that we are testing what we think we are testing
	assert.Assert(t, len(keys) == 1)
	assert.Assert(t, len(newKeys) == 1)

	key := keys[0]
	newKey := newKeys[0]
	switch k1 := key.(type) {
	case *rsa.PublicKey:
		switch k2 := newKey.(type) {
		case *rsa.PublicKey:
			assert.Assert(t, k1.Equal(k2))
		default:
			panic("Expected key loaded from tmp file to be *rsa.PublicKey, but it wasn't!")
		}
	default:
		panic("Expected key loaded from tmp file to be *rsa.PublicKey, but it wasn't!")
	}
}

func TestPrivateKeyType(t *testing.T) {
	o := config.SomePrivateKey("/not/a/real/path")
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestPrivateKeyString(t *testing.T) {
	path := "../tls/rsa/key.pem"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	o := config.SomePrivateKey(path)
	assert.Equal(t, abs, o.String())
}

func TestPrivateKeyMarshalText(t *testing.T) {
	path := "../tls/rsa/key.pem"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	o := config.SomePrivateKey(path)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, abs, string(s))
}

func TestPrivateKeyUnmarshalText(t *testing.T) {
	path := "../tls/rsa/key.pem"
	nullPrivateKey := "null"
	intPrivateKey := "42"
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)
	intAbs, err := filepath.Abs(intPrivateKey)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)


	// Text sucessful unmarshaling
	o := config.NoPrivateKey()
	err = o.UnmarshalText([]byte(path))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, abs, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullPrivateKey))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-path. This will work because it should interpret this as a file named "41", which you could
	// theoretically have.
	err = o.UnmarshalText([]byte(intPrivateKey))
	assert.NilError(t, err)

	ret, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, intAbs, ret)
}

func TestPrivateKeyFilePermsValid(t *testing.T) {
	valid_path := "../tls/rsa/key.pem"
	invalid_path := "../tls/rsa/key_bad_perms.pem"
	v := config.SomePrivateKey(valid_path)
	i := config.SomePrivateKey(invalid_path)

	good, err := v.FilePermsValid()
	assert.NilError(t, err)

	bad, err := i.FilePermsValid()
	assert.NilError(t, err)

	assert.Assert(t, good)
	assert.Assert(t, !bad)
}

func TestPrivateKeySetFilePerms(t *testing.T) {
	f, err := os.CreateTemp("", "key")
	assert.NilError(t, err)

	path := f.Name()
	defer os.Remove(path)
	f.Chmod(0666)

	o := config.SomePrivateKey(path)
	o.SetFilePerms()
	s, err := f.Stat()
	assert.NilError(t, err)
	assert.Equal(t, config.KeyFilePerms, s.Mode())
}

func TestPrivateKeyReadCert(t *testing.T) {
	keyPath := "../tls/rsa/key.pem"
	certPath := "../tls/rsa/cert.pem"
	ko := config.SomePrivateKey(keyPath)
	co := config.SomeCert(certPath)

	key, err := ko.ReadPrivateKey()
	certificate, err := ko.ReadCert(co)
	assert.NilError(t, err)

	// Check that the cert loaded correctly
	certs, err := co.ReadCerts()
	assert.Assert(t, len(certs) == 1)
	assert.Assert(t, len(certificate.Certificate) == 1)
	assert.Assert(t, bytes.Equal(certs[0].Raw, certificate.Certificate[0]))

	// Check that the private key loaded correctly
	type privKeyInter interface {
		Public() crypto.PublicKey
    	Equal(x crypto.PrivateKey) bool
	}

	switch k1 := key.(type) {
	case privKeyInter:
		switch k2 := certificate.PrivateKey.(type) {
		case privKeyInter:
			assert.Assert(t, k1.Equal(k2))
		default:
			panic("Expected *rsa.PrivateKey to implement the crypto.PrivateKey interface, but it didn't!")
		}
	default:
		panic("Expected key from tls.Certificate to implement the crypto.PrivateKey interface, but it didn't!")
	}
}

func TestPrivateKeyReadPrivateKeyRSA(t *testing.T) {
	keyPath := "../tls/rsa/key.pem"
	p := config.SomePrivateKey(keyPath)
	key, err := p.ReadPrivateKey()
	assert.NilError(t, err)

	switch key.(type) {
	case *rsa.PrivateKey:
	default:
		panic("rsa pem file was not read into the correct format key!")
	}
}

func TestPrivateKeyReadPrivateKeyECDSA(t *testing.T) {
	keyPath := "../tls/ecdsa/key.pem"

	p := config.SomePrivateKey(keyPath)
	key, err := p.ReadPrivateKey()
	assert.NilError(t, err)

	switch key.(type) {
	case *ecdsa.PrivateKey:
	default:
		panic("ecdsa pem file was not read into the correct format key!")
	}
}

func TestPrivateKeyReadPrivateKeyED25519(t *testing.T) {
	keyPath := "../tls/ed25519/key.pem"

	p := config.SomePrivateKey(keyPath)
	key, err := p.ReadPrivateKey()
	assert.NilError(t, err)

	switch key.(type) {
	case ed25519.PrivateKey:
	default:
		panic("ed25519 pem file was not read into the correct format key!")
	}
}

func TestPrivateKeyWritePrivateKey(t *testing.T) {
	// Read in valid certificates (tested above)
	keyPath := "../tls/rsa/key.pem"
	c := config.SomePrivateKey(keyPath)
	key, err := c.ReadPrivateKey()
	assert.NilError(t, err)

	f, err := os.CreateTemp("", "key")
	assert.NilError(t, err)

	// Create temporary file (empty)
	path := f.Name()
	defer os.Remove(path)
	f.Chmod(0666)
	f.Close()

	tc := config.SomePrivateKey(path)
	// Write certificates to new file
	err = tc.WritePrivateKey(key)
	assert.NilError(t, err)

	// Check that reading the certs back gives us the expected values
	newKey, err := tc.ReadPrivateKey()
	assert.NilError(t, err)

	switch k1 := key.(type) {
	case *rsa.PrivateKey:
		switch k2 := newKey.(type) {
		case *rsa.PrivateKey:
			assert.Assert(t, k1.Equal(k2))
		default:
			panic("Expected key loaded from tmp file to be *rsa.PrivateKey, but it wasn't!")
		}
	default:
		panic("Expected key loaded from tmp file to be *rsa.PrivateKey, but it wasn't!")
	}
}
