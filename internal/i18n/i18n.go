// Package i18n holds UI strings and exercise names for the 7 supported languages.
package i18n

import "strings"

// Lang is a BCP-47-ish language code we support.
type Lang string

const (
	RU Lang = "ru"
	EN Lang = "en"
	TR Lang = "tr"
	DE Lang = "de"
	ES Lang = "es"
	FR Lang = "fr"
	IT Lang = "it"
)

// Default language when nothing else is known.
const Default = EN

// Order is the canonical column order used by the translation tables.
var Order = []Lang{RU, EN, TR, DE, ES, FR, IT}

// Meta describes a language for the picker.
type Meta struct {
	Code Lang
	Name string // endonym
	// BCP47 is the tag passed to the browser SpeechSynthesis / lang attribute.
	BCP47 string
}

// Languages is the ordered list shown in the UI.
var Languages = []Meta{
	{RU, "Русский", "ru-RU"},
	{EN, "English", "en-US"},
	{TR, "Türkçe", "tr-TR"},
	{DE, "Deutsch", "de-DE"},
	{ES, "Español", "es-ES"},
	{FR, "Français", "fr-FR"},
	{IT, "Italiano", "it-IT"},
}

func index(l Lang) int {
	for i, x := range Order {
		if x == l {
			return i
		}
	}
	return -1
}

// Valid reports whether code is a supported language.
func Valid(code string) bool {
	return index(Lang(code)) >= 0
}

// Parse normalizes a code to a supported Lang, falling back to Default.
func Parse(code string) Lang {
	code = strings.ToLower(strings.TrimSpace(code))
	if i := strings.IndexAny(code, "-_"); i >= 0 {
		code = code[:i]
	}
	if index(Lang(code)) >= 0 {
		return Lang(code)
	}
	return Default
}

// FromAcceptLanguage picks the best supported language from an Accept-Language header.
func FromAcceptLanguage(header string) Lang {
	for _, part := range strings.Split(header, ",") {
		tag := part
		if i := strings.Index(tag, ";"); i >= 0 {
			tag = tag[:i]
		}
		tag = strings.ToLower(strings.TrimSpace(tag))
		if i := strings.IndexAny(tag, "-_"); i >= 0 {
			tag = tag[:i]
		}
		if index(Lang(tag)) >= 0 {
			return Lang(tag)
		}
	}
	return Default
}

// BCP47 returns the speech/lang tag for a language.
func BCP47(l Lang) string {
	for _, m := range Languages {
		if m.Code == l {
			return m.BCP47
		}
	}
	return "en-US"
}

// T looks up a UI string by key, falling back to English then the key itself.
func T(l Lang, key string) string {
	row, ok := ui[key]
	if !ok {
		return key
	}
	if i := index(l); i >= 0 && row[i] != "" {
		return row[i]
	}
	if row[1] != "" { // English column
		return row[1]
	}
	return key
}

// Exercise returns the localized name for an exercise ID.
func Exercise(l Lang, id string) string {
	row, ok := exerciseNames[id]
	if !ok {
		return id
	}
	if i := index(l); i >= 0 && row[i] != "" {
		return row[i]
	}
	return row[1]
}

// Cue returns the localized spoken cue text for a key (used by the voice engine).
func Cue(l Lang, key string) string {
	return T(l, "cue."+key)
}
