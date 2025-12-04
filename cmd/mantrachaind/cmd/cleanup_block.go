package cmd

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtconfig "github.com/cometbft/cometbft/config"
	sm "github.com/cometbft/cometbft/state"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
)

// CleanupBlockEventsCmd creates a command to cleanup bloated block events
func CleanupBlockEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup-block-events [height] [height2] ...",
		Short: "Clean up bloated event data from specific block heights",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config
			var heights []int64
			for _, arg := range args {
				var height int64
				if _, err := fmt.Sscan(arg, &height); err != nil {
					return fmt.Errorf("invalid height '%s': %w", arg, err)
				}
				heights = append(heights, height)
			}
			return cleanupBlockEvents(cfg, heights)
		},
	}
	cmd.Flags().String(flags.FlagHome, "", "The application home directory")
	return cmd
}

func cleanupBlockEvents(cfg *cmtconfig.Config, heights []int64) error {
	stateDB, err := cmtconfig.DefaultDBProvider(&cmtconfig.DBContext{
		ID:     "state",
		Config: cfg,
	})
	if err != nil {
		return fmt.Errorf("failed to open state database: %w", err)
	}
	defer stateDB.Close()
	stateStore := sm.NewStore(stateDB, sm.StoreOptions{
		DiscardABCIResponses: false,
	})
	for _, height := range heights {
		resp, err := stateStore.LoadFinalizeBlockResponse(height)
		if err != nil {
			fmt.Printf("Failed to load block at height %d: %v\n", height, err)
			continue
		}
		if resp == nil {
			fmt.Printf("No finalize block response found at height %d\n", height)
			continue
		}
		resp.Events = []abci.Event{}
		for i := range resp.TxResults {
			if resp.TxResults[i] != nil {
				resp.TxResults[i].Events = []abci.Event{}
			}
		}
		if err := stateStore.SaveFinalizeBlockResponse(height, resp); err != nil {
			return fmt.Errorf("failed to save modified finalize block response at height %d: %w", height, err)
		}
	}
	return nil
}
