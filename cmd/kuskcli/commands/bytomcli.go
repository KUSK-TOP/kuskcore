package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"kuskcore/util"
)

// kuskcli usage template
var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:
    {{range .Commands}}{{if (and .IsAvailableCommand (.Name | WalletDisable))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

  available with wallet enable:
    {{range .Commands}}{{if (and .IsAvailableCommand (.Name | WalletEnable))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// commandError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (c commandError) Error() string {
	return c.s
}

func (c commandError) isUserError() bool {
	return c.userError
}

func newUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

func newSystemError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: false}
}

func newSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: false}
}

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}

// KuskcliCmd is Kuskcli's root command.
// Every other command attached to KuskcliCmd is a child command to it.
var KuskcliCmd = &cobra.Command{
	Use:   "kuskcli",
	Short: "Kuskcli is a commond line client for kusk core (a.k.a. kuskd)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.SetUsageTemplate(usageTemplate)
			cmd.Usage()
		}
	},
}

// Execute adds all child commands to the root command KuskcliCmd and sets flags appropriately.
func Execute() {

	AddCommands()
	AddTemplateFunc()

	if _, err := KuskcliCmd.ExecuteC(); err != nil {
		os.Exit(util.ErrLocalExe)
	}
}

// AddCommands adds child commands to the root command KuskcliCmd.
func AddCommands() {
	KuskcliCmd.AddCommand(createAccessTokenCmd)
	KuskcliCmd.AddCommand(listAccessTokenCmd)
	KuskcliCmd.AddCommand(deleteAccessTokenCmd)
	KuskcliCmd.AddCommand(checkAccessTokenCmd)

	KuskcliCmd.AddCommand(createAccountCmd)
	KuskcliCmd.AddCommand(deleteAccountCmd)
	KuskcliCmd.AddCommand(listAccountsCmd)
	KuskcliCmd.AddCommand(updateAccountAliasCmd)
	KuskcliCmd.AddCommand(createAccountReceiverCmd)
	KuskcliCmd.AddCommand(listAddressesCmd)
	KuskcliCmd.AddCommand(validateAddressCmd)
	KuskcliCmd.AddCommand(listPubKeysCmd)

	KuskcliCmd.AddCommand(createAssetCmd)
	KuskcliCmd.AddCommand(getAssetCmd)
	KuskcliCmd.AddCommand(listAssetsCmd)
	KuskcliCmd.AddCommand(updateAssetAliasCmd)

	KuskcliCmd.AddCommand(getTransactionCmd)
	KuskcliCmd.AddCommand(listTransactionsCmd)

	KuskcliCmd.AddCommand(getUnconfirmedTransactionCmd)
	KuskcliCmd.AddCommand(listUnconfirmedTransactionsCmd)
	KuskcliCmd.AddCommand(decodeRawTransactionCmd)

	KuskcliCmd.AddCommand(listUnspentOutputsCmd)
	KuskcliCmd.AddCommand(listBalancesCmd)

	KuskcliCmd.AddCommand(rescanWalletCmd)
	KuskcliCmd.AddCommand(walletInfoCmd)

	KuskcliCmd.AddCommand(buildTransactionCmd)
	KuskcliCmd.AddCommand(signTransactionCmd)
	KuskcliCmd.AddCommand(submitTransactionCmd)
	KuskcliCmd.AddCommand(estimateTransactionGasCmd)

	KuskcliCmd.AddCommand(getBlockCountCmd)
	KuskcliCmd.AddCommand(getBlockHashCmd)
	KuskcliCmd.AddCommand(getBlockCmd)
	KuskcliCmd.AddCommand(getBlockHeaderCmd)

	KuskcliCmd.AddCommand(createKeyCmd)
	KuskcliCmd.AddCommand(deleteKeyCmd)
	KuskcliCmd.AddCommand(listKeysCmd)
	KuskcliCmd.AddCommand(updateKeyAliasCmd)
	KuskcliCmd.AddCommand(resetKeyPwdCmd)
	KuskcliCmd.AddCommand(checkKeyPwdCmd)

	KuskcliCmd.AddCommand(signMsgCmd)
	KuskcliCmd.AddCommand(verifyMsgCmd)
	KuskcliCmd.AddCommand(decodeProgCmd)

	KuskcliCmd.AddCommand(createTransactionFeedCmd)
	KuskcliCmd.AddCommand(listTransactionFeedsCmd)
	KuskcliCmd.AddCommand(deleteTransactionFeedCmd)
	KuskcliCmd.AddCommand(getTransactionFeedCmd)
	KuskcliCmd.AddCommand(updateTransactionFeedCmd)

	KuskcliCmd.AddCommand(netInfoCmd)
	KuskcliCmd.AddCommand(gasRateCmd)

	KuskcliCmd.AddCommand(versionCmd)
}

// AddTemplateFunc adds usage template to the root command KuskcliCmd.
func AddTemplateFunc() {
	walletEnableCmd := []string{
		createAccountCmd.Name(),
		listAccountsCmd.Name(),
		deleteAccountCmd.Name(),
		updateAccountAliasCmd.Name(),
		createAccountReceiverCmd.Name(),
		listAddressesCmd.Name(),
		validateAddressCmd.Name(),
		listPubKeysCmd.Name(),

		createAssetCmd.Name(),
		getAssetCmd.Name(),
		listAssetsCmd.Name(),
		updateAssetAliasCmd.Name(),

		createKeyCmd.Name(),
		deleteKeyCmd.Name(),
		listKeysCmd.Name(),
		resetKeyPwdCmd.Name(),
		checkKeyPwdCmd.Name(),
		signMsgCmd.Name(),

		buildTransactionCmd.Name(),
		signTransactionCmd.Name(),

		getTransactionCmd.Name(),
		listTransactionsCmd.Name(),
		listUnspentOutputsCmd.Name(),
		listBalancesCmd.Name(),

		rescanWalletCmd.Name(),
		walletInfoCmd.Name(),
	}

	cobra.AddTemplateFunc("WalletEnable", func(cmdName string) bool {
		for _, name := range walletEnableCmd {
			if name == cmdName {
				return true
			}
		}
		return false
	})

	cobra.AddTemplateFunc("WalletDisable", func(cmdName string) bool {
		for _, name := range walletEnableCmd {
			if name == cmdName {
				return false
			}
		}
		return true
	})
}
