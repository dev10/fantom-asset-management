package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "assetmanagement/IssueToken", nil)
	cdc.RegisterConcrete(MsgMintCoins{}, "assetmanagement/MintCoins", nil)
	cdc.RegisterConcrete(MsgBurnCoins{}, "assetmanagement/BurnCoins", nil)
}
