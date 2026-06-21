# frozen_string_literal: true

# Ruby MRI counterpart of go/number_test.go: an N-step integer accumulation
# on native Integer.
require_relative "common"

bench("integer arithmetic") do
  acc = 0
  N.times { |j| acc += j }
end
