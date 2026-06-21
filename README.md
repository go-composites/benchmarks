<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/benchmarks" width="720"></p>

# benchmarks

[![ci](https://github.com/go-composites/benchmarks/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/benchmarks/actions/workflows/ci.yml)

Honest, measured benchmarks of the [go-composites](https://github.com/go-composites)
composite style. For every operation we triangulate three implementations:

1. **go-composites** — the composite (`Array`, `Dictionary`, `Number`, `String`, …).
2. **idiomatic raw Go** — what a Go programmer would actually write (`[]int`,
   `map[string]int`, `int64`, `strings.Builder`).
3. **Ruby MRI** — the reference dynamic language, native `Array`/`Hash`/`Integer`/`String`.

The point is to quantify, not to flatter: the composite style trades raw speed
for uniformity, null-safety and a value-oriented `Result` API. This repo says by
**how much**.

## TL;DR

go-composites is consistently slower than idiomatic raw Go — from **~1.7×** on
hash insert up to **~130×** on a tight integer-arithmetic loop — because every
operation **boxes its operands into `interface{}`** and **allocates a `Result`**
(and often a fresh composite) per call. Where the composite is doing real
container work (map insert, hash lookup) it generally still **beats Ruby MRI**;
where the operation is a primitive the compiler would otherwise keep in a
register (integer add, array iterate), the per-op allocation makes it **as slow
as, or slower than, Ruby**.

## Results

Machine: Apple M4 Max, macOS (`arm64-darwin`). Go 1.26.4. Ruby 4.0.5.
Working-set size `N = 1000` for every operation (one "op" = one full N-sized
build / iterate / lookup pass). Go figures are `go test -bench -benchmem`
ns/op; Ruby figures are wall time from the stdlib `Benchmark` module divided by
the repetition count (see [methodology](#methodology)). Numbers are a single
representative run — re-run `./run.sh` to reproduce.

| Operation (N=1000)        | go-composites | raw Go    | Ruby MRI  | composite ÷ raw Go | composite allocs |
|---------------------------|--------------:|----------:|----------:|-------------------:|-----------------:|
| array push                |    24 267 ns  |  2 310 ns | 37 030 ns |        **10.5×**   | 1 756 allocs/op  |
| array iterate (sum)       |    15 958 ns  |    260 ns | 25 980 ns |        **61×**     | 1 001 allocs/op  |
| dictionary insert         |    70 776 ns  | 42 933 ns | 115 211 ns|        **1.65×**   | 1 767 allocs/op  |
| dictionary lookup         |    28 609 ns  |  6 268 ns | 59 901 ns |        **4.6×**    | 1 000 allocs/op  |
| integer arithmetic loop   |    34 458 ns  |    256 ns | 27 251 ns |        **134×**    | 3 001 allocs/op  |
| string build (vs Builder) |   507 245 ns  |  5 572 ns | 35 875 ns |        **91×**     | 3 001 allocs/op  |
| string build (vs naive +=)|   507 245 ns  | 614 673 ns| 359 720 ns|        **0.83×**   | 3 001 allocs/op  |

Raw Go `allocs/op` for reference: array push 12, iterate 0, dict insert 20,
dict lookup 0, integer loop 0, string-Builder 15, string-naive 999.

## Honest analysis

**Where the cost comes from.** Two structural taxes, neither of which raw Go
pays:

- **`interface{}` boxing.** Every element handed to `Array.Push`, every key and
  value in `Dictionary`, every operand of `Number.Add` is stored as
  `interface{}`. Boxing an `int` heap-allocates it; the raw-Go `[]int` /
  `map[string]int` / `int64` keep the value unboxed in place. This is why the
  composite columns show ~1 000–3 000 allocs/op where raw Go shows 0–20.
- **Per-op `Result` allocation.** Fallible composite operations return a
  `Result` value (so failures are values, not panics). That is one heap
  allocation on the hot path of *every* call. `Number.Add` is the worst case:
  the loop allocates a `Result` **and** a fresh `Number` per step (≈3 allocs ×
  1000 = 3001).

**How big is the tax?** It tracks how much real work the underlying operation
does:

- The closer the operation is to a single machine instruction, the worse the
  composite looks. **Integer arithmetic (134×)** and **array iterate (61×)** are
  the extremes: raw Go keeps the accumulator in a register and the loop body is
  essentially free (256 ns / 260 ns for 1000 iterations, 0 allocs), so the
  composite's boxing+`Result` overhead is *the entire runtime*.
- The more genuine container work per call, the smaller the relative tax.
  **Dictionary insert is only 1.65×** because Go's `map` insert (hashing,
  probing, possible growth) is itself expensive enough to dominate, leaving the
  boxing as a minority cost. Array push (10.5×) and dict lookup (4.6×) sit in
  between.

**Versus Ruby MRI.** This is the more flattering comparison, and it splits
cleanly:

- On real container work, go-composites **beats** interpreted Ruby: dict insert
  70 776 ns vs 115 211 ns, dict lookup 28 609 ns vs 59 901 ns, array push
  24 267 ns vs 37 030 ns. Even paying the composite tax, compiled Go with a
  native runtime stays ahead of the MRI interpreter.
- On primitive-bound loops the picture flips: the **integer loop is *slower*
  than Ruby** (34 458 ns vs 27 251 ns). MRI's fixnums avoid allocation for small
  integers, while the composite allocates a `Number` + `Result` every step — so
  here the composite gives up Go's usual advantage and lands behind the
  reference dynamic language. Array iterate (15 958 ns vs 25 980 ns) stays
  ahead, but by less than 2× rather than the orders of magnitude raw Go enjoys.

**The string outlier.** `String.Concat` rebuilds the whole accumulated string
each call, so building by repeated `Concat` is **O(n²)** — 507 µs. Compared
against the idiomatic `strings.Builder` (linear, 5.6 µs) that is a **91×**
blow-up, but most of that is *algorithm*, not composite overhead: the raw-Go
**naive `+=`** variant, which uses the same O(n²) shape, is actually *slower*
(615 µs) — the composite comes out at **0.83× of naive `+=`**. Ruby's mutating
`String#<<` (linear, 35.9 µs) and naive `+` (O(n²), 359.7 µs) bracket the same
two regimes. Takeaway: the composite `String` has no linear "builder" path, so
accumulation patterns that are cheap in idiomatic Go become quadratic — choose
the operation, not just the type.

**Bottom line.** The composite style is not free and we don't pretend it is.
Budget **roughly 2–10× over idiomatic Go for container-shaped work** and **one
to two orders of magnitude for primitive-shaped work** (integer math, tight
iteration), driven entirely by `interface{}` boxing and per-op `Result`
allocation. It usually still beats Ruby MRI on container work, and can fall
behind it on allocation-bound primitive loops. Reach for go-composites when you
want its uniformity, null-safety and value-error ergonomics — not when a tight
numeric inner loop is the bottleneck.

## Methodology

- **Same sizes both sides.** Every operation uses `N = 1000`
  ([`go/sizes.go`](go/sizes.go), `N` in [`ruby/common.rb`](ruby/common.rb)). One
  reported "op" is one full N-sized pass (e.g. pushing 1000 elements).
- **Go.** `go test -bench=. -benchmem ./go/...`. `testing.B` auto-scales the
  iteration count and reports ns for one op directly. Paired benchmarks
  (`_Composite` vs `_RawGo`) live side by side in each `*_test.go`.
- **Ruby.** Each script runs the operation `REPS = 2000` times under
  `Benchmark.measure`, then divides the real (wall) time by `REPS` and converts
  to ns — yielding the same "ns per one N-sized op" unit as Go, so the columns
  are directly comparable. Ruby timings vary a few percent run to run.
- **Honesty notes.** No warm-up cherry-picking; numbers are a single run. The
  string row is reported against *both* idiomatic (`strings.Builder`) and
  same-algorithm (`+=`) raw Go so the algorithmic vs overhead split is explicit.

## Run it yourself

```sh
export GOWORK=off CGO_ENABLED=0 GOPRIVATE=github.com/go-composites \
       GOPROXY=direct GOSUMDB=off GOFLAGS=-mod=mod
./run.sh
```

or piecemeal:

```sh
go test -bench=. -benchmem ./go/...
ruby ruby/array.rb ; ruby ruby/hash.rb ; ruby ruby/number.rb ; ruby ruby/string.rb
```

## License

BSD-3-Clause — see [LICENSE](LICENSE).
