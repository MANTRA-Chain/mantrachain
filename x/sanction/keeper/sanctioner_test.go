package keeper_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/MANTRA-Chain/mantrachain/v8/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type mockAccountKeeper struct {
	accs map[string]sdk.AccountI
}

func (m mockAccountKeeper) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	return m.accs[addr.String()]
}

type fixture struct {
	k   keeper.Keeper
	ctx sdk.Context
	srv types.MsgServer
	qs  types.QueryServer

	gov       string
	full      string
	full2     string
	addOnly   string
	moduleAcc string // a module account known to x/auth
	reserved  string // a module address passed to the keeper, not in x/auth
}

func addr(seed string) string {
	b := []byte(seed + strings.Repeat("_", 20))
	return sdk.AccAddress(b[:20]).String()
}

func setup(t *testing.T) *fixture {
	t.Helper()
	modAcc := authtypes.NewEmptyModuleAccount("somemodule")
	reserved := authtypes.NewModuleAddress("notyetcreated").String()
	ak := mockAccountKeeper{accs: map[string]sdk.AccountI{modAcc.GetAddress().String(): modAcc}}

	k, ctx, _ := keepertest.SanctionKeeperWithAccounts(t, ak, []string{reserved})
	f := &fixture{
		k:         k,
		ctx:       ctx,
		srv:       keeper.NewMsgServerImpl(k),
		qs:        keeper.NewQueryServerImpl(k),
		gov:       k.GetAuthority(),
		full:      addr("full"),
		full2:     addr("full2"),
		addOnly:   addr("addonly"),
		moduleAcc: modAcc.GetAddress().String(),
		reserved:  reserved,
	}
	_, err := f.srv.UpdateSanctioners(ctx, types.NewMsgUpdateSanctioners(f.gov, []types.Sanctioner{
		{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
		{Address: f.full2, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
		{Address: f.addOnly, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
	}))
	require.NoError(t, err)
	return f
}

// freshCtx returns the fixture context with a new event manager.
func (f *fixture) freshCtx() sdk.Context {
	return f.ctx.WithEventManager(sdk.NewEventManager())
}

func (f *fixture) sanction(signer string, addrs ...string) error {
	_, err := f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(signer, addrs, "rule 7, tx ABC"))
	return err
}

func (f *fixture) unsanction(signer string, addrs ...string) error {
	_, err := f.srv.UnsanctionAddresses(f.ctx, types.NewMsgUnsanctionAddresses(signer, addrs, "false positive"))
	return err
}

func (f *fixture) listed(t *testing.T, a string) bool {
	t.Helper()
	has, err := f.k.BlacklistAccounts.Has(f.ctx, a)
	require.NoError(t, err)
	return has
}

func listChangedEvents(t *testing.T, ctx sdk.Context) []*types.EventSanctionListChanged {
	t.Helper()
	var out []*types.EventSanctionListChanged
	for _, e := range ctx.EventManager().ABCIEvents() {
		msg, err := sdk.ParseTypedEvent(e)
		require.NoError(t, err)
		if ev, ok := msg.(*types.EventSanctionListChanged); ok {
			out = append(out, ev)
		}
	}
	return out
}

func TestUpdateSanctioners(t *testing.T) {
	f := setup(t)

	_, err := f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.full, nil))
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest, "only governance can update the list")

	got, err := f.k.GetSanctioners(f.ctx)
	require.NoError(t, err)
	require.Len(t, got, 3)

	// Replacing the list drops anyone left out.
	_, err = f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, []types.Sanctioner{
		{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
	}))
	require.NoError(t, err)
	got, err = f.k.GetSanctioners(f.ctx)
	require.NoError(t, err)
	require.Equal(t, []types.Sanctioner{{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL}}, got)
	require.ErrorIs(t, f.sanction(f.addOnly, addr("x")), types.ErrNotSanctioner)

	// An empty list turns the fast path off.
	_, err = f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, nil))
	require.NoError(t, err)
	require.ErrorIs(t, f.sanction(f.full, addr("x")), types.ErrNotSanctioner)
}

func TestSanctionAddresses(t *testing.T) {
	f := setup(t)
	a, b := addr("attacker"), addr("attacker2")

	t.Run("non-sanctioner rejected", func(t *testing.T) {
		require.ErrorIs(t, f.sanction(addr("rando"), a), types.ErrNotSanctioner)
		require.False(t, f.listed(t, a))
	})

	t.Run("add-only sanctioner adds, with record and event", func(t *testing.T) {
		f.ctx = f.freshCtx().WithBlockHeight(77)
		require.NoError(t, f.sanction(f.addOnly, a))
		require.True(t, f.listed(t, a))

		rec, err := f.k.SanctionRecords.Get(f.ctx, a)
		require.NoError(t, err)
		require.Equal(t, types.SanctionRecord{
			Actor: f.addOnly, ActorRole: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY, Reason: "rule 7, tx ABC", Height: 77,
		}, rec)

		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1)
		require.Equal(t, &types.EventSanctionListChanged{
			Address: a, Sanctioned: true, Actor: f.addOnly,
			ActorRole: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY, Reason: "rule 7, tx ABC", Height: 77,
		}, evs[0])
	})

	t.Run("already-listed and in-batch duplicates are no-ops", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.sanction(f.full, a, b, b))
		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1, "only b is new")
		require.Equal(t, b, evs[0].Address)

		// The original record is kept, not overwritten.
		rec, err := f.k.SanctionRecords.Get(f.ctx, a)
		require.NoError(t, err)
		require.Equal(t, f.addOnly, rec.Actor)
	})

	t.Run("mixed-case input is stored in canonical form", func(t *testing.T) {
		c := addr("uppercase")
		require.NoError(t, f.sanction(f.full, strings.ToUpper(c)))
		require.True(t, f.listed(t, c))
	})
}

func TestSanctionProtections(t *testing.T) {
	f := setup(t)
	gone := addr("innocent")

	cases := map[string]string{
		"full sanctioner":                f.full,
		"add-only sanctioner":            f.addOnly,
		"gov authority":                  f.gov,
		"module account in x/auth":       f.moduleAcc,
		"module address not yet created": f.reserved,
	}
	for name, target := range cases {
		t.Run(name, func(t *testing.T) {
			for _, signer := range []string{f.full, f.addOnly} {
				err := f.sanction(signer, gone, target)
				require.ErrorIs(t, err, types.ErrProtectedAddress)
			}
			// Governance can still sanction it.
			_, err := f.srv.AddBlacklistAccounts(f.ctx, types.NewMsgAddBlacklistAccounts(f.gov, []string{target}))
			require.NoError(t, err)
			_, err = f.srv.RemoveBlacklistAccounts(f.ctx, types.NewMsgRemoveBlacklistAccounts(f.gov, []string{target}))
			require.NoError(t, err)
		})
	}
	// The message is atomic at the tx level, and the keeper checks each
	// protected address before writing it. Here the innocent address came
	// first in the batch, so the caller's tx would be reverted in full.
}

func TestLimits(t *testing.T) {
	f := setup(t)
	_, err := f.srv.UpdateParams(f.ctx, types.NewMsgUpdateParams(f.gov, types.NewParams(2, 20)))
	require.NoError(t, err)

	err = f.sanction(f.full, addr("a"), addr("b"), addr("c"))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)
	require.NoError(t, f.sanction(f.full, addr("a"), addr("b")))

	err = f.unsanction(f.full, addr("a"), addr("b"), addr("c"))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)

	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, []string{addr("d")}, strings.Repeat("x", 21)))
	require.ErrorIs(t, err, types.ErrReasonTooLong)

	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, []string{addr("d")}, ""))
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)

	// Governance messages aren't capped.
	many := []string{addr("g1"), addr("g2"), addr("g3")}
	_, err = f.srv.AddBlacklistAccounts(f.ctx, types.NewMsgAddBlacklistAccounts(f.gov, many))
	require.NoError(t, err)
}

func TestUpdateParams(t *testing.T) {
	f := setup(t)
	_, err := f.srv.UpdateParams(f.ctx, types.NewMsgUpdateParams(f.full, types.NewParams(5, 5)))
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	_, err = f.srv.UpdateParams(f.ctx, types.NewMsgUpdateParams(f.gov, types.NewParams(0, 5)))
	require.Error(t, err)
	_, err = f.srv.UpdateParams(f.ctx, types.NewMsgUpdateParams(f.gov, types.NewParams(5, 5)))
	require.NoError(t, err)
	res, err := f.qs.Params(f.ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, types.NewParams(5, 5), res.Params)
}

func TestUnsanctionAddresses(t *testing.T) {
	f := setup(t)
	byBot, byFull, byGov, legacy, unlisted := addr("bybot"), addr("byfull"), addr("bygov"), addr("legacy"), addr("unlisted")
	require.NoError(t, f.sanction(f.addOnly, byBot))
	require.NoError(t, f.sanction(f.full2, byFull))
	_, err := f.srv.AddBlacklistAccounts(f.ctx, types.NewMsgAddBlacklistAccounts(f.gov, []string{byGov}))
	require.NoError(t, err)
	require.NoError(t, f.k.BlacklistAccounts.Set(f.ctx, legacy)) // predates sanction records

	t.Run("add-only sanctioner cannot remove", func(t *testing.T) {
		require.ErrorIs(t, f.unsanction(f.addOnly, byBot), types.ErrRoleNotPermitted)
		require.True(t, f.listed(t, byBot))
	})

	t.Run("non-sanctioner cannot remove", func(t *testing.T) {
		require.ErrorIs(t, f.unsanction(addr("rando"), byBot), types.ErrNotSanctioner)
	})

	t.Run("full sanctioner removes entries from a bot and another full sanctioner", func(t *testing.T) {
		f.ctx = f.freshCtx().WithBlockHeight(90)
		require.NoError(t, f.unsanction(f.full, byBot, byFull, unlisted))
		require.False(t, f.listed(t, byBot))
		require.False(t, f.listed(t, byFull))
		has, err := f.k.SanctionRecords.Has(f.ctx, byBot)
		require.NoError(t, err)
		require.False(t, has, "record is dropped with the entry")

		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 2, "unlisted address is a no-op")
		require.Equal(t, &types.EventSanctionListChanged{
			Address: byBot, Sanctioned: false, Actor: f.full,
			ActorRole: types.SanctionerRole_SANCTIONER_ROLE_FULL, Reason: "false positive", Height: 90,
		}, evs[0])
	})

	t.Run("governance and legacy entries are protected", func(t *testing.T) {
		require.ErrorIs(t, f.unsanction(f.full, byGov), types.ErrGovernanceEntry)
		require.ErrorIs(t, f.unsanction(f.full, legacy), types.ErrGovernanceEntry)
		require.True(t, f.listed(t, byGov))
		require.True(t, f.listed(t, legacy))
	})

	t.Run("governance can reverse anything", func(t *testing.T) {
		require.NoError(t, f.sanction(f.addOnly, byBot))
		_, err := f.srv.RemoveBlacklistAccounts(f.ctx, types.NewMsgRemoveBlacklistAccounts(f.gov, []string{byBot, byGov, legacy}))
		require.NoError(t, err)
		for _, a := range []string{byBot, byGov, legacy} {
			require.False(t, f.listed(t, a))
		}
	})
}

func TestRevokeSanctioner(t *testing.T) {
	f := setup(t)
	revoke := func(signer, target string) error {
		_, err := f.srv.RevokeSanctioner(f.ctx, types.NewMsgRevokeSanctioner(signer, target, "key leaked"))
		return err
	}

	require.ErrorIs(t, revoke(f.addOnly, f.addOnly), types.ErrRoleNotPermitted, "add-only cannot revoke")
	require.ErrorIs(t, revoke(f.full, f.full2), types.ErrRoleNotPermitted, "full sanctioners are removed by governance only")
	require.ErrorIs(t, revoke(f.full, addr("rando")), types.ErrNotSanctioner)
	require.ErrorIs(t, revoke(addr("rando"), f.addOnly), types.ErrNotSanctioner)

	f.ctx = f.freshCtx()
	require.NoError(t, revoke(f.full, f.addOnly))
	has, err := f.k.Sanctioners.Has(f.ctx, f.addOnly)
	require.NoError(t, err)
	require.False(t, has)

	var revoked *types.EventSanctionerRevoked
	for _, e := range f.ctx.EventManager().ABCIEvents() {
		msg, err := sdk.ParseTypedEvent(e)
		require.NoError(t, err)
		if ev, ok := msg.(*types.EventSanctionerRevoked); ok {
			revoked = ev
		}
	}
	require.NotNil(t, revoked)
	require.Equal(t, f.addOnly, revoked.Address)
	require.Equal(t, f.full, revoked.RevokedBy)

	require.ErrorIs(t, f.sanction(f.addOnly, addr("x")), types.ErrNotSanctioner, "revoked key can no longer sanction")
}

// TestStolenAddOnlyKey walks through the stolen-key example from the spec.
func TestStolenAddOnlyKey(t *testing.T) {
	f := setup(t)
	attacker, victim := addr("attacker"), addr("victim")
	require.NoError(t, f.sanction(f.full, attacker))

	// The attacker holds the add-only key: they can block others but can't
	// unblock themselves.
	require.NoError(t, f.sanction(f.addOnly, victim))
	require.ErrorIs(t, f.unsanction(f.addOnly, attacker), types.ErrRoleNotPermitted)

	// A full sanctioner revokes the key and unblocks the victim.
	_, err := f.srv.RevokeSanctioner(f.ctx, types.NewMsgRevokeSanctioner(f.full, f.addOnly, "key stolen"))
	require.NoError(t, err)
	require.NoError(t, f.unsanction(f.full, victim))

	require.True(t, f.listed(t, attacker))
	require.False(t, f.listed(t, victim))
	require.ErrorIs(t, f.sanction(f.addOnly, addr("another")), types.ErrNotSanctioner)
}

func TestGovernanceRecordsActor(t *testing.T) {
	f := setup(t)
	a := addr("govadded")
	f.ctx = f.freshCtx().WithBlockHeight(12)
	_, err := f.srv.AddBlacklistAccounts(f.ctx, types.NewMsgAddBlacklistAccounts(f.gov, []string{a}))
	require.NoError(t, err)

	rec, err := f.k.SanctionRecords.Get(f.ctx, a)
	require.NoError(t, err)
	require.Equal(t, types.SanctionRecord{Actor: f.gov, ActorRole: types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, Height: 12}, rec)

	evs := listChangedEvents(t, f.ctx)
	require.Len(t, evs, 1)
	require.Equal(t, f.gov, evs[0].Actor)
	require.True(t, evs[0].Sanctioned)

	// Existing governance behaviour is unchanged: duplicates still fail.
	_, err = f.srv.AddBlacklistAccounts(f.ctx, types.NewMsgAddBlacklistAccounts(f.gov, []string{a}))
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
}

func TestQueries(t *testing.T) {
	f := setup(t)
	byBot, legacy, clean := addr("bybot"), addr("legacy"), addr("clean")
	f.ctx = f.ctx.WithBlockHeight(5)
	require.NoError(t, f.sanction(f.addOnly, byBot))
	require.NoError(t, f.k.BlacklistAccounts.Set(f.ctx, legacy))

	sres, err := f.qs.Sanctioners(f.ctx, &types.QuerySanctionersRequest{})
	require.NoError(t, err)
	require.Len(t, sres.Sanctioners, 3)

	ares, err := f.qs.SanctionedAddresses(f.ctx, &types.QuerySanctionedAddressesRequest{})
	require.NoError(t, err)
	require.Len(t, ares.SanctionedAddresses, 2)
	got := map[string]*types.SanctionRecord{}
	for _, e := range ares.SanctionedAddresses {
		got[e.Address] = e.Record
	}
	require.Nil(t, got[legacy])
	require.Equal(t, f.addOnly, got[byBot].Actor)
	require.Equal(t, int64(5), got[byBot].Height)

	ires, err := f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{Address: byBot})
	require.NoError(t, err)
	require.True(t, ires.Sanctioned)
	require.Equal(t, "rule 7, tx ABC", ires.Record.Reason)

	ires, err = f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{Address: legacy})
	require.NoError(t, err)
	require.True(t, ires.Sanctioned)
	require.Nil(t, ires.Record)

	ires, err = f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{Address: clean})
	require.NoError(t, err)
	require.False(t, ires.Sanctioned)

	_, err = f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{Address: "nope"})
	require.Error(t, err)
}

func TestMigrate1to2(t *testing.T) {
	k, ctx, _ := keepertest.SanctionKeeper(t)
	// Version 1 state: empty params and bare blacklist entries.
	require.NoError(t, k.Params.Set(ctx, types.Params{}))
	old := []string{addr("old1"), addr("old2")}
	for _, a := range old {
		require.NoError(t, k.BlacklistAccounts.Set(ctx, a))
	}

	require.NoError(t, keeper.NewMigrator(k).Migrate1to2(ctx))

	params, err := k.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), params)
	for _, a := range old {
		has, err := k.BlacklistAccounts.Has(ctx, a)
		require.NoError(t, err)
		require.True(t, has)
		hasRec, err := k.SanctionRecords.Has(ctx, a)
		require.NoError(t, err)
		require.False(t, hasRec)
	}
	sanctioners, err := k.GetSanctioners(ctx)
	require.NoError(t, err)
	require.Empty(t, sanctioners)
}
