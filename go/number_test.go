package benchmarks

import (
	"testing"

	Number "github.com/go-composites/number/src"
)

// BenchmarkNumberArithmetic_Composite runs an N-step accumulation through
// the composite Number. Each Add allocates a Result wrapping a fresh
// Number, then the payload is type-asserted back out.
func BenchmarkNumberArithmetic_Composite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		acc := Number.New(Number.WithInt(0))
		for j := 0; j < N; j++ {
			acc = acc.Add(Number.New(Number.WithInt(int64(j)))).Payload().(Number.Interface)
		}
		_ = acc.ToGoInt()
	}
}

// BenchmarkNumberArithmetic_RawGo runs the same accumulation with a plain
// int64.
func BenchmarkNumberArithmetic_RawGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var acc int64
		for j := 0; j < N; j++ {
			acc += int64(j)
		}
		_ = acc
	}
}
