# frozen_string_literal: true

# Shared harness for the Ruby MRI benchmarks. Mirrors the Go suite:
#   - N is the working-set size (must match go/sizes.go and N below).
#   - REPS is how many times the whole N-sized operation is repeated when
#     timing, so we measure a meaningful interval and divide back out.
#
# Methodology: Go's `testing.B` reports ns for ONE operation (one full
# N-sized build/iterate). Ruby's Benchmark times a loop; we run the
# operation REPS times, take the real (wall) time, and divide by REPS to
# get the per-operation time, then convert to ns so the columns line up.
require "benchmark"

N = 1000
REPS = 2000

# bench runs `block` REPS times under Benchmark, prints the per-operation
# time in nanoseconds (one operation == one full N-sized pass).
def bench(label)
  tm = Benchmark.measure do
    REPS.times { yield }
  end
  per_op_ns = (tm.real / REPS) * 1_000_000_000.0
  printf("%-28s %12.1f ns/op  (N=%d, reps=%d)\n", label, per_op_ns, N, REPS)
end
