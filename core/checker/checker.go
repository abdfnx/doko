package checker

import (
	"fmt"
	"strings"

	"github.com/abdfnx/doko/tools"
	"github.com/abdfnx/doko/core/api"
	"github.com/abdfnx/doko/cmd/factory"

	"github.com/mgutz/ansi"
	"github.com/tidwall/gjson"
	tcexe "github.com/Timothee-Cardoso/tc-exe"
)

func Check(buildVersion string) {
	cmdFactory := factory.New()
	stderr := cmdFactory.IOStreams.ErrOut

	latestVersion := api.GetLatest()
	isFromHomebrewTap := isUnderHomebrew()
	isFromGo := isUnderGo()
	isFromUsrBinDir := isUnderUsr()
	isFromAppData := isUnderAppData()

	var command = func() string {
		if isFromHomebrewTap {
			return "brew upgrade doko"
		} else if isFromGo {
			return "go get -u github.com/abdfnx/doko"
		} else if isFromUsrBinDir {
			return "curl -fsSL https://git.io/doko | bash"
		} else if isFromAppData {
			return "iwr -useb https://git.io/doko-win | iex"
		}

		return ""
	}

	if buildVersion != latestVersion && gjson.Get(tools.SettingsContent(), "rs_settings.show_update").String() != "false" {
		fmt.Fprintf(stderr, "\n%s %s → %s\n",
		ansi.Color("There's a new version of ", "yellow") + ansi.Color("doko", "cyan") + ansi.Color(" is avalaible:", "yellow"),
		ansi.Color(buildVersion, "cyan"),
		ansi.Color(latestVersion, "cyan"))

		if command() != "" {
			fmt.Fprintf(stderr, ansi.Color("To upgrade, run: %s\n", "yellow"), ansi.Color(command(), "black:white"))
		}
	}
}

var dokoExe, _ = tcexe.LookPath("doko")

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
