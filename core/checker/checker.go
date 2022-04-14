package checker

import (
	"fmt"
	"strings"

	"github.com/abdfnx/doko/tools"
	"github.com/abdfnx/doko/core/api"
	"github.com/abdfnx/doko/cli/factory"

	"github.com/mgutz/ansi"
	"github.com/tidwall/gjson"
	"github.com/abdfnx/looker"
)

func Check(buildVersion string) {
	cmdFactory := factory.New()
	stderr := cmdFactory.IOStreams.ErrOut

	if buildVersion == "" {
		buildVersion = "unknown"
	}

	latestVersion := api.GetLatest()
	isFromHomebrewTap := isUnderHomebrew()
	isFromGo := isUnderGo()
	isFromUsr := isUnderUsr()
	isFromAppData := isUnderAppData()

	var command = func() string {
		if isFromHomebrewTap {
			return "brew upgrade doko"
		} else if isFromGo {
			return "go get -u github.com/abdfnx/doko"
		} else if isFromUsr {
			return "curl -fsSL https://git.io/doko | bash"
		} else if isFromAppData {
			return "iwr -useb https://git.io/doko-win | iex"
		}

		return ""
	}

	if buildVersion != "unknown" {
		if buildVersion != latestVersion && gjson.Get(tools.SettingsContent(), "dk_settings.show_update").String() != "false" {
			fmt.Fprintf(stderr, "%s %s → %s\n",
			ansi.Color("There's a new version of ", "yellow") + ansi.Color("doko", "cyan") + ansi.Color(" is avalaible:", "yellow"),
			ansi.Color(buildVersion, "cyan"),
			ansi.Color(latestVersion, "cyan"))
	
			if command() != "" {
				fmt.Fprintf(stderr, ansi.Color("To upgrade, run: %s\n", "yellow"), ansi.Color(command(), "black:white"))
			}
		}
	}
}

var dokoExe, _ = looker.LookPath("doko")

func isUnderHomebrew() bool {
	return strings.Contains(dokoExe, "brew")
}

func isUnderGo() bool {
	return strings.Contains(dokoExe, "go")
}

func isUnderUsr() bool {
	return strings.Contains(dokoExe, "usr")
}

func isUnderAppData() bool {
	return strings.Contains(dokoExe, "AppData")
}
