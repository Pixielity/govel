package traits

import (
	"strings"
	"sync"

	"govel/application/helpers"
	traitInterfaces "govel/types/src/interfaces/application/traits"
)

/**
 * Localizable provides localization functionality through composition.
 * This struct implements the LocalizableInterface and contains its own locale data,
 * following the self-contained trait pattern.
 *
 * Unlike dependency injection, this trait owns and manages its own state
 * for locale, fallback locale, and timezone settings.
 */
type Localizable struct {
	locale         string       // Current locale (e.g., "en-US")
	fallbackLocale string       // Fallback locale (e.g., "en")
	timezone       string       // Timezone (e.g., "UTC")
	mutex          sync.RWMutex // Thread safety for trait operations
}

/**
 * NewLocalizable creates a new localizable trait with optional parameters.
 * If values are not provided, they will be read from environment variables.
 *
 * Parameters:
 *   options[0]: Optional locale setting (string)
 *   options[1]: Optional fallback locale setting (string)
 *   options[2]: Optional timezone setting (string)
 *
 * Returns:
 *   *Localizable: The newly created trait instance
 *
 * Example:
 *   // Using environment variables
 *   localizable := NewLocalizable()
 *   // Providing explicit values (locale="en-US", fallback="en", timezone="UTC")
 *   localizable := NewLocalizable("en-US", "en", "UTC")
 */
func NewLocalizable(options ...string) *Localizable {
	envHelper := helpers.NewEnvHelper()

	// Use provided options or fallback to environment variables
	appLocale := envHelper.GetAppLocale()                 // Default from environment
	appFallbackLocale := envHelper.GetAppFallbackLocale() // Default from environment
	appTimezone := envHelper.GetAppTimezone()             // Default from environment

	// If options are provided:
	// First is locale
	// Second is fallback locale
	// Third is timezone
	if len(options) > 0 && options[0] != "" {
		appLocale = options[0]
	}
	if len(options) > 1 && options[1] != "" {
		appFallbackLocale = options[1]
	}
	if len(options) > 2 && options[2] != "" {
		appTimezone = options[2]
	}

	return &Localizable{
		locale:         appLocale,
		fallbackLocale: appFallbackLocale,
		timezone:       appTimezone,
	}
}

// GetLocale returns the current locale setting.
//
// Returns:
//
//	string: The current locale identifier (e.g., "en-US", "fr-FR")
func (l *Localizable) GetLocale() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.locale
}

// SetLocale updates the locale setting.
//
// Parameters:
//
//	locale: The new locale identifier to apply
func (l *Localizable) SetLocale(locale string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.IsValidLocale(locale) {
		l.locale = locale
	}
}

// GetFallbackLocale returns the fallback locale setting.
//
// Returns:
//
//	string: The fallback locale identifier
func (l *Localizable) GetFallbackLocale() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.fallbackLocale
}

// SetFallbackLocale updates the fallback locale setting.
//
// Parameters:
//
//	fallbackLocale: The new fallback locale identifier
func (l *Localizable) SetFallbackLocale(fallbackLocale string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.IsValidLocale(fallbackLocale) {
		l.fallbackLocale = fallbackLocale
	}
}

// GetTimezone returns the current timezone setting.
//
// Returns:
//
//	string: The current timezone identifier
func (l *Localizable) GetTimezone() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.timezone
}

// SetTimezone updates the timezone setting.
//
// Parameters:
//
//	timezone: The new timezone identifier
func (l *Localizable) SetTimezone(timezone string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.timezone = timezone
}

// IsLocale checks if the current locale matches the given locale.
//
// Parameters:
//
//	locale: The locale to check against
//
// Returns:
//
//	bool: true if the current locale matches
func (l *Localizable) IsLocale(locale string) bool {
	return l.GetLocale() == locale
}

// IsTimezone checks if the current timezone matches the given timezone.
//
// Parameters:
//
//	timezone: The timezone to check against
//
// Returns:
//
//	bool: true if the current timezone matches
func (l *Localizable) IsTimezone(timezone string) bool {
	return l.GetTimezone() == timezone
}

// GetLocaleWithFallback returns the locale or fallback if locale is empty.
//
// Returns:
//
//	string: The current locale or fallback locale
func (l *Localizable) GetLocaleWithFallback() string {
	locale := l.GetLocale()
	if locale == "" {
		return l.GetFallbackLocale()
	}
	return locale
}

// IsValidLocale checks if a locale string is valid.
//
// Parameters:
//
//	locale: The locale string to validate
//
// Returns:
//
//	bool: true if locale appears valid
func (l *Localizable) IsValidLocale(locale string) bool {
	if locale == "" {
		return false
	}

	// Support common locale patterns: "en", "en-US", "en_US"
	parts := strings.Split(strings.Replace(locale, "_", "-", -1), "-")

	// Must have at least language code
	if len(parts) == 0 || len(parts[0]) < 2 {
		return false
	}

	// Language code should be 2-3 characters
	langCode := strings.ToLower(parts[0])
	if len(langCode) < 2 || len(langCode) > 3 {
		return false
	}

	// If country code is provided, it should be 2 characters
	if len(parts) > 1 {
		countryCode := strings.ToUpper(parts[1])
		if len(countryCode) != 2 {
			return false
		}
	}

	return true
}

// GetLanguageCode extracts the language code from current locale.
//
// Returns:
//
//	string: The language code portion of the locale
func (l *Localizable) GetLanguageCode() string {
	locale := l.GetLocale()
	if len(locale) >= 2 {
		if len(locale) == 2 {
			return locale
		}
		// Look for separator like "-" or "_"
		for i, char := range locale {
			if char == '-' || char == '_' {
				return locale[:i]
			}
		}
	}
	return locale
}

// GetCountryCode extracts the country code from current locale.
//
// Returns:
//
//	string: The country code portion of the locale, or empty if none
func (l *Localizable) GetCountryCode() string {
	locale := l.GetLocale()
	// Look for separator like "-" or "_"
	for i, char := range locale {
		if char == '-' || char == '_' {
			if i+1 < len(locale) {
				return locale[i+1:]
			}
		}
	}
	return ""
}

// LocaleInfo returns comprehensive locale information.
//
// Returns:
//
//	map[string]string: A map containing locale details
func (l *Localizable) LocaleInfo() map[string]string {
	return map[string]string{
		"locale":      l.GetLocale(),
		"fallback":    l.GetFallbackLocale(),
		"timezone":    l.GetTimezone(),
		"language":    l.GetLanguageCode(),
		"country":     l.GetCountryCode(),
		"safe_locale": l.GetLocaleWithFallback(),
	}
}

// SetLocaleInfo sets multiple locale properties at once.
//
// Parameters:
//
//	locale: The primary locale to set
//	fallback: The fallback locale to set
//	timezone: The timezone to set
func (l *Localizable) SetLocaleInfo(locale, fallback, timezone string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.IsValidLocale(locale) {
		l.locale = locale
	}
	if l.IsValidLocale(fallback) {
		l.fallbackLocale = fallback
	}
	l.timezone = timezone
}

// IsRTL returns whether the current locale uses right-to-left writing.
//
// Returns:
//
//	bool: true if the locale uses RTL writing
func (l *Localizable) IsRTL() bool {
	language := l.GetLanguageCode()
	rtlLanguages := []string{"ar", "he", "fa", "ur", "yi"}

	for _, rtlLang := range rtlLanguages {
		if language == rtlLang {
			return true
		}
	}
	return false
}

// GetTextDirection returns the text direction for the current locale.
//
// Returns:
//
//	string: "rtl" for right-to-left languages, "ltr" for left-to-right
func (l *Localizable) GetTextDirection() string {
	if l.IsRTL() {
		return "rtl"
	}
	return "ltr"
}

// Compile-time interface compliance check
var _ traitInterfaces.LocalizableInterface = (*Localizable)(nil)
