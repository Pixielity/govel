module govel/support

go 1.25

require (
	github.com/dromara/carbon/v2 v2.6.11
	github.com/gobeam/stringy v0.0.7
	github.com/google/uuid v1.6.0
	github.com/shopspring/decimal v1.4.0
)



replace (
    // Map the main module to the src directory
    govel/support => ./src

    // Support subpackages
    govel/support/carbon => ./src/carbon
    govel/support/env => ./src/env
    govel/support/facades => ./src/facades
    govel/support/money => ./src/money
    govel/support/number => ./src/number
    govel/support/reflector => ./src/reflector
    govel/support/sleep => ./src/sleep
    govel/support/str => ./src/str
    govel/support/stringable => ./src/stringable
    govel/support/symbol => ./src/symbol
    govel/support/traits => ./src/traits
)
