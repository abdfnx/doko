package tools

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

var homeDir, _ = homedir.Dir()
var dotDoko = path.Join(homeDir, "./.doko")
var settingsFile = path.Join(dotDoko, "settings.json")

func CheckDotDoko() {
	if _, err := os.Stat(dotDoko); os.IsNotExist(err) {
		os.Mkdir(dotDoko, 0755)
		Files()
	}

	Files()
}

func Files() {
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		os.Create(settingsFile)
		SetDefaultSettings()
	}
}

func SettingsFile() string {
	return settingsFile
}
