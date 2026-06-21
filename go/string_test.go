package benchmarks

import (
	"strings"
	"testing"

	String "github.com/go-composites/string/src"
)

// chunk is the piece appended on every step of the string-build benchmarks.
const chunk = "abcdefghij"

// BenchmarkStringBuild_Composite builds a string of N chunks by repeated
// Concat through the composite String. Each Concat allocates a Result and a
// fresh String holding the whole accumulated value, so this is quadratic by
// construction — exactly the cost the composite hides behind a message send.
func BenchmarkStringBuild_Composite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := String.New(String.WithGoString(""))
		piece := String.New(String.WithGoString(chunk))
		for j := 0; j < N; j++ {
			s = s.Concat(piece).Payload().(String.Interface)
		}
		_ = s.ToGoString()
	}
}

// BenchmarkStringBuild_RawGo builds the same string with strings.Builder,
// the idiomatic linear-time approach a Go programmer would actually write.
func BenchmarkStringBuild_RawGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < N; j++ {
			sb.WriteString(chunk)
		}
		_ = sb.String()
	}
}

// BenchmarkStringBuild_RawGoNaive builds the string with naive "+="
// concatenation. This mirrors the composite Concat algorithm (O(n^2),
// re-copying the whole accumulator each step) so the composite-vs-raw factor
// here isolates the boxing + Result overhead rather than the algorithm.
func BenchmarkStringBuild_RawGoNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < N; j++ {
			s += chunk
		}
		_ = s
	}
}
