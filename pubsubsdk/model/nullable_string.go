package model

type NullableString string

func (s NullableString) Get() *string {
	var str = string(s)
	return &str
}

func (s *NullableString) Set(str *string) {
	*s = NullableString(*str)
}

func (s NullableString) IsSet() bool {
	return true
}

func (s *NullableString) Unset() {
	*s = ""
}
