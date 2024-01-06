package config

import (
	"time"
	"encoding/json"
	"github.com/brnsampson/optional"
)

const DEFAULT_TIME_FORMAT string = time.RFC3339

type Time struct {
	optional.Option[time.Time]
	formats []string
}

func SomeTime(value time.Time, formats ...string) Time {
	if len(formats) == 0 {
		formats = append(formats, DEFAULT_TIME_FORMAT)
	}
	return Time{optional.Some(value), formats}
}

func NoTime(formats ...string) Time {
	if len(formats) == 0 {
		formats = append(formats, DEFAULT_TIME_FORMAT)
	}
	return Time{optional.None[time.Time](), formats}
}

func (o Time) WithFormats(formats ...string) Time {
	if len(formats) > 0 {
		o.formats = formats
	}

	return o
}

func (o *Time) defaultFormatsIfEmpty() {
	if len(o.formats) == 0 {
		o.formats = append(o.formats, DEFAULT_TIME_FORMAT)
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
		tmp, err := o.Get()
		if err != nil {
			return "Error[Time]"
		}
		return tmp.Format(o.formats[0])
	}
}

func (o Time) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, err := o.Get()
		return []byte(tmp.Format(o.formats[0])), err
	}
}

func (o *Time) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.defaultFormatsIfEmpty()
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
				o.SetVal(t)
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

// Unmarshaller interface
func (o *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Clear()
		o.defaultFormatsIfEmpty()
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
