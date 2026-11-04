package cmd

import (
	"fmt"
	"strings"

	"github.com/MANTRA-Chain/mantrachain/v8/app"
	"github.com/spf13/cobra"
)

// upgradeHandlersCommand lets governance tooling check a proposal's plan name
// against the binary before the chain halts on a missing handler.
func upgradeHandlersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade-handlers",
		Short: "List the upgrade plan names this binary has handlers for",
		Args:  cobra.NoArgs,
		// only reads this process, so skip the root's config and home interception
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		RunE: func(cmd *cobra.Command, _ []string) error {
			names := make([]string, len(app.Upgrades))
			for i, u := range app.Upgrades {
				names[i] = u.UpgradeName
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return err
		},
	}
}
