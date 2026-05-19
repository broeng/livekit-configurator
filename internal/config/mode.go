package config

import (
	"fmt"
	"strings"
)

type Mode int

const (
	Overwrite Mode = iota
	Merge
	Assert
)

var modeName = map[Mode]string{
	Overwrite: "overwrite",
	Merge:     "merge",
	Assert:    "assert",
}

func (m Mode) String() string {
	return modeName[m]
}

func (m *Mode) UnmarshalText(text []byte) error {
	s := strings.ToLower(string(text))
	for k, v := range modeName {
		if v == s {
			*m = k
			return nil
		}
	}
	return fmt.Errorf("invalid Mode: %q", s)
}
