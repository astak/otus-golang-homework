package tokenizer

import "unicode"

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

func Scan(str string, accept func(token Token) error) error {
	var prev rune
	for _, r := range str {
		if r == EscapeSymbol && prev != EscapeSymbol {
			prev = r
			continue
		}

		var token Token
		switch {
		case prev == EscapeSymbol:
			token = Token{
				Kind:  TokenEscaped,
				Value: r,
			}
		case unicode.IsDigit(r):
			token = Token{
				Kind:  TokenDigit,
				Value: r,
			}
		default:
			token = Token{
				Kind:  TokenRune,
				Value: r,
			}
		}

		prev = 0
		err := accept(token)
		if err != nil {
			return err
		}
	}
	return nil
}
