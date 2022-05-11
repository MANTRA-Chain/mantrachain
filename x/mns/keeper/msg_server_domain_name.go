package keeper

import (
	"context"
	"fmt"

	"github.com/LimeChain/mantrachain/x/mns/types"
	utils "github.com/LimeChain/mantrachain/x/mns/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDomainName(goCtx context.Context, msg *types.MsgCreateDomainName) (*types.MsgCreateDomainNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert owner address string to sdk.AccAddress
	owner, _ := sdk.AccAddressFromBech32(msg.Creator)
	// Check if the value already exists
	domain, isDomainFound := k.GetDomain(ctx, msg.Domain)
	// If a domain is found in store
	if !isDomainFound {
		// Throw an error
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Domain not set")
	}
	_, isDomainNameFound := k.GetDomainName(
		ctx,
		msg.Domain,
		msg.DomainName,
	)
	if isDomainNameFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Domain name already set")
	}

	id := utils.GetDomainNameIndex(msg.Domain, msg.DomainName)
	didCtrl := NewDidController(ctx, id).
		GenDidId().
		GenDidVerMethod(msg.PubKeyHex, owner.String(), msg.VmType).
		GenDidAuth()

	err := didCtrl.SetDid()
	if err != nil {
		return nil, err
	}
	didDoc, err := didCtrl.GetDidDoc()
	if err != nil {
		return nil, err
	}
	strDoc := string(didDoc)
	fmt.Println(strDoc)

	var domainName = types.DomainName{
		Creator:    msg.Creator,
		Index:      id,
		Domain:     domain.Index,
		DomainName: msg.DomainName,
		Did:        "",
		Owner:      owner.String(),
	}

	k.SetDomainName(
		ctx,
		domainName,
	)
	return &types.MsgCreateDomainNameResponse{}, nil
}
