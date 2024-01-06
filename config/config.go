package config

import (
	"time"
	"github.com/brnsampson/optional"
)

type primatives interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~bool | ~string | time.Time
}

// ConfigOptional is an extension of the Optional interface meant to make it more useful for loading configurations.
type ConfigOptional[T primatives] interface {
	optional.MutableOptional[T]

	// Along with String() and Set(string) error from, implements pflag.Value
	Type() string
	Set(string) error
	// Satisfies fmt.Stringer interface
	String() string
	// Satisfies encoding.TextUnmarshaler
	UnmarshalText(text []byte) error
	// Satisfies encoding.TextMarshaler
	MarshalText() (text []byte, err error)
}
