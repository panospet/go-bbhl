package util

import (
	"fmt"
	"strconv"
	"strings"
)

var ErrBadTitle = fmt.Errorf("bad title")
var ErrStrToIntConversion = fmt.Errorf("conversion from string error")

func ExtractElRound(title string) (int, error) {
	parts := strings.Split(strings.ToLower(title), "round")
	if len(parts) != 2 {
		return 0, ErrBadTitle
	}
	secondPart := parts[1]

	spaced := strings.Split(strings.Trim(secondPart, " "), " ")
	intVar, err := strconv.Atoi(spaced[0])
	if err != nil {
		return 0, ErrStrToIntConversion
	}

	return intVar, err
}
