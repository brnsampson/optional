package optional

import (
	"encoding/json"
)

// Option is a generic way to make a field or parameter optional. Instantiating an Optional value through the Some()
// or None() methods are prefered since it is easier for the reader of your code to see what the expected value of the
// Option is, but to avoid leaking options across api boundries a FromPointer() method is given. This allows you to accept
// a parameter as a pointer and immediately convert it into an Option value without having to do the nil check yourself.
//
// While this can be useful on its own, this can also be used to create a specialized option for your use case. See the
// config package for a comprehensive example of this.
//
// Special care should be taken when creating Options of pointer types. All the same concerns around passing structs with
// pointer fields apply, since copying the pointer by value will create many copies of the same pointer. There is nothing
// to stop you from doing this, but I'm not sure what they use case is and it may lead to less understandable code (which
// is what this library was created to avoid in the first place!)
type Option[T comparable] struct {
	inner T
	none  bool
}

// Some returns an Option with an inferred type and specified value.
func Some[T comparable](value T) Option[T] {
	return Option[T]{inner: value, none: false}
}

// None returns an Option with no value loaded. Note that since there is no value to infer type from, None must be
// instanciated with the desired type like optional.None[string]().
func None[T comparable]() Option[T] {
	var tmp T
	return Option[T]{inner: tmp, none: true}
}

// FromPointer creates an option with an inferred type from a pointer. nil pointers are mapped to a None value and non-nil
// pointers have their value copied into a new option. The pointer itself is not retained and can be modified later without
// affecting the Option value.
func FromPointer[T comparable](p *T) Option[T] {
	if p == nil {
		return None[T]()
	} else {
		return Some(*p)
	}
}

func (o Option[T]) IsSome() bool {
	return !o.none
}

// IsSomeAnd returns true if the Option has a value of Some(x) and f(x) == true
func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	return !o.none && f(o.inner)
}

func (o Option[T]) IsNone() bool {
	return o.none
}

// Clone creates a copy of the Option by value. This means the new Option may be unwrapped or modified without affecting
// the old Option. NOTE: the exception to this would be if you create an Option for a pointer type. Because the pointer
// is copied by value, it will still refer to the same value.
func (o Option[T]) Clone() Optional[T] {
	if o.none {
		return None[T]()
	} else {
		return Some(o.inner)
	}
}

// MutableClone creates a copy of the Option by value the same as Clone. The only differnce is that the returned type
// is a pointer cast as a MutableOptional so that the returned value can be further modified.
func (o Option[T]) MutableClone() MutableOptional[T] {
	if o.none {
		tmp := None[T]()
		return &tmp
	} else {
		tmp := Some(o.inner)
		return &tmp
	}
}

// Clear converts a Some(x) or None type Option into a None value.
func (o *Option[T]) Clear() {
	o.none = true
}

// Set converts a Some(x) or None type Option into a Some(value) value
func (o *Option[T]) Set(value T) {
	o.inner = value
	o.none = false
}

// Get returns the current wrapped value of a Some value Option and returns an error if the Option is None.
func (o Option[T]) Get() (T, error) {
	if o.IsSome() {
		return o.inner, nil
	}
	return o.inner, optionalError("Attempted to Get Option with None value")
}

// GetOr is the same as Get, but will return the passed value instead of an error if the Option is None.
func (o Option[T]) GetOr(val T) T {
	res, err := o.Get()
	if err != nil {
		return val
	} else {
		return res
	}
}

// GetOrInsert is the same as Get, but will call Set on the passed value first if the Option is None
func (o *Option[T]) GetOrInsert(val T) T {
	res, err := o.Get()

	if err != nil {
		o.Set(val)
		return val
	} else {
		return res
	}
}

// Must is like Get, but panic instead of producing an error.
func (o Option[T]) Must() T {
	res, err := o.Get()
	if err != nil {
		panic("Attempted to call Must on an Option with None value")
	} else {
		return res
	}
}

// Unwrap returns the current wrapped value of a Some value Option and returns an error if the Option is None. In either
// case, the Option is set to None to indicate that it has been consumed.
func (o *Option[T]) Unwrap() (T, error) {
	res, err := o.Get()
	o.Clear()
	return res, err
}

// MustUnwrap is like Unwrap, but panic instead of producing an error.
func (o *Option[T]) MustUnwrap() T {
	res, err := o.Unwrap()
	if err != nil {
		panic("Attempted to UnsafeUnwrap an Option with None value")
	} else {
		return res
	}
}

// UnwrapOr is like Unwrap, but return the passed value instead of producing an error.
func (o *Option[T]) UnwrapOr(val T) T {
	res, err := o.Unwrap()
	if err != nil {
		return val
	} else {
		return res
	}
}

// UnwrapOrElse is like Unwrap, but run the passed function and return the result instead of producing an error.
func (o *Option[T]) UnwrapOrElse(f func() T) T {
	res, err := o.Unwrap()
	if err != nil {
		return f()
	} else {
		return res
	}
}

// Match tests if the inner value of Option == the passed value
func (o Option[T]) Match(probe T) bool {
	if o.none {
		return false
	} else {
		return o.inner == probe
	}
}

// Eq tests if two values implementing Optional are equal. They do not need to be the same concrete type.
func (o Option[T]) Eq(other Optional[T]) bool {
	if o.none && other.IsNone() {
		return true
	} else if !o.none && other.IsSome() {
		// We do not know if other is a pointer type or not, so play it safe
		return other.Match(o.inner)
	} else {
		// one is none and the other is some
		return false
	}
}

// And returns None if the Option is None, and other if the original Option is Some. Conceptually, think o && other
func (o Option[T]) And(other Optional[T]) Optional[T] {
	if o.none {
		return &o
	} else {
		return other
	}
}

// Or returns the first Option if it is Some, and other if it is None. Conceptually, it is o || other
func (o Option[T]) Or(other Optional[T]) Optional[T] {
	if o.none {
		return other
	} else {
		return &o
	}
}

// Transform applies function f(T) to the inner value of the Option. If the Option is None, then the Option will remain
// None.
func (o *Option[T]) Transform(f func(inner T) T) {
	if o.IsSome() {
		tmp := f(o.inner)
		o.Set(tmp)
	}
}

// TransformOr is like Transform, except None values are mapped to backup
func (o *Option[T]) TransformOr(f func(T) T, backup T) {
	if o.IsSome() {
		tmp := f(o.inner)
		o.Set(tmp)
	} else {
		o.Set(backup)
	}
}

// TransformOrError is like Transform, except the passed function can return an error for invalid values
func (o *Option[T]) TransformOrError(f func(T) (T, error)) error {
	if o.IsSome() {
		tmp, err := f(o.inner)
		if err != nil {
			return err
		}

		o.Set(tmp)
	}
	return nil
}

// BinaryTransform set the value of the Option to f(inner, second). None options always map to None.
func (o *Option[T]) BinaryTransform(second T, f func(T, T) T) {
	if o.IsSome() {
		tmp := f(o.inner, second)
		o.Set(tmp)
	}
}

// MarshalJSON implements the encoding.json.Marshaler interface. None values are marshaled to json null, while Some values are
// passed into json.Marshal directly.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return json.Marshal(nil)
	} else {
		return json.Marshal(o.inner)
	}
}

// UnmarshalJSON implements the encoding.json.Unmarshaller interface. Json nulls are unmarshaled into None values, while
// any other value is attempted to unmarshal as normal. Any error encountered is returned without modification. There are
// some protections included to make sure that unmarshaling an uninitialized Option[T] does not break the Option invariants.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Clear()
		return nil
	}

	var tmp T
	tmp2 := &tmp
	if err := json.Unmarshal(data, tmp2); err != nil {
		return err
	}

	if tmp2 != nil {
		// According to the spec, this should be nil if we recieved a json null, but I've found that you
		// actually will get a valid pointer going to the zero value.
		o.Set(*tmp2)
	} else {
		o.Clear()
	}

	return nil
}
