package common

import "strings"

type String struct {
	value string
}

type StringSlice struct {
	values []string
}

func NewString(value string) *String {
	return &String{
		value: value,
	}
}

func (s *String) ToLower() *String {
	s.value = strings.ToLower(s.value)
	return s
}

func (s *String) TrimPrefix(prefix string) *String {
	s.value = strings.TrimPrefix(s.value, prefix)
	return s
}

func (s *String) Trim(prefix string) *String {
	s.value = strings.Trim(s.value, prefix)
	return s
}

func (s *String) Split(separator string) *StringSlice {
	values := &StringSlice{
		values: strings.Split(s.value, separator),
	}
	return values
}

func (s *String) Value() string {
	return s.value
}

func (s *StringSlice) Join(separator string) *String {
	return &String{
		value: strings.Join(s.values, separator),
	}
}
