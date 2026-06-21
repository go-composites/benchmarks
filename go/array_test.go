package benchmarks

import (
	"testing"

	Array "github.com/go-composites/array/src"
	Result "github.com/go-composites/result/src"
)

// BenchmarkArrayPush_Composite builds an array of N elements via the
// go-composites Array, which boxes every element into interface{} and
// allocates a Result per Push.
func BenchmarkArrayPush_Composite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := Array.New()
		for j := 0; j < N; j++ {
			a.Push(j)
		}
	}
}

// BenchmarkArrayPush_RawGo builds the same array with a plain []int.
func BenchmarkArrayPush_RawGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var a []int
		for j := 0; j < N; j++ {
			a = append(a, j)
		}
	}
}

// BenchmarkArrayIterate_Composite sums N elements through the composite
// Each iterator (each callback returns a Result).
func BenchmarkArrayIterate_Composite(b *testing.B) {
	a := Array.New()
	for j := 0; j < N; j++ {
		a.Push(j)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		a.Each(func(_ int, item interface{}) Result.Interface {
			sum += item.(int)
			return Result.New()
		})
		_ = sum
	}
}

// BenchmarkArrayIterate_RawGo sums N elements with a plain range loop.
func BenchmarkArrayIterate_RawGo(b *testing.B) {
	a := make([]int, 0, N)
	for j := 0; j < N; j++ {
		a = append(a, j)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for _, v := range a {
			sum += v
		}
		_ = sum
	}
}
