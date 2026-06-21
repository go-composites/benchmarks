#!/usr/bin/env sh
# Run the whole comparison: Go benchmarks (go-composites vs raw Go) then the
# Ruby MRI counterparts. Numbers are ns/op for one N-sized operation on both
# sides (see ruby/common.rb for the Ruby normalisation methodology).
set -eu

export GOWORK=off CGO_ENABLED=0 GOPRIVATE=github.com/go-composites \
       GOPROXY=direct GOSUMDB=off GOFLAGS=-mod=mod

cd "$(dirname "$0")"

echo "==================== Go (testing.B) ===================="
go test -bench=. -benchmem ./go/...

echo
echo "==================== Ruby MRI (Benchmark) ===================="
ruby --version
for f in array hash number string; do
  ruby "ruby/${f}.rb"
done
