# frozen_string_literal: true

# Ruby MRI counterpart of go/dictionary_test.go: insert N entries into a
# native Hash, then look every key up.
require_relative "common"

KEYS = Array.new(N) { |i| "key-#{i}" }

bench("hash insert") do
  h = {}
  N.times { |j| h[KEYS[j]] = j }
end

h = {}
N.times { |j| h[KEYS[j]] = j }
bench("hash lookup") do
  sum = 0
  N.times { |j| sum += h[KEYS[j]] }
end
