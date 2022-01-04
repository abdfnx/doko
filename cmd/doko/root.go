package doko

import (
	"os"
	"fmt"
	"context"
	"runtime"

	"github.com/abdfnx/doko/log"
	"github.com/abdfnx/doko/core"
	"github.com/abdfnx/doko/docker"
	"github.com/abdfnx/doko/core/opts"
	"github.com/abdfnx/doko/cmd/factory"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/mattn/go-runewidth"
	"github.com/MakeNowJust/heredoc"
	"github.com/docker/docker/client"
)

var dokoOpts = &opts.Options{
	Endpoint: "",
	CertPath: "",
	KeyPath:  "",
	CaPath:   "",
	ApiVersion: "",
	LogFilePath: "",
	LogLevelPath: "",
}

// Execute start the CLI
func Execute(f *factory.Factory, version string, buildDate string) *cobra.Command {
	const desc = `🐳 docker you know but with console user interface.`

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "doko <subcommand> [flags]",
		Short:  desc,
		Long: desc,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			# Open doko
			doko
		`),
		Annotations: map[string]string{
			"help:tellus": heredoc.Doc(`
				Open an issue at https://github.com/abdfnx/doko/issues
			`),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS == "windows" && runewidth.IsEastAsian() {
				tview.Borders.Horizontal = '-'
				tview.Borders.Vertical = '|'
				tview.Borders.TopLeft = '+'
				tview.Borders.TopRight = '+'
				tview.Borders.BottomLeft = '+'
				tview.Borders.BottomRight = '+'
				tview.Borders.LeftT = '|'
				tview.Borders.RightT = '|'
				tview.Borders.TopT = '-'
				tview.Borders.BottomT = '-'
				tview.Borders.Cross = '+'
				tview.Borders.HorizontalFocus = '='
				tview.Borders.VerticalFocus = '|'
				tview.Borders.TopLeftFocus = '+'
				tview.Borders.TopRightFocus = '+'
				tview.Borders.BottomLeftFocus = '+'
				tview.Borders.BottomRightFocus = '+'
			}

			logger.NewLogger(string(dokoOpts.LogFilePath), string(dokoOpts.LogLevelPath))

			docker.NewDocker(
				docker.NewClientConfig(
					string(dokoOpts.Endpoint),
					string(dokoOpts.CertPath),
					string(dokoOpts.KeyPath),
					string(dokoOpts.CaPath),
					string(dokoOpts.ApiVersion),
				),
			)

			if _, err := docker.Client.Info(context.TODO()); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}

			ui := core.New()

			if err := ui.Start(); err != nil {
				logger.Logger.Errorf("cannot start `doko`: %s", err)
				return err
			}

			return nil
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Aliases: []string{"ver"},
		Short: "Print the version of your doko binary.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("doko version " + version + " " + buildDate)
		},
	}

	rootCmd.SetOut(f.IOStreams.Out)
	rootCmd.SetErr(f.IOStreams.ErrOut)

	cs := f.IOStreams.ColorScheme()

	helpHelper := func(command *cobra.Command, args []string) {
		rootHelpFunc(cs, command, args)
	}

	rootCmd.PersistentFlags().Bool("help", false, "Help for doko")
	rootCmd.SetHelpFunc(helpHelper)
	rootCmd.SetUsageFunc(rootUsageFunc)
	rootCmd.SetFlagErrorFunc(rootFlagErrorFunc)

	// add `versionCmd` to root command
	rootCmd.AddCommand(versionCmd)

	// add flags
	rootCmd.Flags().StringVarP(&dokoOpts.Endpoint, "endpoint", "e", client.DefaultDockerHost, "The docker endpoint to use")
	rootCmd.Flags().StringVarP(&dokoOpts.CertPath, "cert", "c", "", "The path to the TLS certificate")
	rootCmd.Flags().StringVarP(&dokoOpts.KeyPath, "key", "k", "", "The path to the TLS key")
	rootCmd.Flags().StringVarP(&dokoOpts.CaPath, "ca", "", "", "The path to the TLS CA")
	rootCmd.Flags().StringVarP(&dokoOpts.ApiVersion, "api", "a", "", "The docker api version")
	rootCmd.Flags().StringVarP(&dokoOpts.LogFilePath, "logfile", "l", "", "The path to the log file")
	rootCmd.Flags().StringVarP(&dokoOpts.LogLevelPath, "loglevel", "o", "info", "The log level")

	return rootCmd
}
