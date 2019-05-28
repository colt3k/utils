package webut

import (
	"fmt"
	"strconv"
	"strings"
)

type ConvertibleBoolean bool

func (bit *ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)

	if strings.HasPrefix(asString, "\"") {
		asString = strings.TrimPrefix(asString, "\"")
	}
	if strings.HasSuffix(asString, "\"") {
		asString = strings.TrimSuffix(asString, "\"")
	}
	b, err := strconv.ParseBool(asString)
	if err != nil {
		return fmt.Errorf("boolean unmarshal error: invalid input %s", asString)
	}
	if b {
		*bit = true
	} else {
		*bit = false
	}

	return nil
}
