package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	b := strings.Builder{}

	var nextIsDigit bool
	var curDigit bool
	rStr := []rune(s)

	for ind, val := range rStr {
		nextIsDigit = false
		curDigit = unicode.IsDigit(val)
		if ind+1 < len(rStr) {
			nextIsDigit = unicode.IsDigit(rStr[ind+1])
		}

		if curDigit && (ind == 0 || nextIsDigit) {
			return "", ErrInvalidString
		}

		if !curDigit {
			if nextIsDigit {
				cnt, _ := strconv.Atoi(string(rStr[ind+1]))
				if cnt > 0 {
					b.WriteString(strings.Repeat(string(val), cnt))
				}
			} else {
				b.WriteRune(val)
			}
		}
	}

	return b.String(), nil
}
