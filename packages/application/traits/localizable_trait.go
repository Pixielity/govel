package traits

import (
	"strings"
	"sync"

	traitInterfaces "govel/packages/application/interfaces/traits"
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
 * NewLocalizable creates a new Localizable instance with the specified locale settings.
 *
 * @param locale string The initial locale setting
 * @param fallbackLocale string The fallback locale setting
 * @param timezone string The timezone setting
 * @return *Localizable The newly created trait instance
 */
func NewLocalizable(locale, fallbackLocale, timezone string) *Localizable {
	return &Localizable{
		locale:         locale,
		fallbackLocale: fallbackLocale,
		timezone:       timezone,
	}
}

/**
 * GetLocale returns the current locale setting.
 * This method implements the LocalizableInterface contract for locale access.
 *
 * @return string The current locale identifier
 */
func (l *Localizable) GetLocale() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.locale
}

/**
 * SetLocale updates the locale setting.
 * This method implements the LocalizableInterface contract for locale modification.
 *
 * @param locale string The new locale identifier to apply
 */
func (l *Localizable) SetLocale(locale string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.IsValidLocale(locale) {
		l.locale = locale
	}
}

/**
 * GetFallbackLocale returns the fallback locale setting.
 *
 * @return string The fallback locale identifier
 */
func (l *Localizable) GetFallbackLocale() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.fallbackLocale
}

/**
 * SetFallbackLocale updates the fallback locale setting.
 *
 * @param fallbackLocale string The new fallback locale identifier
 */
func (l *Localizable) SetFallbackLocale(fallbackLocale string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.IsValidLocale(fallbackLocale) {
		l.fallbackLocale = fallbackLocale
	}
}

/**
 * GetTimezone returns the current timezone setting.
 *
 * @return string The current timezone identifier
 */
func (l *Localizable) GetTimezone() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.timezone
}

/**
 * SetTimezone updates the timezone setting.
 *
 * @param timezone string The new timezone identifier
 */
func (l *Localizable) SetTimezone(timezone string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.timezone = timezone
}

/**
 * IsLocale checks if the current locale matches the given locale.
 *
 * @param locale string The locale to check against
 * @return bool true if the current locale matches
 */
func (l *Localizable) IsLocale(locale string) bool {
	return l.GetLocale() == locale
}

/**
 * IsTimezone checks if the current timezone matches the given timezone.
 *
 * @param timezone string The timezone to check against
 * @return bool true if the current timezone matches
 */
func (l *Localizable) IsTimezone(timezone string) bool {
	return l.GetTimezone() == timezone
}

/**
 * GetLocaleWithFallback returns the locale or fallback if locale is empty.
 *
 * @return string The current locale or fallback locale
 */
func (l *Localizable) GetLocaleWithFallback() string {
	locale := l.GetLocale()
	if locale == "" {
		return l.GetFallbackLocale()
	}
	return locale
}

/**
 * IsValidLocale checks if a locale string is valid (basic validation).
 *
 * @param locale string The locale string to validate
 * @return bool true if locale appears valid
 */
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

/**
 * GetLanguageCode extracts the language code from current locale.
 *
 * @return string The language code portion of the locale
 */
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

/**
 * GetCountryCode extracts the country code from current locale.
 *
 * @return string The country code portion of the locale, or empty if none
 */
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

/**
 * LocaleInfo returns comprehensive locale information.
 *
 * @return map[string]string A map containing locale details
 */
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

/**
 * SetLocaleInfo sets multiple locale properties at once.
 *
 * @param locale string The primary locale to set
 * @param fallback string The fallback locale to set
 * @param timezone string The timezone to set
 */
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

/**
 * IsRTL returns whether the current locale uses right-to-left writing.
 *
 * @return bool true if the locale uses RTL writing
 */
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

/**
 * GetTextDirection returns the text direction for the current locale.
 *
 * @return string "rtl" for right-to-left languages, "ltr" for left-to-right
 */
func (l *Localizable) GetTextDirection() string {
	if l.IsRTL() {
		return "rtl"
	}
	return "ltr"
}

// Compile-time interface compliance check
// This ensures that Localizable implements the LocalizableInterface at compile time
// If Localizable doesn't implement all required methods, compilation will fail
var _ traitInterfaces.LocalizableInterface = (*Localizable)(nil)
