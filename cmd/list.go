package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/ramiawar/superpet/config"
	"github.com/ramiawar/superpet/envvar"
	"github.com/ramiawar/superpet/snippet"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all snippets",
	Long:  `Show all snippets`,
	RunE:  list,
}

var listenvCmd = &cobra.Command{
	Use:   "listenv",
	Short: "Show all env vars",
	Long:  `Show all env vars`,
	RunE:  listenv,
}

func list(cmd *cobra.Command, args []string) error {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, snippet := range snippets.Snippets {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(snippet.Description, col, "..."), col)
			command := runewidth.Truncate(snippet.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), snippet.Description)
			if strings.Contains(snippet.Command, "\n") {
				lines := strings.Split(snippet.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), snippet.Command)
			}
			if snippet.Tag != nil {
				tag := strings.Join(snippet.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if snippet.Output != "" {
				output := strings.Replace(snippet.Output, "\n", "\n             ", -1)
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.RedString("     Output:"), output)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}

func listenv(cmd *cobra.Command, args []string) error {
	var envvars envvar.EnvVar
	if err := envvars.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	fmt.Println("")
	for _, envvar := range envvars.EnvVars {
		var variables []string
		for _, value := range envvar.Variables {
			values := strings.Split(value, "=")
			variables = append(variables, values[0])
		}
		vars := strings.Join(variables, ", ")

		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(envvar.Description, col, "..."), col)
			vars = runewidth.Truncate(vars, 100-4-col, "...")
			// make sure multiline vars printed as oneline
			vars = strings.Replace(vars, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(vars))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), envvar.Description)
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.YellowString("    Variables:"), vars)

			if envvar.Tag != nil {
				tag := strings.Join(envvar.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(listenvCmd)
	listCmd.Flags().BoolVarP(&config.Flag.OneLine, "oneline", "", false,
		`Display snippets in one line`)
}
