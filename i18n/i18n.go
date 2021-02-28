package i18n

import (
	"fmt"
)

func T(lang, key string, p ...interface{}) string {
	return fmt.Sprintf(languageDictionary[ToLanguage(lang)][key], p...)
}

func CommonMessage(key string) string {
	return commonDictionary[key]
}

func HelpBasicCommand(lang string) string {
	var ret string
	for _, c := range helpBasicCommand[ToLanguage(lang)] {
		if c.command == "" {
			ret += fmt.Sprintf("\n%v\n", c.description)
		} else {
			ret += fmt.Sprintf("`%v` %v\n", c.command, c.description)
		}
	}
	return ret
}

func ToLanguage(lang string) string {
	for _, l := range Languages {
		if l == lang {
			return lang
		}
	}
	return DefaultLanguage
}
