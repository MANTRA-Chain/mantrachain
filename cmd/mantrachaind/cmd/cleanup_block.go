package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtconfig "github.com/cometbft/cometbft/config"
	sm "github.com/cometbft/cometbft/state"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
)

type BlockEventsExport struct {
	Height      int64          `json:"height"`
	BlockEvents []abci.Event   `json:"block_events"`
	TxEvents    [][]abci.Event `json:"tx_events"`
}

// ExportBlockEventsCmd creates a command to export block events to a file
func ExportBlockEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-block-events [height] [height2] ...",
		Short: "Export block event data from specific block heights to a JSON file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config
			outputFile, _ := cmd.Flags().GetString("output")
			if outputFile == "" {
				return fmt.Errorf("output file is required")
			}

			var heights []int64
			for _, arg := range args {
				var height int64
				if _, err := fmt.Sscan(arg, &height); err != nil {
					return fmt.Errorf("invalid height '%s': %w", arg, err)
				}
				heights = append(heights, height)
			}
			return exportBlockEvents(cfg, heights, outputFile)
		},
	}
	cmd.Flags().String(flags.FlagHome, "", "The application home directory")
	cmd.Flags().String("output", "block-events-export.json", "Output file for exported events")
	return cmd
}

// ImportBlockEventsCmd creates a command to import block events from a file
func ImportBlockEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-block-events",
		Short: "Import block event data from a JSON file and restore events",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config
			inputFile, _ := cmd.Flags().GetString("input")
			if inputFile == "" {
				return fmt.Errorf("input file is required")
			}
			return importBlockEvents(cfg, inputFile)
		},
	}
	cmd.Flags().String(flags.FlagHome, "", "The application home directory")
	cmd.Flags().String("input", "block-events-export.json", "Input file with exported events")
	return cmd
}

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

func exportBlockEvents(cfg *cmtconfig.Config, heights []int64, outputFile string) error {
	stateStore, closeFunc, err := openStateStore(cfg)
	if err != nil {
		return err
	}
	defer closeFunc()

	var exports []BlockEventsExport

	for _, height := range heights {
		export, err := loadBlockEventsExport(stateStore, height)
		if err != nil {
			fmt.Printf("Failed to load block at height %d: %v\n", height, err)
			continue
		}

		exports = append(exports, *export)
		fmt.Printf("Exported %d block events and %d tx event groups from height %d\n",
			len(export.BlockEvents), len(export.TxEvents), height)
	}

	data, err := json.MarshalIndent(exports, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	fmt.Printf("Successfully exported events from %d heights to %s\n", len(exports), outputFile)
	return nil
}

func importBlockEvents(cfg *cmtconfig.Config, inputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	var exports []BlockEventsExport
	if err := json.Unmarshal(data, &exports); err != nil {
		return fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	stateStore, closeFunc, err := openStateStore(cfg)
	if err != nil {
		return err
	}
	defer closeFunc()

	for _, export := range exports {
		if err := applyBlockEventsExport(stateStore, export); err != nil {
			fmt.Printf("Failed to restore events at height %d: %v\n", export.Height, err)
			continue
		}

		fmt.Printf("Restored %d block events and %d tx event groups at height %d\n",
			len(export.BlockEvents), len(export.TxEvents), export.Height)
	}

	fmt.Printf("Successfully imported (restored) events for %d heights\n", len(exports))
	return nil
}

func cleanupBlockEvents(cfg *cmtconfig.Config, heights []int64) error {
	stateStore, closeFunc, err := openStateStore(cfg)
	if err != nil {
		return err
	}
	defer closeFunc()

	for _, height := range heights {
		export, err := loadBlockEventsExport(stateStore, height)
		if err != nil {
			fmt.Printf("Failed to load block at height %d: %v\n", height, err)
			continue
		}

		filteredBlockEvents := []abci.Event{}
		for _, event := range export.BlockEvents {
			if event.Type == "oracle_prices" {
				filteredBlockEvents = append(filteredBlockEvents, event)
			}
		}

		cleared := BlockEventsExport{
			Height:      height,
			BlockEvents: filteredBlockEvents,
			TxEvents:    make([][]abci.Event, len(export.TxEvents)),
		}
		for i := range cleared.TxEvents {
			cleared.TxEvents[i] = []abci.Event{}
		}

		if err := applyBlockEventsExport(stateStore, cleared); err != nil {
			return fmt.Errorf("failed to save modified finalize block response at height %d: %w", height, err)
		}

		removedCount := len(export.BlockEvents) - len(filteredBlockEvents)
		fmt.Printf("Cleaned up %d block events (kept %d oracle_prices events) and %d tx event groups at height %d\n",
			removedCount, len(filteredBlockEvents), len(export.TxEvents), height)
	}
	return nil
}

func openStateStore(cfg *cmtconfig.Config) (sm.Store, func(), error) {
	stateDB, err := cmtconfig.DefaultDBProvider(&cmtconfig.DBContext{
		ID:     "state",
		Config: cfg,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open state database: %w", err)
	}

	stateStore := sm.NewStore(stateDB, sm.StoreOptions{
		DiscardABCIResponses: false,
	})

	closeFunc := func() {
		stateDB.Close()
	}

	return stateStore, closeFunc, nil
}

func loadBlockEventsExport(stateStore sm.Store, height int64) (*BlockEventsExport, error) {
	resp, err := stateStore.LoadFinalizeBlockResponse(height)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no finalize block response found at height %d", height)
	}

	export := &BlockEventsExport{
		Height:      height,
		BlockEvents: resp.Events,
		TxEvents:    make([][]abci.Event, len(resp.TxResults)),
	}

	for i := range resp.TxResults {
		if resp.TxResults[i] != nil {
			export.TxEvents[i] = resp.TxResults[i].Events
		}
	}

	return export, nil
}

func applyBlockEventsExport(stateStore sm.Store, export BlockEventsExport) error {
	resp, err := stateStore.LoadFinalizeBlockResponse(export.Height)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("no finalize block response found at height %d", export.Height)
	}

	resp.Events = export.BlockEvents
	for i := range resp.TxResults {
		if resp.TxResults[i] != nil && i < len(export.TxEvents) {
			resp.TxResults[i].Events = export.TxEvents[i]
		}
	}

	return stateStore.SaveFinalizeBlockResponse(export.Height, resp)
}
