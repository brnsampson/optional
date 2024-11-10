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
	some  bool
}

// Some returns an Option with an inferred type and specified value.
func Some[T comparable](value T) Option[T] {
	return Option[T]{inner: value, some: true}
}

// None returns an Option with no value loaded. Note that since there is no value to infer type from, None must be
// instanciated with the desired type like optional.None[string]().
func None[T comparable]() Option[T] {
	var tmp T
	return Option[T]{inner: tmp, some: false}
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
	return o.some
}

func (o Option[T]) IsNone() bool {
	return !o.some
}

// Clone creates a copy of the Option by value. This means the new Option may be unwrapped or modified without affecting
// the old Option. NOTE: the exception to this would be if you create an Option for a pointer type. Because the pointer
// is copied by value, it will still refer to the same value.
func (o Option[T]) Clone() Optional[T] {
	if o.some {
		return Some(o.inner)
	} else {
		return None[T]()
	}
}

// MutableClone creates a copy of the Option by value the same as Clone. The only differnce is that the returned type
// is a pointer cast as a MutableOptional so that the returned value can be further modified.
func (o Option[T]) MutableClone() MutableOptional[T] {
	if o.some {
		tmp := Some(o.inner)
		return &tmp
	} else {
		tmp := None[T]()
		return &tmp
	}
}

// Clear converts a Some(x) or None type Option into a None value.
func (o *Option[T]) Clear() {
	o.some = false
}

// Replace converts a Some(x) or None type Option into a Some(value) value.
// The base Option struct can never return an error from Replace, so it is generally safe to ignore the returned values
// from this, e.g. calling o.Replace() instead of _, _ = o.Replace().
func (o *Option[T]) Replace(value T) (Optional[T], error) {
	tmp, err := o.Get()
	o.inner = value
	o.some = true

	if err != nil {
		// it was None
		return None[T](), nil
	} else {
		return Some(tmp), nil
	}
}

// Get returns the current wrapped value of a Some value Option and returns an error if the Option is None.
func (o Option[T]) Get() (T, error) {
	if o.IsSome() {
		return o.inner, nil
	}
	return o.inner, optionalError("Attempted to Get Option with None value")
}

// Match tests if the inner value of Option == the passed value
func (o Option[T]) Match(probe T) bool {
	if o.some {
		return o.inner == probe
	} else {
		return false
	}
}

// Transform applies function f(T) to the inner value of the Option. If the Option is None, then the Option will remain
// None.
func (o *Option[T]) Transform(t Transformer[T]) error {
	if o.IsSome() {
		tmp, err := t(o.inner)
		if err != nil {
			return err
		}

		o.Replace(tmp)
	}
	return nil
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

	// WARNING: this will probably trigger a linter warning because not nil != nil, but leave it anyways.
	if tmp2 != nil {
		// According to the spec, this should be nil if we recieved a json null, but I've found that you
		// actually will get a valid pointer going to the zero value. This unfortunately means we have no
		// good way to determine if the original json value was null, which is why we checked the string earlier.
		o.Replace(*tmp2)
	} else {
		o.Clear()
	}

	return nil
}
