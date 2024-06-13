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
var LocaleFS embed.FS

func Init_i18n(path string) *i18n.Localizer {
	lang := os.Getenv("LANG")
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, path+".en.toml")
	bundle.LoadMessageFileFS(LocaleFS, path+".zh.toml")
	return i18n.NewLocalizer(bundle, strings.Split(lang, ".")[0])
}

func Query_i18n(localizer *i18n.Localizer, mesID string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: mesID})
}

func Query_template_i18n(localizer *i18n.Localizer, mesID string, tpData map[string]interface{}) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: mesID, TemplateData: tpData})
}
