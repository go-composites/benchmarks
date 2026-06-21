package benchmarks

import (
	"strconv"
	"testing"

	Dictionary "github.com/go-composites/dictionary/src"
)

// keys is a fixed set of N string keys, generated once, so the benchmarks
// time the map operations rather than the key construction.
var keys = func() []string {
	ks := make([]string, N)
	for i := range ks {
		ks[i] = "key-" + strconv.Itoa(i)
	}
	return ks
}()

// BenchmarkDictInsert_Composite inserts N entries through the composite
// Dictionary (interface{} keys/values, internal map[interface{}]interface{}).
func BenchmarkDictInsert_Composite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := Dictionary.New()
		for j := 0; j < N; j++ {
			d.Set(keys[j], j)
		}
	}
}

// BenchmarkDictInsert_RawGo inserts N entries into a plain map[string]int.
func BenchmarkDictInsert_RawGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[string]int)
		for j := 0; j < N; j++ {
			m[keys[j]] = j
		}
	}
}

// BenchmarkDictLookup_Composite looks up all N keys through the composite
// Dictionary, each Get returning a Result.
func BenchmarkDictLookup_Composite(b *testing.B) {
	d := Dictionary.New()
	for j := 0; j < N; j++ {
		d.Set(keys[j], j)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum int
		for j := 0; j < N; j++ {
			if v := d.Get(keys[j]); !v.HasError() {
				sum += v.Payload().(int)
			}
		}
		_ = sum
	}
}

// BenchmarkDictLookup_RawGo looks up all N keys in a plain map[string]int.
func BenchmarkDictLookup_RawGo(b *testing.B) {
	m := make(map[string]int, N)
	for j := 0; j < N; j++ {
		m[keys[j]] = j
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum int
		for j := 0; j < N; j++ {
			if v, ok := m[keys[j]]; ok {
				sum += v
			}
		}
		_ = sum
	}
}
