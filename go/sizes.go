// Package benchmarks holds paired micro-benchmarks comparing the
// go-composites composite style against idiomatic raw Go. The same
// operations are timed in Ruby MRI under ../ruby with identical sizes,
// so the three columns in the README are directly comparable.
package benchmarks

// N is the working-set size used by every benchmark (and mirrored by the
// Ruby scripts). Keep this in sync with the `N` constant in ruby/*.rb.
const N = 1000
