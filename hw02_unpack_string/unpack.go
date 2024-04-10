package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const escapeSymbol = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	next := newTokenizer(str)
	st := state{}

	for r, ok := next(); ok; r, ok = next() {
		if err := st.Visit(r); err != nil {
			return "", err
		}
	}

	return st.String(), nil
}

const (
	tokenRune = iota
	tokenDigit
	tokenEscaped
)

type token struct {
	kind  int
	value rune
}

func newTokenizer(str string) func() (*token, bool) {
	data := str
	position := 0

	return func() (*token, bool) {
		r, size := utf8.DecodeRuneInString(data[position:])
		if r == utf8.RuneError {
			return nil, false
		}
		position += size

		if r == escapeSymbol {
			r, size = utf8.DecodeRuneInString(data[position:])
			position += size

			return &token{
				kind:  tokenEscaped,
				value: r,
			}, true
		}

		if unicode.IsDigit(r) {
			return &token{
				kind:  tokenDigit,
				value: r,
			}, true
		}

		return &token{
			kind:  tokenRune,
			value: r,
		}, true
	}
}

type state struct {
	builder strings.Builder
	last    rune
}

func (st *state) visitDigit(r rune) error {
	if st.last == 0 {
		return ErrInvalidString
	}

	n, err := strconv.Atoi(string(r))
	if err != nil {
		return err
	}

	st.builder.WriteString(strings.Repeat(string(st.last), n))
	st.last = 0
	return nil
}

func (st *state) visitRune(r rune) error {
	if st.last != 0 {
		st.builder.WriteRune(st.last)
	}

	st.last = r
	return nil
}

func (st *state) visitEscaped(r rune) error {
	if !(unicode.IsDigit(r) || r == escapeSymbol) {
		return ErrInvalidString
	}

	return st.visitRune(r)
}

func (st *state) Visit(t *token) error {
	if t == nil {
		return ErrInvalidString
	}

	switch t.kind {
	case tokenRune:
		return st.visitRune(t.value)
	case tokenDigit:
		return st.visitDigit(t.value)
	case tokenEscaped:
		return st.visitEscaped(t.value)
	default:
		return ErrInvalidString
	}
}

func (st *state) String() string {
	if st.last != 0 {
		st.builder.WriteRune(st.last)
	}
	return st.builder.String()
}
