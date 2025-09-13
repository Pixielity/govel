package stringable

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	interfaces "govel/types/src/interfaces/support"
	"govel/support/src/str"
)

// StringableImpl implements the Stringable interface with all Laravel-style string manipulation methods.
// This implementation provides a fluent API for string operations, JSON serialization,
// and array-like access to string characters.
type StringableImpl struct {
	value string
}

// NewStringable creates a new Stringable instance with the given string value.
//
// Parameters:
//
//	value: The initial string value
//
// Returns:
//
//	interfaces.Stringable: A new Stringable instance
//
// Example:
//
//	str := NewStringable("Hello World")
//	result := str.Lower().Snake()
//	fmt.Println(result.ToString()) // "hello_world"
func NewStringable(value string) interfaces.Stringable {
	return &StringableImpl{value: value}
}

// JsonSerializable implementation

// JsonSerialize returns the string value for JSON serialization.
// This allows the Stringable to be directly serialized to JSON as a string value.
//
// Returns:
//
//	interface{}: The underlying string value
func (s *StringableImpl) JsonSerialize() interface{} {
	return s.value
}

// ArrayAccess implementation

// OffsetExists checks if a character exists at the given index.
// Supports both positive and negative indices (negative counts from end).
//
// Parameters:
//
//	offset: The index to check (should be an int)
//
// Returns:
//
//	bool: True if the index is valid, false otherwise
func (s *StringableImpl) OffsetExists(offset interface{}) bool {
	index, ok := offset.(int)
	if !ok {
		return false
	}

	length := utf8.RuneCountInString(s.value)
	if index < 0 {
		index = length + index
	}

	return index >= 0 && index < length
}

// OffsetGet retrieves the character at the given index.
// Supports both positive and negative indices.
//
// Parameters:
//
//	offset: The index to retrieve (should be an int)
//
// Returns:
//
//	interface{}: The character at the index, or nil if invalid
func (s *StringableImpl) OffsetGet(offset interface{}) interface{} {
	index, ok := offset.(int)
	if !ok {
		return nil
	}

	runes := []rune(s.value)
	length := len(runes)

	if index < 0 {
		index = length + index
	}

	if index >= 0 && index < length {
		return string(runes[index])
	}

	return nil
}

// OffsetSet is not supported for strings as they are immutable.
// This method panics if called to maintain string immutability.
//
// Parameters:
//
//	offset: The index (ignored)
//	value: The value (ignored)
func (s *StringableImpl) OffsetSet(offset, value interface{}) {
	panic("Cannot modify string characters directly. Strings are immutable. Use string manipulation methods instead.")
}

// OffsetUnset is not supported for strings as they are immutable.
// This method panics if called to maintain string immutability.
//
// Parameters:
//
//	offset: The index (ignored)
func (s *StringableImpl) OffsetUnset(offset interface{}) {
	panic("Cannot remove string characters directly. Strings are immutable. Use string manipulation methods instead.")
}

// String manipulation methods

// After returns the remainder of a string after the first occurrence of a given value.
func (s *StringableImpl) After(search string) interfaces.Stringable {
	return NewStringable(str.After(s.value, search))
}

// AfterLast returns the remainder of a string after the last occurrence of a given value.
func (s *StringableImpl) AfterLast(search string) interfaces.Stringable {
	return NewStringable(str.AfterLast(s.value, search))
}

// Append appends the given values to the string.
func (s *StringableImpl) Append(values ...string) interfaces.Stringable {
	return NewStringable(s.value + strings.Join(values, ""))
}

// NewLine appends one or more new lines to the string.
func (s *StringableImpl) NewLine(count int) interfaces.Stringable {
	return s.Append(strings.Repeat("\n", count))
}

// Ascii transliterates a UTF-8 value to ASCII.
func (s *StringableImpl) Ascii(language string) interfaces.Stringable {
	return NewStringable(str.Ascii(s.value))
}

// Basename gets the trailing name component of the path.
func (s *StringableImpl) Basename(suffix string) interfaces.Stringable {
	return NewStringable(str.Basename(s.value, suffix))
}

// Before gets the portion of a string before the first occurrence of a given value.
func (s *StringableImpl) Before(search string) interfaces.Stringable {
	return NewStringable(str.Before(s.value, search))
}

// BeforeLast gets the portion of a string before the last occurrence of a given value.
func (s *StringableImpl) BeforeLast(search string) interfaces.Stringable {
	return NewStringable(str.BeforeLast(s.value, search))
}

// Between gets the portion of a string between two given values.
func (s *StringableImpl) Between(from, to string) interfaces.Stringable {
	return NewStringable(str.Between(s.value, from, to))
}

// BetweenFirst gets the smallest possible portion of a string between two given values.
func (s *StringableImpl) BetweenFirst(from, to string) interfaces.Stringable {
	return NewStringable(str.BetweenFirst(s.value, from, to))
}

// Case conversion methods

// Camel converts a value to camel case.
func (s *StringableImpl) Camel() interfaces.Stringable {
	return NewStringable(str.Camel(s.value))
}

// Kebab converts a string to kebab case.
func (s *StringableImpl) Kebab() interfaces.Stringable {
	return NewStringable(str.Kebab(s.value))
}

// Lower converts the given string to lower-case.
func (s *StringableImpl) Lower() interfaces.Stringable {
	return NewStringable(strings.ToLower(s.value))
}

// Snake converts a string to snake case.
func (s *StringableImpl) Snake(delimiter string) interfaces.Stringable {
	if delimiter == "" {
		delimiter = "_"
	}
	return NewStringable(str.Snake(s.value))
}

// Studly converts a value to studly caps case.
func (s *StringableImpl) Studly() interfaces.Stringable {
	return NewStringable(str.Studly(s.value))
}

// Pascal converts the string to Pascal case.
func (s *StringableImpl) Pascal() interfaces.Stringable {
	return NewStringable(str.Pascal(s.value))
}

// Title converts the given string to proper case.
func (s *StringableImpl) Title() interfaces.Stringable {
	return NewStringable(strings.Title(s.value))
}

// Upper converts the given string to upper-case.
func (s *StringableImpl) Upper() interfaces.Stringable {
	return NewStringable(strings.ToUpper(s.value))
}

// Headline converts the given string to proper case for each word.
func (s *StringableImpl) Headline() interfaces.Stringable {
	return NewStringable(str.Headline(s.value))
}

// Apa converts the given string to APA-style title case.
func (s *StringableImpl) Apa() interfaces.Stringable {
	return NewStringable(str.Apa(s.value))
}

// Character and position methods

// CharAt gets the character at the specified index.
func (s *StringableImpl) CharAt(index int) string {
	runes := []rune(s.value)
	if index >= 0 && index < len(runes) {
		return string(runes[index])
	}
	return ""
}

// ChopStart removes the given string if it exists at the start of the current string.
func (s *StringableImpl) ChopStart(needle string) interfaces.Stringable {
	return NewStringable(str.ChopStart(s.value, needle))
}

// ChopEnd removes the given string if it exists at the end of the current string.
func (s *StringableImpl) ChopEnd(needle string) interfaces.Stringable {
	return NewStringable(str.ChopEnd(s.value, needle))
}

// Length returns the length of the given string.
func (s *StringableImpl) Length(encoding string) int {
	if encoding == "" || encoding == "UTF-8" {
		return utf8.RuneCountInString(s.value)
	}
	// For other encodings, return byte length as approximation
	return len(s.value)
}

// Position finds the multi-byte safe position of the first occurrence of the given substring.
func (s *StringableImpl) Position(needle string, offset int, encoding string) int {
	if offset < 0 {
		return -1
	}

	runes := []rune(s.value)
	needleRunes := []rune(needle)

	if offset >= len(runes) {
		return -1
	}

	for i := offset; i <= len(runes)-len(needleRunes); i++ {
		match := true
		for j, needleRune := range needleRunes {
			if i+j >= len(runes) || runes[i+j] != needleRune {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}

	return -1
}

// Content checking methods

// Contains determines if a given string contains a given substring.
func (s *StringableImpl) Contains(needles []string, ignoreCase bool) bool {
	value := s.value
	if ignoreCase {
		value = strings.ToLower(value)
	}

	for _, needle := range needles {
		searchNeedle := needle
		if ignoreCase {
			searchNeedle = strings.ToLower(needle)
		}
		if strings.Contains(value, searchNeedle) {
			return true
		}
	}
	return false
}

// ContainsAll determines if a given string contains all array values.
func (s *StringableImpl) ContainsAll(needles []string, ignoreCase bool) bool {
	value := s.value
	if ignoreCase {
		value = strings.ToLower(value)
	}

	for _, needle := range needles {
		searchNeedle := needle
		if ignoreCase {
			searchNeedle = strings.ToLower(needle)
		}
		if !strings.Contains(value, searchNeedle) {
			return false
		}
	}
	return true
}

// StartsWith determines if a given string starts with a given substring.
func (s *StringableImpl) StartsWith(needles []string) bool {
	for _, needle := range needles {
		if strings.HasPrefix(s.value, needle) {
			return true
		}
	}
	return false
}

// DoesntStartWith determines if a given string doesn't start with a given substring.
func (s *StringableImpl) DoesntStartWith(needles []string) bool {
	return !s.StartsWith(needles)
}

// EndsWith determines if a given string ends with a given substring.
func (s *StringableImpl) EndsWith(needles []string) bool {
	for _, needle := range needles {
		if strings.HasSuffix(s.value, needle) {
			return true
		}
	}
	return false
}

// DoesntEndWith determines if a given string doesn't end with a given substring.
func (s *StringableImpl) DoesntEndWith(needles []string) bool {
	return !s.EndsWith(needles)
}

// Exactly determines if the string is an exact match with the given value.
func (s *StringableImpl) Exactly(value string) bool {
	return s.value == value
}

// Is determines if a given string matches a given pattern.
func (s *StringableImpl) Is(pattern string, ignoreCase bool) bool {
	return str.Is(pattern, s.value)
}

// IsMatch determines if a given string matches a given pattern.
func (s *StringableImpl) IsMatch(pattern string) bool {
	matched, _ := regexp.MatchString(pattern, s.value)
	return matched
}

// Test determines if the string matches the given pattern.
func (s *StringableImpl) Test(pattern string) bool {
	return s.IsMatch(pattern)
}

// Validation methods

// IsAscii determines if a given string is 7 bit ASCII.
func (s *StringableImpl) IsAscii() bool {
	for _, r := range s.value {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsJson determines if a given string is valid JSON.
func (s *StringableImpl) IsJson() bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s.value), &js) == nil
}

// IsUrl determines if a given value is a valid URL.
func (s *StringableImpl) IsUrl(protocols []string) bool {
	return str.IsUrl(s.value, protocols)
}

// IsUuid determines if a given string is a valid UUID.
func (s *StringableImpl) IsUuid(version int) bool {
	return str.IsUuid(s.value)
}

// IsUlid determines if a given string is a valid ULID.
func (s *StringableImpl) IsUlid() bool {
	return str.IsUlid(s.value)
}

// IsEmpty determines if the given string is empty.
func (s *StringableImpl) IsEmpty() bool {
	return s.value == ""
}

// IsNotEmpty determines if the given string is not empty.
func (s *StringableImpl) IsNotEmpty() bool {
	return s.value != ""
}

// String modification methods

// ConvertCase converts the case of a string.
func (s *StringableImpl) ConvertCase(mode int, encoding string) interfaces.Stringable {
	return NewStringable(str.ConvertCase(s.value, mode, encoding))
}

// Deduplicate replaces consecutive instances of a given character with a single character.
func (s *StringableImpl) Deduplicate(character string) interfaces.Stringable {
	if character == "" {
		return NewStringable(s.value)
	}

	pattern := regexp.QuoteMeta(character) + "+"
	re := regexp.MustCompile(pattern)
	result := re.ReplaceAllString(s.value, character)
	return NewStringable(result)
}

// Excerpt extracts an excerpt from text that matches the first instance of a phrase.
func (s *StringableImpl) Excerpt(phrase string, options map[string]interface{}) string {
	return str.Excerpt(s.value, phrase, options)
}

// Finish caps a string with a single instance of a given value.
func (s *StringableImpl) Finish(cap string) interfaces.Stringable {
	return NewStringable(str.Finish(s.value, cap))
}

// Limit limits the number of characters in a string.
func (s *StringableImpl) Limit(limit int, end string, preserveWords bool) interfaces.Stringable {
	return NewStringable(str.Limit(s.value, limit, end))
}

// Mask masks a portion of a string with a repeated character.
func (s *StringableImpl) Mask(character string, index int, length int, encoding string) interfaces.Stringable {
	return NewStringable(str.Mask(s.value, character, index, length))
}

// PadBoth pads both sides of the string with another.
func (s *StringableImpl) PadBoth(length int, pad string) interfaces.Stringable {
	return NewStringable(str.PadBoth(s.value, length, pad))
}

// PadLeft pads the left side of the string with another.
func (s *StringableImpl) PadLeft(length int, pad string) interfaces.Stringable {
	return NewStringable(str.PadLeft(s.value, length, pad))
}

// PadRight pads the right side of the string with another.
func (s *StringableImpl) PadRight(length int, pad string) interfaces.Stringable {
	return NewStringable(str.PadRight(s.value, length, pad))
}

// Prepend prepends the given values to the string.
func (s *StringableImpl) Prepend(values ...string) interfaces.Stringable {
	return NewStringable(strings.Join(values, "") + s.value)
}

// Remove removes any occurrence of the given string in the subject.
func (s *StringableImpl) Remove(search []string, caseSensitive bool) interfaces.Stringable {
	result := s.value
	for _, needle := range search {
		if caseSensitive {
			result = strings.ReplaceAll(result, needle, "")
		} else {
			// Case insensitive replacement
			re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(needle))
			result = re.ReplaceAllString(result, "")
		}
	}
	return NewStringable(result)
}

// Repeat repeats the string.
func (s *StringableImpl) Repeat(times int) interfaces.Stringable {
	return NewStringable(strings.Repeat(s.value, times))
}

// Replace replaces the given value in the given string.
func (s *StringableImpl) Replace(search, replace []string, caseSensitive bool) interfaces.Stringable {
	result := s.value

	minLen := len(search)
	if len(replace) < minLen {
		minLen = len(replace)
	}

	for i := 0; i < minLen; i++ {
		if caseSensitive {
			result = strings.ReplaceAll(result, search[i], replace[i])
		} else {
			re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(search[i]))
			result = re.ReplaceAllString(result, replace[i])
		}
	}

	return NewStringable(result)
}

// ReplaceArray replaces a given value in the string sequentially with an array.
func (s *StringableImpl) ReplaceArray(search string, replace []string) interfaces.Stringable {
	result := s.value
	for _, replacement := range replace {
		index := strings.Index(result, search)
		if index == -1 {
			break
		}
		result = result[:index] + replacement + result[index+len(search):]
	}
	return NewStringable(result)
}

// ReplaceFirst replaces the first occurrence of a given value in the string.
func (s *StringableImpl) ReplaceFirst(search, replace string) interfaces.Stringable {
	index := strings.Index(s.value, search)
	if index == -1 {
		return NewStringable(s.value)
	}
	result := s.value[:index] + replace + s.value[index+len(search):]
	return NewStringable(result)
}

// ReplaceLast replaces the last occurrence of a given value in the string.
func (s *StringableImpl) ReplaceLast(search, replace string) interfaces.Stringable {
	index := strings.LastIndex(s.value, search)
	if index == -1 {
		return NewStringable(s.value)
	}
	result := s.value[:index] + replace + s.value[index+len(search):]
	return NewStringable(result)
}

// ReplaceStart replaces the first occurrence of the given value if it appears at the start of the string.
func (s *StringableImpl) ReplaceStart(search, replace string) interfaces.Stringable {
	if strings.HasPrefix(s.value, search) {
		return NewStringable(replace + s.value[len(search):])
	}
	return NewStringable(s.value)
}

// ReplaceEnd replaces the last occurrence of a given value if it appears at the end of the string.
func (s *StringableImpl) ReplaceEnd(search, replace string) interfaces.Stringable {
	if strings.HasSuffix(s.value, search) {
		return NewStringable(s.value[:len(s.value)-len(search)] + replace)
	}
	return NewStringable(s.value)
}

// ReplaceMatches replaces the patterns matching the given regular expression.
func (s *StringableImpl) ReplaceMatches(pattern, replace string, limit int) interfaces.Stringable {
	re := regexp.MustCompile(pattern)
	if limit <= 0 {
		return NewStringable(re.ReplaceAllString(s.value, replace))
	}

	result := s.value
	for i := 0; i < limit; i++ {
		newResult := re.ReplaceAllString(result, replace)
		if newResult == result {
			break
		}
		result = newResult
	}

	return NewStringable(result)
}

// Reverse reverses the string.
func (s *StringableImpl) Reverse() interfaces.Stringable {
	runes := []rune(s.value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return NewStringable(string(runes))
}

// Squish removes all "extra" blank space from the given string.
func (s *StringableImpl) Squish() interfaces.Stringable {
	// Replace multiple whitespace with single space and trim
	re := regexp.MustCompile(`\s+`)
	result := re.ReplaceAllString(strings.TrimSpace(s.value), " ")
	return NewStringable(result)
}

// Start begins a string with a single instance of a given value.
func (s *StringableImpl) Start(prefix string) interfaces.Stringable {
	return NewStringable(str.Start(s.value, prefix))
}

// StripTags strips HTML and PHP tags from the given string.
func (s *StringableImpl) StripTags(allowedTags string) interfaces.Stringable {
	return NewStringable(str.StripTags(s.value, allowedTags))
}

// Substr returns the portion of the string specified by the start and length parameters.
func (s *StringableImpl) Substr(start int, length int, encoding string) interfaces.Stringable {
	runes := []rune(s.value)
	runeLen := len(runes)

	if start < 0 {
		start = runeLen + start
	}

	if start < 0 {
		start = 0
	}

	if start >= runeLen {
		return NewStringable("")
	}

	end := start + length
	if length < 0 {
		end = runeLen + length
	}

	if end > runeLen {
		end = runeLen
	}

	if end <= start {
		return NewStringable("")
	}

	return NewStringable(string(runes[start:end]))
}

// SubstrCount returns the number of substring occurrences.
func (s *StringableImpl) SubstrCount(needle string, offset int, length int) int {
	return str.SubstrCount(s.value, needle, offset, length)
}

// SubstrReplace replaces text within a portion of a string.
func (s *StringableImpl) SubstrReplace(replace string, offset int, length int) interfaces.Stringable {
	return NewStringable(str.SubstrReplace(s.value, replace, offset, length))
}

// Swap swaps multiple keywords in a string with other keywords.
func (s *StringableImpl) Swap(replacements map[string]string) interfaces.Stringable {
	result := s.value
	for search, replace := range replacements {
		result = strings.ReplaceAll(result, search, replace)
	}
	return NewStringable(result)
}

// Take takes the first or last {limit} characters.
func (s *StringableImpl) Take(limit int) interfaces.Stringable {
	runes := []rune(s.value)
	runeLen := len(runes)

	if limit == 0 {
		return NewStringable("")
	}

	if limit > 0 {
		if limit >= runeLen {
			return NewStringable(s.value)
		}
		return NewStringable(string(runes[:limit]))
	} else {
		// Negative limit means take from the end
		limit = -limit
		if limit >= runeLen {
			return NewStringable(s.value)
		}
		return NewStringable(string(runes[runeLen-limit:]))
	}
}

// Trim trims the string of the given characters.
func (s *StringableImpl) Trim(characters string) interfaces.Stringable {
	if characters == "" {
		return NewStringable(strings.TrimSpace(s.value))
	}
	return NewStringable(strings.Trim(s.value, characters))
}

// Ltrim left trims the string of the given characters.
func (s *StringableImpl) Ltrim(characters string) interfaces.Stringable {
	if characters == "" {
		return NewStringable(strings.TrimLeftFunc(s.value, unicode.IsSpace))
	}
	return NewStringable(strings.TrimLeft(s.value, characters))
}

// Rtrim right trims the string of the given characters.
func (s *StringableImpl) Rtrim(characters string) interfaces.Stringable {
	if characters == "" {
		return NewStringable(strings.TrimRightFunc(s.value, unicode.IsSpace))
	}
	return NewStringable(strings.TrimRight(s.value, characters))
}

// Lcfirst makes a string's first character lowercase.
func (s *StringableImpl) Lcfirst() interfaces.Stringable {
	if s.value == "" {
		return NewStringable("")
	}
	runes := []rune(s.value)
	runes[0] = unicode.ToLower(runes[0])
	return NewStringable(string(runes))
}

// Ucfirst makes a string's first character uppercase.
func (s *StringableImpl) Ucfirst() interfaces.Stringable {
	if s.value == "" {
		return NewStringable("")
	}
	runes := []rune(s.value)
	runes[0] = unicode.ToUpper(runes[0])
	return NewStringable(string(runes))
}

// Wrap wraps the string with the given strings.
func (s *StringableImpl) Wrap(before, after string) interfaces.Stringable {
	return NewStringable(before + s.value + after)
}

// Unwrap unwraps the string with the given strings.
func (s *StringableImpl) Unwrap(before, after string) interfaces.Stringable {
	result := s.value
	if strings.HasPrefix(result, before) {
		result = result[len(before):]
	}
	if strings.HasSuffix(result, after) {
		result = result[:len(result)-len(after)]
	}
	return NewStringable(result)
}

// String array and collection methods

// Explode explodes the string into a slice.
func (s *StringableImpl) Explode(delimiter string, limit int) []string {
	if limit <= 0 {
		return strings.Split(s.value, delimiter)
	}
	return strings.SplitN(s.value, delimiter, limit)
}

// Split splits a string using a regular expression or by length.
func (s *StringableImpl) Split(pattern interface{}, limit int, flags int) []string {
	switch p := pattern.(type) {
	case string:
		// Treat as regular expression
		re := regexp.MustCompile(p)
		if limit <= 0 {
			return re.Split(s.value, -1)
		}
		return re.Split(s.value, limit)
	case int:
		// Split by length
		if p <= 0 {
			return []string{s.value}
		}

		var result []string
		runes := []rune(s.value)
		for i := 0; i < len(runes); i += p {
			end := i + p
			if end > len(runes) {
				end = len(runes)
			}
			result = append(result, string(runes[i:end]))
		}
		return result
	default:
		return []string{s.value}
	}
}

// Ucsplit splits a string by uppercase characters.
func (s *StringableImpl) Ucsplit() []string {
	return str.Ucsplit(s.value)
}

// Words limits the number of words in a string.
func (s *StringableImpl) Words(words int, end string) interfaces.Stringable {
	return NewStringable(str.Words(s.value, words, end))
}

// WordCount gets the number of words a string contains.
func (s *StringableImpl) WordCount(characters string) int {
	return str.WordCount(s.value, characters)
}

// WordWrap wraps a string to a given number of characters.
func (s *StringableImpl) WordWrap(characters int, breakStr string, cutLongWords bool) interfaces.Stringable {
	return NewStringable(str.WordWrap(s.value, characters, breakStr))
}

// Output methods

// ToString gets the underlying string value.
func (s *StringableImpl) ToString() string {
	return s.value
}

// Value gets the underlying string value.
func (s *StringableImpl) Value() string {
	return s.value
}

// ToInteger gets the underlying string value as an integer.
func (s *StringableImpl) ToInteger(base int) int {
	if base <= 0 {
		base = 10
	}
	result, _ := strconv.ParseInt(s.value, base, 64)
	return int(result)
}

// ToFloat gets the underlying string value as a float.
func (s *StringableImpl) ToFloat() float64 {
	result, _ := strconv.ParseFloat(s.value, 64)
	return result
}

// ToBoolean gets the underlying string value as a boolean.
func (s *StringableImpl) ToBoolean() bool {
	lower := strings.ToLower(strings.TrimSpace(s.value))
	return lower == "true" || lower == "1" || lower == "yes" || lower == "on"
}

// Conditional methods - simplified implementations

// When executes the given callback if the condition is true.
func (s *StringableImpl) When(condition bool, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	if condition && callback != nil {
		return callback(s)
	}
	if !condition && defaultCallback != nil {
		return defaultCallback(s)
	}
	return s
}

// Unless executes the given callback if the condition is false.
func (s *StringableImpl) Unless(condition bool, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(!condition, callback, defaultCallback)
}

// WhenContains executes the given callback if the string contains a given substring.
func (s *StringableImpl) WhenContains(needles []string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.Contains(needles, false), callback, defaultCallback)
}

// WhenContainsAll executes the given callback if the string contains all array values.
func (s *StringableImpl) WhenContainsAll(needles []string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.ContainsAll(needles, false), callback, defaultCallback)
}

// WhenEmpty executes the given callback if the string is empty.
func (s *StringableImpl) WhenEmpty(callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.IsEmpty(), callback, defaultCallback)
}

// WhenNotEmpty executes the given callback if the string is not empty.
func (s *StringableImpl) WhenNotEmpty(callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.IsNotEmpty(), callback, defaultCallback)
}

// WhenStartsWith executes the given callback if the string starts with a given substring.
func (s *StringableImpl) WhenStartsWith(needles []string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.StartsWith(needles), callback, defaultCallback)
}

// WhenEndsWith executes the given callback if the string ends with a given substring.
func (s *StringableImpl) WhenEndsWith(needles []string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.EndsWith(needles), callback, defaultCallback)
}

// WhenExactly executes the given callback if the string is an exact match with the given value.
func (s *StringableImpl) WhenExactly(value string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.Exactly(value), callback, defaultCallback)
}

// WhenIs executes the given callback if the string matches a given pattern.
func (s *StringableImpl) WhenIs(pattern string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.Is(pattern, false), callback, defaultCallback)
}

// WhenIsAscii executes the given callback if the string is 7 bit ASCII.
func (s *StringableImpl) WhenIsAscii(callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.IsAscii(), callback, defaultCallback)
}

// WhenIsUuid executes the given callback if the string is a valid UUID.
func (s *StringableImpl) WhenIsUuid(callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.IsUuid(0), callback, defaultCallback)
}

// WhenIsUlid executes the given callback if the string is a valid ULID.
func (s *StringableImpl) WhenIsUlid(callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.IsUlid(), callback, defaultCallback)
}

// WhenTest executes the given callback if the string matches the given pattern.
func (s *StringableImpl) WhenTest(pattern string, callback func(interfaces.Stringable) interfaces.Stringable, defaultCallback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return s.When(s.Test(pattern), callback, defaultCallback)
}

// Utility methods

// Pipe calls the given callback and returns a new string.
func (s *StringableImpl) Pipe(callback func(interfaces.Stringable) interfaces.Stringable) interfaces.Stringable {
	return callback(s)
}

// Tap calls the given callback with the string and returns the string.
func (s *StringableImpl) Tap(callback func(interfaces.Stringable)) interfaces.Stringable {
	callback(s)
	return s
}

// Dump dumps the string and returns the string.
func (s *StringableImpl) Dump() interfaces.Stringable {
	fmt.Printf("Stringable: %q\n", s.value)
	return s
}

// ClassBasename gets the basename of the class path.
func (s *StringableImpl) ClassBasename() interfaces.Stringable {
	return NewStringable(str.ClassBasename(s.value))
}

// Dirname gets the parent directory's path.
func (s *StringableImpl) Dirname(levels int) interfaces.Stringable {
	return NewStringable(str.Dirname(s.value, levels))
}

// Scan parses input from a string according to a format.
func (s *StringableImpl) Scan(format string) []interface{} {
	return str.Scan(s.value, format)
}

// Additional methods that need implementation stubs for interface compliance

// Plural gets the plural form of an English word.
func (s *StringableImpl) Plural(count int, prependCount bool) interfaces.Stringable {
	return NewStringable(str.Plural(s.value, count))
}

// PluralStudly pluralizes the last word of an English, studly caps case string.
func (s *StringableImpl) PluralStudly(count int) interfaces.Stringable {
	return NewStringable(str.PluralStudly(s.value, count))
}

// PluralPascal pluralizes the last word of an English, Pascal caps case string.
func (s *StringableImpl) PluralPascal(count int) interfaces.Stringable {
	return NewStringable(str.PluralPascal(s.value, count))
}

// Singular gets the singular form of an English word.
func (s *StringableImpl) Singular() interfaces.Stringable {
	return NewStringable(str.Singular(s.value))
}

// Slug generates a URL friendly "slug" from a given string.
func (s *StringableImpl) Slug(separator, language string, dictionary map[string]string) interfaces.Stringable {
	return NewStringable(str.Slug(s.value, separator))
}

// Match gets the string matching the given pattern.
func (s *StringableImpl) Match(pattern string) interfaces.Stringable {
	return NewStringable(str.Match(s.value, pattern))
}

// MatchAll gets all strings matching the given pattern.
func (s *StringableImpl) MatchAll(pattern string) []string {
	return str.MatchAll(s.value, pattern)
}

// Numbers removes all non-numeric characters from a string.
func (s *StringableImpl) Numbers() interfaces.Stringable {
	numbers := str.Numbers(s.value)
	if len(numbers) > 0 {
		return NewStringable(strings.Join(numbers, ""))
	}
	return NewStringable("")
}

// ParseCallback parses a Class@method style callback into class and method.
func (s *StringableImpl) ParseCallback(defaultMethod string) []string {
	class, method := str.ParseCallback(s.value, defaultMethod)
	return []string{class, method}
}

// ToBase64 converts the string to Base64 encoding.
func (s *StringableImpl) ToBase64() interfaces.Stringable {
	return NewStringable(str.ToBase64(s.value))
}

// FromBase64 decodes the Base64 encoded string.
func (s *StringableImpl) FromBase64(strict bool) interfaces.Stringable {
	result, err := str.FromBase64(s.value, strict)
	if err != nil {
		return NewStringable("")
	}
	return NewStringable(result)
}

// Hash hashes the string using the given algorithm.
func (s *StringableImpl) Hash(algorithm string) interfaces.Stringable {
	return NewStringable(str.Hash(s.value, algorithm))
}

// Transliterate transliterates a string to its closest ASCII representation.
func (s *StringableImpl) Transliterate(unknown string, strict bool) interfaces.Stringable {
	return NewStringable(str.Transliterate(s.value, unknown))
}
