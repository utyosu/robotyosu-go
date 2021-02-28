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

func HelpBasicCommands(lang string) string {
	return buildCommands(ToLanguage(lang), helpBasicCommands)
}

func HelpRecruitmentCommands(lang string) string {
	return buildCommands(ToLanguage(lang), helpRecruitmentCommands)
}

func buildCommands(lang string, commands map[string][]commandSet) string {
	var ret string
	for _, c := range commands[ToLanguage(lang)] {
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
