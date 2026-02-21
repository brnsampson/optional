// Package optional provides a convenient and safe way to deal with options which may or may
// not have a designated value without using pointers and risking nil pointer dereferencing panics.
package optional

import "time"

type OptionalError struct {
	msg string
}

func optionalError(msg string) *OptionalError {
	return &OptionalError{msg}
}

func (e OptionalError) Error() string {
	return e.msg
}

type Transformer[T comparable] func(T) (T, error)

type primatives interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~bool | ~string | time.Time
}

// Optional defines the functionality needed to provide good ergonimics around optional fields and values. In general,
// code should not declare variables or parameters as Optionals and instead prefer using concrete types like Option.
// This interface is meant to ensure compatablility between different concrete option types and for the rare cases where
// the abstraction is actually necessary.
type Optional[T comparable] interface {
	IsSome() bool
	IsNone() bool
	Clone() Optional[T]
	Get() (val T, ok bool)
	MustGet() T
	Match(T) bool
	// Satisfies encoding.json.Marshaler
	MarshalJSON() ([]byte, error)
}

// MutableOptional is a superset of Optional which allows mutating and transforming the wrapped value.
type MutableOptional[T comparable] interface {
	Optional[T]

	MutableClone() MutableOptional[T]
	Clear()
	Default(T) (replaced bool)
	// Replace updates the current value in the MutableOptional and returns the previous value
	Replace(T) Optional[T]
	// Transform only applies the func to the values of Some valued Optionals. Any mapping of None is None.
	Transform(f Transformer[T]) error
	// Satisfies encoding.json.UnMarshaler
	UnmarshalJSON([]byte) error
}

// StorableOptional is an extension of the MutableOptional interface that allows golang's
// database/sql package to correctly store and retrieve values.
type StorableOptional[T comparable] interface {
	MutableOptional[T]

	// Implements database/sql.Scanner interface
	Scan(src any) error
	// Implements the database/sqlc.Valuer interface
	Value() (any, error)
}

// LoadableOptional is an extension of the Optional interface meant to make it more useful for
// loading from a variety of sources.
type LoadableOptional[T primatives] interface {
	MutableOptional[T]

	// Satisfies fmt.Stringer interface
	String() string
	// Along with String(), implements flag.Value and pflag.Value
	Set(string) error
	// Along with String() and Set(), implements pflag.Value
	Type() string
	// Satisfies the encoding.TextMarshaler
	MarshalText() (text []byte, err error)
	// Satisfies encoding.TextUnmarshaler
	UnmarshalText(text []byte) error
	// json.Marshaler and json.Unmarshaler are implemented by the embedded MutableOptional interface
}

// IsSomeAnd returns true if the Option has a value of Some(x) and f(x) == true
func IsSomeAnd[T comparable](opt Option[T], f func(T) bool) bool {
	tmp, ok := opt.Get()
	if ok {
		return f(tmp)
	} else {
		return false
	}
}

// Equal is a convenience function for checking if the contents of two Optional types are equivilent.
// Note that Get() and Match() may be overridden by more complex types which wrap a vanilla Option.
// In these situations, the writer is responsible for making sure that the invariant Some(x).Match(Some(x).Get())
// is always true.
func Equal[T comparable, O Optional[T]](left, right O) bool {
	if left.IsNone() && right.IsNone() {
		return true
	} else if left.IsSome() && right.IsSome() {
		tmp, ok := left.Get()
		if !ok {
			return false
		}

		return right.Match(tmp)
	} else {
		// one is none and the other is some
		return false
	}
}

// And returns None if the first Optional is None, and the second Optional otherwise. Conceptually similar to
// left && right. This is a convenience function for Option selection. Convenient for merging configs, implementing
// builder patterns, etc.
func And[T comparable, O Optional[T]](left, right O) O {
	if left.IsNone() {
		return left
	} else {
		return right
	}
}

// Or returns the first Optional if it contains a value. Otherwise, return the second Optional. This is conceptually
// similar to left || right. This is a convenience function for situations like merging configs or implementing
// builder patterns.
func Or[T comparable, O Optional[T]](left, right O) O {
	if left.IsSome() {
		return left
	} else {
		return right
	}
}

// ClearIfMatch calls clear if Optional.Match(probe) == true. This is a convenience for situations where you need to convert
// from a value of T with known "magic value" which implies None. A good example of this is if you have an int loaded
// from command line flags and you know that any flag omitted by the user will be assigned to 0. This can be done like this:
// o := Some(x)
// o.ClearIfMatch(0)
func ClearIfMatch[T comparable](opt MutableOptional[T], probe T) {
	if opt.Match(probe) {
		opt.Clear()
	}
}

// GetOr is the same as Get, but will return the passed value instead of an error if the Option is None. Another convenience
// function
func GetOr[T comparable](opt Optional[T], val T) T {
	res, ok := opt.Get()
	if ok {
		return res
	} else {
		return val
	}
}

// GetOrElse calls Get(), but run the passed function and return the result instead of producing an error if the option
// is None.
func GetOrElse[T comparable](opt Optional[T], f func() T) T {
	res, ok := opt.Get()
	if ok {
		return res
	} else {
		return f()
	}
}

// GetOrInsert calls Get, but will call Default on the passed value then return it if the Option is None
func GetOrInsert[T comparable](opt MutableOptional[T], val T) (T, error) {
	_ = opt.Default(val)
	res, ok := opt.Get()

	if ok {
		return res, nil
	} else {
		panic("Optional did not have value immediately after being set")
	}
}

// TransformOr just calls Transform(), except None values are mapped to backup before being transformed.
func TransformOr[T comparable](opt MutableOptional[T], t Transformer[T], backup T) error {
	_ = opt.Default(backup)
	return opt.Transform(t)
}
