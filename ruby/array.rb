# frozen_string_literal: true

# Ruby MRI counterpart of go/array_test.go: push N elements onto a native
# Array, then iterate summing them.
require_relative "common"

bench("array push") do
  a = []
  N.times { |j| a.push(j) }
end

a = []
N.times { |j| a.push(j) }
bench("array iterate") do
  sum = 0
  a.each { |v| sum += v }
end
