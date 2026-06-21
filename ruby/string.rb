# frozen_string_literal: true

# Ruby MRI counterpart of go/string_test.go: build a string of N chunks.
# String#<< mutates in place (linear, like strings.Builder); "+" rebuilds
# (quadratic, like the composite Concat and the Go naive "+=" variant).
require_relative "common"

CHUNK = "abcdefghij"

bench("string build (<<)") do
  s = +""
  N.times { s << CHUNK }
end

bench("string build (+ naive)") do
  s = +""
  N.times { s += CHUNK }
end
