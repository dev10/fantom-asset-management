package assetmanagement

import "github.com/dev10/fantom-asset-management/x/assetmanagement/internal/keeper"
import "github.com/dev10/fantom-asset-management/x/assetmanagement/internal/types"

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier

	// messages
	NewMsgBurnCoins     = types.NewMsgBurnCoins
	NewMsgFreezeCoins   = types.NewMsgFreezeCoins
	NewMsgIssueToken    = types.NewMsgIssueToken
	NewMsgMintCoins     = types.NewMsgMintCoins
	NewMsgUnfreezeCoins = types.NewMsgUnfreezeCoins

	NewToken = types.NewToken

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper

	// messages
	MsgBurnCoins     = types.MsgBurnCoins
	MsgFreezeCoins   = types.MsgFreezeCoins
	MsgIssueToken    = types.MsgIssueToken
	MsgMintCoins     = types.MsgMintCoins
	MsgUnfreezeCoins = types.MsgUnfreezeCoins

	// queries
	QueryResultSymbol = types.QueryResultSymbol

	// state/stored types
	CustomAccount = types.CustomAccount
	Token         = types.Token
)
