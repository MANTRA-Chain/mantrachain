package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/MANTRA-Chain/mantrachain/v8/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type fixture struct {
	k   keeper.Keeper
	ctx sdk.Context
	srv types.MsgServer
	qs  types.QueryServer

	gov     string
	full    string
	full2   string
	addOnly string
}

func addr(seed string) string {
	b := []byte(seed + strings.Repeat("_", 20))
	return sdk.AccAddress(b[:20]).String()
}

func setup(t *testing.T) *fixture {
	t.Helper()
	k, ctx, _ := keepertest.SanctionKeeper(t)
	f := &fixture{
		k:       k,
		ctx:     ctx,
		srv:     keeper.NewMsgServerImpl(k),
		qs:      keeper.NewQueryServerImpl(k),
		gov:     k.GetAuthority(),
		full:    addr("full"),
		full2:   addr("full2"),
		addOnly: addr("addonly"),
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

	t.Run("add-only sanctioner adds, with event", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.sanction(f.addOnly, a))
		require.True(t, f.listed(t, a))

		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1)
		require.Equal(t, &types.EventSanctionListChanged{
			Addresses: []string{a}, Sanctioned: true, Actor: f.addOnly,
			ActorRole: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY, Reason: "rule 7, tx ABC",
		}, evs[0])
	})

	t.Run("already-listed and in-batch duplicates are no-ops", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.sanction(f.full, a, b, b))
		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1, "only b is new")
		require.Equal(t, []string{b}, evs[0].Addresses)
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
		"full sanctioner":     f.full,
		"add-only sanctioner": f.addOnly,
	}
	for name, target := range cases {
		t.Run(name, func(t *testing.T) {
			for _, signer := range []string{f.full, f.addOnly} {
				err := f.sanction(signer, gone, target)
				require.ErrorIs(t, err, types.ErrProtectedAddress)
			}
			// Governance can still sanction it.
			require.NoError(t, f.sanction(f.gov, target))
			require.True(t, f.listed(t, target))
			require.NoError(t, f.unsanction(f.gov, target))
		})
	}
	// The message is atomic at the tx level, and the keeper checks each
	// protected address before writing it. Here the innocent address came
	// first in the batch, so the caller's tx would be reverted in full.
}

// TestSanctionModuleAccountRejected checks nobody, governance included, can
// sanction a module account, since the send restriction would stop it
// sending. The whole batch fails and nothing is written.
func TestSanctionModuleAccountRejected(t *testing.T) {
	f := setup(t)
	for _, name := range keepertest.SanctionModuleAccounts {
		module := authtypes.NewModuleAddress(name).String()
		for _, signer := range []string{f.gov, f.full, f.addOnly} {
			err := f.sanction(signer, addr("innocent"), module)
			require.ErrorIs(t, err, types.ErrProtectedAddress, "%s by %s", name, signer)
			require.ErrorContains(t, err, "module account")
			require.False(t, f.listed(t, module))
			require.False(t, f.listed(t, addr("innocent")))
		}
	}
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

	// Governance is capped the same way and also needs a reason.
	many := []string{addr("g1"), addr("g2"), addr("g3")}
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.gov, many, "r"))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)
	_, err = f.srv.UnsanctionAddresses(f.ctx, types.NewMsgUnsanctionAddresses(f.gov, many, "r"))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.gov, many[:2], strings.Repeat("x", 21)))
	require.ErrorIs(t, err, types.ErrReasonTooLong)
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.gov, many[:2], ""))
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	require.NoError(t, f.sanction(f.gov, many[:2]...))
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
	require.NoError(t, f.sanction(f.gov, byGov))
	require.NoError(t, f.k.BlacklistAccounts.Set(f.ctx, legacy)) // predates sanctioners

	t.Run("add-only sanctioner cannot remove", func(t *testing.T) {
		require.ErrorIs(t, f.unsanction(f.addOnly, byBot), types.ErrRoleNotPermitted)
		require.True(t, f.listed(t, byBot))
	})

	t.Run("non-sanctioner cannot remove", func(t *testing.T) {
		require.ErrorIs(t, f.unsanction(addr("rando"), byBot), types.ErrNotSanctioner)
	})

	t.Run("full sanctioner removes any entry, whoever added it", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.unsanction(f.full, byBot, byFull, byGov, legacy, unlisted))
		for _, a := range []string{byBot, byFull, byGov, legacy} {
			require.False(t, f.listed(t, a))
		}

		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1, "one event per message")
		require.Equal(t, &types.EventSanctionListChanged{
			Addresses: []string{byBot, byFull, byGov, legacy}, Sanctioned: false, Actor: f.full,
			ActorRole: types.SanctionerRole_SANCTIONER_ROLE_FULL, Reason: "false positive",
		}, evs[0], "unlisted address is left out")
	})

	t.Run("governance removes any entry, whoever added it", func(t *testing.T) {
		require.NoError(t, f.sanction(f.addOnly, byBot))
		require.NoError(t, f.sanction(f.full, byFull))
		require.NoError(t, f.sanction(f.gov, byGov))
		require.NoError(t, f.k.BlacklistAccounts.Set(f.ctx, legacy))
		require.NoError(t, f.unsanction(f.gov, byBot, byFull, byGov, legacy))
		for _, a := range []string{byBot, byFull, byGov, legacy} {
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

func TestGovernanceUsesSameMessages(t *testing.T) {
	f := setup(t)
	a := addr("govadded")
	f.ctx = f.freshCtx()
	require.NoError(t, f.sanction(f.gov, a))
	require.True(t, f.listed(t, a))

	evs := listChangedEvents(t, f.ctx)
	require.Len(t, evs, 1)
	require.Equal(t, &types.EventSanctionListChanged{
		Addresses: []string{a}, Sanctioned: true, Actor: f.gov,
		ActorRole: types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, Reason: "rule 7, tx ABC",
	}, evs[0])

	// Duplicates are a no-op for governance too.
	f.ctx = f.freshCtx()
	require.NoError(t, f.sanction(f.gov, a))
	require.Empty(t, listChangedEvents(t, f.ctx))

	// Governance removes, with an event naming it.
	f.ctx = f.freshCtx()
	require.NoError(t, f.unsanction(f.gov, a))
	evs = listChangedEvents(t, f.ctx)
	require.Len(t, evs, 1)
	require.False(t, evs[0].Sanctioned)
	require.Equal(t, f.gov, evs[0].Actor)
}

func TestQueries(t *testing.T) {
	f := setup(t)
	byBot, clean := addr("bybot"), addr("clean")
	require.NoError(t, f.sanction(f.addOnly, byBot))

	sres, err := f.qs.Sanctioners(f.ctx, &types.QuerySanctionersRequest{})
	require.NoError(t, err)
	require.Len(t, sres.Sanctioners, 3)

	ires, err := f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{Address: byBot})
	require.NoError(t, err)
	require.True(t, ires.Sanctioned)

	// Mixed case resolves to the stored canonical form.
	ires, err = f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{Address: strings.ToUpper(byBot)})
	require.NoError(t, err)
	require.True(t, ires.Sanctioned)

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
	// A module account listed under version 1 is dropped.
	module := authtypes.NewModuleAddress(distrtypes.ModuleName).String()
	require.NoError(t, k.BlacklistAccounts.Set(ctx, module))

	require.NoError(t, keeper.NewMigrator(k).Migrate1to2(ctx))

	params, err := k.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), params)
	for _, a := range old {
		has, err := k.BlacklistAccounts.Has(ctx, a)
		require.NoError(t, err)
		require.True(t, has)
	}
	has, err := k.BlacklistAccounts.Has(ctx, module)
	require.NoError(t, err)
	require.False(t, has, "module account removed")
	sanctioners, err := k.GetSanctioners(ctx)
	require.NoError(t, err)
	require.Empty(t, sanctioners)
}

func revokedEvents(t *testing.T, ctx sdk.Context) []*types.EventSanctionerRevoked {
	t.Helper()
	var out []*types.EventSanctionerRevoked
	for _, e := range ctx.EventManager().ABCIEvents() {
		msg, err := sdk.ParseTypedEvent(e)
		require.NoError(t, err)
		if ev, ok := msg.(*types.EventSanctionerRevoked); ok {
			out = append(out, ev)
		}
	}
	return out
}

func (f *fixture) revoke(signer, target, reason string) error {
	_, err := f.srv.RevokeSanctioner(f.ctx, types.NewMsgRevokeSanctioner(signer, target, reason))
	return err
}

func (f *fixture) setParams(t *testing.T, maxAddrs, maxReason uint32) {
	t.Helper()
	_, err := f.srv.UpdateParams(f.ctx, types.NewMsgUpdateParams(f.gov, types.NewParams(maxAddrs, maxReason)))
	require.NoError(t, err)
}

// TestSanctionIsAllOrNothing checks a failing address anywhere in the batch
// leaves the whole list untouched, even without the tx-level revert.
func TestSanctionIsAllOrNothing(t *testing.T) {
	f := setup(t)
	innocent, later := addr("innocent"), addr("later")

	f.ctx = f.freshCtx()
	require.ErrorIs(t, f.sanction(f.full, innocent, f.addOnly, later), types.ErrProtectedAddress)
	require.False(t, f.listed(t, innocent))
	require.False(t, f.listed(t, later))
	require.Empty(t, listChangedEvents(t, f.ctx))

	// An invalid address reaching the keeper directly, bypassing Validate.
	err := f.k.SanctionAddresses(f.ctx, []string{innocent, "not-an-address"}, f.full, "r")
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	require.False(t, f.listed(t, innocent))

	require.NoError(t, f.sanction(f.full, innocent))
	err = f.k.UnsanctionAddresses(f.ctx, []string{innocent, "not-an-address"}, f.full, "r")
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	require.True(t, f.listed(t, innocent), "nothing removed")
}

// TestLimitBoundaries checks each cap accepts exactly the limit and rejects
// one more, for every message that has a reason.
func TestLimitBoundaries(t *testing.T) {
	f := setup(t)
	f.setParams(t, 3, 10)
	atCap := []string{addr("c1"), addr("c2"), addr("c3")}
	overCap := append([]string{addr("c4")}, atCap...)
	reasonAtCap, reasonOver := strings.Repeat("r", 10), strings.Repeat("r", 11)

	_, err := f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, overCap, reasonAtCap))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, atCap, reasonOver))
	require.ErrorIs(t, err, types.ErrReasonTooLong)
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, atCap, reasonAtCap))
	require.NoError(t, err)

	_, err = f.srv.UnsanctionAddresses(f.ctx, types.NewMsgUnsanctionAddresses(f.full, overCap, reasonAtCap))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)
	_, err = f.srv.UnsanctionAddresses(f.ctx, types.NewMsgUnsanctionAddresses(f.full, atCap, reasonOver))
	require.ErrorIs(t, err, types.ErrReasonTooLong)
	_, err = f.srv.UnsanctionAddresses(f.ctx, types.NewMsgUnsanctionAddresses(f.full, atCap, reasonAtCap))
	require.NoError(t, err)

	require.ErrorIs(t, f.revoke(f.full, f.addOnly, reasonOver), types.ErrReasonTooLong)
	require.NoError(t, f.revoke(f.full, f.addOnly, reasonAtCap))

	// Duplicates count toward the cap: it limits message size, not changes.
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, []string{atCap[0], atCap[0], atCap[0], atCap[0]}, "r"))
	require.ErrorIs(t, err, types.ErrTooManyAddresses)

	// A params change takes effect on the next message.
	f.setParams(t, 4, 10)
	_, err = f.srv.SanctionAddresses(f.ctx, types.NewMsgSanctionAddresses(f.full, overCap, reasonAtCap))
	require.NoError(t, err)
}

func TestUnsanctionEvents(t *testing.T) {
	f := setup(t)
	a, b := addr("ev-a"), addr("ev-b")
	require.NoError(t, f.sanction(f.full, a, b))

	t.Run("only unlisted addresses emits nothing", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.unsanction(f.full, addr("never1"), addr("never2")))
		require.Empty(t, listChangedEvents(t, f.ctx))
	})

	t.Run("mixed-case input is reported in canonical form", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.unsanction(f.full, strings.ToUpper(a)))
		require.False(t, f.listed(t, a))
		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1)
		require.Equal(t, []string{a}, evs[0].Addresses)
	})

	t.Run("governance removal has no role", func(t *testing.T) {
		f.ctx = f.freshCtx()
		require.NoError(t, f.unsanction(f.gov, b))
		require.Equal(t, []*types.EventSanctionListChanged{{
			Addresses: []string{b}, Sanctioned: false, Actor: f.gov,
			ActorRole: types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, Reason: "false positive",
		}}, listChangedEvents(t, f.ctx))
	})

	t.Run("in-batch duplicates are removed once", func(t *testing.T) {
		require.NoError(t, f.sanction(f.full, a))
		f.ctx = f.freshCtx()
		require.NoError(t, f.unsanction(f.full, a, a))
		evs := listChangedEvents(t, f.ctx)
		require.Len(t, evs, 1)
		require.Equal(t, []string{a}, evs[0].Addresses)
	})
}

// TestRoleChanges checks permissions follow the current list, so a role
// change by governance applies to the next message.
func TestRoleChanges(t *testing.T) {
	f := setup(t)
	target := addr("rolechange")
	require.NoError(t, f.sanction(f.full, target))

	// Demote full to add-only: it can still add, but no longer remove or
	// revoke.
	_, err := f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, []types.Sanctioner{
		{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
		{Address: f.addOnly, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
	}))
	require.NoError(t, err)
	require.NoError(t, f.sanction(f.full, addr("still-adds")))
	require.ErrorIs(t, f.unsanction(f.full, target), types.ErrRoleNotPermitted)
	require.ErrorIs(t, f.revoke(f.full, f.addOnly, "r"), types.ErrRoleNotPermitted)

	// Promote add-only to full: it can now remove and revoke.
	_, err = f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, []types.Sanctioner{
		{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
		{Address: f.addOnly, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
	}))
	require.NoError(t, err)
	require.NoError(t, f.unsanction(f.addOnly, target))
	require.NoError(t, f.revoke(f.addOnly, f.full, "r"))
}

func TestRevokeSanctionerEdges(t *testing.T) {
	f := setup(t)
	flagged := addr("flaggedbybot")
	require.NoError(t, f.sanction(f.addOnly, flagged))

	require.ErrorIs(t, f.revoke(f.gov, f.addOnly, "r"), types.ErrNotSanctioner, "governance uses MsgUpdateSanctioners instead")
	require.ErrorIs(t, f.revoke(f.full, f.full, "r"), types.ErrRoleNotPermitted, "a full sanctioner can't revoke itself")

	// Mixed-case target resolves to the stored entry.
	f.ctx = f.freshCtx()
	require.NoError(t, f.revoke(f.full, strings.ToUpper(f.addOnly), "key rotated"))
	require.Equal(t, []*types.EventSanctionerRevoked{{
		Address: f.addOnly, RevokedBy: f.full, Reason: "key rotated", Height: f.ctx.BlockHeight(),
	}}, revokedEvents(t, f.ctx))

	require.ErrorIs(t, f.revoke(f.full, f.addOnly, "r"), types.ErrNotSanctioner, "already revoked")
	require.True(t, f.listed(t, flagged), "revoking leaves the sanctioner's entries in place")

	// The revoked address is no longer protected, so it can be sanctioned.
	require.NoError(t, f.sanction(f.full, f.addOnly))

	// Governance can appoint it again.
	_, err := f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, []types.Sanctioner{
		{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
		{Address: f.addOnly, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
	}))
	require.NoError(t, err)
	require.NoError(t, f.sanction(f.addOnly, addr("back")))
}

func TestUpdateSanctionersMsg(t *testing.T) {
	f := setup(t)
	before, err := f.k.GetSanctioners(f.ctx)
	require.NoError(t, err)

	for name, list := range map[string][]types.Sanctioner{
		"duplicate":        {{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL}, {Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY}},
		"unspecified role": {{Address: f.full}},
		"invalid address":  {{Address: "nope", Role: types.SanctionerRole_SANCTIONER_ROLE_FULL}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, list))
			require.Error(t, err)
			after, err := f.k.GetSanctioners(f.ctx)
			require.NoError(t, err)
			require.Equal(t, before, after, "list unchanged")
		})
	}

	t.Run("stores canonical addresses and emits the new list", func(t *testing.T) {
		f.ctx = f.freshCtx()
		_, err := f.srv.UpdateSanctioners(f.ctx, types.NewMsgUpdateSanctioners(f.gov, []types.Sanctioner{
			{Address: strings.ToUpper(f.full), Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
		}))
		require.NoError(t, err)
		want := []types.Sanctioner{{Address: f.full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL}}
		got, err := f.k.GetSanctioners(f.ctx)
		require.NoError(t, err)
		require.Equal(t, want, got)
		require.NoError(t, f.sanction(f.full, addr("canon")), "signer matches the stored form")

		var ev *types.EventSanctionersUpdated
		for _, e := range f.ctx.EventManager().ABCIEvents() {
			msg, err := sdk.ParseTypedEvent(e)
			require.NoError(t, err)
			if u, ok := msg.(*types.EventSanctionersUpdated); ok {
				ev = u
			}
		}
		require.NotNil(t, ev)
		require.Equal(t, want, ev.Sanctioners)
	})
}

// TestKeeperDirectCalls covers another module calling the keeper with its
// own address as actor. The keeper still applies the limits and protections.
func TestKeeperDirectCalls(t *testing.T) {
	f := setup(t)
	module := addr("othermodule")
	target := addr("frommodule")

	f.ctx = f.freshCtx()
	require.NoError(t, f.k.SanctionAddresses(f.ctx, []string{target}, module, "hook"))
	require.True(t, f.listed(t, target))
	require.Equal(t, []*types.EventSanctionListChanged{{
		Addresses: []string{target}, Sanctioned: true, Actor: module,
		ActorRole: types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, Reason: "hook",
	}}, listChangedEvents(t, f.ctx))

	require.ErrorIs(t, f.k.SanctionAddresses(f.ctx, []string{f.full}, module, "hook"), types.ErrProtectedAddress)
	require.ErrorIs(t, f.k.SanctionAddresses(f.ctx, []string{target}, module, strings.Repeat("x", int(types.DefaultMaxReasonLength)+1)), types.ErrReasonTooLong)
	require.NoError(t, f.k.UnsanctionAddresses(f.ctx, []string{target}, module, "hook"))
	require.False(t, f.listed(t, target))

	// The authority skips the protections, as a gov proposal does.
	require.NoError(t, f.k.SanctionAddresses(f.ctx, []string{f.full}, f.gov, "gov"))
	require.True(t, f.listed(t, f.full))
}

// TestSanctionedSanctionerKeepsRole checks sanctioning a sanctioner doesn't
// change the sanctioner list. Its transactions are stopped by the ante
// handler instead (see app/sanction_integration_test.go).
func TestSanctionedSanctionerKeepsRole(t *testing.T) {
	f := setup(t)
	require.NoError(t, f.sanction(f.gov, f.addOnly))
	has, err := f.k.Sanctioners.Has(f.ctx, f.addOnly)
	require.NoError(t, err)
	require.True(t, has)

	// A full sanctioner can lift it; the address is protected only from being
	// sanctioned.
	require.NoError(t, f.unsanction(f.full, f.addOnly))
	require.False(t, f.listed(t, f.addOnly))
}

func TestQueryPaginationAndErrors(t *testing.T) {
	f := setup(t)
	listed := []string{addr("p1"), addr("p2"), addr("p3")}
	require.NoError(t, f.sanction(f.full, listed...))

	sres, err := f.qs.Sanctioners(f.ctx, &types.QuerySanctionersRequest{Pagination: &query.PageRequest{Limit: 2, CountTotal: true}})
	require.NoError(t, err)
	require.Len(t, sres.Sanctioners, 2)
	require.EqualValues(t, 3, sres.Pagination.Total)
	require.NotEmpty(t, sres.Pagination.NextKey)

	bres, err := f.qs.Blacklist(f.ctx, &types.QueryBlacklistRequest{})
	require.NoError(t, err)
	require.ElementsMatch(t, listed, bres.BlacklistedAccounts)

	_, err = f.qs.Sanctioners(f.ctx, nil)
	require.Error(t, err)
	_, err = f.qs.IsSanctioned(f.ctx, nil)
	require.Error(t, err)
	_, err = f.qs.IsSanctioned(f.ctx, &types.QueryIsSanctionedRequest{})
	require.Error(t, err)
}
