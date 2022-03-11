package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	isPreviousChar := true
	isPreviousEscape := false
	var charCodeToRepeat rune
	var responseBuilder strings.Builder
	for index, inputCharCode := range inputString {
		switch {
		case string(inputCharCode) == `\` && !isPreviousEscape:
			isPreviousEscape = true
		case !unicode.IsDigit(inputCharCode) && isPreviousEscape && string(inputCharCode) != `\`:
			return "", ErrInvalidString
		case !unicode.IsDigit(inputCharCode) || isPreviousEscape:
			if index != 0 && isPreviousChar {
				responseBuilder.WriteRune(charCodeToRepeat)
			}
			if index == len(inputString)-1 {
				responseBuilder.WriteRune(inputCharCode)
			}
			charCodeToRepeat = inputCharCode
			isPreviousChar = true
			isPreviousEscape = false
		case unicode.IsDigit(inputCharCode) && isPreviousChar:
			if index == 0 {
				return "", ErrInvalidString
			}
			count, err := strconv.Atoi(string(inputCharCode))
			if err != nil {
				return "", ErrInvalidString
			}
			responseBuilder.WriteString(strings.Repeat(string(charCodeToRepeat), count))
			isPreviousChar = false
		default:
			return "", ErrInvalidString
		}
	}
	return responseBuilder.String(), nil
}
