package config

import (
	"time"
	"github.com/brnsampson/optional"
	"crypto/tls"
	"crypto/x509"
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/ecdh"
)

type primatives interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~bool | ~string | time.Time
}

type PublicKeys interface {
	*rsa.PublicKey | *ecdsa.PublicKey | ed25519.PublicKey | *ecdh.PublicKey
}

type PrivateKeys interface {
	*rsa.PrivateKey | *ecdsa.PrivateKey | ed25519.PrivateKey | *ecdh.PrivateKey
}

// ConfigOptional is an extension of the Optional interface meant to make it more useful for loading configurations.
type ConfigOptional[T primatives] interface {
	optional.Optional[T]

	// Satisfies fmt.Stringer interface
	String() string
	// Satisfies encoding.TextUnmarshaler
	UnmarshalText(text []byte) error
	// Satisfies encoding.TextMarshaler
	MarshalText() (text []byte, err error)
}

type CertOptional interface {
	ConfigOptional[string]

	FilePermsValid() (bool, error)
	SetFilePerms() error
	// See https://pkg.go.dev/crypto/x509@go1.21.5#ParsePKIXPublicKey for the type of pub
	ReadCerts() (certs []*x509.Certificate, err error)
	// See https://pkg.go.dev/crypto/x509@go1.21.5#MarshalPKIXPublicKey for the accepted types for pub
	WriteCerts(certs []*x509.Certificate) error
}

type PubKeyOptional interface {
	ConfigOptional[string]

	FilePermsValid() (bool, error)
	SetFilePerms() error
	// See https://pkg.go.dev/crypto/x509@go1.21.5#ParsePKIXPublicKey for the type of pub
	ReadPublicKeys() (pub []any, err error)
	// See https://pkg.go.dev/crypto/x509@go1.21.5#MarshalPKIXPublicKey for the accepted types for pub
	WritePublicKeys(pubs []any) error
}


type PrivKeyOptional interface {
	ConfigOptional[string]

	FilePermsValid() bool
	SetFilePerms() error
	// See https://pkg.go.dev/crypto/x509@go1.21.5#ParsePKCS8PrivateKey for details about the type of key
	ReadPrivateKey() (key any, err error)
	ReadCert(pub PubKeyOptional) (cert tls.Certificate, err error)
	// See https://pkg.go.dev/crypto/x509@go1.21.5#MarshalPKIXPublicKey for the accepted type for key
	WritePrivateKey(key any) error
}
