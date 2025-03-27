package optional

// Str implements Configand for the string type.
type Str struct {
	Option[string]
}

func SomeStr(value string) Str {
	return Str{Some(value)}
}

func NoStr() Str {
	return Str{None[string]()}
}

func (o Str) Type() string {
	return "Str"
}

func (o *Str) Set(str string) error {
	return o.UnmarshalText([]byte(str))
}

func (o Str) String() string {
	if o.IsNone() {
		return "None[Str]"
	} else {
		tmp, ok := o.Get()
		if !ok {
			return "Error[Str]"
		}
		return tmp
	}
}

func (o Str) MarshalText() (text []byte, err error) {
	if o.IsNone() {
		return []byte("None"), nil
	} else {
		tmp, ok := o.Get()
		var err error
		if !ok {
			err = optionalError("Attempted to Get Option with None value")
		}
		return []byte(tmp), err
	}
}

func (o *Str) UnmarshalText(text []byte) error {
	tmp := string(text)
	if tmp == "None" || tmp == "none" || tmp == "null" || tmp == "nil" {
		o.Clear()
	} else {
		o.Replace(tmp)
	}
	return nil
}
