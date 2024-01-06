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
	// Satisfies encoding.json.Marshaler
	MarshalJSON() ([]byte, error)
}

// MutableOptional is a superset of Optional which allows mutating and transforming the wrapped value.
type MutableOptional[T comparable] interface {
	Optional[T]

	MutableClone() MutableOptional[T]
	Clear()
	ClearIfMatch(T)
	SetVal(T)
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

// Equal is a convenience function for checking if the contents of two Optional types are equivilent.
// Note that Get() and Match() may be overridden by more complex types which wrap a vanilla Option.
// In these situations, the writer is responsible for making sure that the invariant Some(x).Match(Some(x).Get())
// is always true.
func Equal[T comparable, O Optional[T]](left, right O) bool {
	if left.IsNone() && right.IsNone() {
		return true
	} else if left.IsSome() && right.IsSome() {
		tmp, err := left.Get()
		if err != nil {
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
