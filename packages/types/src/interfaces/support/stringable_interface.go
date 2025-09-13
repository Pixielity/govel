package interfaces

/**
 * StringableInterface defines the contract for Laravel-style string manipulation functionality.
 *
 * This interface provides a fluent API for comprehensive string operations including:
 * - JSON serialization support through JsonSerializable interface
 * - Array-like access to characters through ArrayAccess interface
 * - Extensive string manipulation and case conversion methods
 * - Content validation and checking methods
 * - Pattern matching and regular expression operations
 * - Conditional execution based on string properties
 * - String transformation and modification capabilities
 *
 * The interface follows the Laravel Stringable contract, enabling fluent chaining
 * of string operations while maintaining immutability through method chaining.
 *
 * Example Usage:
 *   str := NewStringable("Hello World")
 *   result := str.Lower().Snake().Upper()  // Fluent chaining
 *   fmt.Println(result.ToString())         // "HELLO_WORLD"
 *
 * Thread Safety:
 *   Implementations should be immutable and thread-safe. Each method should
 *   return a new instance rather than modifying the existing one.
 *
 * Performance Considerations:
 *   Methods should be optimized for common use cases while maintaining
 *   compatibility with Laravel's Stringable behavior.
 */
type StringableInterface interface {
	/**
	 * JsonSerializable interface methods
	 * These methods support JSON serialization of the string value
	 */

	/**
	 * JsonSerialize returns the string value for JSON serialization.
	 *
	 * @return interface{} The underlying string value for JSON encoding
	 */
	JsonSerialize() interface{}

	/**
	 * ArrayAccess interface methods
	 * These methods provide array-like access to string characters
	 */

	/**
	 * OffsetExists checks if a character exists at the given index.
	 *
	 * @param offset interface{} The index to check (should be an int)
	 * @return bool True if the index is valid, false otherwise
	 */
	OffsetExists(offset interface{}) bool

	/**
	 * OffsetGet retrieves the character at the given index.
	 *
	 * @param offset interface{} The index to retrieve (should be an int)
	 * @return interface{} The character at the index, or nil if invalid
	 */
	OffsetGet(offset interface{}) interface{}

	/**
	 * OffsetSet is not supported for strings as they are immutable.
	 *
	 * @param offset interface{} The index (ignored)
	 * @param value interface{} The value (ignored)
	 * @panic Always panics to maintain string immutability
	 */
	OffsetSet(offset, value interface{})

	/**
	 * OffsetUnset is not supported for strings as they are immutable.
	 *
	 * @param offset interface{} The index (ignored)
	 * @panic Always panics to maintain string immutability
	 */
	OffsetUnset(offset interface{})

	/**
	 * String manipulation and positioning methods
	 */

	/**
	 * After returns the remainder of a string after the first occurrence of a given value.
	 *
	 * @param search string The value to search for
	 * @return StringableInterface New stringable with content after the search value
	 */
	After(search string) StringableInterface

	/**
	 * AfterLast returns the remainder of a string after the last occurrence of a given value.
	 *
	 * @param search string The value to search for
	 * @return StringableInterface New stringable with content after the last occurrence
	 */
	AfterLast(search string) StringableInterface

	/**
	 * Append appends the given values to the string.
	 *
	 * @param values ...string Variable number of strings to append
	 * @return StringableInterface New stringable with appended content
	 */
	Append(values ...string) StringableInterface

	/**
	 * NewLine appends one or more new lines to the string.
	 *
	 * @param count int Number of new lines to append
	 * @return StringableInterface New stringable with appended new lines
	 */
	NewLine(count int) StringableInterface

	/**
	 * Ascii transliterates a UTF-8 value to ASCII.
	 *
	 * @param language string Language code for transliteration rules
	 * @return StringableInterface New stringable with ASCII representation
	 */
	Ascii(language string) StringableInterface

	/**
	 * Basename gets the trailing name component of the path.
	 *
	 * @param suffix string Optional suffix to remove
	 * @return StringableInterface New stringable with basename
	 */
	Basename(suffix string) StringableInterface

	/**
	 * Before gets the portion of a string before the first occurrence of a given value.
	 *
	 * @param search string The value to search for
	 * @return StringableInterface New stringable with content before the search value
	 */
	Before(search string) StringableInterface

	/**
	 * BeforeLast gets the portion of a string before the last occurrence of a given value.
	 *
	 * @param search string The value to search for
	 * @return StringableInterface New stringable with content before the last occurrence
	 */
	BeforeLast(search string) StringableInterface

	/**
	 * Between gets the portion of a string between two given values.
	 *
	 * @param from string Starting delimiter
	 * @param to string Ending delimiter
	 * @return StringableInterface New stringable with content between delimiters
	 */
	Between(from, to string) StringableInterface

	/**
	 * BetweenFirst gets the smallest possible portion of a string between two given values.
	 *
	 * @param from string Starting delimiter
	 * @param to string Ending delimiter
	 * @return StringableInterface New stringable with smallest content between delimiters
	 */
	BetweenFirst(from, to string) StringableInterface

	/**
	 * Case conversion methods
	 */

	/**
	 * Camel converts a value to camel case.
	 *
	 * @return StringableInterface New stringable in camel case
	 */
	Camel() StringableInterface

	/**
	 * Kebab converts a string to kebab case.
	 *
	 * @return StringableInterface New stringable in kebab case
	 */
	Kebab() StringableInterface

	/**
	 * Lower converts the given string to lower-case.
	 *
	 * @return StringableInterface New stringable in lowercase
	 */
	Lower() StringableInterface

	/**
	 * Snake converts a string to snake case.
	 *
	 * @param delimiter string Delimiter to use (defaults to underscore)
	 * @return StringableInterface New stringable in snake case
	 */
	Snake(delimiter string) StringableInterface

	/**
	 * Studly converts a value to studly caps case.
	 *
	 * @return StringableInterface New stringable in studly caps case
	 */
	Studly() StringableInterface

	/**
	 * Pascal converts the string to Pascal case.
	 *
	 * @return StringableInterface New stringable in Pascal case
	 */
	Pascal() StringableInterface

	/**
	 * Title converts the given string to proper case.
	 *
	 * @return StringableInterface New stringable in title case
	 */
	Title() StringableInterface

	/**
	 * Upper converts the given string to upper-case.
	 *
	 * @return StringableInterface New stringable in uppercase
	 */
	Upper() StringableInterface

	/**
	 * Headline converts the given string to proper case for each word.
	 *
	 * @return StringableInterface New stringable with headline formatting
	 */
	Headline() StringableInterface

	/**
	 * Apa converts the given string to APA-style title case.
	 *
	 * @return StringableInterface New stringable in APA title case
	 */
	Apa() StringableInterface

	/**
	 * Character and position methods
	 */

	/**
	 * CharAt gets the character at the specified index.
	 *
	 * @param index int The character index
	 * @return string The character at the index, or empty string if invalid
	 */
	CharAt(index int) string

	/**
	 * ChopStart removes the given string if it exists at the start of the current string.
	 *
	 * @param needle string The string to remove from the start
	 * @return StringableInterface New stringable with chopped start
	 */
	ChopStart(needle string) StringableInterface

	/**
	 * ChopEnd removes the given string if it exists at the end of the current string.
	 *
	 * @param needle string The string to remove from the end
	 * @return StringableInterface New stringable with chopped end
	 */
	ChopEnd(needle string) StringableInterface

	/**
	 * Length returns the length of the given string.
	 *
	 * @param encoding string Character encoding (defaults to UTF-8)
	 * @return int The string length in characters
	 */
	Length(encoding string) int

	/**
	 * Position finds the multi-byte safe position of the first occurrence of the given substring.
	 *
	 * @param needle string The substring to find
	 * @param offset int Starting position for search
	 * @param encoding string Character encoding
	 * @return int Position of the substring, or -1 if not found
	 */
	Position(needle string, offset int, encoding string) int

	/**
	 * Content checking methods
	 */

	/**
	 * Contains determines if a given string contains a given substring.
	 *
	 * @param needles []string Array of substrings to check for
	 * @param ignoreCase bool Whether to ignore case in comparison
	 * @return bool True if any needle is found
	 */
	Contains(needles []string, ignoreCase bool) bool

	/**
	 * ContainsAll determines if a given string contains all array values.
	 *
	 * @param needles []string Array of substrings to check for
	 * @param ignoreCase bool Whether to ignore case in comparison
	 * @return bool True if all needles are found
	 */
	ContainsAll(needles []string, ignoreCase bool) bool

	/**
	 * StartsWith determines if a given string starts with a given substring.
	 *
	 * @param needles []string Array of prefixes to check
	 * @return bool True if string starts with any needle
	 */
	StartsWith(needles []string) bool

	/**
	 * DoesntStartWith determines if a given string doesn't start with a given substring.
	 *
	 * @param needles []string Array of prefixes to check
	 * @return bool True if string doesn't start with any needle
	 */
	DoesntStartWith(needles []string) bool

	/**
	 * EndsWith determines if a given string ends with a given substring.
	 *
	 * @param needles []string Array of suffixes to check
	 * @return bool True if string ends with any needle
	 */
	EndsWith(needles []string) bool

	/**
	 * DoesntEndWith determines if a given string doesn't end with a given substring.
	 *
	 * @param needles []string Array of suffixes to check
	 * @return bool True if string doesn't end with any needle
	 */
	DoesntEndWith(needles []string) bool

	/**
	 * Exactly determines if the string is an exact match with the given value.
	 *
	 * @param value string The value to compare exactly
	 * @return bool True if strings are exactly equal
	 */
	Exactly(value string) bool

	/**
	 * Is determines if a given string matches a given pattern.
	 *
	 * @param pattern string The pattern to match (supports wildcards)
	 * @param ignoreCase bool Whether to ignore case in comparison
	 * @return bool True if string matches the pattern
	 */
	Is(pattern string, ignoreCase bool) bool

	/**
	 * IsMatch determines if a given string matches a given pattern.
	 *
	 * @param pattern string The regular expression pattern
	 * @return bool True if string matches the regex pattern
	 */
	IsMatch(pattern string) bool

	/**
	 * Test determines if the string matches the given pattern.
	 *
	 * @param pattern string The regular expression pattern
	 * @return bool True if string matches the regex pattern
	 */
	Test(pattern string) bool

	/**
	 * Validation methods
	 */

	/**
	 * IsAscii determines if a given string is 7 bit ASCII.
	 *
	 * @return bool True if string contains only ASCII characters
	 */
	IsAscii() bool

	/**
	 * IsJson determines if a given string is valid JSON.
	 *
	 * @return bool True if string is valid JSON
	 */
	IsJson() bool

	/**
	 * IsUrl determines if a given value is a valid URL.
	 *
	 * @param protocols []string Allowed protocols (optional)
	 * @return bool True if string is a valid URL
	 */
	IsUrl(protocols []string) bool

	/**
	 * IsUuid determines if a given string is a valid UUID.
	 *
	 * @param version int UUID version to validate (0 for any)
	 * @return bool True if string is a valid UUID
	 */
	IsUuid(version int) bool

	/**
	 * IsUlid determines if a given string is a valid ULID.
	 *
	 * @return bool True if string is a valid ULID
	 */
	IsUlid() bool

	/**
	 * IsEmpty determines if the given string is empty.
	 *
	 * @return bool True if string is empty
	 */
	IsEmpty() bool

	/**
	 * IsNotEmpty determines if the given string is not empty.
	 *
	 * @return bool True if string is not empty
	 */
	IsNotEmpty() bool

	/**
	 * String modification methods
	 */

	/**
	 * ConvertCase converts the case of a string.
	 *
	 * @param mode int Case conversion mode constant
	 * @param encoding string Character encoding
	 * @return StringableInterface New stringable with converted case
	 */
	ConvertCase(mode int, encoding string) StringableInterface

	/**
	 * Deduplicate replaces consecutive instances of a given character with a single character.
	 *
	 * @param character string The character to deduplicate
	 * @return StringableInterface New stringable with deduplicated characters
	 */
	Deduplicate(character string) StringableInterface

	/**
	 * Excerpt extracts an excerpt from text that matches the first instance of a phrase.
	 *
	 * @param phrase string The phrase to search for
	 * @param options map[string]interface{} Extraction options
	 * @return string The extracted excerpt
	 */
	Excerpt(phrase string, options map[string]interface{}) string

	/**
	 * Finish caps a string with a single instance of a given value.
	 *
	 * @param cap string The string to append if not already present
	 * @return StringableInterface New stringable with finish cap
	 */
	Finish(cap string) StringableInterface

	/**
	 * Limit limits the number of characters in a string.
	 *
	 * @param limit int Maximum number of characters
	 * @param end string String to append when truncated
	 * @param preserveWords bool Whether to preserve word boundaries
	 * @return StringableInterface New stringable with limited length
	 */
	Limit(limit int, end string, preserveWords bool) StringableInterface

	/**
	 * Mask masks a portion of a string with a repeated character.
	 *
	 * @param character string The masking character
	 * @param index int Starting position to mask
	 * @param length int Number of characters to mask
	 * @param encoding string Character encoding
	 * @return StringableInterface New stringable with masked content
	 */
	Mask(character string, index int, length int, encoding string) StringableInterface

	/**
	 * PadBoth pads both sides of the string with another.
	 *
	 * @param length int Total desired length
	 * @param pad string Padding string
	 * @return StringableInterface New stringable with both sides padded
	 */
	PadBoth(length int, pad string) StringableInterface

	/**
	 * PadLeft pads the left side of the string with another.
	 *
	 * @param length int Total desired length
	 * @param pad string Padding string
	 * @return StringableInterface New stringable with left padding
	 */
	PadLeft(length int, pad string) StringableInterface

	/**
	 * PadRight pads the right side of the string with another.
	 *
	 * @param length int Total desired length
	 * @param pad string Padding string
	 * @return StringableInterface New stringable with right padding
	 */
	PadRight(length int, pad string) StringableInterface

	/**
	 * Prepend prepends the given values to the string.
	 *
	 * @param values ...string Variable number of strings to prepend
	 * @return StringableInterface New stringable with prepended content
	 */
	Prepend(values ...string) StringableInterface

	/**
	 * Remove removes any occurrence of the given string in the subject.
	 *
	 * @param search []string Array of strings to remove
	 * @param caseSensitive bool Whether removal is case sensitive
	 * @return StringableInterface New stringable with removed content
	 */
	Remove(search []string, caseSensitive bool) StringableInterface

	/**
	 * Repeat repeats the string.
	 *
	 * @param times int Number of times to repeat
	 * @return StringableInterface New stringable with repeated content
	 */
	Repeat(times int) StringableInterface

	/**
	 * Replace replaces the given value in the given string.
	 *
	 * @param search []string Array of strings to search for
	 * @param replace []string Array of replacement strings
	 * @param caseSensitive bool Whether replacement is case sensitive
	 * @return StringableInterface New stringable with replaced content
	 */
	Replace(search, replace []string, caseSensitive bool) StringableInterface

	/**
	 * ReplaceArray replaces a given value in the string sequentially with an array.
	 *
	 * @param search string The string to search for
	 * @param replace []string Array of replacement strings
	 * @return StringableInterface New stringable with array replacements
	 */
	ReplaceArray(search string, replace []string) StringableInterface

	/**
	 * ReplaceFirst replaces the first occurrence of a given value in the string.
	 *
	 * @param search string The string to search for
	 * @param replace string The replacement string
	 * @return StringableInterface New stringable with first occurrence replaced
	 */
	ReplaceFirst(search, replace string) StringableInterface

	/**
	 * ReplaceLast replaces the last occurrence of a given value in the string.
	 *
	 * @param search string The string to search for
	 * @param replace string The replacement string
	 * @return StringableInterface New stringable with last occurrence replaced
	 */
	ReplaceLast(search, replace string) StringableInterface

	/**
	 * ReplaceStart replaces the first occurrence of the given value if it appears at the start of the string.
	 *
	 * @param search string The string to search for at start
	 * @param replace string The replacement string
	 * @return StringableInterface New stringable with start replaced
	 */
	ReplaceStart(search, replace string) StringableInterface

	/**
	 * ReplaceEnd replaces the last occurrence of a given value if it appears at the end of the string.
	 *
	 * @param search string The string to search for at end
	 * @param replace string The replacement string
	 * @return StringableInterface New stringable with end replaced
	 */
	ReplaceEnd(search, replace string) StringableInterface

	/**
	 * ReplaceMatches replaces the patterns matching the given regular expression.
	 *
	 * @param pattern string The regular expression pattern
	 * @param replace string The replacement string
	 * @param limit int Maximum number of replacements (0 for unlimited)
	 * @return StringableInterface New stringable with regex replacements
	 */
	ReplaceMatches(pattern, replace string, limit int) StringableInterface

	/**
	 * Reverse reverses the string.
	 *
	 * @return StringableInterface New stringable with reversed content
	 */
	Reverse() StringableInterface

	/**
	 * Squish removes all "extra" blank space from the given string.
	 *
	 * @return StringableInterface New stringable with normalized whitespace
	 */
	Squish() StringableInterface

	/**
	 * Start begins a string with a single instance of a given value.
	 *
	 * @param prefix string The prefix to ensure at start
	 * @return StringableInterface New stringable with ensured start
	 */
	Start(prefix string) StringableInterface

	/**
	 * StripTags strips HTML and PHP tags from the given string.
	 *
	 * @param allowedTags string Tags to allow (in HTML format)
	 * @return StringableInterface New stringable with tags stripped
	 */
	StripTags(allowedTags string) StringableInterface

	/**
	 * Substr returns the portion of the string specified by the start and length parameters.
	 *
	 * @param start int Starting position
	 * @param length int Length of substring (-1 for to end)
	 * @param encoding string Character encoding
	 * @return StringableInterface New stringable with substring
	 */
	Substr(start int, length int, encoding string) StringableInterface

	/**
	 * SubstrCount returns the number of substring occurrences.
	 *
	 * @param needle string The substring to count
	 * @param offset int Starting position for counting
	 * @param length int Length of search area
	 * @return int Number of occurrences found
	 */
	SubstrCount(needle string, offset int, length int) int

	/**
	 * SubstrReplace replaces text within a portion of a string.
	 *
	 * @param replace string The replacement string
	 * @param offset int Starting position for replacement
	 * @param length int Length to replace
	 * @return StringableInterface New stringable with substring replaced
	 */
	SubstrReplace(replace string, offset int, length int) StringableInterface

	/**
	 * Swap swaps multiple keywords in a string with other keywords.
	 *
	 * @param replacements map[string]string Map of search=>replacement pairs
	 * @return StringableInterface New stringable with swapped content
	 */
	Swap(replacements map[string]string) StringableInterface

	/**
	 * Take takes the first or last {limit} characters.
	 *
	 * @param limit int Number of characters to take (negative for from end)
	 * @return StringableInterface New stringable with taken characters
	 */
	Take(limit int) StringableInterface

	/**
	 * Trim trims the string of the given characters.
	 *
	 * @param characters string Characters to trim (empty for whitespace)
	 * @return StringableInterface New stringable with trimmed content
	 */
	Trim(characters string) StringableInterface

	/**
	 * Ltrim left trims the string of the given characters.
	 *
	 * @param characters string Characters to trim from left
	 * @return StringableInterface New stringable with left-trimmed content
	 */
	Ltrim(characters string) StringableInterface

	/**
	 * Rtrim right trims the string of the given characters.
	 *
	 * @param characters string Characters to trim from right
	 * @return StringableInterface New stringable with right-trimmed content
	 */
	Rtrim(characters string) StringableInterface

	/**
	 * Lcfirst makes a string's first character lowercase.
	 *
	 * @return StringableInterface New stringable with lowercase first character
	 */
	Lcfirst() StringableInterface

	/**
	 * Ucfirst makes a string's first character uppercase.
	 *
	 * @return StringableInterface New stringable with uppercase first character
	 */
	Ucfirst() StringableInterface

	/**
	 * Wrap wraps the string with the given strings.
	 *
	 * @param before string String to prepend
	 * @param after string String to append
	 * @return StringableInterface New stringable wrapped with given strings
	 */
	Wrap(before, after string) StringableInterface

	/**
	 * Unwrap unwraps the string with the given strings.
	 *
	 * @param before string String to remove from start
	 * @param after string String to remove from end
	 * @return StringableInterface New stringable with unwrapped content
	 */
	Unwrap(before, after string) StringableInterface

	/**
	 * String array and collection methods
	 */

	/**
	 * Explode explodes the string into a slice.
	 *
	 * @param delimiter string The delimiter to split on
	 * @param limit int Maximum number of pieces (0 for unlimited)
	 * @return []string Array of string pieces
	 */
	Explode(delimiter string, limit int) []string

	/**
	 * Split splits a string using a regular expression or by length.
	 *
	 * @param pattern interface{} Regex pattern (string) or length (int)
	 * @param limit int Maximum number of pieces
	 * @param flags int Regex flags (when using regex pattern)
	 * @return []string Array of split pieces
	 */
	Split(pattern interface{}, limit int, flags int) []string

	/**
	 * Ucsplit splits a string by uppercase characters.
	 *
	 * @return []string Array of pieces split by uppercase characters
	 */
	Ucsplit() []string

	/**
	 * Words limits the number of words in a string.
	 *
	 * @param words int Maximum number of words
	 * @param end string String to append when truncated
	 * @return StringableInterface New stringable with limited words
	 */
	Words(words int, end string) StringableInterface

	/**
	 * WordCount gets the number of words a string contains.
	 *
	 * @param characters string Additional word separator characters
	 * @return int Number of words in the string
	 */
	WordCount(characters string) int

	/**
	 * WordWrap wraps a string to a given number of characters.
	 *
	 * @param characters int Maximum line length
	 * @param breakStr string Line break string
	 * @param cutLongWords bool Whether to cut long words
	 * @return StringableInterface New stringable with wrapped lines
	 */
	WordWrap(characters int, breakStr string, cutLongWords bool) StringableInterface

	/**
	 * Output and conversion methods
	 */

	/**
	 * ToString gets the underlying string value.
	 *
	 * @return string The underlying string value
	 */
	ToString() string

	/**
	 * Value gets the underlying string value.
	 *
	 * @return string The underlying string value
	 */
	Value() string

	/**
	 * ToInteger gets the underlying string value as an integer.
	 *
	 * @param base int Numeric base for conversion
	 * @return int The parsed integer value
	 */
	ToInteger(base int) int

	/**
	 * ToFloat gets the underlying string value as a float.
	 *
	 * @return float64 The parsed floating point value
	 */
	ToFloat() float64

	/**
	 * ToBoolean gets the underlying string value as a boolean.
	 *
	 * @return bool The interpreted boolean value
	 */
	ToBoolean() bool

	/**
	 * Conditional execution methods
	 */

	/**
	 * When executes the given callback if the condition is true.
	 *
	 * @param condition bool The condition to check
	 * @param callback func(StringableInterface) StringableInterface Callback for true condition
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback for false condition
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	When(condition bool, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * Unless executes the given callback if the condition is false.
	 *
	 * @param condition bool The condition to check
	 * @param callback func(StringableInterface) StringableInterface Callback for false condition
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback for true condition
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	Unless(condition bool, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenContains executes the given callback if the string contains a given substring.
	 *
	 * @param needles []string Array of substrings to check for
	 * @param callback func(StringableInterface) StringableInterface Callback when condition is met
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when condition is not met
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenContains(needles []string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenContainsAll executes the given callback if the string contains all array values.
	 *
	 * @param needles []string Array of substrings that must all be present
	 * @param callback func(StringableInterface) StringableInterface Callback when condition is met
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when condition is not met
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenContainsAll(needles []string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenEmpty executes the given callback if the string is empty.
	 *
	 * @param callback func(StringableInterface) StringableInterface Callback when string is empty
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when string is not empty
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenEmpty(callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenNotEmpty executes the given callback if the string is not empty.
	 *
	 * @param callback func(StringableInterface) StringableInterface Callback when string is not empty
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when string is empty
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenNotEmpty(callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenStartsWith executes the given callback if the string starts with a given substring.
	 *
	 * @param needles []string Array of prefixes to check
	 * @param callback func(StringableInterface) StringableInterface Callback when condition is met
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when condition is not met
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenStartsWith(needles []string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenEndsWith executes the given callback if the string ends with a given substring.
	 *
	 * @param needles []string Array of suffixes to check
	 * @param callback func(StringableInterface) StringableInterface Callback when condition is met
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when condition is not met
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenEndsWith(needles []string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenExactly executes the given callback if the string is an exact match with the given value.
	 *
	 * @param value string The value to match exactly
	 * @param callback func(StringableInterface) StringableInterface Callback when condition is met
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when condition is not met
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenExactly(value string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenIs executes the given callback if the string matches a given pattern.
	 *
	 * @param pattern string The pattern to match (supports wildcards)
	 * @param callback func(StringableInterface) StringableInterface Callback when condition is met
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when condition is not met
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenIs(pattern string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenIsAscii executes the given callback if the string is 7 bit ASCII.
	 *
	 * @param callback func(StringableInterface) StringableInterface Callback when string is ASCII
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when string is not ASCII
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenIsAscii(callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenIsUuid executes the given callback if the string is a valid UUID.
	 *
	 * @param callback func(StringableInterface) StringableInterface Callback when string is UUID
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when string is not UUID
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenIsUuid(callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenIsUlid executes the given callback if the string is a valid ULID.
	 *
	 * @param callback func(StringableInterface) StringableInterface Callback when string is ULID
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when string is not ULID
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenIsUlid(callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * WhenTest executes the given callback if the string matches the given pattern.
	 *
	 * @param pattern string The regular expression pattern
	 * @param callback func(StringableInterface) StringableInterface Callback when pattern matches
	 * @param defaultCallback func(StringableInterface) StringableInterface Callback when pattern doesn't match
	 * @return StringableInterface Result of callback execution or original stringable
	 */
	WhenTest(pattern string, callback func(StringableInterface) StringableInterface, defaultCallback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * Utility and pipeline methods
	 */

	/**
	 * Pipe calls the given callback and returns a new string.
	 *
	 * @param callback func(StringableInterface) StringableInterface Transformation callback
	 * @return StringableInterface Result of the pipeline transformation
	 */
	Pipe(callback func(StringableInterface) StringableInterface) StringableInterface

	/**
	 * Tap calls the given callback with the string and returns the string.
	 *
	 * @param callback func(StringableInterface) Side-effect callback
	 * @return StringableInterface Original stringable for method chaining
	 */
	Tap(callback func(StringableInterface)) StringableInterface

	/**
	 * Dump dumps the string and returns the string.
	 *
	 * @return StringableInterface Original stringable for method chaining
	 */
	Dump() StringableInterface

	/**
	 * ClassBasename gets the basename of the class path.
	 *
	 * @return StringableInterface New stringable with class basename
	 */
	ClassBasename() StringableInterface

	/**
	 * Dirname gets the parent directory's path.
	 *
	 * @param levels int Number of parent levels to go up
	 * @return StringableInterface New stringable with directory path
	 */
	Dirname(levels int) StringableInterface

	/**
	 * Scan parses input from a string according to a format.
	 *
	 * @param format string Scanf-style format string
	 * @return []interface{} Array of parsed values
	 */
	Scan(format string) []interface{}

	/**
	 * Pluralization and language methods
	 */

	/**
	 * Plural gets the plural form of an English word.
	 *
	 * @param count int Count determining singular/plural
	 * @param prependCount bool Whether to prepend the count
	 * @return StringableInterface New stringable with plural form
	 */
	Plural(count int, prependCount bool) StringableInterface

	/**
	 * PluralStudly pluralizes the last word of an English, studly caps case string.
	 *
	 * @param count int Count determining singular/plural
	 * @return StringableInterface New stringable with studly plural form
	 */
	PluralStudly(count int) StringableInterface

	/**
	 * PluralPascal pluralizes the last word of an English, Pascal caps case string.
	 *
	 * @param count int Count determining singular/plural
	 * @return StringableInterface New stringable with Pascal plural form
	 */
	PluralPascal(count int) StringableInterface

	/**
	 * Singular gets the singular form of an English word.
	 *
	 * @return StringableInterface New stringable with singular form
	 */
	Singular() StringableInterface

	/**
	 * Slug generates a URL friendly "slug" from a given string.
	 *
	 * @param separator string Separator character for slug
	 * @param language string Language for slug generation rules
	 * @param dictionary map[string]string Custom character replacements
	 * @return StringableInterface New stringable with URL slug
	 */
	Slug(separator, language string, dictionary map[string]string) StringableInterface

	/**
	 * Pattern matching and extraction methods
	 */

	/**
	 * Match gets the string matching the given pattern.
	 *
	 * @param pattern string Regular expression pattern
	 * @return StringableInterface New stringable with first match
	 */
	Match(pattern string) StringableInterface

	/**
	 * MatchAll gets all strings matching the given pattern.
	 *
	 * @param pattern string Regular expression pattern
	 * @return []string Array of all matches
	 */
	MatchAll(pattern string) []string

	/**
	 * Numbers removes all non-numeric characters from a string.
	 *
	 * @return StringableInterface New stringable containing only numbers
	 */
	Numbers() StringableInterface

	/**
	 * ParseCallback parses a Class@method style callback into class and method.
	 *
	 * @param defaultMethod string Default method name if not specified
	 * @return []string Array containing [class, method]
	 */
	ParseCallback(defaultMethod string) []string

	/**
	 * Encoding and hashing methods
	 */

	/**
	 * ToBase64 converts the string to Base64 encoding.
	 *
	 * @return StringableInterface New stringable with Base64 encoded content
	 */
	ToBase64() StringableInterface

	/**
	 * FromBase64 decodes the Base64 encoded string.
	 *
	 * @param strict bool Whether to use strict decoding
	 * @return StringableInterface New stringable with decoded content
	 */
	FromBase64(strict bool) StringableInterface

	/**
	 * Hash hashes the string using the given algorithm.
	 *
	 * @param algorithm string Hash algorithm name (e.g., "md5", "sha1", "sha256")
	 * @return StringableInterface New stringable with hashed content
	 */
	Hash(algorithm string) StringableInterface

	/**
	 * Transliterate transliterates a string to its closest ASCII representation.
	 *
	 * @param unknown string Replacement for untranslatable characters
	 * @param strict bool Whether to use strict transliteration
	 * @return StringableInterface New stringable with transliterated content
	 */
	Transliterate(unknown string, strict bool) StringableInterface
}
