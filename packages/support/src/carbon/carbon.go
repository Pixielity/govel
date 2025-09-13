// Package carbon provides a comprehensive Carbon date/time library for Go,
// inspired by PHP's Carbon library with Laravel-style methods and functionality.
//
// This package wraps the powerful dromara/carbon/v2 library and provides familiar
// PHP Carbon method names and patterns, making it easy for Laravel developers to
// work with dates and times in Go applications.
//
// Key Features:
//   - Laravel/PHP Carbon-style API
//   - Full timezone and locale support
//   - Extensive date manipulation methods
//   - Human-readable date differences
//   - Multiple output formats (ISO8601, RFC3339, custom formats)
//   - Testing utilities with fixed time capabilities
//   - Thread-safe operations
//
// Example usage:
//
//	now := carbon.Now()
//	tomorrow := carbon.Tomorrow()
//	diff := now.DiffForHumans(tomorrow)
//	formatted := now.Format("2006-01-02 15:04:05")
package carbon

import (
	"fmt"
	stdtime "time"

	"github.com/dromara/carbon/v2"
)

// Carbon is an alias to the dromara/carbon Carbon struct for convenience.
// This allows us to use Carbon as a type while leveraging all the functionality
// from the underlying dromara/carbon library.
type Carbon = carbon.Carbon

// Constants for common date/time formats (PHP Carbon style)
// These format constants follow Go's time formatting convention using the reference time
// "Mon Jan 2 15:04:05 MST 2006", which is Unix time 1136239445.
const (
	// RFC3339 formats - Internet date/time format (ISO 8601 profile)
	// RFC3339Format represents the standard RFC3339 format without nanoseconds
	RFC3339Format = "2006-01-02T15:04:05Z07:00"
	// RFC3339NanoFormat represents RFC3339 format with nanosecond precision
	RFC3339NanoFormat = "2006-01-02T15:04:05.999999999Z07:00"

	// Common date and time formats for everyday use
	// DateFormat represents date in YYYY-MM-DD format (ISO 8601 date)
	DateFormat = "2006-01-02"
	// TimeFormat represents time in HH:MM:SS format (24-hour)
	TimeFormat = "15:04:05"
	// DateTimeFormat combines date and time in YYYY-MM-DD HH:MM:SS format
	DateTimeFormat = "2006-01-02 15:04:05"
	// ISO8601Format represents ISO 8601 format with 'Z' UTC indicator
	ISO8601Format = "2006-01-02T15:04:05Z"

	// Human readable formats for display purposes
	// DefaultFormat is the standard date-time format for general use
	DefaultFormat = "2006-01-02 15:04:05"
	// CookieFormat follows the HTTP cookie date format specification
	CookieFormat = "Monday, 02-Jan-2006 15:04:05 MST"
	// RSSFormat follows the RSS feed date format (RFC 2822)
	RSSFormat = "Mon, 02 Jan 2006 15:04:05 -0700"
	// W3CFormat follows the W3C date-time format specification
	W3CFormat = "2006-01-02T15:04:05-07:00"

	// Common time units for duration calculations and conversions
	// These constants help with time arithmetic and provide semantic meaning
	SecondsPerMinute = 60       // 1 minute = 60 seconds
	SecondsPerHour   = 3600     // 1 hour = 3,600 seconds
	SecondsPerDay    = 86400    // 1 day = 86,400 seconds
	SecondsPerWeek   = 604800   // 1 week = 604,800 seconds
	SecondsPerMonth  = 2629746  // Average month ≈ 30.44 days (365.25/12 days)
	SecondsPerYear   = 31556952 // Average year = 365.25 days (accounting for leap years)
)

// Global Configuration Functions (Laravel-style)
// These functions provide global configuration for the Carbon library,
// allowing you to set default behaviors that affect all Carbon instances.

// SetTimezone sets the global default timezone for all new Carbon instances.
// This timezone will be used when creating new Carbon instances without
// explicitly specifying a timezone.
//
// Parameters:
//   - timezone: A timezone identifier (e.g., "UTC", "America/New_York", "Asia/Tokyo")
//
// Example:
//
//	carbon.SetTimezone("America/New_York")
//	now := carbon.Now() // Will use America/New_York timezone
func SetTimezone(timezone string) {
	// Delegate to the underlying carbon library's global timezone setting
	carbon.SetTimezone(timezone)
}

// SetLocale sets the global language locale for date formatting and display.
// This affects how dates are displayed in human-readable formats and
// localized string representations.
//
// Parameters:
//   - locale: A locale identifier (e.g., "en", "fr", "es", "de", "zh")
//
// Example:
//
//	carbon.SetLocale("fr")
//	now := carbon.Now()
//	fmt.Println(now.DiffForHumans()) // Will display in French
func SetLocale(locale string) {
	// Delegate to the underlying carbon library's global locale setting
	carbon.SetLocale(locale)
}

// SetTestNow sets a fixed "now" time for testing purposes (Laravel-style).
// This is extremely useful for testing time-dependent code by freezing
// the current time to a specific moment.
//
// Parameters:
//   - c: The Carbon instance to use as the fixed "now" time
//
// Note: After calling this function, all calls to Now() will return the fixed time
// until ClearTestNow() is called.
//
// Example:
//
//	fixedTime := carbon.Parse("2023-01-01 12:00:00")
//	carbon.SetTestNow(fixedTime)
//	now := carbon.Now() // Will always return 2023-01-01 12:00:00
func SetTestNow(c *Carbon) {
	// Set the fixed test time in the underlying carbon library
	carbon.SetTestNow(c)
}

// ClearTestNow clears any previously set fixed test time.
// After calling this function, Now() will return the actual current time again.
//
// This should be called in test cleanup or after testing scenarios that
// required a fixed time.
//
// Example:
//
//	defer carbon.ClearTestNow() // Ensure test time is cleared after test
//	carbon.SetTestNow(carbon.Parse("2023-01-01 12:00:00"))
//	// ... run tests ...
//	// ClearTestNow() is called automatically by defer
func ClearTestNow() {
	// Clear the fixed test time in the underlying carbon library
	carbon.ClearTestNow()
}

// IsTestNow determines if a test time is currently set.
// This function helps you check whether the Carbon library is currently
// using a fixed time for testing purposes.
//
// Returns:
//   - bool: true if a test time is set, false if using actual current time
//
// Example:
//
//	if carbon.IsTestNow() {
//		fmt.Println("Warning: Using fixed test time")
//	}
func IsTestNow() bool {
	// Check if test time is set in the underlying carbon library
	return carbon.IsTestNow()
}

// Creation Functions (Laravel-style static methods)
// These functions provide various ways to create Carbon instances,
// mimicking Laravel's Carbon static methods for familiar usage patterns.

// Now returns a new Carbon instance representing the current date and time.
// This is one of the most commonly used functions for getting the current moment.
//
// Parameters:
//   - timezone: Optional timezone identifier. If not provided, uses the global default timezone.
//
// Returns:
//   - *Carbon: A new Carbon instance representing the current moment
//
// Example:
//
//	now := carbon.Now()                    // Current time in default timezone
//	nyTime := carbon.Now("America/New_York") // Current time in New York timezone
func Now(timezone ...string) *Carbon {
	// Delegate to the underlying carbon library's Now function
	// This will respect any test time that might be set via SetTestNow
	return carbon.Now(timezone...)
}

// Today returns a new Carbon instance for today at 00:00:00 (midnight).
// This is useful when you need to work with date boundaries or compare dates
// without considering the time component.
//
// Parameters:
//   - timezone: Optional timezone identifier. If not provided, uses the global default timezone.
//
// Returns:
//   - *Carbon: A new Carbon instance representing today at midnight
//
// Example:
//
//	today := carbon.Today()                    // Today at 00:00:00 in default timezone
//	todayUTC := carbon.Today("UTC")            // Today at 00:00:00 UTC
func Today(timezone ...string) *Carbon {
	// Get current time and set it to the start of the day (00:00:00)
	return carbon.Now(timezone...).StartOfDay()
}

// Tomorrow returns a new Carbon instance for tomorrow at 00:00:00 (midnight).
// Useful for scheduling, date comparisons, and working with future date boundaries.
//
// Parameters:
//   - timezone: Optional timezone identifier. If not provided, uses the global default timezone.
//
// Returns:
//   - *Carbon: A new Carbon instance representing tomorrow at midnight
//
// Example:
//
//	tomorrow := carbon.Tomorrow()               // Tomorrow at 00:00:00
//	tomorrowLA := carbon.Tomorrow("America/Los_Angeles") // Tomorrow at 00:00:00 LA time
func Tomorrow(timezone ...string) *Carbon {
	// Get current time, add one day, and set to start of day
	return carbon.Now(timezone...).AddDay().StartOfDay()
}

// Yesterday returns a new Carbon instance for yesterday at 00:00:00 (midnight).
// Useful for working with past date boundaries and historical data processing.
//
// Parameters:
//   - timezone: Optional timezone identifier. If not provided, uses the global default timezone.
//
// Returns:
//   - *Carbon: A new Carbon instance representing yesterday at midnight
//
// Example:
//
//	yesterday := carbon.Yesterday()             // Yesterday at 00:00:00
//	yesterdayTokyo := carbon.Yesterday("Asia/Tokyo") // Yesterday at 00:00:00 Tokyo time
func Yesterday(timezone ...string) *Carbon {
	// Get current time, subtract one day, and set to start of day
	return carbon.Now(timezone...).SubDay().StartOfDay()
}

// Parse parses a date string into a Carbon instance using automatic format detection.
// This function attempts to intelligently parse various common date/time string formats.
// It's very flexible but for better performance and reliability, consider using ParseByFormat
// when you know the exact format of your input strings.
//
// Parameters:
//   - value: The date/time string to parse (e.g., "2023-01-15", "2023-01-15 14:30:00", "Jan 15, 2023")
//   - timezone: Optional timezone identifier for the parsed time
//
// Returns:
//   - *Carbon: A new Carbon instance representing the parsed date/time, or zero value if parsing fails
//
// Example:
//
//	date1 := carbon.Parse("2023-01-15")              // Parses ISO date
//	date2 := carbon.Parse("2023-01-15 14:30:00")     // Parses datetime
//	date3 := carbon.Parse("Jan 15, 2023", "UTC")     // Parses with timezone
func Parse(value string, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's intelligent parsing
	return carbon.Parse(value, timezone...)
}

// ParseByFormat parses a date string using one or more specific formats.
// This function provides more control and better performance than Parse() when you know
// the exact format of your input strings. It supports both single format and multiple
// format fallbacks using Go generics.
//
// Parameters:
//   - value: The date/time string to parse
//   - format: Either a single format string or slice of format strings to try
//   - timezone: Optional timezone identifier for the parsed time
//
// Returns:
//   - *Carbon: A new Carbon instance representing the parsed date/time
//
// Format uses Go's reference time: "Mon Jan 2 15:04:05 MST 2006"
//
// Example:
//
//	// Single format
//	date1 := carbon.ParseByFormat("15/01/2023", "02/01/2006")
//	
//	// Multiple formats (tries each until one succeeds)
//	formats := []string{"2006-01-02", "02/01/2006", "Jan 2, 2006"}
//	date2 := carbon.ParseByFormat("Jan 15, 2023", formats)
func ParseByFormat[T string | []string](value string, format T, timezone ...string) *Carbon {
	// Use type assertion to handle both single format and multiple formats
	switch v := any(format).(type) {
	case string:
		// Single format - use ParseByFormat
		return carbon.ParseByFormat(value, v, timezone...)
	case []string:
		// Multiple formats - use ParseByFormats (tries each format until success)
		return carbon.ParseByFormats(value, v, timezone...)
	}
	// Fallback to general parsing if type assertion fails
	return carbon.Parse(value, timezone...)
}

// CreateFromFormat is an alias for ParseByFormat to provide PHP Carbon compatibility.
// This function maintains the same parameter order as PHP's Carbon::createFromFormat().
// Note: The parameter order is reversed compared to ParseByFormat for PHP compatibility.
//
// Parameters:
//   - format: The format string to use for parsing
//   - value: The date/time string to parse
//   - timezone: Optional timezone identifier for the parsed time
//
// Returns:
//   - *Carbon: A new Carbon instance representing the parsed date/time
//
// Example:
//
//	// PHP-style usage
//	date := carbon.CreateFromFormat("d/m/Y", "15/01/2023")
func CreateFromFormat(format, value string, timezone ...string) *Carbon {
	// Swap parameter order to match PHP Carbon's createFromFormat method
	return ParseByFormat(value, format, timezone...)
}

// CreateFromDate creates a Carbon instance from individual date components.
// The time will be set to 00:00:00 (midnight). This is useful when you want to
// work with dates without caring about the specific time.
//
// Parameters:
//   - year: The year (e.g., 2023)
//   - month: The month (1-12)
//   - day: The day of the month (1-31)
//   - timezone: Optional timezone identifier
//
// Returns:
//   - *Carbon: A new Carbon instance representing the specified date at midnight
//
// Example:
//
//	date := carbon.CreateFromDate(2023, 1, 15)        // 2023-01-15 00:00:00
//	dateUTC := carbon.CreateFromDate(2023, 1, 15, "UTC") // 2023-01-15 00:00:00 UTC
func CreateFromDate(year, month, day int, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's date creation function
	return carbon.CreateFromDate(year, month, day, timezone...)
}

// CreateFromTime creates a Carbon instance from individual time components.
// The date will be set to today's date. This is useful when you need to work
// with specific times while using the current date.
//
// Parameters:
//   - hour: The hour (0-23)
//   - minute: The minute (0-59)
//   - second: The second (0-59)
//   - timezone: Optional timezone identifier
//
// Returns:
//   - *Carbon: A new Carbon instance representing today with the specified time
//
// Example:
//
//	time := carbon.CreateFromTime(14, 30, 45)         // Today at 14:30:45
//	timeUTC := carbon.CreateFromTime(14, 30, 45, "UTC") // Today at 14:30:45 UTC
func CreateFromTime(hour, minute, second int, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's time creation function
	return carbon.CreateFromTime(hour, minute, second, timezone...)
}

// CreateFromDateTime creates a Carbon instance from individual date and time components.
// This provides the most granular control over Carbon instance creation, allowing you
// to specify every component of the date and time.
//
// Parameters:
//   - year: The year (e.g., 2023)
//   - month: The month (1-12)
//   - day: The day of the month (1-31)
//   - hour: The hour (0-23)
//   - minute: The minute (0-59)
//   - second: The second (0-59)
//   - timezone: Optional timezone identifier
//
// Returns:
//   - *Carbon: A new Carbon instance representing the specified date and time
//
// Example:
//
//	dateTime := carbon.CreateFromDateTime(2023, 1, 15, 14, 30, 45)
//	// Creates: 2023-01-15 14:30:45
//	
//	dateTimeLA := carbon.CreateFromDateTime(2023, 1, 15, 14, 30, 45, "America/Los_Angeles")
//	// Creates: 2023-01-15 14:30:45 in Los Angeles timezone
func CreateFromDateTime(year, month, day, hour, minute, second int, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's datetime creation function
	return carbon.CreateFromDateTime(year, month, day, hour, minute, second, timezone...)
}

// CreateFromTimestamp creates a Carbon instance from a Unix timestamp.
// Unix timestamps represent the number of seconds that have elapsed since
// January 1, 1970, 00:00:00 UTC (the Unix epoch).
//
// Parameters:
//   - timestamp: Unix timestamp in seconds since epoch
//   - timezone: Optional timezone identifier to apply to the created instance
//
// Returns:
//   - *Carbon: A new Carbon instance representing the timestamp
//
// Example:
//
//	// Create from current Unix timestamp
//	timestamp := time.Now().Unix()
//	date := carbon.CreateFromTimestamp(timestamp)
//	
//	// Create from specific timestamp with timezone
//	date2 := carbon.CreateFromTimestamp(1640995200, "America/New_York") // 2022-01-01 00:00:00 EST
func CreateFromTimestamp(timestamp int64, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's timestamp creation function
	return carbon.CreateFromTimestamp(timestamp, timezone...)
}

// CreateFromTimestampMilli creates a Carbon instance from a millisecond timestamp.
// Millisecond timestamps are commonly used in JavaScript and web APIs.
// They represent milliseconds since the Unix epoch.
//
// Parameters:
//   - timestamp: Timestamp in milliseconds since epoch
//   - timezone: Optional timezone identifier to apply to the created instance
//
// Returns:
//   - *Carbon: A new Carbon instance representing the timestamp
//
// Example:
//
//	// JavaScript-style timestamp (common in APIs)
//	date := carbon.CreateFromTimestampMilli(1640995200000) // 2022-01-01 00:00:00
//	
//	// From Go's time.Time UnixMilli()
//	milliTimestamp := time.Now().UnixMilli()
//	date2 := carbon.CreateFromTimestampMilli(milliTimestamp, "UTC")
func CreateFromTimestampMilli(timestamp int64, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's millisecond timestamp creation function
	return carbon.CreateFromTimestampMilli(timestamp, timezone...)
}

// CreateFromTimestampMicro creates a Carbon instance from a microsecond timestamp.
// Microsecond timestamps provide higher precision than milliseconds and are used
// in high-precision timing applications and some database systems.
//
// Parameters:
//   - timestamp: Timestamp in microseconds since epoch
//   - timezone: Optional timezone identifier to apply to the created instance
//
// Returns:
//   - *Carbon: A new Carbon instance representing the timestamp
//
// Example:
//
//	// High-precision timestamp
//	microTimestamp := time.Now().UnixMicro()
//	date := carbon.CreateFromTimestampMicro(microTimestamp)
func CreateFromTimestampMicro(timestamp int64, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's microsecond timestamp creation function
	return carbon.CreateFromTimestampMicro(timestamp, timezone...)
}

// CreateFromTimestampNano creates a Carbon instance from a nanosecond timestamp.
// Nanosecond timestamps provide the highest precision and are used in Go's
// internal time representation and high-frequency trading systems.
//
// Parameters:
//   - timestamp: Timestamp in nanoseconds since epoch
//   - timezone: Optional timezone identifier to apply to the created instance
//
// Returns:
//   - *Carbon: A new Carbon instance representing the timestamp
//
// Example:
//
//	// Highest precision timestamp (Go's default internal representation)
//	nanoTimestamp := time.Now().UnixNano()
//	date := carbon.CreateFromTimestampNano(nanoTimestamp, "UTC")
func CreateFromTimestampNano(timestamp int64, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's nanosecond timestamp creation function
	return carbon.CreateFromTimestampNano(timestamp, timezone...)
}

// CreateFromStdTime creates a Carbon instance from Go's standard time.Time.
// This function provides seamless interoperability between Go's native time
// package and the Carbon library, allowing easy conversion of existing time.Time
// values to Carbon instances.
//
// Parameters:
//   - time: A Go time.Time instance to convert
//   - timezone: Optional timezone identifier to convert the time to
//
// Returns:
//   - *Carbon: A new Carbon instance representing the same moment as the input time.Time
//
// Example:
//
//	// Convert from Go's time.Time
//	stdTime := time.Now()
//	carbonTime := carbon.CreateFromStdTime(stdTime)
//	
//	// Convert with timezone change
//	utcTime := time.Now().UTC()
//	nyTime := carbon.CreateFromStdTime(utcTime, "America/New_York")
func CreateFromStdTime(time stdtime.Time, timezone ...string) *Carbon {
	// Delegate to underlying carbon library's standard time conversion function
	return carbon.CreateFromStdTime(time, timezone...)
}

// Helper Functions
// These utility functions provide convenient ways to create special Carbon instances
// and maintain consistency with the underlying carbon library's API.

// NewCarbon creates a new Carbon instance, optionally from a Go time.Time.
// This is an alias function that provides consistency with the underlying carbon library.
// If no time is provided, it creates a zero-value Carbon instance.
//
// Parameters:
//   - time: Optional Go time.Time instances to convert to Carbon
//
// Returns:
//   - *Carbon: A new Carbon instance
//
// Example:
//
//	// Create zero-value Carbon
//	c1 := carbon.NewCarbon()
//	
//	// Create from existing time.Time
//	stdTime := time.Now()
//	c2 := carbon.NewCarbon(stdTime)
func NewCarbon(time ...stdtime.Time) *Carbon {
	// Delegate to underlying carbon library's NewCarbon function
	return carbon.NewCarbon(time...)
}

// ZeroValue returns a zero-value Carbon instance.
// A zero-value Carbon represents January 1, year 1, 00:00:00.000000000 UTC.
// This is useful for initializations, comparisons, and representing "no time" states.
//
// Returns:
//   - *Carbon: A Carbon instance representing the zero time
//
// Example:
//
//	zero := carbon.ZeroValue()
//	if someDate.Eq(zero) {
//		fmt.Println("Date is not set")
//	}
func ZeroValue() *Carbon {
	// Return the zero-value Carbon instance from underlying library
	return carbon.ZeroValue()
}

// EpochValue returns a Carbon instance representing the Unix epoch.
// The Unix epoch is January 1, 1970, 00:00:00 UTC, which is the reference point
// for Unix timestamps. This is useful for timestamp calculations and comparisons.
//
// Returns:
//   - *Carbon: A Carbon instance representing the Unix epoch (1970-01-01 00:00:00 UTC)
//
// Example:
//
//	epoch := carbon.EpochValue()
//	fmt.Println(epoch.Format("2006-01-02 15:04:05")) // Output: 1970-01-01 00:00:00
//	
//	// Calculate seconds since epoch
//	now := carbon.Now()
//	secondsSinceEpoch := now.DiffInSeconds(epoch)
func EpochValue() *Carbon {
	// Return the Unix epoch Carbon instance from underlying library
	return carbon.EpochValue()
}

// Extension Methods (Laravel-style convenience methods)
// The CarbonHelper type provides additional utility methods that extend Carbon's
// functionality with Laravel-style convenience methods and common date/time operations.

// CarbonHelper provides additional utility methods for Carbon instances.
// This helper struct contains methods that provide Laravel-style convenience functions,
// date/time validations, formatting helpers, and common operations that aren't
// directly available on Carbon instances. It acts as a utility belt for Carbon operations.
//
// The helper pattern allows for:
//   - Extending Carbon functionality without modifying the core Carbon struct
//   - Grouping related utility functions logically
//   - Providing Laravel-familiar method names and behaviors
//   - Maintaining separation between core Carbon functionality and convenience methods
//
// Usage:
//   helper := carbon.NewHelper()
//   if helper.IsWeekend(someDate) {
//       // Handle weekend logic
//   }
type CarbonHelper struct{
	// This struct intentionally has no fields as it only provides methods
	// that operate on Carbon instances passed as parameters
}

// NewHelper returns a new CarbonHelper instance.
// The CarbonHelper provides utility methods for working with Carbon instances,
// offering Laravel-style convenience functions and additional date/time operations.
//
// Returns:
//   - *CarbonHelper: A new helper instance ready to use
//
// Example:
//
//	helper := carbon.NewHelper()
//	now := carbon.Now()
//	
//	if helper.IsWeekend(now) {
//		fmt.Println("It's the weekend!")
//	}
//	
//	age := helper.Age(birthDate)
//	fmt.Printf("Age: %d years\n", age)
func NewHelper() *CarbonHelper {
	// Return a new instance of the helper struct
	return &CarbonHelper{}
}

// IsWeekend checks if the given Carbon instance is a weekend
func (h *CarbonHelper) IsWeekend(c *Carbon) bool {
	return c.IsSaturday() || c.IsSunday()
}

// IsWeekday checks if the given Carbon instance is a weekday
func (h *CarbonHelper) IsWeekday(c *Carbon) bool {
	return !h.IsWeekend(c)
}

// IsPast checks if the given Carbon instance is in the past
func (h *CarbonHelper) IsPast(c *Carbon) bool {
	return c.Lt(Now(c.Timezone()))
}

// IsFuture checks if the given Carbon instance is in the future
func (h *CarbonHelper) IsFuture(c *Carbon) bool {
	return c.Gt(Now(c.Timezone()))
}

// IsToday checks if the given Carbon instance is today
func (h *CarbonHelper) IsToday(c *Carbon) bool {
	now := Now(c.Timezone())
	return c.ToDateString() == now.ToDateString()
}

// IsYesterday checks if the given Carbon instance is yesterday
func (h *CarbonHelper) IsYesterday(c *Carbon) bool {
	yesterday := Now(c.Timezone()).SubDay()
	return c.ToDateString() == yesterday.ToDateString()
}

// IsTomorrow checks if the given Carbon instance is tomorrow
func (h *CarbonHelper) IsTomorrow(c *Carbon) bool {
	tomorrow := Now(c.Timezone()).AddDay()
	return c.ToDateString() == tomorrow.ToDateString()
}

// DiffForHumans returns a human-readable difference string
func (h *CarbonHelper) DiffForHumans(c *Carbon, other ...*Carbon) string {
	if len(other) > 0 {
		return c.DiffForHumans(other[0])
	}
	return c.DiffForHumans()
}

// Age returns the age in years
func (h *CarbonHelper) Age(c *Carbon) int64 {
	return int64(c.Age())
}

// Format Helpers (PHP Carbon style)

// FormatLocalized returns a localized formatted string
func (h *CarbonHelper) FormatLocalized(c *Carbon, format string) string {
	// This would require locale-specific formatting implementation
	// For now, return the standard format
	return c.Format(format)
}

// ToDateString returns date in Y-m-d format
func (h *CarbonHelper) ToDateString(c *Carbon) string {
	return c.ToDateString()
}

// ToTimeString returns time in H:i:s format
func (h *CarbonHelper) ToTimeString(c *Carbon) string {
	return c.ToTimeString()
}

// ToDateTimeString returns datetime in Y-m-d H:i:s format
func (h *CarbonHelper) ToDateTimeString(c *Carbon) string {
	return c.ToDateTimeString()
}

// ToFormattedDateString returns date in readable format (e.g., "Jan 1, 2023")
func (h *CarbonHelper) ToFormattedDateString(c *Carbon) string {
	return c.ToFormattedDateString()
}

// ToISOString returns ISO8601 formatted string
func (h *CarbonHelper) ToISOString(c *Carbon) string {
	return c.ToIso8601String()
}

// ToJSON returns JSON-formatted string
func (h *CarbonHelper) ToJSON(c *Carbon) string {
	return c.String()
}

// ToRFC3339String returns RFC3339 formatted string
func (h *CarbonHelper) ToRFC3339String(c *Carbon) string {
	return c.ToRfc3339String()
}

// Manipulation Helpers

// StartOf returns the start of the given unit
func (h *CarbonHelper) StartOf(c *Carbon, unit string) *Carbon {
	switch unit {
	case "year":
		return c.StartOfYear()
	case "quarter":
		return c.StartOfQuarter()
	case "month":
		return c.StartOfMonth()
	case "week":
		return c.StartOfWeek()
	case "day":
		return c.StartOfDay()
	case "hour":
		return c.StartOfHour()
	case "minute":
		return c.StartOfMinute()
	case "second":
		return c.StartOfSecond()
	default:
		return c
	}
}

// EndOf returns the end of the given unit
func (h *CarbonHelper) EndOf(c *Carbon, unit string) *Carbon {
	switch unit {
	case "year":
		return c.EndOfYear()
	case "quarter":
		return c.EndOfQuarter()
	case "month":
		return c.EndOfMonth()
	case "week":
		return c.EndOfWeek()
	case "day":
		return c.EndOfDay()
	case "hour":
		return c.EndOfHour()
	case "minute":
		return c.EndOfMinute()
	case "second":
		return c.EndOfSecond()
	default:
		return c
	}
}

// Add adds time to the Carbon instance
func (h *CarbonHelper) Add(c *Carbon, value int64, unit string) *Carbon {
	switch unit {
	case "years", "year":
		return c.AddYears(int(value))
	case "quarters", "quarter":
		return c.AddQuarters(int(value))
	case "months", "month":
		return c.AddMonths(int(value))
	case "weeks", "week":
		return c.AddWeeks(int(value))
	case "days", "day":
		return c.AddDays(int(value))
	case "hours", "hour":
		return c.AddHours(int(value))
	case "minutes", "minute":
		return c.AddMinutes(int(value))
	case "seconds", "second":
		return c.AddSeconds(int(value))
	default:
		return c
	}
}

// Sub subtracts time from the Carbon instance
func (h *CarbonHelper) Sub(c *Carbon, value int64, unit string) *Carbon {
	return h.Add(c, -value, unit)
}

// Comparison Helpers

// IsSame checks if two Carbon instances are the same for given unit
func (h *CarbonHelper) IsSame(c1, c2 *Carbon, unit string) bool {
	switch unit {
	case "year":
		return c1.IsSameYear(c2)
	case "quarter":
		return c1.IsSameQuarter(c2)
	case "month":
		return c1.IsSameMonth(c2)
	case "week":
		// Check if both dates are in the same week
		return c1.StartOfWeek().ToDateString() == c2.StartOfWeek().ToDateString()
	case "day":
		return c1.IsSameDay(c2)
	case "hour":
		return c1.IsSameHour(c2)
	case "minute":
		return c1.IsSameMinute(c2)
	case "second":
		return c1.IsSameSecond(c2)
	default:
		return c1.Eq(c2)
	}
}

// Specialized Creation Functions

// CreateMidnightDate creates a Carbon instance for midnight of given date
func CreateMidnightDate(year, month, day int, timezone ...string) *Carbon {
	return CreateFromDateTime(year, month, day, 0, 0, 0, timezone...)
}

// CreateEndOfDay creates a Carbon instance for end of given date (23:59:59)
func CreateEndOfDay(year, month, day int, timezone ...string) *Carbon {
	return CreateFromDateTime(year, month, day, 23, 59, 59, timezone...)
}

// CreateNoon creates a Carbon instance for noon of given date (12:00:00)
func CreateNoon(year, month, day int, timezone ...string) *Carbon {
	return CreateFromDateTime(year, month, day, 12, 0, 0, timezone...)
}

// Range and Iteration Functions

// CarbonPeriod represents a period between two Carbon instances
type CarbonPeriod struct {
	start    *Carbon
	end      *Carbon
	interval *Carbon
	exclude  []*Carbon
}

// NewPeriod creates a new CarbonPeriod
func NewPeriod(start, end *Carbon, interval ...string) *CarbonPeriod {
	period := &CarbonPeriod{
		start: start,
		end:   end,
	}

	if len(interval) > 0 {
		// Parse interval (simplified implementation)
		switch interval[0] {
		case "1 day", "day":
			period.interval = CreateFromTime(24, 0, 0)
		case "1 hour", "hour":
			period.interval = CreateFromTime(1, 0, 0)
		default:
			period.interval = CreateFromTime(24, 0, 0) // Default to 1 day
		}
	} else {
		period.interval = CreateFromTime(24, 0, 0) // Default to 1 day
	}

	return period
}

// ToSlice returns all dates in the period as a slice
func (p *CarbonPeriod) ToSlice() []*Carbon {
	var dates []*Carbon
	current := p.start.Copy()

	for current.Lte(p.end) {
		dates = append(dates, current.Copy())
		current = current.AddDay() // Simplified - should use interval
	}

	return dates
}

// Factory Functions for Common Use Cases

// StartOfWeek returns the start of the current week
func StartOfWeek(timezone ...string) *Carbon {
	return Now(timezone...).StartOfWeek()
}

// EndOfWeek returns the end of the current week
func EndOfWeek(timezone ...string) *Carbon {
	return Now(timezone...).EndOfWeek()
}

// StartOfMonth returns the start of the current month
func StartOfMonth(timezone ...string) *Carbon {
	return Now(timezone...).StartOfMonth()
}

// EndOfMonth returns the end of the current month
func EndOfMonth(timezone ...string) *Carbon {
	return Now(timezone...).EndOfMonth()
}

// StartOfYear returns the start of the current year
func StartOfYear(timezone ...string) *Carbon {
	return Now(timezone...).StartOfYear()
}

// EndOfYear returns the end of the current year
func EndOfYear(timezone ...string) *Carbon {
	return Now(timezone...).EndOfYear()
}

// Laravel-style Helper Functions

// Max returns the maximum Carbon instance from the given instances
func Max(dates ...*Carbon) *Carbon {
	if len(dates) == 0 {
		return nil
	}

	max := dates[0]
	for _, date := range dates[1:] {
		if date.Gt(max) {
			max = date
		}
	}
	return max
}

// Min returns the minimum Carbon instance from the given instances
func Min(dates ...*Carbon) *Carbon {
	if len(dates) == 0 {
		return nil
	}

	min := dates[0]
	for _, date := range dates[1:] {
		if date.Lt(min) {
			min = date
		}
	}
	return min
}

// Create is an alias for CreateFromDateTime for Laravel compatibility
func Create(year, month, day int, hour, minute, second int, timezone ...string) *Carbon {
	return CreateFromDateTime(year, month, day, hour, minute, second, timezone...)
}

// HasFormat checks if a string can be parsed with the given format
func HasFormat(value, format string) bool {
	_, err := stdtime.Parse(format, value)
	return err == nil
}

// Utility Functions

// GetWeekStartsAt returns the day of the week that weeks start at
func GetWeekStartsAt() int {
	return int(stdtime.Sunday) // Can be configured
}

// GetWeekEndsAt returns the day of the week that weeks end at
func GetWeekEndsAt() int {
	return int(stdtime.Saturday) // Can be configured
}

// IsValid checks if the given values form a valid date
func IsValid(year, month, day int) bool {
	date := stdtime.Date(year, stdtime.Month(month), day, 0, 0, 0, 0, stdtime.UTC)
	return date.Year() == year && int(date.Month()) == month && date.Day() == day
}

// DaysInMonth returns the number of days in the given month/year
func DaysInMonth(year, month int) int {
	firstOfMonth := stdtime.Date(year, stdtime.Month(month), 1, 0, 0, 0, 0, stdtime.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth.Day()
}

// IsLeapYear checks if the given year is a leap year
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// CarbonExtensions provides additional methods that can be added to Carbon instances
type CarbonExtensions struct{}

// NewExtensions returns a new CarbonExtensions instance
func NewExtensions() *CarbonExtensions {
	return &CarbonExtensions{}
}

// ToStdTime converts a Carbon instance to Go's time.Time
func (ext *CarbonExtensions) ToStdTime(c *Carbon) stdtime.Time {
	return c.StdTime()
}

// Clone creates a copy of the Carbon instance
func (ext *CarbonExtensions) Clone(c *Carbon) *Carbon {
	return c.Copy()
}

// SetDateFrom sets the date from another Carbon instance
func (ext *CarbonExtensions) SetDateFrom(c *Carbon, date *Carbon) *Carbon {
	return c.SetDate(date.Year(), date.Month(), date.Day())
}

// SetTimeFrom sets the time from another Carbon instance
func (ext *CarbonExtensions) SetTimeFrom(c *Carbon, time *Carbon) *Carbon {
	return c.SetTime(time.Hour(), time.Minute(), time.Second())
}

// Error types for Carbon operations
type CarbonError struct {
	Operation string
	Message   string
	Cause     error
}

func (e *CarbonError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("carbon %s error: %s (cause: %v)", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("carbon %s error: %s", e.Operation, e.Message)
}

// NewCarbonError creates a new CarbonError
func NewCarbonError(operation, message string, cause error) *CarbonError {
	return &CarbonError{
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}
