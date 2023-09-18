package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var ErrInvalidDigitParse = errors.New("invalid digit parse")

func Unpack(str string) (string, error) {
	tableID := [3][5]int{
		{2, 0, 3, 0, 4},
		{2, 1, 3, 0, 4},
		{0, 2, 2, 0, 0},
	}
	state := 1
	backSlash := '\\'
	var unpackString string
	for i, val := range str {
		prevState := state
		if unicode.IsLetter(val) {
			state = tableID[state-1][0]
		} else if unicode.IsDigit(val) {
			state = tableID[state-1][1]
		} else if val == backSlash {
			state = tableID[state-1][2]
		} else {
			state = tableID[state-1][3]
		}

		if state == 0 {
			return "", ErrInvalidString
		}

		if state == 1 {
			if prevState == 2 {
				num, err := strconv.Atoi(string(val))
				if err != nil {
					return "", ErrInvalidDigitParse
				}
				unpackString = unpackString[:len(unpackString)-1]
				unpackString += strings.Repeat(string(str[i-1]), num)
			} else if prevState == 3 {
				unpackString += string(val)
			}
		} else if state == 2 {
			unpackString += string(val)
		}
	}

	state = tableID[state-1][4]
	if state == 0 {
		return "", ErrInvalidString
	}
	return unpackString, nil
}
