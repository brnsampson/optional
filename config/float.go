package config

import (
	"strconv"
	"github.com/brnsampson/optional"
)

// 32bit sized floats
type Float32 struct {
	optional.Option[float32]
}

func SomeFloat32(value float32) Float32 {
	return Float32{optional.Some(value)}
}

func NoFloat32() Float32 {
	return Float32{optional.None[float32]()}
}

func (o Float32) Type() string {
	return "Float32"
}

func (o Float32) String() string {
	if o.IsNone() {
		return "None[Float32]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Float32]"
		}
		return strconv.FormatFloat(float64(tmp), 'g', -1, 32)
	}
}

func (o Float32) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatFloat(float64(tmp), 'g', -1, 32)), err
	}
}

func (o *Float32) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseFloat(tmp, 32)
		if err != nil {
			return err
		}
		o.Set(float32(i))
	}
	return nil
}

// 64bit sized floats
type Float64 struct {
	optional.Option[float64]
}

func SomeFloat64(value float64) Float64 {
	return Float64{optional.Some(value)}
}

func NoFloat64() Float64 {
	return Float64{optional.None[float64]()}
}

func (o Float64) Type() string {
	return "Float64"
}

func (o Float64) String() string {
	if o.IsNone() {
		return "None[Float64]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Float64]"
		}
		return strconv.FormatFloat(tmp, 'g', -1, 64)
	}
}

func (o Float64) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatFloat(tmp, 'g', -1, 64)), err
	}
}

func (o *Float64) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			return err
		}
		o.Set(i)
	}
	return nil
}
