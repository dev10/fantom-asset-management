package assetmanagement

import "github.com/dev10/fantom-asset-management/x/assetmanagement/types"

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	// messages
	NewMsgBurnCoins     = types.NewMsgBurnCoins
	NewMsgFreezeCoins   = types.NewMsgFreezeCoins
	NewMsgIssueToken    = types.NewMsgIssueToken
	NewMsgMintCoins     = types.NewMsgMintCoins
	NewMsgUnfreezeCoins = types.NewMsgUnfreezeCoins

	ModuleCdc    = types.ModuleCdc
	RegisterCode = types.RegisterCodec
)

type (
	// messages
	MsgBurnCoins     = types.MsgBurnCoins
	MsgFreezeCoins   = types.MsgFreezeCoins
	MsgIssueToken    = types.MsgIssueToken
	MsgMintCoins     = types.MsgMintCoins
	MsgUnfreezeCoins = types.MsgUnfreezeCoins

	// queries
	QueryResultSymbol = types.QueryResultSymbol

	// state/stored types
	CustomCoinAccount = types.CustomAccount
	Token             = types.Token
)
