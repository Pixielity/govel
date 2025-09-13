package interfaces

// LocalizableInterface defines the contract for localization functionality.
type LocalizableInterface interface {
	// GetLocale returns the current locale setting
	GetLocale() string
	
	// SetLocale updates the locale setting
	SetLocale(locale string)
	
	// GetFallbackLocale returns the fallback locale setting
	GetFallbackLocale() string
	
	// SetFallbackLocale updates the fallback locale setting
	SetFallbackLocale(fallbackLocale string)
	
	// GetTimezone returns the current timezone setting
	GetTimezone() string
	
	// SetTimezone updates the timezone setting
	SetTimezone(timezone string)
	
	// IsLocale checks if the current locale matches the given locale
	IsLocale(locale string) bool
	
	// IsTimezone checks if the current timezone matches the given timezone
	IsTimezone(timezone string) bool
	
	// GetLocaleWithFallback returns the locale or fallback if locale is empty
	GetLocaleWithFallback() string
	
	// IsValidLocale checks if a locale string is valid
	IsValidLocale(locale string) bool
	
	// GetLanguageCode extracts the language code from current locale
	GetLanguageCode() string
	
	// GetCountryCode extracts the country code from current locale
	GetCountryCode() string
	
	// LocaleInfo returns comprehensive locale information
	LocaleInfo() map[string]string
	
	// SetLocaleInfo sets multiple locale properties at once
	SetLocaleInfo(locale, fallback, timezone string)
	
	// IsRTL returns whether the current locale uses right-to-left writing
	IsRTL() bool
	
	// GetTextDirection returns the text direction for the current locale
	GetTextDirection() string
}
