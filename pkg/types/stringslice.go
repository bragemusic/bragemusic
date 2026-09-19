package types

import (
	"fmt"
	"strings"
)

type StringSlice []string

func (s *StringSlice) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	*s = strings.Split(str, ",")
	return nil
}
