package optional

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Bool struct {
	Option[bool]
}

func SomeBool(value bool) Bool {
	return Bool{Some(value)}
}

func NoBool() Bool {
	return Bool{None[bool]()}
}

func (o Bool) Type() string {
	return "Bool"
}

func (o *Bool) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Bool) String() string {
	if o.IsNone() {
		return "None[Bool]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Bool]"
		}
		return strconv.FormatBool(tmp)
	}
}

func (o Bool) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
		return []byte(strconv.FormatBool(tmp)), err
	}
}

func (o *Bool) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		i, err := strconv.ParseBool(tmp)
		if err != nil {
			return err
		}
		o.Replace(i)
	}
	return nil
}

// Scan implements database/sql.Scanner interface.
func (o *Bool) Scan(src any) error {
	if src == nil {
		// NULL value row
		o.Clear()
		return nil
	}
	switch t := src.(type) {
	case bool:
		_ = o.Replace(t)
	case string:
		tmp, err := strconv.ParseBool(t)
		if err != nil {
			return err
		}
		_ = o.Replace(tmp)
	default:
		return fmt.Errorf("converting driver.Value type %T to %s", src, o.Type())
	}
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (o Bool) Value() (driver.Value, error) {
	val, ok := o.Get()
	if ok {
		return val, nil
	}
	return nil, nil
}

// True returns true iff the value is Some(true). It is a special method exclusive to Bool
// optionals, and is the same as calling Bool.Match(true).
func (o *Bool) True() bool {
	return o.Match(true)
}
