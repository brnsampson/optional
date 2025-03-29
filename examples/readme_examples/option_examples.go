package main

import (
	"encoding/json"
	"fmt"

	"github.com/brnsampson/optional"
)

func DefiningOptionalValues() {
	// There are types for all the primatives we would normally expect
	// Bool
	// Int Int16 Int32 Int64
	// Uint Uint16 Uint32 Uint64
	// Float32 Float64
	// Str
	// Time
	// and the generic Option[T comparable]

	// Create an Optional Int with no initial value
	// The zero value of an optional is None
	var i optional.Int

	// This also works fine
	i = optional.Int{}

	// However, I normally use this functional form for symmetry with creating
	// Options with values, optional.SomeInt(<my int>)
	i = optional.NoInt()

	// Check if i is None (empty)
	if i.IsNone() {
		fmt.Println("i is empty!")
	}

	// Set the value to some default if it was previously unset
	i.Default(42)

	// Update the value and get the previous value back for any comparisons you might need to do
	previous := i.Replace(42)

	// Some methods like Option.Replace() return an Optional interface type. This erases the
	// concrete type and hides all of the methods which could mutate the value,
	// as the previous value is only provided as a reference. Unfortunatly, this also
	// hides some convenient things like the implemntations of TextMarshaler and Stringer
	if previous.IsSome() {
		fmt.Println("Replaced previous value:")
		fmt.Println(previous.MustGet())
	}

	// Overwrite the previous value without care
	i = optional.SomeInt(42)
}

func InspectingValues() {
	i := optional.SomeInt(42)

	// Check if i has a value
	if i.IsSome() {
		fmt.Println("i has a value!")
	}

	// We can check to make sure i is 42
	if i.Match(42) {
		fmt.Println("i was indeed 42!")
	} else {
		panic("wtf?")
	}

	// Get i's value along with an 'ok' boolean telling us if the read is valid
	val, ok := i.Get()
	if ok {
		fmt.Println("Got i's value:")
		fmt.Println(val)
	}

	// Get i's value, but just panic if i is None
	val = i.MustGet()

	// Get i's value or a default value if i is None
	tmp := optional.GetOr(i, 123)
	fmt.Println("Got i's value or 123:")
	fmt.Println(tmp)

	// Get i's value or a default value AND set i to the default value if it is used
	// Note that helper functions require a MutableOptional interface, which only Option
	// Pointers fulful. That should be a given, since it's just like passing an int;
	// you can't expect a function to modify an int, it can only return a new int.
	tmp, err := optional.GetOrInsert(&i, 42)
	if err != nil {
		fmt.Println("Error while replacing i's value with default")
	} else {
		fmt.Println("Got i's value which should DEFINITELY be 42:")
		fmt.Println(tmp)
	}

	// For functions that automatically convert types into their string representation,
	// the Option can be used directly:
	fmt.Println("Printing i directly:")
	fmt.Println(i)
}

func MarshalingExamples() {
	i := optional.SomeInt(42)
	f := optional.SomeFloat32(12.34)
	s := optional.SomeStr("Hello!")
	nope := optional.NoStr()

	// Let's create a text string first using Sprintf. We can't use more specific verbs like
	// %d or %f because we have no way to represent None.
	newString := fmt.Sprintf("i: %s, f: %s, s: %s, nothing: %s", i, f, s, nope)
	fmt.Println(newString)

	// Options do have TextMarshaler and String methods implemented though, so we can equally well use %v
	newString = fmt.Sprintf("i: %v, f: %v, s: %v, nothing: %v", i, f, s, nope)
	fmt.Println(newString)

	// Now let's marshal a json string
	type MyStruct struct {
		Int        optional.Int
		Float      optional.Float32
		GoodString optional.Str
		BadString  optional.Str
	}

	myStruct := MyStruct{i, f, s, nope}
	jsonEncoded, err := json.Marshal(myStruct)
	if err != nil {
		fmt.Println("Failed to marshal json from struct!")
	} else {
		fmt.Println(string(jsonEncoded))
	}
}

func TransformationExamples() {
	// Define our value and transformation first
	i := optional.SomeInt(42)
	transform := func(x int) (int, error) { return x + 5, nil }

	// Modify the value in an Option without unpacking it
	err := i.Transform(transform)
	if err != nil {
		fmt.Println("The transform function returned an error!")
	}

	// Apply our transform to a slice of options, while modifying None values to be their index in the slice.
	// Remember, the zero value is None
	opts := make([]optional.Int, 10)
	for i, opt := range opts {
		// Functions which modify options in place should accept the MutableOptional interface which
		// is implemented by Option pointer types, such as this helper function. Try to use optional.TransformOr
		// with opt instead of &opt. It doesn't work, just in the same way that passing an int into a function
		// and expecting the integer to be changed in place doesn't work.
		err = optional.TransformOr(&opt, transform, i)
		if err != nil {
			fmt.Println("The transform function returned an error!")
		}
	}
}
