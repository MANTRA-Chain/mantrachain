package app

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

const sanctionTestEVMChainID = 5888

// sweepContract stands in for a contract an attacker delegates to.
var sweepContract = common.HexToAddress("0x5eeb5eeb5eeb5eeb5eeb5eeb5eeb5eeb5eeb5eeb")

func ethAccAddr(key *ecdsa.PrivateKey) string {
	return sdk.AccAddress(crypto.PubkeyToAddress(key.PublicKey).Bytes()).String()
}

// setCodeTx builds a type-4 tx sent by relayer that carries an EIP-7702
// authorization signed by authority, delegating its code to sweepContract
// and calling it in the same tx.
func setCodeTx(t *testing.T, app *App, relayer, authority *ecdsa.PrivateKey) sdk.Tx {
	t.Helper()
	chainID := big.NewInt(sanctionTestEVMChainID)
	auth, err := ethtypes.SignSetCode(authority, ethtypes.SetCodeAuthorization{
		ChainID: *uint256.NewInt(sanctionTestEVMChainID),
		Address: sweepContract,
		Nonce:   0,
	})
	require.NoError(t, err)

	authorityAddr := crypto.PubkeyToAddress(authority.PublicKey)
	signer := ethtypes.LatestSignerForChainID(chainID)
	signed, err := ethtypes.SignNewTx(relayer, signer, &ethtypes.SetCodeTx{
		ChainID:   uint256.NewInt(sanctionTestEVMChainID),
		Nonce:     0,
		GasTipCap: uint256.NewInt(1),
		GasFeeCap: uint256.NewInt(1_000_000_000),
		Gas:       200_000,
		To:        authorityAddr,
		Value:     uint256.NewInt(0),
		AuthList:  []ethtypes.SetCodeAuthorization{auth},
	})
	require.NoError(t, err)
	require.Equal(t, uint8(ethtypes.SetCodeTxType), signed.Type())

	msg := &evmtypes.MsgEthereumTx{}
	require.NoError(t, msg.FromSignedEthereumTx(signed, signer))
	tx, err := msg.BuildTxWithEvmParams(app.TxConfig().NewTxBuilder(), evmtypes.Params{
		EvmDenom:             FeeDenom,
		ExtendedDenomOptions: &evmtypes.ExtendedDenomOptions{ExtendedDenom: FeeDenom},
	})
	require.NoError(t, err)
	return tx
}

// TestSanctionBlocksEIP7702Authorization sends a type-4 tx whose sender is
// clean but which carries an authorization signed by a sanctioned account.
// The app's ante handler must reject it, or the relayer could run code as
// the sanctioned account and move its funds.
func TestSanctionBlocksEIP7702Authorization(t *testing.T) {
	app, ctx, _, addOnly := sanctionAppCtx(t)
	relayer, err := crypto.GenerateKey()
	require.NoError(t, err)
	sanctioned, err := crypto.GenerateKey()
	require.NoError(t, err)
	clean, err := crypto.GenerateKey()
	require.NoError(t, err)

	msg := sanctiontypes.NewMsgSanctionAddresses(addOnly, []string{ethAccAddr(sanctioned)}, "rule drain-v1")
	_, err = app.MsgServiceRouter().Handler(msg)(ctx, msg)
	require.NoError(t, err)

	// Through the app's full ante handler.
	_, err = app.AnteHandler()(ctx, setCodeTx(t, app, relayer, sanctioned), false)
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)
	require.ErrorContains(t, err, "authorization signer")

	// The same tx with a clean authority passes the sanction check.
	dec := sanctionkeeper.NewEVMBlacklistCheckDecorator(app.SanctionKeeper)
	next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }
	_, err = dec.AnteHandle(ctx, setCodeTx(t, app, relayer, clean), false, next)
	require.NoError(t, err)

	// A sanctioned relayer is still rejected as the sender.
	_, err = dec.AnteHandle(ctx, setCodeTx(t, app, sanctioned, clean), false, next)
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)
	require.ErrorContains(t, err, "sender")
}

// setDelegation installs an EIP-7702 delegation on addr, as a type-4 tx
// would have before the account was sanctioned.
func setDelegation(app *App, ctx sdk.Context, addr common.Address) {
	code := ethtypes.AddressToDelegation(sweepContract)
	hash := crypto.Keccak256(code)
	app.EVMKeeper.SetCode(ctx, hash, code)
	app.EVMKeeper.SetCodeHash(ctx, addr.Bytes(), hash)
}

func hasDelegation(app *App, ctx sdk.Context, addr common.Address) bool {
	_, ok := ethtypes.ParseDelegation(app.EVMKeeper.GetCode(ctx, app.EVMKeeper.GetCodeHash(ctx, addr)))
	return ok
}

// TestSanctionClearsEIP7702Delegation checks sanctioning an account removes a
// delegation it set up earlier, so nobody can call into it to run the
// delegated code. Real contract code is left alone.
func TestSanctionClearsEIP7702Delegation(t *testing.T) {
	app, ctx, _, addOnly := sanctionAppCtx(t)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	delegated := crypto.PubkeyToAddress(key.PublicKey)
	setDelegation(app, ctx, delegated)
	require.True(t, hasDelegation(app, ctx, delegated))

	contract := common.HexToAddress("0xc0de00000000000000000000000000000000c0de")
	contractCode := []byte{0x60, 0x80, 0x60, 0x40, 0x52}
	contractHash := crypto.Keccak256(contractCode)
	app.EVMKeeper.SetCode(ctx, contractHash, contractCode)
	app.EVMKeeper.SetCodeHash(ctx, contract.Bytes(), contractHash)

	msg := sanctiontypes.NewMsgSanctionAddresses(addOnly, []string{
		sdk.AccAddress(delegated.Bytes()).String(),
		sdk.AccAddress(contract.Bytes()).String(),
	}, "rule drain-v1")
	_, err = app.MsgServiceRouter().Handler(msg)(ctx, msg)
	require.NoError(t, err)

	require.False(t, hasDelegation(app, ctx, delegated))
	require.Equal(t, common.BytesToHash(evmtypes.EmptyCodeHash), app.EVMKeeper.GetCodeHash(ctx, delegated))
	require.Equal(t, common.BytesToHash(contractHash), app.EVMKeeper.GetCodeHash(ctx, contract), "contract code untouched")

	// Governance sanctions go through the same path.
	key2, err := crypto.GenerateKey()
	require.NoError(t, err)
	delegated2 := crypto.PubkeyToAddress(key2.PublicKey)
	setDelegation(app, ctx, delegated2)
	gmsg := sanctiontypes.NewMsgSanctionAddresses(authtypes.NewModuleAddress(govtypes.ModuleName).String(), []string{sdk.AccAddress(delegated2.Bytes()).String()}, "gov")
	_, err = app.MsgServiceRouter().Handler(gmsg)(ctx, gmsg)
	require.NoError(t, err)
	require.False(t, hasDelegation(app, ctx, delegated2))
}

// TestMigrate1to2ClearsEIP7702Delegation checks the v8.9.0 migration clears
// delegations from accounts that were blacklisted before the upgrade.
func TestMigrate1to2ClearsEIP7702Delegation(t *testing.T) {
	app, ctx, _, _ := sanctionAppCtx(t)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	listed := crypto.PubkeyToAddress(key.PublicKey)
	key2, err := crypto.GenerateKey()
	require.NoError(t, err)
	notListed := crypto.PubkeyToAddress(key2.PublicKey)

	setDelegation(app, ctx, listed)
	setDelegation(app, ctx, notListed)
	require.NoError(t, app.SanctionKeeper.BlacklistAccounts.Set(ctx, sdk.AccAddress(listed.Bytes()).String()))

	require.NoError(t, sanctionkeeper.NewMigrator(app.SanctionKeeper).Migrate1to2(ctx))

	require.False(t, hasDelegation(app, ctx, listed))
	require.True(t, hasDelegation(app, ctx, notListed), "accounts not on the blacklist keep their delegation")
	params, err := app.SanctionKeeper.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, sanctiontypes.DefaultParams(), params)
}
