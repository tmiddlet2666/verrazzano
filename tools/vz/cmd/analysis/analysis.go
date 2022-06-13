package analysis

import (
	"fmt"
	"github.com/spf13/cobra"
	cmdhelpers "github.com/verrazzano/verrazzano/tools/vz/cmd/helpers"
	analysis "github.com/verrazzano/verrazzano/tools/vz/pkg/analysis/main_pkg"
	"github.com/verrazzano/verrazzano/tools/vz/pkg/constants"
	"github.com/verrazzano/verrazzano/tools/vz/pkg/helpers"
)

const (
	CommandName = "analyze"
	helpShort   = "Verrazzano Analysis Tool"
	helpLong    = "Verrazzano Analysis Tool"
	helpExample = ``
)

func NewCmdAnalysis(vzHelper helpers.VZHelper) *cobra.Command {
	cmd := cmdhelpers.NewCommand(vzHelper, CommandName, helpShort, helpLong)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runCmdAnalysis(cmd, args, vzHelper)
	}

	cmd.PersistentFlags().String(constants.DirectoryFlagName, constants.DirectoryFlagValue, constants.DirectoryFlagUsage)
	cmd.PersistentFlags().String(constants.ReportFileFlagName, constants.ReportFileFlagValue, constants.ReportFileFlagUsage)
	cmd.PersistentFlags().String(constants.ReportFormatFlagName, constants.ReportFormatFlagValue, constants.ReportFormatFlagUsage)
	return cmd
}

func runCmdAnalysis(cmd *cobra.Command, args []string, helper helpers.VZHelper) error {
	fmt.Println("ran command analysis")
	reportFile, err := cmd.PersistentFlags().GetString(constants.ReportFileFlagName)
	if err != nil {
		fmt.Println("error fetching flag: %s", constants.ReportFileFlagName)
	}
	analysis.MainExecAnalysis(reportFile)
	return nil
}
