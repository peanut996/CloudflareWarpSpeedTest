package i18n

import (
	"embed"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed resources/*.toml
var localeFS embed.FS

var localizer *i18n.Localizer

var folderPrefix = "resources/"

func init() {
	initLocalizer()
}

func initLocalizer() {
	lang := os.Getenv("LANG")
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFileFS(localeFS, folderPrefix+"en.toml")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	_, err = bundle.LoadMessageFileFS(localeFS, folderPrefix+"zh.toml")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	localizer = i18n.NewLocalizer(bundle, strings.Split(lang, ".")[0])
}

func QueryI18n(mesID string) string {
	if localizer == nil {
		initLocalizer()
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: mesID})
}

func QueryTemplateI18n(mesID string, tpData map[string]interface{}) string {
	if localizer == nil {
		initLocalizer()
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: mesID, TemplateData: tpData})
}
