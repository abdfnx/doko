package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/abdfnx/doko/tools"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

func SettingsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Update Settings Or Change it",
		Long:  `Update Doko settings like enable mouse`,
	}

	cmd.AddCommand(SettingsSet())

	return cmd
}

func SettingsSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set new or update settings",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 0 {
				var value bool

				if string(args[1]) == "true" {
					value = true
				} else {
					value = false
				}

				if strings.Contains(args[0], "show_update") {
					if string(args[1]) == "true" || string(args[1]) == "false" {
						update, err := sjson.Set(tools.SettingsContent(), "dk_settings.show_update", value)

						if err != nil {
							panic(err)
						}

						uerr := ioutil.WriteFile(tools.SettingsFile(), []byte(string(update)), 0644)

						if uerr != nil {
							panic(uerr)
						}
					} else {
						fmt.Println(ansi.Color("dk_settings.show_update must be `true` or `false`", "red"))
						os.Exit(1)
					}
				} else if strings.Contains(args[0], "mouse") {
					if string(args[1]) == "true" || string(args[1]) == "false" {
						mouse, err := sjson.Set(tools.SettingsContent(), "dk_settings.enable_mouse", value)

						if err != nil {
							panic(err)
						}

						merr := ioutil.WriteFile(tools.SettingsFile(), []byte(string(mouse)), 0644)

						if merr != nil {
							panic(merr)
						}
					} else {
						fmt.Println(ansi.Color("dk_settings.enable_mouse must be `true` or `false`", "red"))
						os.Exit(1)
					}
				}

				fmt.Println(ansi.Color("Settings updated", "green"))
			}
		},
	}

	return cmd
}
