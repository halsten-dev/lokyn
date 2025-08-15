// Package lokyn is a small and lightweight library to help using lower level libs like jeandeaual/go-locale and nicksnyder/go-i18n.
// It's heavily inspired by the Fyne lang package, but in a standalone philosophy. The real deal comes when Lokyn app is used to help
// translating the application.
package lokyn

import (
	"embed"
	"encoding/json"
	"log"
	"sync"

	"github.com/jeandeaual/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	once      sync.Once

	currentLang language.Tag
	translated  []language.Tag
)

// Init inits the package, only once.
func Init() {
	once.Do(initBundle)
}

// AddTranslationFS registers all the managed languages to lokyn.
func AddTranslationFS(fs embed.FS, dir string) error {
	files, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		name := f.Name()
		data, err := fs.ReadFile(dir + "/" + name)
		if err != nil {
			continue
		}

		err = addLanguage(data, name)
		if err != nil {
			continue
		}
	}

	initLanguage()

	return nil
}

// GetCurrentLanguage returns the current language as string.
func GetCurrentLanguage() string {
	return currentLang.String()
}

// SetLanguage helps defining the current language.
func SetLanguage(lang string) {
	setupLang(lang)
}

// L returns translation of the given key.
func L(key string) string {
	return getKey(key, key)
}

// P returns translation with plural management.
func P(key string, count int) string {
	return getPluralKey(key, key, count)
}

// initBundle initialize lokyn with english language.
func initBundle() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	translated = []language.Tag{language.Make("en")}
}

// initLanguage init the language based on the system language.
func initLanguage() {
	all, err := locale.GetLocales()
	if err != nil {
		all = []string{"en"}
	}

	setupLang(closestSupportedLocale(all).String())
}

// setupLang initialize a new localizer for the requested language.
func setupLang(lang string) {
	currentLang = language.Make(lang)
	localizer = i18n.NewLocalizer(bundle, lang)
}

// addLanguage adds a language to lokyn managed languages.
func addLanguage(data []byte, name string) error {
	f, err := bundle.ParseMessageFileBytes(data, name)
	if err != nil {
		return err
	}

	translated = append(translated, f.Tag)
	return nil
}

// getKey gets the requested key, manage a fallback key.
func getKey(key, fallback string) string {
	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    key,
			Other: fallback,
		},
	})

	if err != nil {
		log.Println("Error in translation")
	}

	return ret
}

// getPluralKey gets the requested key, manage a fallback and plural.
func getPluralKey(key, fallback string, count int) string {
	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    key,
			Other: fallback,
		},
		PluralCount: count,
	})

	if err != nil {
		log.Println("Error in translation")
	}

	return ret
}

// closestSupportedLocale helps to determine the closest language tag based on the locales given in parameter.
func closestSupportedLocale(locs []string) language.Tag {
	matcher := language.NewMatcher(translated)

	tags := make([]language.Tag, len(locs))
	for i, loc := range locs {
		tag, err := language.Parse(loc)
		if err != nil {
			log.Println("Error in parsing tags")
		}
		tags[i] = tag
	}
	best, _, _ := matcher.Match(tags...)
	return best
}
