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

// BenchmarkSmallIntArithmetic_Composite is small-value-dominated work (operands
// and results stay within the fixnum cache band), where a small-Number cache
// can eliminate the per-Number allocation — unlike the growing-accumulator loop
// above, whose results escape the band.
func BenchmarkSmallIntArithmetic_Composite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			a := Number.New(Number.WithInt(int64(j % 100)))
			c := Number.New(Number.WithInt(int64((j + 1) % 100)))
			_ = a.Add(c) // result <= 198, in band
		}
	}
}

// BenchmarkComparison_Composite produces a Boolean per call (GreaterThan), the
// path that benefits from interned true/false.
func BenchmarkComparison_Composite(b *testing.B) {
	x := Number.New(Number.WithInt(5))
	y := Number.New(Number.WithInt(3))
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			_ = x.GreaterThan(y)
		}
	}
}
