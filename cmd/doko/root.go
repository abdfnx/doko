package doko

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/abdfnx/doko/cli"
	"github.com/abdfnx/doko/cmd/factory"
	"github.com/abdfnx/doko/core"
	"github.com/abdfnx/doko/core/opts"
	"github.com/abdfnx/doko/docker"
	"github.com/abdfnx/doko/log"
	"github.com/abdfnx/doko/tools"

	"github.com/MakeNowJust/heredoc"
	"github.com/docker/docker/client"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var dokoOpts = &opts.Options{
	Endpoint:      "",
	CertPath:      "",
	KeyPath:       "",
	CaPath:        "",
	EngineVersion: "",
	LogFilePath:   "",
	LogLevelPath:  "",
}

// Execute start the CLI
func Execute(f *factory.Factory, version string, buildDate string) *cobra.Command {
	tools.CheckDotDoko()

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

			# With specific endpoint
			doko --endpoint <DOCKER_ENDPOINT>

			# Use another docker engine version
			doko --engine "1.40"

			# Log file path
			doko --log-file /home/doko/my-log.log
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

			logger.NewLogger(dokoOpts.LogLevelPath, dokoOpts.LogFilePath)

			docker.NewDocker(
				docker.NewClientConfig(
					dokoOpts.Endpoint,
					dokoOpts.CertPath,
					dokoOpts.KeyPath,
					dokoOpts.CaPath,
					dokoOpts.EngineVersion,
				),
			)

			if _, err := docker.Client.Info(context.TODO()); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}

			ui := core.New()

			if err := ui.Start(version); err != nil {
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
	rootCmd.AddCommand(cli.SettingsCMD(), versionCmd)

	// add flags
	rootCmd.Flags().StringVarP(&dokoOpts.Endpoint, "endpoint", "e", client.DefaultDockerHost, "The docker endpoint to use")
	rootCmd.Flags().StringVarP(&dokoOpts.CertPath, "cert", "c", "", "The path to the TLS certificate (cert.pem)")
	rootCmd.Flags().StringVarP(&dokoOpts.KeyPath, "key", "k", "", "The path to the TLS key (key.pem)")
	rootCmd.Flags().StringVarP(&dokoOpts.CaPath, "ca", "", "", "The path to the TLS CA (ca.pem)")
	rootCmd.Flags().StringVarP(&dokoOpts.EngineVersion, "engine", "g", "1.41", "The docker engine version")
	rootCmd.Flags().StringVarP(&dokoOpts.LogFilePath, "log-file", "l", "", "The path to the log file")
	rootCmd.Flags().StringVarP(&dokoOpts.LogLevelPath, "log-level", "o", "info", "The log level")

	return rootCmd
}
