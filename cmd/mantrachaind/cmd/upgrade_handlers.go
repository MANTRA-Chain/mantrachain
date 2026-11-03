package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

// UpgradeHandlerRegistry exposes the upgrade handlers registered by the application. The upgrade
// keeper is the source of truth for it, since it is the keeper that matches handler names against
// the on-chain upgrade plan.
type UpgradeHandlerRegistry interface {
	// UpgradeHandlerNames returns the sorted names of the registered upgrade handlers.
	UpgradeHandlerNames() []string
	// HasUpgradeHandler returns true if a handler is registered for the given name.
	HasUpgradeHandler(name string) bool
}

// UpgradeHandlersCmd returns a command that prints the upgrade handler names registered by this
// binary. Operators can use it to verify that the binary they are about to run contains the
// handler of a pending upgrade proposal.
func UpgradeHandlersCmd(registry UpgradeHandlerRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade-handlers [plan-name]",
		Short: "Print the upgrade handler names registered by this binary",
		Long: `Print the upgrade handler names registered by this binary, one per line.

Upgrade handlers are matched by name against the on-chain upgrade plan. Use this command to
verify that the binary you are about to run contains the handler of a pending upgrade proposal.

When a plan name is provided, the command exits with an error if this binary does not register
a handler for that name.`,
		Example: `  mantrachaind upgrade-handlers
  mantrachaind upgrade-handlers v8.0.0`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		// The command only reads the handlers registered in-process, so override the root
		// persistent pre-run and avoid requiring a home directory, config files or a client context.
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return checkUpgradeHandler(cmd.OutOrStdout(), registry, args[0])
			}

			return printUpgradeHandlers(cmd.OutOrStdout(), registry.UpgradeHandlerNames())
		},
	}

	return cmd
}

// printUpgradeHandlers writes the given upgrade handler names, one per line.
func printUpgradeHandlers(w io.Writer, names []string) error {
	for _, name := range names {
		if _, err := fmt.Fprintln(w, name); err != nil {
			return err
		}
	}

	return nil
}

// checkUpgradeHandler returns an error when the application has no handler registered for planName.
func checkUpgradeHandler(w io.Writer, registry UpgradeHandlerRegistry, planName string) error {
	if !registry.HasUpgradeHandler(planName) {
		return fmt.Errorf(
			"upgrade handler %s is not registered by this binary (registered handlers: %s)",
			planName,
			strings.Join(registry.UpgradeHandlerNames(), ", "),
		)
	}

	_, err := fmt.Fprintf(w, "upgrade handler %s is registered\n", planName)

	return err
}
