package model

type NullableString string

func (s NullableString) Get() *string {
	var str = string(s)
	return &str
}

func (s *NullableString) Set(str *string) {
	*s = NullableString(*str)
}

// IsSet reports whether the value is non-empty. An empty string counts as not set, so
// MarshalJSON leaves the field out instead of sending "" (the API rejects "" for time).
func (s NullableString) IsSet() bool {
	return s != ""
}

func (s *NullableString) Unset() {
	*s = ""
}
