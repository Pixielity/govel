package money

import (
	"fmt"
	"strings"

	"govel/packages/support/src/number"

	"github.com/shopspring/decimal"
)

// Currency represents different currency codes
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	CAD Currency = "CAD"
	AUD Currency = "AUD"
	CHF Currency = "CHF"
	CNY Currency = "CNY"
	SEK Currency = "SEK"
	NZD Currency = "NZD"
)

// CurrencyInfo holds information about a currency
type CurrencyInfo struct {
	Code          Currency
	Symbol        string
	DecimalPlaces int
	Name          string
}

// Common currency information
var currencyInfoMap = map[Currency]CurrencyInfo{
	USD: {USD, "$", 2, "US Dollar"},
	EUR: {EUR, "€", 2, "Euro"},
	GBP: {GBP, "£", 2, "British Pound"},
	JPY: {JPY, "¥", 0, "Japanese Yen"},
	CAD: {CAD, "C$", 2, "Canadian Dollar"},
	AUD: {AUD, "A$", 2, "Australian Dollar"},
	CHF: {CHF, "CHF", 2, "Swiss Franc"},
	CNY: {CNY, "¥", 2, "Chinese Yuan"},
	SEK: {SEK, "kr", 2, "Swedish Krona"},
	NZD: {NZD, "NZ$", 2, "New Zealand Dollar"},
}

// Money represents a monetary amount with currency
// Uses precise decimal arithmetic to avoid floating-point errors
type Money struct {
	amount   decimal.Decimal
	currency Currency
}

// MoneyError represents errors in money operations
type MoneyError struct {
	Op       string
	Amount   string
	Currency Currency
	Err      error
}

func (e *MoneyError) Error() string {
	return fmt.Sprintf("money %s error for %s %s: %v", e.Op, e.Amount, e.Currency, e.Err)
}

func (e *MoneyError) Unwrap() error {
	return e.Err
}

// NewMoney creates a new Money instance
func NewMoney(amount decimal.Decimal, currency Currency) *Money {
	return &Money{
		amount:   amount,
		currency: currency,
	}
}

// NewMoneyFromString creates Money from string amount and currency
func NewMoneyFromString(amount string, currency Currency) (*Money, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, &MoneyError{
			Op:       "parse",
			Amount:   amount,
			Currency: currency,
			Err:      err,
		}
	}
	return &Money{amount: d, currency: currency}, nil
}

// NewMoneyFromFloat creates Money from float64 (use with caution)
func NewMoneyFromFloat(amount float64, currency Currency) *Money {
	return &Money{
		amount:   decimal.NewFromFloat(amount),
		currency: currency,
	}
}

// NewMoneyFromInt creates Money from integer cents/pence/etc
func NewMoneyFromInt(cents int64, currency Currency) *Money {
	info := currencyInfoMap[currency]
	divisor := decimal.NewFromInt(1)

	// Calculate divisor based on decimal places
	for i := 0; i < info.DecimalPlaces; i++ {
		divisor = divisor.Mul(decimal.NewFromInt(10))
	}

	amount := decimal.NewFromInt(cents).Div(divisor)
	return &Money{amount: amount, currency: currency}
}

// Amount returns the decimal amount
func (m *Money) Amount() decimal.Decimal {
	return m.amount
}

// Currency returns the currency
func (m *Money) Currency() Currency {
	return m.currency
}

// String returns the string representation
func (m *Money) String() string {
	info, exists := currencyInfoMap[m.currency]
	if !exists {
		return fmt.Sprintf("%s %s", m.amount.String(), m.currency)
	}

	rounded := m.amount.Round(int32(info.DecimalPlaces))
	return fmt.Sprintf("%s%s", info.Symbol, rounded.String())
}

// StringWithCode returns string representation with currency code
func (m *Money) StringWithCode() string {
	info, exists := currencyInfoMap[m.currency]
	if !exists {
		return fmt.Sprintf("%s %s", m.amount.String(), m.currency)
	}

	rounded := m.amount.Round(int32(info.DecimalPlaces))
	return fmt.Sprintf("%s %s", rounded.String(), m.currency)
}

// Format returns formatted money string with custom formatting
func (m *Money) Format(useSymbol bool, thousandsSep string) string {
	info, exists := currencyInfoMap[m.currency]
	if !exists {
		return fmt.Sprintf("%s %s", m.amount.String(), m.currency)
	}

	rounded := m.amount.Round(int32(info.DecimalPlaces))
	amountStr := rounded.String()

	// Add thousands separators
	if thousandsSep != "" {
		parts := strings.Split(amountStr, ".")
		integerPart := parts[0]

		// Handle negative numbers
		negative := strings.HasPrefix(integerPart, "-")
		if negative {
			integerPart = integerPart[1:]
		}

		// Add thousands separators
		if len(integerPart) > 3 {
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
			amountStr = integerPart + "." + parts[1]
		} else {
			amountStr = integerPart
		}
	}

	if useSymbol {
		return fmt.Sprintf("%s%s", info.Symbol, amountStr)
	}
	return fmt.Sprintf("%s %s", amountStr, m.currency)
}

// Cents returns the amount in smallest currency unit (cents, pence, etc.)
func (m *Money) Cents() int64 {
	info := currencyInfoMap[m.currency]
	multiplier := decimal.NewFromInt(1)

	// Calculate multiplier based on decimal places
	for i := 0; i < info.DecimalPlaces; i++ {
		multiplier = multiplier.Mul(decimal.NewFromInt(10))
	}

	return m.amount.Mul(multiplier).IntPart()
}

// Arithmetic operations

// Add adds another money amount (must be same currency)
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, &MoneyError{
			Op:       "add",
			Amount:   other.amount.String(),
			Currency: other.currency,
			Err:      fmt.Errorf("cannot add different currencies: %s + %s", m.currency, other.currency),
		}
	}

	return &Money{
		amount:   m.amount.Add(other.amount),
		currency: m.currency,
	}, nil
}

// Subtract subtracts another money amount (must be same currency)
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, &MoneyError{
			Op:       "subtract",
			Amount:   other.amount.String(),
			Currency: other.currency,
			Err:      fmt.Errorf("cannot subtract different currencies: %s - %s", m.currency, other.currency),
		}
	}

	return &Money{
		amount:   m.amount.Sub(other.amount),
		currency: m.currency,
	}, nil
}

// Multiply multiplies by a number (returns new Money)
func (m *Money) Multiply(multiplier *number.Number) *Money {
	return &Money{
		amount:   m.amount.Mul(multiplier.Value()),
		currency: m.currency,
	}
}

// MultiplyByDecimal multiplies by a decimal
func (m *Money) MultiplyByDecimal(multiplier decimal.Decimal) *Money {
	return &Money{
		amount:   m.amount.Mul(multiplier),
		currency: m.currency,
	}
}

// MultiplyByFloat multiplies by a float64
func (m *Money) MultiplyByFloat(multiplier float64) *Money {
	return &Money{
		amount:   m.amount.Mul(decimal.NewFromFloat(multiplier)),
		currency: m.currency,
	}
}

// Divide divides by a number (returns new Money)
func (m *Money) Divide(divisor *number.Number) (*Money, error) {
	if divisor.IsZero() {
		return nil, &MoneyError{
			Op:       "divide",
			Amount:   m.amount.String(),
			Currency: m.currency,
			Err:      fmt.Errorf("division by zero"),
		}
	}

	return &Money{
		amount:   m.amount.Div(divisor.Value()),
		currency: m.currency,
	}, nil
}

// DivideByDecimal divides by a decimal
func (m *Money) DivideByDecimal(divisor decimal.Decimal) (*Money, error) {
	if divisor.IsZero() {
		return nil, &MoneyError{
			Op:       "divide",
			Amount:   m.amount.String(),
			Currency: m.currency,
			Err:      fmt.Errorf("division by zero"),
		}
	}

	return &Money{
		amount:   m.amount.Div(divisor),
		currency: m.currency,
	}, nil
}

// Percentage calculates a percentage of the money amount
func (m *Money) Percentage(percent *number.Number) *Money {
	multiplier := percent.Value().Div(decimal.NewFromInt(100))
	return &Money{
		amount:   m.amount.Mul(multiplier),
		currency: m.currency,
	}
}

// Abs returns the absolute value
func (m *Money) Abs() *Money {
	return &Money{
		amount:   m.amount.Abs(),
		currency: m.currency,
	}
}

// Negate returns the negated amount
func (m *Money) Negate() *Money {
	return &Money{
		amount:   m.amount.Neg(),
		currency: m.currency,
	}
}

// Round rounds to the currency's standard decimal places
func (m *Money) Round() *Money {
	info := currencyInfoMap[m.currency]
	return &Money{
		amount:   m.amount.Round(int32(info.DecimalPlaces)),
		currency: m.currency,
	}
}

// RoundToDecimalPlaces rounds to specific decimal places
func (m *Money) RoundToDecimalPlaces(decimals int32) *Money {
	return &Money{
		amount:   m.amount.Round(decimals),
		currency: m.currency,
	}
}

// Comparison methods

// Equal checks if two money amounts are equal (same currency and amount)
func (m *Money) Equal(other *Money) bool {
	return m.currency == other.currency && m.amount.Equal(other.amount)
}

// GreaterThan checks if this money is greater than another (same currency)
func (m *Money) GreaterThan(other *Money) (bool, error) {
	if m.currency != other.currency {
		return false, &MoneyError{
			Op:       "compare",
			Amount:   other.amount.String(),
			Currency: other.currency,
			Err:      fmt.Errorf("cannot compare different currencies: %s vs %s", m.currency, other.currency),
		}
	}
	return m.amount.GreaterThan(other.amount), nil
}

// GreaterThanOrEqual checks if this money is greater than or equal to another
func (m *Money) GreaterThanOrEqual(other *Money) (bool, error) {
	if m.currency != other.currency {
		return false, &MoneyError{
			Op:       "compare",
			Amount:   other.amount.String(),
			Currency: other.currency,
			Err:      fmt.Errorf("cannot compare different currencies: %s vs %s", m.currency, other.currency),
		}
	}
	return m.amount.GreaterThanOrEqual(other.amount), nil
}

// LessThan checks if this money is less than another
func (m *Money) LessThan(other *Money) (bool, error) {
	if m.currency != other.currency {
		return false, &MoneyError{
			Op:       "compare",
			Amount:   other.amount.String(),
			Currency: other.currency,
			Err:      fmt.Errorf("cannot compare different currencies: %s vs %s", m.currency, other.currency),
		}
	}
	return m.amount.LessThan(other.amount), nil
}

// LessThanOrEqual checks if this money is less than or equal to another
func (m *Money) LessThanOrEqual(other *Money) (bool, error) {
	if m.currency != other.currency {
		return false, &MoneyError{
			Op:       "compare",
			Amount:   other.amount.String(),
			Currency: other.currency,
			Err:      fmt.Errorf("cannot compare different currencies: %s vs %s", m.currency, other.currency),
		}
	}
	return m.amount.LessThanOrEqual(other.amount), nil
}

// State checking methods

// IsZero checks if the amount is zero
func (m *Money) IsZero() bool {
	return m.amount.IsZero()
}

// IsPositive checks if the amount is positive
func (m *Money) IsPositive() bool {
	return m.amount.IsPositive()
}

// IsNegative checks if the amount is negative
func (m *Money) IsNegative() bool {
	return m.amount.IsNegative()
}

// Utility methods for collections of money

// Sum calculates the sum of multiple money amounts (must all be same currency)
func Sum(monies ...*Money) (*Money, error) {
	if len(monies) == 0 {
		return nil, &MoneyError{
			Op:       "sum",
			Amount:   "empty slice",
			Currency: "",
			Err:      fmt.Errorf("cannot sum empty slice"),
		}
	}

	currency := monies[0].currency
	result := NewMoney(decimal.Zero, currency)

	for _, money := range monies {
		var err error
		result, err = result.Add(money)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// SplitEvenly divides money evenly among n parts
func (m *Money) SplitEvenly(parts int) ([]*Money, *Money, error) {
	if parts <= 0 {
		return nil, nil, &MoneyError{
			Op:       "split",
			Amount:   m.amount.String(),
			Currency: m.currency,
			Err:      fmt.Errorf("cannot split into %d parts", parts),
		}
	}

	partsDecimal := decimal.NewFromInt(int64(parts))
	perPart := m.amount.Div(partsDecimal)

	info := currencyInfoMap[m.currency]
	roundedPerPart := perPart.Round(int32(info.DecimalPlaces))

	result := make([]*Money, parts)
	totalDistributed := decimal.Zero

	// Create parts with rounded amounts
	for i := 0; i < parts; i++ {
		result[i] = &Money{amount: roundedPerPart, currency: m.currency}
		totalDistributed = totalDistributed.Add(roundedPerPart)
	}

	// Calculate remainder
	remainder := m.amount.Sub(totalDistributed)
	remainderMoney := &Money{amount: remainder, currency: m.currency}

	return result, remainderMoney, nil
}

// Utility functions for creating common money amounts

// Zero returns zero money for the given currency
func ZeroMoney(currency Currency) *Money {
	return &Money{amount: decimal.Zero, currency: currency}
}

// GetCurrencyInfo returns information about a currency
func GetCurrencyInfo(currency Currency) (CurrencyInfo, bool) {
	info, exists := currencyInfoMap[currency]
	return info, exists
}

// ParseMoney attempts to parse a money string (e.g., "$123.45", "123.45 USD")
func ParseMoney(value string, defaultCurrency Currency) (*Money, error) {
	value = strings.TrimSpace(value)

	// Try to detect currency from symbols or codes
	for currency, info := range currencyInfoMap {
		// Check for symbol at the beginning
		if strings.HasPrefix(value, info.Symbol) {
			amountStr := strings.TrimSpace(value[len(info.Symbol):])
			return NewMoneyFromString(amountStr, currency)
		}

		// Check for currency code at the end
		if strings.HasSuffix(value, " "+string(currency)) {
			amountStr := strings.TrimSpace(value[:len(value)-4])
			return NewMoneyFromString(amountStr, currency)
		}
	}

	// If no currency detected, use default
	cleaned := strings.ReplaceAll(value, ",", "") // Remove thousands separators
	return NewMoneyFromString(cleaned, defaultCurrency)
}
