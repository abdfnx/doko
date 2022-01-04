package tools

import (
	"io/ioutil"
	"strconv"

	"github.com/tidwall/pretty"
)

func SettingsContent() string {
	stgFile, err := ioutil.ReadFile(settingsFile)

	if err != nil {
		panic(err)
	}

	return string(stgFile)
}

func UpdateSettings(ssu, sem bool) {
	settings := `
		{
			"dk_settings": {
				"show_update":` + strconv.FormatBool(ssu) + `,
				"enable_mouse":` + strconv.FormatBool(sem) + `
			}
		}
	`

	prettySettings := pretty.Pretty([]byte(settings))

	err := ioutil.WriteFile(settingsFile, []byte(string(prettySettings)), 0644)

	if err != nil {
		panic(err)
	}
}

func SetDefaultSettings() {
	defaultSettings := `
		{
			"dk_settings": {
				"show_update": true,
				"enable_mouse": true
			}
		}
	`

	prettySettings := pretty.Pretty([]byte(defaultSettings))

	err := ioutil.WriteFile(settingsFile, []byte(string(prettySettings)), 0644)

	if err != nil {
		panic(err)
	}
}
