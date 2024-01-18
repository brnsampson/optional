package confopt

import (
	"github.com/brnsampson/optional"
	"strconv"
)

// default sized uint
type Uint struct {
	optional.Option[uint]
}

func SomeUint(value uint) Uint {
	return Uint{optional.Some(value)}
}

func NoUint() Uint {
	return Uint{optional.None[uint]()}
}

func (o Uint) Type() string {
	return "Uint"
}

func (o *Uint) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Uint) String() string {
	if o.IsNone() {
		return "None[Uint]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Uint]"
		}
		return strconv.FormatUint(uint64(tmp), 10)
	}
}

func (o Uint) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatUint(uint64(tmp), 10)), err
	}
}

func (o *Uint) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseUint(tmp, 10, 0)
		if err != nil {
			return err
		}
		o.Replace(uint(i))
	}
	return nil
}

// 8bit sized uint
type Uint8 struct {
	optional.Option[uint8]
}

func SomeUint8(value uint8) Uint8 {
	return Uint8{optional.Some(value)}
}

func NoUint8() Uint8 {
	return Uint8{optional.None[uint8]()}
}

func (o Uint8) Type() string {
	return "Uint8"
}

func (o *Uint8) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Uint8) String() string {
	if o.IsNone() {
		return "None[Uint8]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Uint8]"
		}
		return strconv.FormatUint(uint64(tmp), 10)
	}
}

func (o Uint8) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatUint(uint64(tmp), 10)), err
	}
}

func (o *Uint8) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseUint(tmp, 10, 8)
		if err != nil {
			return err
		}
		o.Replace(uint8(i))
	}
	return nil
}

// 16bit sized uint
type Uint16 struct {
	optional.Option[uint16]
}

func SomeUint16(value uint16) Uint16 {
	return Uint16{optional.Some(value)}
}

func NoUint16() Uint16 {
	return Uint16{optional.None[uint16]()}
}

func (o Uint16) Type() string {
	return "Uint16"
}

func (o *Uint16) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Uint16) String() string {
	if o.IsNone() {
		return "None[Uint16]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Uint16]"
		}
		return strconv.FormatUint(uint64(tmp), 10)
	}
}

func (o Uint16) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatUint(uint64(tmp), 10)), err
	}
}

func (o *Uint16) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseUint(tmp, 10, 16)
		if err != nil {
			return err
		}
		o.Replace(uint16(i))
	}
	return nil
}

// 32bit sized uint
type Uint32 struct {
	optional.Option[uint32]
}

func SomeUint32(value uint32) Uint32 {
	return Uint32{optional.Some(value)}
}

func NoUint32() Uint32 {
	return Uint32{optional.None[uint32]()}
}

func (o Uint32) Type() string {
	return "Uint32"
}

func (o *Uint32) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Uint32) String() string {
	if o.IsNone() {
		return "None[Uint32]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Uint32]"
		}
		return strconv.FormatUint(uint64(tmp), 10)
	}
}

func (o Uint32) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatUint(uint64(tmp), 10)), err
	}
}

func (o *Uint32) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseUint(tmp, 10, 32)
		if err != nil {
			return err
		}
		o.Replace(uint32(i))
	}
	return nil
}

// 64bit sized uint
type Uint64 struct {
	optional.Option[uint64]
}

func SomeUint64(value uint64) Uint64 {
	return Uint64{optional.Some(value)}
}

func NoUint64() Uint64 {
	return Uint64{optional.None[uint64]()}
}

func (o Uint64) Type() string {
	return "Uint64"
}

func (o *Uint64) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Uint64) String() string {
	if o.IsNone() {
		return "None[Uint64]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[Uint64]"
		}
		return strconv.FormatUint(tmp, 10)
	}
}

func (o Uint64) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(strconv.FormatUint(tmp, 10)), err
	}
}

func (o *Uint64) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseUint(tmp, 10, 64)
		if err != nil {
			return err
		}
		o.Replace(i)
	}
	return nil
}
