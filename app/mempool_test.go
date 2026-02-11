package app

import (
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/require"
)

func TestEVMMempool_ProcessProposalRejectsInvalidTxBytes(t *testing.T) {
	require.NoError(t, evmtypes.SetChainConfig(nil))
	app := SetupWithEmptyStore(t)
	require.NotNil(t, app.EVMMempool)
	req := &abci.RequestProcessProposal{
		Height: app.LastBlockHeight() + 1,
		Time:   time.Now().UTC(),
		Txs:    [][]byte{[]byte("invalid")},
	}
	resp, err := app.ProcessProposal(req)
	require.NoError(t, err)
	require.Equal(t, abci.ResponseProcessProposal_REJECT, resp.Status)
}
