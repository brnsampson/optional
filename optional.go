package optional

type OptionalError struct {
	msg string
}

func optionalError(msg string) *OptionalError {
	return &OptionalError{msg}
}

func (e OptionalError) Error() string {
	return e.msg
}

// Optional defines the functionality needed to provide good ergonimics around optional fields and values. In general,
// code should not declare variables or parameters as Optionals and instead prefer using concrete types like Option.
// This interface is meant to ensure compatablility between different concrete option types and for the rare cases where
// the abstraction is actually necessary.
type Optional[T comparable] interface {
	IsSome() bool
	IsSomeAnd(func(T) bool) bool
	IsNone() bool
	Clone() Optional[T]
	Get() (T, error)
	GetOr(T) T
	Must() T
	Match(T) bool
	Eq(Optional[T]) bool
	And(Optional[T]) Optional[T]
	Or(Optional[T]) Optional[T]
	// Satisfies encoding.json.Marshaler
	MarshalJSON() ([]byte, error)
}

// MutableOptional is a superset of Optional which allows mutating and transforming the wrapped value.
type MutableOptional[T comparable] interface {
	Optional[T]

	MutableClone() MutableOptional[T]
	Clear()
	Set(T)
	GetOrInsert(T) T
	Unwrap() (T, error)
	MustUnwrap() T
	UnwrapOr(T) T
	UnwrapOrElse(func() T) T
	// Transform only applies the func to the values of Some valued Optionals. Any mapping of None is None.
	Transform(f func(T) T)
	// TransformOr works just like Transform, but maps None -> backup
	TransformOr(f func(T) T, backup T)
	// TransformOrError works just like Transform, but the transform function can return an error which is returned as-is
	TransformOrError(f func(T) (T, error)) error
	BinaryTransform(second T, f func(T, T) T)
	// Satisfies encoding.json.UnMarshaler
	UnmarshalJSON([]byte) error
}
