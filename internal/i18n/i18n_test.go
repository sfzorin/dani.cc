package i18n

import (
	"morning/internal/content"
	"testing"
)

// Every UI string must at least have an English value (the fallback column).
func TestUIStringsHaveEnglish(t *testing.T) {
	enCol := index(EN)
	for key, row := range ui {
		if row[enCol] == "" {
			t.Errorf("ui key %q has no English value", key)
		}
	}
}

// Every exercise in the pool must have a localized name in all 7 languages.
func TestExerciseNamesCoverPool(t *testing.T) {
	for id := range content.Pool {
		row, ok := exerciseNames[id]
		if !ok {
			t.Errorf("exercise %q has no translations", id)
			continue
		}
		for i, name := range row {
			if name == "" {
				t.Errorf("exercise %q missing translation for %s", id, Order[i])
			}
		}
	}
}

func TestParseAndFallback(t *testing.T) {
	if Parse("ru-RU") != RU {
		t.Errorf("Parse(ru-RU) != RU")
	}
	if Parse("xx") != Default {
		t.Errorf("Parse(xx) should fall back to default")
	}
	if FromAcceptLanguage("fr-FR,fr;q=0.9,en;q=0.8") != FR {
		t.Errorf("Accept-Language parse failed")
	}
	if T(RU, "nonexistent.key") != "nonexistent.key" {
		t.Errorf("unknown key should return the key")
	}
	if Month(RU, 6) == "" || Weekday(EN, 0) != "Mon" {
		t.Errorf("calendar localization failed")
	}
}
