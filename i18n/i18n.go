package i18n

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
)

// Init 初始化i18n
func Init(lang string) error {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载所有语言文件
	languages := []string{"zh", "en"}
	for _, l := range languages {
		if _, err := bundle.LoadMessageFileFS(localeFS, "locales/"+l+".toml"); err != nil {
			return err
		}
	}

	// 设置当前语言
	SetLanguage(lang)
	return nil
}

// SetLanguage 设置当前语言
func SetLanguage(lang string) {
	var tag language.Tag
	switch lang {
	case "en", "english":
		tag = language.English
	default:
		tag = language.Chinese
	}
	localizer = i18n.NewLocalizer(bundle, tag.String())
}

// T 翻译函数
func T(messageID string, templateData ...map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}

	cfg := &i18n.LocalizeConfig{
		MessageID: messageID,
	}

	if len(templateData) > 0 {
		cfg.TemplateData = templateData[0]
	}

	msg, err := localizer.Localize(cfg)
	if err != nil {
		return messageID
	}
	return msg
}

// Tf 带格式化参数的翻译函数
func Tf(messageID string, count int, templateData map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}

	if templateData == nil {
		templateData = make(map[string]interface{})
	}
	templateData["Count"] = count

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  count,
	})
	if err != nil {
		return messageID
	}
	return msg
}
