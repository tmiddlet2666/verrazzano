package analysis

import (
	"fmt"
	"github.com/spf13/cobra"
	cmdhelpers "github.com/verrazzano/verrazzano/tools/vz/cmd/helpers"
	"github.com/verrazzano/verrazzano/tools/vz/pkg/constants"
	"github.com/verrazzano/verrazzano/tools/vz/pkg/helpers"
)

const (
	CommandName = "analysis"
	helpShort   = "analysis"
	helpLong    = "analysis"
	helpExample = ``
)

func NewCmdAnalysis(vzHelper helpers.VZHelper) *cobra.Command {
	cmd := cmdhelpers.NewCommand(vzHelper, CommandName, helpShort, helpLong)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runCmdAnalysis(cmd, args, vzHelper)
	}

	//TODO add flags
	cmd.PersistentFlags().String(constants.ActionsFlagName, constants.ActionsFlagValue, constants.ActionsFlagUsage)
	cmd.PersistentFlags().String(constants.AnalysisFlagName, constants.AnalysisFlagValue, constants.AnalysisFlagUsage)
	//cmd.PersistentFlags().String(constants.HelpFlagName, constants.HelpFlagValue, constants.HelpFlagUsage) //TODO
	cmd.PersistentFlags().String(constants.InfoFlagName, constants.InfoFlagValue, constants.InfoFlagUsage)
	cmd.PersistentFlags().String(constants.MinConfidenceFlagName, constants.MinConfidenceFlagValue, constants.MinConfidenceFlagUsage)
	cmd.PersistentFlags().String(constants.MinImpactFlagName, constants.MinImpactFlagValue, constants.MinImpactFlagUsage)
	cmd.PersistentFlags().String(constants.ReportFileFlagName, constants.ReportFileFlagValue, constants.ReportFileFlagUsage)
	cmd.PersistentFlags().String(constants.SupportFlagName, constants.SupportFlagValue, constants.SupportFlagUsage)
	cmd.PersistentFlags().String(constants.VersionFlagName, constants.VersionFlagValue, constants.VersionFlagUsage)
	cmd.PersistentFlags().String(constants.ZapDevelFlagName, constants.ZapDevelFlagValue, constants.ZapDevelFlagUsage)
	cmd.PersistentFlags().String(constants.ZapEncoderFlagName, constants.ZapEncoderFlagValue, constants.ZapEncoderFlagUsage)
	cmd.PersistentFlags().String(constants.ZapLogLevelFlagName, constants.ZapLogLevelFlagValue, constants.ZapLogLevelFlagUsage)
	cmd.PersistentFlags().String(constants.ZapStackTraceLevelFlagName, constants.ZapStackTraceLevelFlagValue, constants.ZapStackTraceLevelFlagUsage)
	cmd.PersistentFlags().String(constants.ZapTimeEncodingFlagName, constants.ZapTimeEncodingFlagValue, constants.ZapTimeEncodingFlagUsage)
	return cmd
}

func runCmdAnalysis(cmd *cobra.Command, args []string, helper helpers.VZHelper) error {
	fmt.Println("ran command analysis")
	return nil
}
