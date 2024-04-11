package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	tokenizer "github.com/astak/otus-golang-homework/hw02_unpack_string/loop"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	builder := builder{}
	if err := tokenizer.Scan(str, builder.Visit); err != nil {
		return "", err
	}

	return builder.String(), nil
}

type builder struct {
	builder strings.Builder
	last    rune
}

func (st *builder) visitDigit(r rune) error {
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

func (st *builder) visitRune(r rune) error {
	if st.last != 0 {
		st.builder.WriteRune(st.last)
	}

	st.last = r
	return nil
}

func (st *builder) visitEscaped(r rune) error {
	if !(unicode.IsDigit(r) || r == tokenizer.EscapeSymbol) {
		return ErrInvalidString
	}

	return st.visitRune(r)
}

func (st *builder) Visit(t tokenizer.Token) error {
	switch t.Kind {
	case tokenizer.TokenRune:
		return st.visitRune(t.Value)
	case tokenizer.TokenDigit:
		return st.visitDigit(t.Value)
	case tokenizer.TokenEscaped:
		return st.visitEscaped(t.Value)
	default:
		return ErrInvalidString
	}
}

func (st *builder) String() string {
	if st.last != 0 {
		st.builder.WriteRune(st.last)
	}
	return st.builder.String()
}
