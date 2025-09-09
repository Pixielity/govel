package interfaces

/**
 * LocalizableInterface defines the comprehensive contract for components that support
 * full localization functionality. This interface goes beyond the basic Localizable
 * interface to provide advanced localization features including validation,
 * parsing, and RTL support.
 *
 * Following the Interface Segregation Principle (ISP), this interface focuses
 * specifically on comprehensive localization capabilities, allowing for clean
 * separation of concerns while providing rich i18n functionality.
 */
type LocalizableInterface interface {
	/**
	 * GetLocale returns the current locale setting.
	 *
	 * @return string The current locale identifier (e.g., "en-US", "fr-FR")
	 */
	GetLocale() string

	/**
	 * SetLocale updates the locale setting.
	 *
	 * @param locale string The new locale identifier to apply
	 */
	SetLocale(locale string)

	/**
	 * GetFallbackLocale returns the fallback locale setting.
	 *
	 * @return string The fallback locale identifier
	 */
	GetFallbackLocale() string

	/**
	 * SetFallbackLocale updates the fallback locale setting.
	 *
	 * @param fallbackLocale string The new fallback locale identifier
	 */
	SetFallbackLocale(fallbackLocale string)

	/**
	 * GetTimezone returns the current timezone setting.
	 *
	 * @return string The current timezone identifier
	 */
	GetTimezone() string

	/**
	 * SetTimezone updates the timezone setting.
	 *
	 * @param timezone string The new timezone identifier
	 */
	SetTimezone(timezone string)

	/**
	 * IsLocale checks if the current locale matches the given locale.
	 *
	 * @param locale string The locale to check against
	 * @return bool true if the current locale matches
	 */
	IsLocale(locale string) bool

	/**
	 * IsTimezone checks if the current timezone matches the given timezone.
	 *
	 * @param timezone string The timezone to check against
	 * @return bool true if the current timezone matches
	 */
	IsTimezone(timezone string) bool

	/**
	 * GetLocaleWithFallback returns the locale or fallback if locale is empty.
	 *
	 * @return string The current locale or fallback locale
	 */
	GetLocaleWithFallback() string

	/**
	 * IsValidLocale checks if a locale string is valid.
	 *
	 * @param locale string The locale string to validate
	 * @return bool true if locale appears valid
	 */
	IsValidLocale(locale string) bool

	/**
	 * GetLanguageCode extracts the language code from current locale.
	 *
	 * @return string The language code portion of the locale
	 */
	GetLanguageCode() string

	/**
	 * GetCountryCode extracts the country code from current locale.
	 *
	 * @return string The country code portion of the locale, or empty if none
	 */
	GetCountryCode() string

	/**
	 * LocaleInfo returns comprehensive locale information.
	 *
	 * @return map[string]string A map containing locale details
	 */
	LocaleInfo() map[string]string

	/**
	 * SetLocaleInfo sets multiple locale properties at once.
	 *
	 * @param locale string The primary locale to set
	 * @param fallback string The fallback locale to set
	 * @param timezone string The timezone to set
	 */
	SetLocaleInfo(locale, fallback, timezone string)

	/**
	 * IsRTL returns whether the current locale uses right-to-left writing.
	 *
	 * @return bool true if the locale uses RTL writing
	 */
	IsRTL() bool

	/**
	 * GetTextDirection returns the text direction for the current locale.
	 *
	 * @return string "rtl" for right-to-left languages, "ltr" for left-to-right
	 */
	GetTextDirection() string
}
