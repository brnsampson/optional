package optional

import (
	"encoding/hex"
	"fmt"
)

// Byte implements LoadableOptional for the byte type.
type Byte struct {
	Option[byte]
}

func SomeByte(value byte) Byte {
	return Byte{Some(value)}
}

func NoByte() Byte {
	return Byte{None[byte]()}
}

func (o Byte) Type() string {
	return "Byte"
}

func (o *Byte) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Byte) String() string {
	if o.IsNone() {
		return "None[Byte]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Byte]"
		}
		return hex.EncodeToString([]byte{tmp})
	}
}

func (o Byte) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
		in := []byte{tmp}
		encoded := make([]byte, hex.EncodedLen(1))
		hex.Encode(encoded, in)
		return encoded, err
	}
}

func (o *Byte) UnmarshalText(text []byte) error {
	if len(text) == 1 {
		// just a byte
		o.Replace(text[0])
		return nil
	}

	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		// hex code?
		bytes, err := hex.DecodeString(tmp)
		if err != nil {
			return err
		} else if len(bytes) == 1 {
			o.Replace(bytes[0])
			return nil
		} else {
			return fmt.Errorf("could not unmarshal text into byte: %x does not fit into byte", bytes)
		}
	}
	return nil
}

// Scan implements database/sql.Scanner interface.
func (o *Byte) Scan(src any) error {
	if src == nil {
		// NULL value row
		o.Clear()
		return nil
	}
	switch src.(type) {
	case string:
		o.UnmarshalText([]byte(src.(string)))
	case []byte:
		o.UnmarshalText(src.([]byte))
	default:
		return fmt.Errorf("converting driver.Value type %T to %s", src, o.Type())
	}
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (o Byte) Value() (any, error) {
	val, ok := o.Get()
	if ok {
		return val, nil
	}
	return nil, nil
}
