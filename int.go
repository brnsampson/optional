package optional

import (
	"strconv"
)

// default sized int
type Int struct {
	Option[int]
}

func SomeInt(value int) Int {
	return Int{Some(value)}
}

func NoInt() Int {
	return Int{None[int]()}
}

func (o Int) Type() string {
	return "Int"
}

func (o *Int) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Int) String() string {
	if o.IsNone() {
		return "None[Int]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Int]"
		}
		return strconv.Itoa(tmp)
	}
}

func (o Int) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.Itoa(tmp)), err
	}
}

func (o *Int) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.Atoi(tmp)
		if err != nil {
			return err
		}
		o.Replace(i)
	}
	return nil
}

// 8bit sized int
type Int8 struct {
	Option[int8]
}

func SomeInt8(value int8) Int8 {
	return Int8{Some(value)}
}

func NoInt8() Int8 {
	return Int8{None[int8]()}
}

func (o Int8) Type() string {
	return "Int8"
}

func (o *Int8) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Int8) String() string {
	if o.IsNone() {
		return "None[Int8]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Int8]"
		}
		return strconv.FormatInt(int64(tmp), 10)
	}
}

func (o Int8) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatInt(int64(tmp), 10)), err
	}
}

func (o *Int8) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseInt(tmp, 10, 8)
		if err != nil {
			return err
		}
		o.Replace(int8(i))
	}
	return nil
}

// 16bit sized int
type Int16 struct {
	Option[int16]
}

func SomeInt16(value int16) Int16 {
	return Int16{Some(value)}
}

func NoInt16() Int16 {
	return Int16{None[int16]()}
}

func (o Int16) Type() string {
	return "Int16"
}

func (o *Int16) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Int16) String() string {
	if o.IsNone() {
		return "None[Int16]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Int16]"
		}
		return strconv.FormatInt(int64(tmp), 10)
	}
}

func (o Int16) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatInt(int64(tmp), 10)), err
	}
}

func (o *Int16) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseInt(tmp, 10, 16)
		if err != nil {
			return err
		}
		o.Replace(int16(i))
	}
	return nil
}

// 32bit sized int
type Int32 struct {
	Option[int32]
}

func SomeInt32(value int32) Int32 {
	return Int32{Some(value)}
}

func NoInt32() Int32 {
	return Int32{None[int32]()}
}

func (o Int32) Type() string {
	return "Int32"
}

func (o *Int32) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Int32) String() string {
	if o.IsNone() {
		return "None[Int32]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Int32]"
		}
		return strconv.FormatInt(int64(tmp), 10)
	}
}

func (o Int32) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatInt(int64(tmp), 10)), err
	}
}

func (o *Int32) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseInt(tmp, 10, 32)
		if err != nil {
			return err
		}
		o.Replace(int32(i))
	}
	return nil
}

// 64bit sized int
type Int64 struct {
	Option[int64]
}

func SomeInt64(value int64) Int64 {
	return Int64{Some(value)}
}

func NoInt64() Int64 {
	return Int64{None[int64]()}
}

func (o Int64) Type() string {
	return "Int64"
}

func (o *Int64) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Int64) String() string {
	if o.IsNone() {
		return "None[Int64]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Int64]"
		}
		return strconv.FormatInt(tmp, 10)
	}
}

func (o Int64) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatInt(tmp, 10)), err
	}
}

func (o *Int64) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseInt(tmp, 10, 64)
		if err != nil {
			return err
		}
		o.Replace(i)
	}
	return nil
}
