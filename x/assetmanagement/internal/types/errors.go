package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeTokenSymbolDoesNotExist sdk.CodeType = 101
)

func ErrTokenSymbolDoesNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenSymbolDoesNotExist, "Token symbol does not exist")
}
