package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

// Localizer is for localizing texts.
type Localizer struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

// NewLocalizer creates a new instance of Localizer.
func NewLocalizer() *Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	localizer := i18n.NewLocalizer(bundle, "en")
	return &Localizer{bundle: bundle, localizer: localizer}
}

// AddLocalization loads message files into the Localizer bundle.
func (l *Localizer) AddLocalization(messageFilePath ...string) {
	for _, path := range messageFilePath {
		_, err := l.bundle.LoadMessageFile(path)
		if err != nil {
			log.Fatalf("could not load message file: %v", err)
		}
	}
}

// Localize translates the given key into a localized string using the Localizer bundle.
// It either returns the localized string or a value indicating that something went wrong.
func (l *Localizer) Localize(key string) string {
	localized, err := l.localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		log.Printf("could not get text for '%s'", key)
		return "☒"
	}
	return localized
}
