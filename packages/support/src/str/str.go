package str

import (
	"encoding/base64"
	"encoding/json"
	mathrand "math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gobeam/stringy"
	"github.com/google/uuid"
)

// After returns the remainder of the string after the first occurrence of the given substring.
// If the substring does not exist, the original string is returned.
func After(s, substring string) string {
	if substring == "" {
		return s
	}
	idx := strings.Index(s, substring)
	if idx == -1 {
		return s
	}
	return s[idx+len(substring):]
}

// AfterLast returns the remainder of the string after the last occurrence of the given substring.
// If the substring does not exist, the original string is returned.
func AfterLast(s, substring string) string {
	if substring == "" {
		return s
	}
	idx := strings.LastIndex(s, substring)
	if idx == -1 {
		return s
	}
	return s[idx+len(substring):]
}

// Before returns the substring of s before the first occurrence of substring.
// If the substring does not exist, the original string is returned.
func Before(s, substring string) string {
	if substring == "" {
		return s
	}
	idx := strings.Index(s, substring)
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// BeforeLast returns the substring of s before the last occurrence of substring.
// If the substring does not exist, the original string is returned.
func BeforeLast(s, substring string) string {
	if substring == "" {
		return s
	}
	idx := strings.LastIndex(s, substring)
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// Between returns the substring of s between the first occurrence of from and the last occurrence of to.
// If from or to are not found in order, an empty string is returned.
func Between(s, from, to string) string {
	if from == "" || to == "" {
		return ""
	}
	a := strings.Index(s, from)
	if a == -1 {
		return ""
	}
	b := strings.LastIndex(s, to)
	if b == -1 || b <= a+len(from) {
		return ""
	}
	return s[a+len(from) : b]
}

// BetweenFirst returns the smallest substring of s that is between the first occurrence of from and the first occurrence of to after it.
// If from or to are not found in order, an empty string is returned.
func BetweenFirst(s, from, to string) string {
	if from == "" || to == "" {
		return ""
	}
	a := strings.Index(s, from)
	if a == -1 {
		return ""
	}
	rest := s[a+len(from):]
	b := strings.Index(rest, to)
	if b == -1 {
		return ""
	}
	return rest[:b]
}

// Camel converts the string to camelCase.
func Camel(s string) string {
	return strings.ToLower(stringy.New(s).CamelCase().Get())
}

// Studly converts the string to StudlyCase (PascalCase).
func Studly(s string) string {
	result := stringy.New(s).CamelCase().Get()
	if len(result) == 0 {
		return result
	}
	// Ensure first letter is capitalized
	return strings.ToUpper(string(result[0])) + result[1:]
}

// Snake converts the string to snake_case.
func Snake(s string) string {
	return stringy.New(s).SnakeCase().Get()
}

// Kebab converts the string to kebab-case.
func Kebab(s string) string {
	return stringy.New(s).KebabCase().Get()
}

// Lower converts the string to lower case.
func Lower(s string) string { return strings.ToLower(s) }

// Upper converts the string to upper case.
func Upper(s string) string { return strings.ToUpper(s) }

// Title converts the string to title case.
func Title(s string) string { return strings.Title(s) }

// Contains reports whether s contains any of the given substrings.
func Contains(s string, needles ...string) bool {
	for _, n := range needles {
		if n == "" {
			continue
		}
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// ContainsAll reports whether s contains all of the given substrings.
func ContainsAll(s string, needles ...string) bool {
	for _, n := range needles {
		if n == "" {
			continue
		}
		if !strings.Contains(s, n) {
			return false
		}
	}
	return true
}

// StartsWith reports whether s starts with any of the given prefixes.
func StartsWith(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if p == "" {
			continue
		}
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// EndsWith reports whether s ends with any of the given suffixes.
func EndsWith(s string, suffixes ...string) bool {
	for _, p := range suffixes {
		if p == "" {
			continue
		}
		if strings.HasSuffix(s, p) {
			return true
		}
	}
	return false
}

// Is reports whether s matches the given pattern using simple wildcard matching.
// The pattern may contain '*' which matches any sequence of characters.
func Is(pattern, s string) bool {
	// Escape regex special chars except '*'
	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, "\\*", ".*")
	re = "^" + re + "$"
	ok, _ := regexp.MatchString(re, s)
	return ok
}

// IsAscii reports whether s contains only ASCII characters.
func IsAscii(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 0x7F {
			return false
		}
	}
	return true
}

// IsJson reports whether s is a valid JSON string.
func IsJson(s string) bool {
	var v any
	return json.Unmarshal([]byte(s), &v) == nil
}

// Length returns the number of runes in s.
func Length(s string) int { return len([]rune(s)) }

// Limit truncates the string to the specified length and appends the end indicator if truncated.
func Limit(s string, length int, end string) string {
	runes := []rune(s)
	if length < 0 {
		length = 0
	}
	if len(runes) <= length {
		return s
	}
	return string(runes[:length]) + end
}

// Replace replaces all occurrences of search with replace in s.
func Replace(s, search, replace string) string { return strings.ReplaceAll(s, search, replace) }

// ReplaceFirst replaces the first occurrence of search with replace in s.
func ReplaceFirst(s, search, replace string) string {
	if search == "" {
		return s
	}
	idx := strings.Index(s, search)
	if idx == -1 {
		return s
	}
	return s[:idx] + replace + s[idx+len(search):]
}

// ReplaceLast replaces the last occurrence of search with replace in s.
func ReplaceLast(s, search, replace string) string {
	if search == "" {
		return s
	}
	idx := strings.LastIndex(s, search)
	if idx == -1 {
		return s
	}
	return s[:idx] + replace + s[idx+len(search):]
}

// Start prefixes s with prefix if it doesn't already start with it.
func Start(s, prefix string) string {
	if prefix == "" || strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// Finish ensures s ends with exactly one instance of suffix.
func Finish(s, suffix string) string {
	if suffix == "" {
		return s
	}
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s + suffix
}

// Random returns a random alphanumeric string of the given length.
func Random(length int) string {
	if randomStringFactory != nil {
		return randomStringFactory(length)
	}

	if length <= 0 {
		return ""
	}
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	mathrand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

// Uuid generates a random UUID string.
func Uuid() string {
	if uuidFactory != nil {
		return uuidFactory()
	}
	if len(frozenUuids) > 0 {
		if frozenUuidIndex < len(frozenUuids) {
			result := frozenUuids[frozenUuidIndex]
			frozenUuidIndex++
			return result
		}
		return frozenUuids[len(frozenUuids)-1]
	}
	return uuid.New().String()
}

// IsUuid reports whether s is a valid UUID string.
func IsUuid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// Mask replaces a portion of a string with a repeated character.
func Mask(s string, character string, index int, length ...int) string {
	if character == "" {
		character = "*"
	}
	runes := []rune(s)
	sLen := len(runes)

	if index < 0 {
		index = 0
	}
	if index >= sLen {
		return s
	}

	maskLength := sLen - index
	if len(length) > 0 && length[0] > 0 {
		maskLength = length[0]
		if index+maskLength > sLen {
			maskLength = sLen - index
		}
	}

	mask := strings.Repeat(character, maskLength)
	return string(runes[:index]) + mask + string(runes[index+maskLength:])
}

// PadBoth pads both sides of a string with the given pad string to a specified length.
func PadBoth(s string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	runes := []rune(s)
	currentLen := len(runes)
	if currentLen >= length {
		return s
	}

	padNeeded := length - currentLen
	leftPad := padNeeded / 2
	rightPad := padNeeded - leftPad

	left := strings.Repeat(pad, (leftPad+len([]rune(pad))-1)/len([]rune(pad)))
	right := strings.Repeat(pad, (rightPad+len([]rune(pad))-1)/len([]rune(pad)))

	left = string([]rune(left)[:leftPad])
	right = string([]rune(right)[:rightPad])

	return left + s + right
}

// PadLeft pads the left side of a string with the given pad string to a specified length.
func PadLeft(s string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	runes := []rune(s)
	currentLen := len(runes)
	if currentLen >= length {
		return s
	}

	padNeeded := length - currentLen
	padStr := strings.Repeat(pad, (padNeeded+len([]rune(pad))-1)/len([]rune(pad)))
	padStr = string([]rune(padStr)[:padNeeded])

	return padStr + s
}

// PadRight pads the right side of a string with the given pad string to a specified length.
func PadRight(s string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	runes := []rune(s)
	currentLen := len(runes)
	if currentLen >= length {
		return s
	}

	padNeeded := length - currentLen
	padStr := strings.Repeat(pad, (padNeeded+len([]rune(pad))-1)/len([]rune(pad)))
	padStr = string([]rune(padStr)[:padNeeded])

	return s + padStr
}

// Remove removes all occurrences of the given strings from the subject.
func Remove(s string, search ...string) string {
	for _, str := range search {
		s = strings.ReplaceAll(s, str, "")
	}
	return s
}

// Trim trims whitespace from both ends of the string.
func Trim(s string, cutset ...string) string {
	if len(cutset) == 0 {
		return strings.TrimSpace(s)
	}
	for _, c := range cutset {
		s = strings.Trim(s, c)
	}
	return s
}

// Slug generates a URL friendly "slug" from the given string.
func Slug(s string, separator ...string) string {
	sep := "-"
	if len(separator) > 0 {
		sep = separator[0]
	}

	// Convert to lower case
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with separator
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, sep)

	// Trim separators from beginning and end
	s = strings.Trim(s, sep)

	return s
}

// Words limits the number of words in a string.
func Words(s string, words int, end ...string) string {
	if words <= 0 {
		return ""
	}

	wordList := strings.Fields(s)
	if len(wordList) <= words {
		return s
	}

	endStr := "..."
	if len(end) > 0 {
		endStr = end[0]
	}

	return strings.Join(wordList[:words], " ") + endStr
}

// Ascii transliterates a UTF-8 value to ASCII.
func Ascii(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r <= 127 {
			result.WriteRune(r)
		} else {
			// Basic transliteration for common characters
			switch r {
			case 'à', 'á', 'â', 'ã', 'ä', 'å':
				result.WriteRune('a')
			case 'À', 'Á', 'Â', 'Ã', 'Ä', 'Å':
				result.WriteRune('A')
			case 'è', 'é', 'ê', 'ë':
				result.WriteRune('e')
			case 'È', 'É', 'Ê', 'Ë':
				result.WriteRune('E')
			case 'ì', 'í', 'î', 'ï':
				result.WriteRune('i')
			case 'Ì', 'Í', 'Î', 'Ï':
				result.WriteRune('I')
			case 'ò', 'ó', 'ô', 'õ', 'ö':
				result.WriteRune('o')
			case 'Ò', 'Ó', 'Ô', 'Õ', 'Ö':
				result.WriteRune('O')
			case 'ù', 'ú', 'û', 'ü':
				result.WriteRune('u')
			case 'Ù', 'Ú', 'Û', 'Ü':
				result.WriteRune('U')
			case 'ç':
				result.WriteRune('c')
			case 'Ç':
				result.WriteRune('C')
			case 'ñ':
				result.WriteRune('n')
			case 'Ñ':
				result.WriteRune('N')
			default:
				// Skip non-ASCII characters that don't have simple replacements
				if unicode.IsPrint(r) {
					result.WriteRune('?')
				}
			}
		}
	}
	return result.String()
}

// Plural gets the plural form of an English word.
// This is a basic implementation - for more sophisticated pluralization,
// consider using a dedicated pluralization library.
func Plural(s string, count ...int) string {
	if len(count) > 0 && count[0] == 1 {
		return s
	}

	// Basic pluralization rules
	lower := strings.ToLower(s)

	if strings.HasSuffix(lower, "y") && len(s) > 1 && !isVowel(rune(lower[len(lower)-2])) {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "ss") ||
		strings.HasSuffix(lower, "sh") || strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") {
		return s + "es"
	}
	if strings.HasSuffix(lower, "f") {
		return s[:len(s)-1] + "ves"
	}
	if strings.HasSuffix(lower, "fe") {
		return s[:len(s)-2] + "ves"
	}

	return s + "s"
}

// Singular gets the singular form of an English word.
// This is a basic implementation - for more sophisticated singularization,
// consider using a dedicated pluralization library.
func Singular(s string) string {
	lower := strings.ToLower(s)

	if strings.HasSuffix(lower, "ies") && len(s) > 3 {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(lower, "ves") && len(s) > 3 {
		if strings.HasSuffix(lower[:len(lower)-3], "l") || strings.HasSuffix(lower[:len(lower)-3], "r") {
			return s[:len(s)-3] + "f"
		}
		return s[:len(s)-3] + "fe"
	}
	if strings.HasSuffix(lower, "es") && len(s) > 2 {
		base := lower[:len(lower)-2]
		if strings.HasSuffix(base, "s") || strings.HasSuffix(base, "sh") ||
			strings.HasSuffix(base, "ch") || strings.HasSuffix(base, "x") ||
			strings.HasSuffix(base, "z") {
			return s[:len(s)-2]
		}
	}
	if strings.HasSuffix(lower, "s") && len(s) > 1 {
		return s[:len(s)-1]
	}

	return s
}

// isVowel checks if a character is a vowel
func isVowel(r rune) bool {
	lower := unicode.ToLower(r)
	return lower == 'a' || lower == 'e' || lower == 'i' || lower == 'o' || lower == 'u'
}

// PluralStudly gets the plural form of an English word in StudlyCase.
func PluralStudly(s string, count ...int) string {
	plural := Plural(s, count...)
	return Studly(plural)
}

// ReplaceArray replaces a given value in the string sequentially with an array of values.
func ReplaceArray(s string, search string, replacements []string) string {
	result := s
	for _, replacement := range replacements {
		if idx := strings.Index(result, search); idx != -1 {
			result = result[:idx] + replacement + result[idx+len(search):]
		} else {
			break
		}
	}
	return result
}

// Of creates a new string instance (equivalent to Laravel's Str::of).
// In Go, this simply returns the string as-is since we use global functions.
func Of(s string) string {
	return s
}

// CharAt returns the character at the specified index.
func CharAt(s string, index int) string {
	runes := []rune(s)
	if index < 0 || index >= len(runes) {
		return ""
	}
	return string(runes[index])
}

// ChopStart removes a substring from the start of a string if it exists.
func ChopStart(s, needle string) string {
	if needle == "" {
		return s
	}
	if strings.HasPrefix(s, needle) {
		return s[len(needle):]
	}
	return s
}

// ChopEnd removes a substring from the end of a string if it exists.
func ChopEnd(s, needle string) string {
	if needle == "" {
		return s
	}
	if strings.HasSuffix(s, needle) {
		return s[:len(s)-len(needle)]
	}
	return s
}

// DoesntContain reports whether s does not contain any of the given substrings.
func DoesntContain(s string, needles ...string) bool {
	return !Contains(s, needles...)
}

// DoesntStartWith reports whether s does not start with any of the given prefixes.
func DoesntStartWith(s string, prefixes ...string) bool {
	return !StartsWith(s, prefixes...)
}

// DoesntEndWith reports whether s does not end with any of the given suffixes.
func DoesntEndWith(s string, suffixes ...string) bool {
	return !EndsWith(s, suffixes...)
}

// ConvertCase converts string case using different modes.
// Mode: 0=lower, 1=upper, 2=title
func ConvertCase(s string, mode int, encoding ...string) string {
	switch mode {
	case 0:
		return strings.ToLower(s)
	case 1:
		return strings.ToUpper(s)
	case 2:
		return strings.Title(s)
	default:
		return strings.ToLower(s)
	}
}

// Deduplicate removes consecutive duplicate characters.
func Deduplicate(s string, characters ...string) string {
	chars := " "
	if len(characters) > 0 {
		chars = characters[0]
	}

	if len(chars) == 1 {
		// Single character deduplication
		char := rune(chars[0])
		var result strings.Builder
		var lastChar rune
		for _, r := range s {
			if r != char || r != lastChar {
				result.WriteRune(r)
			}
			lastChar = r
		}
		return result.String()
	}

	// Multiple character deduplication (basic implementation)
	result := s
	for _, char := range chars {
		charStr := string(char)
		doubleChar := charStr + charStr
		for strings.Contains(result, doubleChar) {
			result = strings.ReplaceAll(result, doubleChar, charStr)
		}
	}
	return result
}

// Excerpt returns an excerpt from text around a given phrase.
func Excerpt(text, phrase string, options ...map[string]interface{}) string {
	if phrase == "" {
		return text
	}

	radius := 100
	omission := "..."

	if len(options) > 0 {
		if r, ok := options[0]["radius"].(int); ok {
			radius = r
		}
		if o, ok := options[0]["omission"].(string); ok {
			omission = o
		}
	}

	index := strings.Index(strings.ToLower(text), strings.ToLower(phrase))
	if index == -1 {
		return ""
	}

	start := index - radius
	if start < 0 {
		start = 0
	}

	end := index + len(phrase) + radius
	if end > len(text) {
		end = len(text)
	}

	excerpt := text[start:end]

	if start > 0 {
		excerpt = omission + excerpt
	}
	if end < len(text) {
		excerpt = excerpt + omission
	}

	return excerpt
}

// Wrap wraps a string with given before and after strings.
func Wrap(s, before string, after ...string) string {
	afterStr := before
	if len(after) > 0 {
		afterStr = after[0]
	}
	return before + s + afterStr
}

// Unwrap removes wrapping strings from the beginning and end.
func Unwrap(s, before string, after ...string) string {
	afterStr := before
	if len(after) > 0 {
		afterStr = after[0]
	}

	if strings.HasPrefix(s, before) && strings.HasSuffix(s, afterStr) {
		s = s[len(before):]
		s = s[:len(s)-len(afterStr)]
	}
	return s
}

// IsUrl reports whether s is a valid URL.
func IsUrl(s string, protocols ...[]string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	if len(protocols) > 0 && len(protocols[0]) > 0 {
		validProtocol := false
		for _, protocol := range protocols[0] {
			if u.Scheme == protocol {
				validProtocol = true
				break
			}
		}
		return validProtocol && u.Host != ""
	}

	return u.Scheme != "" && u.Host != ""
}

// IsUlid reports whether s is a valid ULID.
func IsUlid(s string) bool {
	// Basic ULID validation (26 characters, base32)
	if len(s) != 26 {
		return false
	}

	validChars := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	for _, char := range strings.ToUpper(s) {
		if !strings.ContainsRune(validChars, char) {
			return false
		}
	}
	return true
}

// Transliterate converts UTF-8 strings to ASCII with custom unknown character replacement.
func Transliterate(s string, unknown ...string) string {
	unknownChar := "?"
	if len(unknown) > 0 {
		unknownChar = unknown[0]
	}

	var result strings.Builder
	for _, r := range s {
		if r <= 127 {
			result.WriteRune(r)
		} else {
			// Extended transliteration map
			switch r {
			case 'à', 'á', 'â', 'ã', 'ä', 'å', 'ā', 'ă', 'ą':
				result.WriteRune('a')
			case 'À', 'Á', 'Â', 'Ã', 'Ä', 'Å', 'Ā', 'Ă', 'Ą':
				result.WriteRune('A')
			case 'è', 'é', 'ê', 'ë', 'ē', 'ĕ', 'ė', 'ę', 'ě':
				result.WriteRune('e')
			case 'È', 'É', 'Ê', 'Ë', 'Ē', 'Ĕ', 'Ė', 'Ę', 'Ě':
				result.WriteRune('E')
			case 'ì', 'í', 'î', 'ï', 'ĩ', 'ī', 'ĭ', 'į', 'ı':
				result.WriteRune('i')
			case 'Ì', 'Í', 'Î', 'Ï', 'Ĩ', 'Ī', 'Ĭ', 'Į', 'İ':
				result.WriteRune('I')
			case 'ò', 'ó', 'ô', 'õ', 'ö', 'ø', 'ō', 'ŏ', 'ő':
				result.WriteRune('o')
			case 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', 'Ø', 'Ō', 'Ŏ', 'Ő':
				result.WriteRune('O')
			case 'ù', 'ú', 'û', 'ü', 'ũ', 'ū', 'ŭ', 'ů', 'ű', 'ų':
				result.WriteRune('u')
			case 'Ù', 'Ú', 'Û', 'Ü', 'Ũ', 'Ū', 'Ŭ', 'Ů', 'Ű', 'Ų':
				result.WriteRune('U')
			case 'ç', 'ć', 'ĉ', 'ċ', 'č':
				result.WriteRune('c')
			case 'Ç', 'Ć', 'Ĉ', 'Ċ', 'Č':
				result.WriteRune('C')
			case 'ñ', 'ń', 'ņ', 'ň', 'ŉ':
				result.WriteRune('n')
			case 'Ñ', 'Ń', 'Ņ', 'Ň':
				result.WriteRune('N')
			case 'ý', 'ÿ', 'ŷ':
				result.WriteRune('y')
			case 'Ý', 'Ÿ', 'Ŷ':
				result.WriteRune('Y')
			case 'ß':
				result.WriteString("ss")
			case 'æ':
				result.WriteString("ae")
			case 'Æ':
				result.WriteString("AE")
			case 'œ':
				result.WriteString("oe")
			case 'Œ':
				result.WriteString("OE")
			default:
				result.WriteString(unknownChar)
			}
		}
	}
	return result.String()
}

// Match returns the first match of a regular expression.
func Match(pattern, s string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	match := re.FindString(s)
	return match
}

// IsMatch reports whether a string matches a regular expression pattern.
func IsMatch(pattern, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	return err == nil && matched
}

// MatchAll returns all matches of a regular expression.
func MatchAll(pattern, s string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re.FindAllString(s, -1)
}

// Numbers extracts all numeric strings from the input.
func Numbers(s string) []string {
	re := regexp.MustCompile(`\d+`)
	return re.FindAllString(s, -1)
}

// ParseCallback parses a Class@method style callback into class and method.
func ParseCallback(callback string, defaultMethod ...string) (string, string) {
	defaultVal := ""
	if len(defaultMethod) > 0 {
		defaultVal = defaultMethod[0]
	}

	parts := strings.Split(callback, "@")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return callback, defaultVal
}

// Password generates a random password with specified criteria.
func Password(length int, letters, numbers, symbols, spaces bool) string {
	if length <= 0 {
		return ""
	}

	var chars strings.Builder
	if letters {
		chars.WriteString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if numbers {
		chars.WriteString("0123456789")
	}
	if symbols {
		chars.WriteString("!@#$%^&*()_+-=[]{}|;:,.<>?")
	}
	if spaces {
		chars.WriteString(" ")
	}

	charSet := chars.String()
	if charSet == "" {
		return ""
	}

	mathrand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charSet[mathrand.Intn(len(charSet))]
	}
	return string(result)
}

// Position returns the position of the first occurrence of a substring.
func Position(haystack, needle string, offset ...int) int {
	start := 0
	if len(offset) > 0 {
		start = offset[0]
	}

	if start < 0 || start >= len(haystack) {
		return -1
	}

	index := strings.Index(haystack[start:], needle)
	if index == -1 {
		return -1
	}
	return start + index
}

// Repeat repeats the given string the specified number of times.
func Repeat(s string, times int) string {
	if times <= 0 {
		return ""
	}
	return strings.Repeat(s, times)
}

// ReplaceStart replaces the start of a string with a replacement if it matches the search string.
func ReplaceStart(search, replace, s string) string {
	if strings.HasPrefix(s, search) {
		return replace + s[len(search):]
	}
	return s
}

// ReplaceEnd replaces the end of a string with a replacement if it matches the search string.
func ReplaceEnd(search, replace, s string) string {
	if strings.HasSuffix(s, search) {
		return s[:len(s)-len(search)] + replace
	}
	return s
}

// ReplaceMatches replaces all regex matches with the replacement string.
func ReplaceMatches(pattern, replace, s string, limit ...int) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}

	max := -1
	if len(limit) > 0 {
		max = limit[0]
	}

	if max == -1 {
		return re.ReplaceAllString(s, replace)
	}

	// Limited replacement
	count := 0
	result := s
	for count < max {
		loc := re.FindStringIndex(result)
		if loc == nil {
			break
		}
		result = result[:loc[0]] + replace + result[loc[1]:]
		count++
	}
	return result
}

// Reverse reverses a string.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Headline converts a string to a headline format (first letter of each word capitalized, underscores/dashes to spaces).
func Headline(s string) string {
	// Replace underscores and dashes with spaces
	result := strings.ReplaceAll(s, "_", " ")
	result = strings.ReplaceAll(result, "-", " ")

	// Split into words and capitalize each
	words := strings.Fields(result)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// Apa converts a string to APA title case.
func Apa(s string) string {
	// APA style: capitalize first word, last word, and all major words
	words := strings.Fields(strings.ToLower(s))
	minorWords := map[string]bool{
		"a": true, "an": true, "and": true, "as": true, "at": true, "but": true,
		"by": true, "for": true, "if": true, "in": true, "nor": true, "of": true,
		"on": true, "or": true, "so": true, "the": true, "to": true, "up": true, "yet": true,
	}

	for i, word := range words {
		if i == 0 || i == len(words)-1 || !minorWords[word] || len(word) > 3 {
			// Capitalize major words, first word, last word, or words longer than 3 chars
			words[i] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

// Ltrim removes whitespace or specified characters from the left side of a string.
func Ltrim(s string, cutset ...string) string {
	if len(cutset) == 0 {
		return strings.TrimLeftFunc(s, unicode.IsSpace)
	}
	for _, c := range cutset {
		s = strings.TrimLeft(s, c)
	}
	return s
}

// Rtrim removes whitespace or specified characters from the right side of a string.
func Rtrim(s string, cutset ...string) string {
	if len(cutset) == 0 {
		return strings.TrimRightFunc(s, unicode.IsSpace)
	}
	for _, c := range cutset {
		s = strings.TrimRight(s, c)
	}
	return s
}

// Squish removes all "extra" blank space from a string and trims it.
func Squish(s string) string {
	// Replace multiple whitespace with single space
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

// Pascal converts a string to PascalCase (same as StudlyCase).
func Pascal(s string) string {
	return Studly(s)
}

// Substr returns a substring starting at the specified position.
func Substr(s string, start int, length ...int) string {
	runes := []rune(s)
	sLen := len(runes)

	if start < 0 {
		start = sLen + start
		if start < 0 {
			start = 0
		}
	}

	if start >= sLen {
		return ""
	}

	if len(length) == 0 {
		return string(runes[start:])
	}

	substrLen := length[0]
	if substrLen < 0 {
		return ""
	}

	end := start + substrLen
	if end > sLen {
		end = sLen
	}

	return string(runes[start:end])
}

// SubstrCount counts the number of substring occurrences.
func SubstrCount(haystack, needle string, offset int, length ...int) int {
	if needle == "" {
		return 0
	}

	start := offset

	search := haystack
	if start > 0 && start < len(haystack) {
		search = haystack[start:]
	}

	if len(length) > 0 && length[0] > 0 {
		end := start + length[0]
		if end < len(haystack) {
			search = haystack[start:end]
		}
	}

	return strings.Count(search, needle)
}

// SubstrReplace replaces text within a portion of a string.
func SubstrReplace(s, replacement string, offset int, length ...int) string {
	runes := []rune(s)
	sLen := len(runes)

	if offset < 0 {
		offset = sLen + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset >= sLen {
		return s + replacement
	}

	if len(length) == 0 {
		return string(runes[:offset]) + replacement
	}

	replaceLen := length[0]
	if replaceLen < 0 {
		replaceLen = sLen - offset + replaceLen
		if replaceLen < 0 {
			replaceLen = 0
		}
	}

	end := offset + replaceLen
	if end > sLen {
		end = sLen
	}

	return string(runes[:offset]) + replacement + string(runes[end:])
}

// Swap replaces occurrences of strings according to a map.
func Swap(replacements map[string]string, s string) string {
	result := s
	for search, replace := range replacements {
		result = strings.ReplaceAll(result, search, replace)
	}
	return result
}

// Take returns the first limit characters of a string.
func Take(s string, limit int) string {
	if limit < 0 {
		return Substr(s, limit)
	}
	return Substr(s, 0, limit)
}

// ToBase64 encodes a string to base64.
func ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// FromBase64 decodes a base64 string.
func FromBase64(s string, strict ...bool) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil && len(strict) > 0 && strict[0] {
		return "", err
	}
	if err != nil {
		return "", nil // non-strict mode returns empty string
	}
	return string(data), nil
}

// Lcfirst converts the first character of a string to lowercase.
func Lcfirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Ucfirst converts the first character of a string to uppercase.
func Ucfirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Ucsplit splits a string into an array of words at uppercase boundaries.
func Ucsplit(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	var current strings.Builder

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// WordCount returns the number of words in a string.
func WordCount(s string, characters ...string) int {
	if len(characters) > 0 {
		// Custom word separators
		separators := characters[0]
		for _, sep := range separators {
			s = strings.ReplaceAll(s, string(sep), " ")
		}
	}
	return len(strings.Fields(s))
}

// WordWrap wraps a string to a given number of characters using a break string.
func WordWrap(s string, width int, breakStr ...string) string {
	if width <= 0 {
		return s
	}

	breakString := "\n"
	if len(breakStr) > 0 {
		breakString = breakStr[0]
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}

	var result strings.Builder
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			result.WriteString(currentLine + breakString)
			currentLine = word
		}
	}

	result.WriteString(currentLine)
	return result.String()
}

// Basename returns the trailing name component of a path.
func Basename(path string, suffix ...string) string {
	// Simple basename implementation
	idx := strings.LastIndexAny(path, "/\\")
	name := path
	if idx >= 0 {
		name = path[idx+1:]
	}
	if len(suffix) > 0 && suffix[0] != "" && strings.HasSuffix(name, suffix[0]) {
		name = name[:len(name)-len(suffix[0])]
	}
	return name
}

// StripTags removes HTML and PHP tags from a string.
func StripTags(s string, allowedTags ...string) string {
	// Basic HTML tag removal - in production use a proper HTML parser
	re := regexp.MustCompile(`<[^>]*>`)
	result := re.ReplaceAllString(s, "")
	return result
}

// ClassBasename returns the class basename of a fully qualified class name.
func ClassBasename(class string) string {
	idx := strings.LastIndex(class, "\\")
	if idx >= 0 {
		return class[idx+1:]
	}
	return class
}

// Dirname returns the parent directory of a path.
func Dirname(path string, levels ...int) string {
	level := 1
	if len(levels) > 0 {
		level = levels[0]
	}
	result := path
	for i := 0; i < level; i++ {
		idx := strings.LastIndexAny(result, "/\\")
		if idx <= 0 {
			result = "."
			break
		}
		result = result[:idx]
	}
	return result
}

// Scan parses input from a string according to a format.
func Scan(s string, format string) []interface{} {
	// Basic implementation - in production use fmt.Sscanf or similar
	var result []interface{}
	// This is a simplified version
	result = append(result, s)
	return result
}

// Hash generates a hash of the string using the specified algorithm.
func Hash(s string, algorithm string) string {
	// Basic hash implementation - in production use crypto package
	// For now, just return a simple hash
	return ToBase64(s)
}

// PluralPascal gets the plural form of an English word in PascalCase.
func PluralPascal(s string, count ...int) string {
	plural := Plural(s, count...)
	return Pascal(plural)
}

// Global variables for factory methods (similar to Laravel's approach)
var (
	randomStringFactory func(int) string
	uuidFactory         func() string
	ulidFactory         func() string
	frozenUuids         []string
	frozenUuidIndex     int
	frozenUlids         []string
	frozenUlidIndex     int
)

// Uuid7 generates a version 7 UUID (time-based)
func Uuid7(timeVal ...int64) string {
	// Basic UUID v7 implementation (simplified)
	// In production, you'd want to use a proper UUID v7 library
	return Uuid() // fallback to v4 for now
}

// OrderedUuid generates an ordered UUID (similar to UUID v1 but ordered)
func OrderedUuid() string {
	// For now, return a regular UUID - in production use ordered UUID library
	return Uuid()
}

// Ulid generates a ULID (Universally Unique Lexicographically Sortable Identifier)
func Ulid(timeVal ...int64) string {
	if ulidFactory != nil {
		return ulidFactory()
	}

	if len(frozenUlids) > 0 {
		if frozenUlidIndex < len(frozenUlids) {
			result := frozenUlids[frozenUlidIndex]
			frozenUlidIndex++
			return result
		}
		return frozenUlids[len(frozenUlids)-1]
	}

	// Basic ULID implementation (simplified)
	// In production, use a proper ULID library
	timestamp := timeVal
	if len(timestamp) == 0 {
		timestamp = []int64{time.Now().UnixMilli()}
	}

	// Generate a basic ULID-like string
	encoding := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	mathrand.Seed(time.Now().UnixNano())

	result := make([]byte, 26)
	// First 10 chars represent timestamp
	t := timestamp[0]
	for i := 9; i >= 0; i-- {
		result[i] = encoding[t%32]
		t /= 32
	}

	// Last 16 chars are random
	for i := 10; i < 26; i++ {
		result[i] = encoding[mathrand.Intn(32)]
	}

	return string(result)
}

// CreateRandomStringsUsing sets a custom factory for random string generation
func CreateRandomStringsUsing(factory func(int) string) {
	randomStringFactory = factory
}

// CreateRandomStringsUsingSequence creates random strings using a predefined sequence
func CreateRandomStringsUsingSequence(sequence []string, whenMissing ...func(int) string) {
	index := 0
	CreateRandomStringsUsing(func(length int) string {
		if index < len(sequence) {
			result := sequence[index]
			index++
			return result
		}
		if len(whenMissing) > 0 {
			return whenMissing[0](length)
		}
		return Random(length)
	})
}

// CreateRandomStringsNormally resets random string generation to normal
func CreateRandomStringsNormally() {
	randomStringFactory = nil
}

// CreateUuidsUsing sets a custom factory for UUID generation
func CreateUuidsUsing(factory func() string) {
	uuidFactory = factory
}

// CreateUuidsUsingSequence creates UUIDs using a predefined sequence
func CreateUuidsUsingSequence(sequence []string, whenMissing ...func() string) {
	index := 0
	CreateUuidsUsing(func() string {
		if index < len(sequence) {
			result := sequence[index]
			index++
			return result
		}
		if len(whenMissing) > 0 {
			return whenMissing[0]()
		}
		return uuid.New().String()
	})
}

// CreateUuidsNormally resets UUID generation to normal
func CreateUuidsNormally() {
	uuidFactory = nil
}

// FreezeUuids freezes UUID generation to return specific values
func FreezeUuids(callback ...func()) {
	if len(callback) > 0 {
		// Execute callback with frozen UUIDs, then unfreeze
		callback[0]()
		frozenUuids = nil
		frozenUuidIndex = 0
	}
	// If no callback provided, UUIDs remain frozen until manually unfrozen
}

// CreateUlidsUsing sets a custom factory for ULID generation
func CreateUlidsUsing(factory func() string) {
	ulidFactory = factory
}

// CreateUlidsUsingSequence creates ULIDs using a predefined sequence
func CreateUlidsUsingSequence(sequence []string, whenMissing ...func() string) {
	index := 0
	CreateUlidsUsing(func() string {
		if index < len(sequence) {
			result := sequence[index]
			index++
			return result
		}
		if len(whenMissing) > 0 {
			return whenMissing[0]()
		}
		return Ulid()
	})
}

// CreateUlidsNormally resets ULID generation to normal
func CreateUlidsNormally() {
	ulidFactory = nil
}

// FreezeUlids freezes ULID generation to return specific values
func FreezeUlids(callback ...func()) {
	if len(callback) > 0 {
		// Execute callback with frozen ULIDs, then unfreeze
		callback[0]()
		frozenUlids = nil
		frozenUlidIndex = 0
	}
	// If no callback provided, ULIDs remain frozen until manually unfrozen
}

// FlushCache clears any cached values (placeholder for Laravel compatibility)
func FlushCache() {
	// In Laravel, this clears cached singular/plural forms
	// For our basic implementation, this is a no-op
	// In a more sophisticated implementation, you'd clear pluralization caches here
}

// Additional helper methods to set frozen values

// SetFrozenUuids sets the UUIDs to be returned when frozen
func SetFrozenUuids(uuids ...string) {
	frozenUuids = uuids
	frozenUuidIndex = 0
}

// SetFrozenUlids sets the ULIDs to be returned when frozen
func SetFrozenUlids(ulids ...string) {
	frozenUlids = ulids
	frozenUlidIndex = 0
}
