package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MANTRA-Chain/mantrachain/v8/app"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
)

// fakeUpgradeHandlerRegistry is an in-memory UpgradeHandlerRegistry.
type fakeUpgradeHandlerRegistry struct {
	registered map[string]bool
}

func newFakeUpgradeHandlerRegistry(names ...string) *fakeUpgradeHandlerRegistry {
	registered := make(map[string]bool, len(names))
	for _, name := range names {
		registered[name] = true
	}

	return &fakeUpgradeHandlerRegistry{registered: registered}
}

func (r *fakeUpgradeHandlerRegistry) UpgradeHandlerNames() []string {
	names := make([]string, 0, len(r.registered))
	for name := range r.registered {
		names = append(names, name)
	}

	return names
}

func (r *fakeUpgradeHandlerRegistry) HasUpgradeHandler(name string) bool {
	return r.registered[name]
}

func TestUpgradeHandlersCmd(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := UpgradeHandlersCmd(newFakeUpgradeHandlerRegistry("v8.0.0", "v8.0.0-rc3"))
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	require.NoError(t, cmd.Execute())
	require.ElementsMatch(t, []string{"v8.0.0", "v8.0.0-rc3"}, strings.Fields(buf.String()))
}

func TestUpgradeHandlersCmdWithPlanName(t *testing.T) {
	t.Run("registered handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cmd := UpgradeHandlersCmd(newFakeUpgradeHandlerRegistry("v8.0.0"))
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"v8.0.0"})

		require.NoError(t, cmd.Execute())
		require.Equal(t, "upgrade handler v8.0.0 is registered\n", buf.String())
	})

	t.Run("unknown handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cmd := UpgradeHandlersCmd(newFakeUpgradeHandlerRegistry("v8.0.0"))
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"v9.0.0"})

		err := cmd.Execute()
		require.ErrorContains(t, err, "upgrade handler v9.0.0 is not registered by this binary")
		require.ErrorContains(t, err, "registered handlers: v8.0.0")
	})
}

// TestUpgradeHandlersCmdUsesUpgradeKeeper makes sure the command reports the handlers registered
// with the upgrade keeper by the application, not a hardcoded list.
func TestUpgradeHandlersCmdUsesUpgradeKeeper(t *testing.T) {
	rootCmd := NewRootCmd()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"upgrade-handlers"})

	require.NoError(t, rootCmd.Execute())

	expNames := upgrades.HandlerNames(app.Upgrades)
	require.NotEmpty(t, expNames)
	require.Equal(t, expNames, strings.Fields(buf.String()))
}
