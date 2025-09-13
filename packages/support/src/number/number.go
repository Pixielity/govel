package number

import (
	"fmt"
	"govel/hashing/src/exceptions"
	"math"
	"strings"

	"github.com/shopspring/decimal"
)

// Number provides enterprise-grade decimal arithmetic operations.
// All operations use precise decimal math to avoid floating-point errors.
type Number struct {
	value decimal.Decimal
}

// NumberError represents errors in number operations
type NumberError struct {
	Op    string
	Value string
	Err   error
}

func (e *NumberError) Error() string {
	return fmt.Sprintf("number %s error for value '%s': %v", e.Op, e.Value, e.Err)
}

func (e *NumberError) Unwrap() error {
	return e.Err
}

// New creates a new Number from a decimal value
func New(value decimal.Decimal) *Number {
	return &Number{value: value}
}

// NewFromString creates a Number from a string representation
func NewFromString(value string) (*Number, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return nil, &NumberError{
			Op:    "parse",
			Value: value,
			Err:   err,
		}
	}
	return &Number{value: d}, nil
}

// NewFromFloat creates a Number from a float64 (use with caution for financial calculations)
func NewFromFloat(value float64) *Number {
	return &Number{value: decimal.NewFromFloat(value)}
}

// NewFromInt creates a Number from an integer
func NewFromInt(value int64) *Number {
	return &Number{value: decimal.NewFromInt(value)}
}

// Value returns the underlying decimal value
func (n *Number) Value() decimal.Decimal {
	return n.value
}

// String returns the string representation of the number
func (n *Number) String() string {
	return n.value.String()
}

// Float64 converts to float64 (precision may be lost)
func (n *Number) Float64() (float64, bool) {
	return n.value.Float64()
}

// Int64 converts to int64 (truncates decimal places)
func (n *Number) Int64() int64 {
	return n.value.IntPart()
}

// Add performs addition with another Number
func (n *Number) Add(other *Number) *Number {
	return &Number{value: n.value.Add(other.value)}
}

// Subtract performs subtraction with another Number
func (n *Number) Subtract(other *Number) *Number {
	return &Number{value: n.value.Sub(other.value)}
}

// Multiply performs multiplication with another Number
func (n *Number) Multiply(other *Number) *Number {
	return &Number{value: n.value.Mul(other.value)}
}

// Divide performs division with another Number
func (n *Number) Divide(other *Number) (*Number, error) {
	if other.value.IsZero() {
		return nil, &NumberError{
			Op:    "divide",
			Value: "0",
			Err:   fmt.Errorf("division by zero"),
		}
	}
	return &Number{value: n.value.Div(other.value)}, nil
}

// Power raises the number to the given power
func (n *Number) Power(exponent *Number) *Number {
	return &Number{value: n.value.Pow(exponent.value)}
}

// Mod returns the remainder of division
func (n *Number) Mod(other *Number) (*Number, error) {
	if other.value.IsZero() {
		return nil, &NumberError{
			Op:    "mod",
			Value: "0",
			Err:   fmt.Errorf("modulo by zero"),
		}
	}
	return &Number{value: n.value.Mod(other.value)}, nil
}

// Abs returns the absolute value
func (n *Number) Abs() *Number {
	return &Number{value: n.value.Abs()}
}

// Negate returns the negated value
func (n *Number) Negate() *Number {
	return &Number{value: n.value.Neg()}
}

// Round rounds to the specified number of decimal places
func (n *Number) Round(decimals int32) *Number {
	return &Number{value: n.value.Round(decimals)}
}

// RoundUp rounds up to the specified number of decimal places
func (n *Number) RoundUp(decimals int32) *Number {
	return &Number{value: n.value.RoundUp(decimals)}
}

// RoundDown rounds down to the specified number of decimal places
func (n *Number) RoundDown(decimals int32) *Number {
	return &Number{value: n.value.RoundDown(decimals)}
}

// Truncate truncates to the specified number of decimal places
func (n *Number) Truncate(decimals int32) *Number {
	return &Number{value: n.value.Truncate(decimals)}
}

// Comparison methods

// Equal checks if two numbers are equal
func (n *Number) Equal(other *Number) bool {
	return n.value.Equal(other.value)
}

// GreaterThan checks if this number is greater than another
func (n *Number) GreaterThan(other *Number) bool {
	return n.value.GreaterThan(other.value)
}

// GreaterThanOrEqual checks if this number is greater than or equal to another
func (n *Number) GreaterThanOrEqual(other *Number) bool {
	return n.value.GreaterThanOrEqual(other.value)
}

// LessThan checks if this number is less than another
func (n *Number) LessThan(other *Number) bool {
	return n.value.LessThan(other.value)
}

// LessThanOrEqual checks if this number is less than or equal to another
func (n *Number) LessThanOrEqual(other *Number) bool {
	return n.value.LessThanOrEqual(other.value)
}

// State checking methods

// IsZero checks if the number is zero
func (n *Number) IsZero() bool {
	return n.value.IsZero()
}

// IsPositive checks if the number is positive
func (n *Number) IsPositive() bool {
	return n.value.IsPositive()
}

// IsNegative checks if the number is negative
func (n *Number) IsNegative() bool {
	return n.value.IsNegative()
}

// IsInteger checks if the number is an integer (no decimal places)
func (n *Number) IsInteger() bool {
	return n.value.Equal(n.value.Truncate(0))
}

// Utility methods

// Min returns the smaller of two numbers
func (n *Number) Min(other *Number) *Number {
	if n.value.LessThan(other.value) {
		return n
	}
	return other
}

// Max returns the larger of two numbers
func (n *Number) Max(other *Number) *Number {
	if n.value.GreaterThan(other.value) {
		return n
	}
	return other
}

// Clamp constrains the number between min and max
func (n *Number) Clamp(min, max *Number) *Number {
	if n.value.LessThan(min.value) {
		return min
	}
	if n.value.GreaterThan(max.value) {
		return max
	}
	return n
}

// Format formats the number with specified decimal places and thousands separator
func (n *Number) Format(decimals int32, thousandsSep string) string {
	rounded := n.value.Round(decimals)
	str := rounded.String()

	// Split integer and decimal parts
	parts := strings.Split(str, ".")
	integerPart := parts[0]

	// Handle negative numbers
	negative := strings.HasPrefix(integerPart, "-")
	if negative {
		integerPart = integerPart[1:]
	}

	// Add thousands separators
	if len(integerPart) > 3 && thousandsSep != "" {
		var result strings.Builder
		for i, r := range integerPart {
			if i > 0 && (len(integerPart)-i)%3 == 0 {
				result.WriteString(thousandsSep)
			}
			result.WriteRune(r)
		}
		integerPart = result.String()
	}

	// Reconstruct the number
	if negative {
		integerPart = "-" + integerPart
	}

	if len(parts) > 1 {
		return integerPart + "." + parts[1]
	}
	return integerPart
}

// Percentage converts the number to a percentage string
func (n *Number) Percentage(decimals int32) string {
	percentage := n.value.Mul(decimal.NewFromInt(100))
	return New(percentage).Format(decimals, "") + "%"
}

// Static utility functions

// Zero returns a Number with value zero
func Zero() *Number {
	return &Number{value: decimal.Zero}
}

// One returns a Number with value one
func One() *Number {
	return &Number{value: decimal.NewFromInt(1)}
}

// Sum calculates the sum of multiple numbers
func Sum(numbers ...*Number) *Number {
	result := Zero()
	for _, num := range numbers {
		result = result.Add(num)
	}
	return result
}

// Average calculates the average of multiple numbers
func Average(numbers ...*Number) (*Number, error) {
	if len(numbers) == 0 {
		return nil, &NumberError{
			Op:    "average",
			Value: "empty slice",
			Err:   fmt.Errorf("cannot calculate average of empty slice"),
		}
	}

	sum := Sum(numbers...)
	count := NewFromInt(int64(len(numbers)))
	return sum.Divide(count)
}

// Parse attempts to parse a string as a Number, handling common formats
func Parse(value string) (*Number, error) {
	// Clean the string
	cleaned := strings.TrimSpace(value)
	cleaned = strings.ReplaceAll(cleaned, ",", "") // Remove thousands separators

	// Try to parse as decimal
	return NewFromString(cleaned)
}

// Type conversion helper methods

// ToUint32 safely converts various numeric types to uint32 with bounds checking.
// This helper method provides type-safe conversion for configuration values that need
// to be uint32, such as Argon2 memory and time parameters.
//
// Parameters:
//   - value: The value to convert (supports int, int64, float64, uint32)
//
// Returns:
//   - uint32: The converted value
//   - error: ErrInvalidOptions if conversion fails or value is out of bounds
//
// Usage:
//
//	memory, err := ToUint32(configValue)
//	if err != nil {
//	    return err
//	}
func ToUint32(value interface{}) (uint32, error) {
	switch v := value.(type) {
	case uint32:
		return v, nil
	case int:
		if v < 0 || v > math.MaxUint32 {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint32(v), nil
	case int64:
		if v < 0 || v > math.MaxUint32 {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint32(v), nil
	case float64:
		if v < 0 || v > math.MaxUint32 || v != float64(int64(v)) {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint32(v), nil
	default:
		return 0, exceptions.ErrInvalidOptions
	}
}

// ToUint8 safely converts various numeric types to uint8 with bounds checking.
// This helper method provides type-safe conversion for configuration values that need
// to be uint8, such as Argon2 thread parameters.
//
// Parameters:
//   - value: The value to convert (supports int, int64, float64, uint8, uint32)
//
// Returns:
//   - uint8: The converted value
//   - error: ErrInvalidOptions if conversion fails or value is out of bounds
//
// Usage:
//
//	threads, err := ToUint8(configValue)
//	if err != nil {
//	    return err
//	}
func ToUint8(value interface{}) (uint8, error) {
	switch v := value.(type) {
	case uint8:
		return v, nil
	case int:
		if v < 0 || v > math.MaxUint8 {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint8(v), nil
	case int64:
		if v < 0 || v > math.MaxUint8 {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint8(v), nil
	case uint32:
		if v > math.MaxUint8 {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint8(v), nil
	case float64:
		if v < 0 || v > math.MaxUint8 || v != float64(int64(v)) {
			return 0, exceptions.ErrInvalidOptions
		}
		return uint8(v), nil
	default:
		return 0, exceptions.ErrInvalidOptions
	}
}

// ToInt safely converts various numeric types to int with bounds checking.
// This helper method provides type-safe conversion for configuration values that need
// to be int, such as bcrypt cost parameters.
//
// Parameters:
//   - value: The value to convert (supports int, int64, float64, uint32)
//
// Returns:
//   - int: The converted value
//   - error: ErrInvalidOptions if conversion fails or value is out of bounds
//
// Usage:
//
//	cost, err := ToInt(configValue)
//	if err != nil {
//	    return err
//	}
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		if v > math.MaxInt || v < math.MinInt {
			return 0, exceptions.ErrInvalidOptions
		}
		return int(v), nil
	case uint32:
		if uint64(v) > math.MaxInt {
			return 0, exceptions.ErrInvalidOptions
		}
		return int(v), nil
	case float64:
		if v > math.MaxInt || v < math.MinInt || v != float64(int64(v)) {
			return 0, exceptions.ErrInvalidOptions
		}
		return int(v), nil
	default:
		return 0, exceptions.ErrInvalidOptions
	}
}
