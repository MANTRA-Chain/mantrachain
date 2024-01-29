package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateNftCollection{}, "token/CreateNftCollection", nil)
	cdc.RegisterConcrete(&MsgMintNfts{}, "token/MintNfts", nil)
	cdc.RegisterConcrete(&MsgBurnNfts{}, "token/BurnNfts", nil)
	cdc.RegisterConcrete(&MsgTransferNfts{}, "token/TransferNfts", nil)
	cdc.RegisterConcrete(&MsgApproveNfts{}, "token/ApproveNfts", nil)
	cdc.RegisterConcrete(&MsgApproveAllNfts{}, "token/ApproveAllNfts", nil)
	cdc.RegisterConcrete(&MsgMintNft{}, "token/MintNft", nil)
	cdc.RegisterConcrete(&MsgBurnNft{}, "token/BurnNft", nil)
	cdc.RegisterConcrete(&MsgTransferNft{}, "token/TransferNft", nil)
	cdc.RegisterConcrete(&MsgApproveNft{}, "token/ApproveNft", nil)
	cdc.RegisterConcrete(&MsgUpdateGuardSoulBondNftImage{}, "token/UpdateGuardSoulBondNftImage", nil)
	cdc.RegisterConcrete(&MsgUpdateRestrictedCollectionNftImage{}, "token/UpdateRestrictedCollectionNftImage", nil)
	cdc.RegisterConcrete(&MsgUpdateRestrictedCollectionNftImageBatch{}, "token/UpdateRestrictedCollectionNftImageBatch", nil)
	cdc.RegisterConcrete(&MsgUpdateRestrictedCollectionNftImageGroupedBatch{}, "token/UpdateRestrictedCollectionNftImageGroupedBatch", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateNftCollection{},
		&MsgMintNfts{},
		&MsgBurnNfts{},
		&MsgTransferNfts{},
		&MsgApproveNfts{},
		&MsgApproveAllNfts{},
		&MsgMintNft{},
		&MsgBurnNft{},
		&MsgTransferNft{},
		&MsgApproveNft{},
		&MsgUpdateGuardSoulBondNftImage{},
		&MsgUpdateRestrictedCollectionNftImage{},
		&MsgUpdateRestrictedCollectionNftImageBatch{},
		&MsgUpdateRestrictedCollectionNftImageGroupedBatch{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
