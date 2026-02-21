package optional

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	DEFAULT_TIME_FORMAT        string = time.RFC3339Nano
	DEFAULT_TIME_STRING_FORMAT string = time.DateTime
)

var DefaultExtraTimeFormats []string = []string{time.RFC3339Nano, time.RFC3339, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z}

func SetDefaultExtraTimeFormats(formats []string) {
	DefaultExtraTimeFormats = formats
}

type Time struct {
	Option[time.Time]
	StringFormat string
	DataFormat   string
	formats      []string
}

func SomeTime(value time.Time, formats ...string) Time {
	return Time{Some(value), DEFAULT_TIME_STRING_FORMAT, DEFAULT_TIME_FORMAT, formats}
}

func NoTime(formats ...string) Time {
	return Time{None[time.Time](), DEFAULT_TIME_STRING_FORMAT, DEFAULT_TIME_FORMAT, formats}
}

func (o Time) WithFormats(formats ...string) Time {
	if len(formats) > 0 {
		o.formats = formats
	}

	return o
}

func (o *Time) defaultFormatsIfEmpty() {
	if o.StringFormat == "" {
		o.StringFormat = DEFAULT_TIME_STRING_FORMAT
	}
	if o.DataFormat == "" {
		o.DataFormat = DEFAULT_TIME_FORMAT
	}
	if len(o.formats) == 0 {
		o.formats = append(o.formats, DefaultExtraTimeFormats...)
	}
}

func (o Time) Formats() []string {
	return o.formats
}

func (o Time) Type() string {
	return "Time"
}

func (o *Time) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Time) String() string {
	if o.IsNone() {
		return "None[Time]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Time]"
		}
		o.defaultFormatsIfEmpty()
		return tmp.Format(o.StringFormat)
	}
}

func (o Time) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
		o.defaultFormatsIfEmpty()
		return []byte(tmp.Format(o.DataFormat)), err
	}
}

func (o *Time) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		l := len(o.formats)
		if l == 0 {
			// No formats?
			o.defaultFormatsIfEmpty()
			l = len(o.formats)
		}
		for i, f := range o.formats {
			t, err := time.Parse(f, tmp)
			if err == nil {
				o.Replace(t)
				break
			}
			if i == l-1 {
				// Tried all formats without success. This only gives the error for the last format, but hopefully that
				// should be enough for now. Maybe give a more informative error message in the future?
				return err
			}
		}
	}

	return nil
}

// Marshaler interface

func (o Time) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return json.Marshal(nil)
	} else {
		tmp, err := o.MarshalText()
		ret := append([]byte(`"`), tmp...)
		ret = append(ret, '"')
		return ret, err
	}
}

// UnmarshalJSON implements encoding/json.Unmarshaller interface
func (o *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Clear()
		return nil
	}

	// TODO: think if a better way to do this.
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	return o.UnmarshalText([]byte(s))
}

// Scan implements database/sql.Scanner interface.
func (o *Time) Scan(src any) error {
	if src == nil {
		// NULL value row
		o.Clear()
		return nil
	}
	// TODO: support an int64 as a unix timestamp?
	switch t := src.(type) {
	case time.Time:
		_ = o.Replace(t)
	case string:
		err := o.UnmarshalText([]byte(t))
		if err != nil {
			return err
		}
	case []byte:
		err := o.UnmarshalText(t)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("converting driver.Value type %T to %s", src, o.Type())
	}
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (o Time) Value() (driver.Value, error) {
	val, ok := o.Get()
	if ok {
		return val, nil
	}
	return nil, nil
}

// Duration is the optional version of time.Duration. Under the hood, time.Duration is an int64
// but we still use an underlying Option[time.Duration] type to get access to all the associated methods.
type Duration struct {
	Option[time.Duration]
}

func SomeDuration(value time.Duration) Duration {
	return Duration{Some(value)}
}

func NoDuration(formats ...string) Duration {
	return Duration{None[time.Duration]()}
}

func (o Duration) Type() string {
	return "Duration"
}

func (o *Duration) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Duration) String() string {
	if o.IsNone() {
		return "None[Duration]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Duration]"
		}
		// TODO: make the number of significant figures printed configurable?
		return tmp.String()
	}
}

func (o Duration) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
		return []byte(tmp.String()), err
	}
}

func (o *Duration) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		d, err := time.ParseDuration(tmp)
		if err != nil {
			return err
		}
		o.Replace(d)
	}
	return nil
}

// Marshaler interface

func (o Duration) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return json.Marshal(nil)
	} else {
		tmp, err := o.MarshalText()
		ret := append([]byte(`"`), tmp...)
		ret = append(ret, '"')
		return ret, err
	}
}

// UnmarshalJSON implements encoding/json.Unmarshaller interface
func (o *Duration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Clear()
		return nil
	}

	// TODO: think if a better way to do this.
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	return o.UnmarshalText([]byte(s))
}

// Scan implements database/sql.Scanner interface.
func (o *Duration) Scan(src any) error {
	if src == nil {
		// NULL value row
		o.Clear()
		return nil
	}
	switch src.(type) {
	case time.Duration:
		_ = o.Replace(src.(time.Duration))
	case int64:
		_ = o.Replace(time.Duration(src.(int64)))
	default:
		return fmt.Errorf("converting driver.Value type %T to %s", src, o.Type())
	}
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (o Duration) Value() (driver.Value, error) {
	val, ok := o.Get()
	if ok {
		// The set of types which sql can accept does not include time.Duration. See https://pkg.go.dev/database/sql/driver#Value
		return int64(val), nil
	}
	return nil, nil
}
