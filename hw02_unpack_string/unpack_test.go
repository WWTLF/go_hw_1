package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebug(t *testing.T) {
	result, err := Unpack(`qwз\\5`)
	require.NoError(t, err)
	require.Equal(t, `qwз\\\\\`, result)
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "а4бц2д5е", expected: "аааабццддддде"},
		{input: "a4bc2d5гг", expected: "aaaabccdddddгг"},
		{input: "a4bc2з5e", expected: "aaaabccзззззe"},
		{input: "abccd", expected: "abccd"},
		{input: "abccdз", expected: "abccdз"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "ззз0b", expected: "ззb"},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\4г\5`, expected: `qwe4г5`},
		{input: `qwe\4г\5г2д`, expected: `qwe4г5ггд`},
		{input: `ггг\4\5`, expected: `ггг45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `гwг\45`, expected: `гwг44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwз\\5`, expected: `qwз\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwз\\\3`, expected: `qwз\3`},
		{input: `абц\\\3`, expected: `абц\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{`3abc`, `45`, `aaa10b`, `qw\ne`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
