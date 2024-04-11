package tokenizer

import "testing"

func BenchmarkScan(b *testing.B) {
	fn := func(t Token) error {
		return nil
	}
	for i := 0; i < b.N; i++ {
		Scan("a4bc2d5e", fn)
	}
}
