package i18n

import (
	"embed"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed *.toml */*.toml
var localeFS embed.FS

func InitI18n(path string) *i18n.Localizer {
	lang := os.Getenv("LANG")
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(localeFS, path+".en.toml")
	bundle.LoadMessageFileFS(localeFS, path+".zh.toml")
	return i18n.NewLocalizer(bundle, strings.Split(lang, ".")[0])
}

func QueryI18n(localizer *i18n.Localizer, mesID string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: mesID})
}

func QueryTemplateI18n(localizer *i18n.Localizer, mesID string, tpData map[string]interface{}) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: mesID, TemplateData: tpData})
}
