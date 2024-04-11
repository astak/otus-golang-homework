package tokenizer

import (
	"unicode"
	"unicode/utf8"
)

const EscapeSymbol = '\\'

const (
	TokenRune = iota
	TokenDigit
	TokenEscaped
)

type Token struct {
	Kind  int
	Value rune
}

func NewTokenizer(str string) func() (*Token, bool) {
	data := []byte(str)
	position := 0

	return func() (*Token, bool) {
		r, size := utf8.DecodeRune(data[position:])
		if r == utf8.RuneError {
			return nil, false
		}
		position += size

		if r == EscapeSymbol {
			r, size = utf8.DecodeRune(data[position:])
			position += size

			return &Token{
				Kind:  TokenEscaped,
				Value: r,
			}, true
		}

		if unicode.IsDigit(r) {
			return &Token{
				Kind:  TokenDigit,
				Value: r,
			}, true
		}

		return &Token{
			Kind:  TokenRune,
			Value: r,
		}, true
	}
}

func Scan(str string, accept func(token *Token) error) error {
	next := NewTokenizer(str)
	for t, ok := next(); ok; t, ok = next() {
		if err := accept(t); err != nil {
			return err
		}
	}
	return nil
}
