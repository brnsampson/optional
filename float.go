package optional

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// 32bit sized floats

type Float32 struct {
	Option[float32]
}

func SomeFloat32(value float32) Float32 {
	return Float32{Some(value)}
}

func NoFloat32() Float32 {
	return Float32{None[float32]()}
}

func (o Float32) Type() string {
	return "Float32"
}

func (o *Float32) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Float32) String() string {
	if o.IsNone() {
		return "None[Float32]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Float32]"
		}
		return strconv.FormatFloat(float64(tmp), 'g', -1, 32)
	}
}

func (o Float32) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
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
		o.Replace(float32(i))
	}
	return nil
}

// 64bit sized floats

type Float64 struct {
	Option[float64]
}

func SomeFloat64(value float64) Float64 {
	return Float64{Some(value)}
}

func NoFloat64() Float64 {
	return Float64{None[float64]()}
}

func (o Float64) Type() string {
	return "Float64"
}

func (o *Float64) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Float64) String() string {
	if o.IsNone() {
		return "None[Float64]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Float64]"
		}
		return strconv.FormatFloat(tmp, 'g', -1, 64)
	}
}

func (o Float64) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
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
		o.Replace(i)
	}
	return nil
}

// Scan implements database/sql.Scanner interface.
func (o *Float64) Scan(src any) error {
	if src == nil {
		// NULL value row
		o.Clear()
		return nil
	}
	switch t := src.(type) {
	case float64:
		_ = o.Replace(t)
	case string:
		if err := o.UnmarshalText([]byte(t)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("converting driver.Value type %T to %s", src, o.Type())
	}
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (o Float64) Value() (driver.Value, error) {
	val, ok := o.Get()
	if ok {
		return val, nil
	}
	return nil, nil
}
