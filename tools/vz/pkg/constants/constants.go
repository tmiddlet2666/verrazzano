// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package constants

// GlobalFlagKubeConfig - global flag for specifying the location of the kube config
const GlobalFlagKubeConfig = "kubeconfig"
const GlobalFlagKubeConfigHelp = "Path to the kubeconfig file to use"

// GlobalFlagContext - global flag for specifying which kube config context to use
const GlobalFlagContext = "context"
const GlobalFlagContextHelp = "The name of the kubeconfig context to use"

// GlobalFlagHelp - global help flag
const GlobalFlagHelp = "help"

// Flags that are common to more than one command
const (
	WaitFlag        = "wait"
	WaitFlagHelp    = "Wait for the command to complete and stream the logs to the console. The wait period is controlled by --timeout."
	WaitFlagDefault = true

	TimeoutFlag     = "timeout"
	TimeoutFlagHelp = "Limits the amount of time a command will wait to complete"

	VersionFlag        = "version"
	VersionFlagDefault = "latest"
	VersionFlagHelp    = "The version of Verrazzano to install or upgrade"

	DryRunFlag = "dry-run"

	SetFlag          = "set"
	SetFlagShorthand = "s"
	SetFlagHelp      = "Override a Verrazzano resource value (e.g. --set profile=dev).  This flag can be specified multiple times."

	OperatorFileFlag     = "operator-file"
	OperatorFileFlagHelp = "The path to the file for installing the Verrazzano platform operator. The default is derived from the version string."

	LogFormatFlag = "log-format"
	LogFormatHelp = "The format of the log output. Valid output formats are \"simple\" and \"json\"."

	FilenameFlag          = "filename"
	FilenameFlagShorthand = "f"
	FilenameFlagHelp      = "Path to file containing Verrazzano custom resource.  This flag can be specified multiple times to overlay multiple files."
)

// VerrazzanoReleaseList - API for getting the list of Verrazzano releases
const VerrazzanoReleaseList = "https://api.github.com/repos/verrazzano/verrazzano/releases"

// VerrazzanoOperatorURL - URL for downloading Verrazzano releases
const VerrazzanoOperatorURL = "https://github.com/verrazzano/verrazzano/releases/download/%s/operator.yaml"

// Analysis tool flags
const (
	ActionsFlagName  = "actions"
	ActionsFlagValue = "actions"
	ActionsFlagUsage = "Include actions in the report, default is true (default true)"

	AnalysisFlagName  = "analysis"
	AnalysisFlagValue = "analysis"
	AnalysisFlagUsage = "Type of analysis: cluster (default \"cluster\")"

	HelpFlagName  = "help"
	HelpFlagValue = "help"
	HelpFlagUsage = "Display usage help"

	InfoFlagName  = "info"
	InfoFlagValue = "info"
	InfoFlagUsage = "Include informational messages, default is true (default true)"

	MinConfidenceFlagName  = "minConfidence"
	MinConfidenceFlagValue = "minConfidence"
	MinConfidenceFlagUsage = "Minimum confidence threshold to report for issues, 0-10, default is 0"

	MinImpactFlagName  = "minImpact"
	MinImpactFlagValue = "minImpact"
	MinImpactFlagUsage = "Minimum impact threshold to report for issues, 0-10, default is 0"

	ReportFileFlagName  = "reportFile"
	ReportFileFlagValue = "reportFile"
	ReportFileFlagUsage = "Name of report output file, default is stdout"

	SupportFlagName  = "support"
	SupportFlagValue = "support"
	SupportFlagUsage = "Include support data in the report, default is true (default true)"

	VersionFlagName  = "version"
	VersionFlagValue = "version"
	VersionFlagUsage = "Display version"

	ZapDevelFlagName  = "zap-devel"
	ZapDevelFlagValue = "zap-devel"
	ZapDevelFlagUsage = "Development Mode defaults(encoder=consoleEncoder,logLevel=Debug,stackTraceLevel=Warn). Production Mode defaults(encoder=jsonEncoder,logLevel=Info,stackTraceLevel=Error)"

	ZapEncoderFlagName  = "zap-encoder"
	ZapEncoderFlagValue = "zap-encoder"
	ZapEncoderFlagUsage = "Zap log encoding (one of 'json' or 'console')"

	ZapLogLevelFlagName  = "zap-log-level"
	ZapLogLevelFlagValue = "zap-log-level"
	ZapLogLevelFlagUsage = "Zap Level to configure the verbosity of logging. Can be one of 'debug', 'info', 'error', or any integer value > 0 which corresponds to custom debug levels of increasing verbosity"

	ZapStackTraceLevelFlagName  = "zap-stacktrace-level"
	ZapStackTraceLevelFlagValue = "zap-stacktrace-level"
	ZapStackTraceLevelFlagUsage = "Zap Level at and above which stacktraces are captured (one of 'info', 'error', 'panic')."

	ZapTimeEncodingFlagName  = "zap-time-encoding"
	ZapTimeEncodingFlagValue = "zap-time-encoding"
	ZapTimeEncodingFlagUsage = "Zap time encoding (one of 'epoch', 'millis', 'nano', 'iso8601', 'rfc3339' or 'rfc3339nano'). Defaults to 'epoch'."
)
