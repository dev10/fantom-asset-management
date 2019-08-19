package assetmanagement

import (
	"fmt"

	"github.com/dev10/fantom-asset-management/x/assetmanagement/rand"
	"github.com/dev10/fantom-asset-management/x/assetmanagement/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "assetmanagement" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgIssueToken:
			return handleMsgIssueToken(ctx, keeper, msg)
		case types.MsgMintCoins:
			return handleMsgMintCoins(ctx, keeper, msg)
		case types.MsgBurnCoins:
			return handleMsgBurnCoins(ctx, keeper, msg)
		case types.MsgFreezeCoins:
			return handleMsgFreezeCoins(ctx, keeper, msg)
		case types.MsgUnfreezeCoins:
			return handleMsgUnfreezeCoins(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized assetmanagement Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handle message to issue token
func handleMsgIssueToken(ctx sdk.Context, keeper Keeper, msg types.MsgIssueToken) sdk.Result {
	var newRandomSymbol = rand.GenerateNewSymbol(msg.Symbol)
	token := types.NewToken(msg.Name, newRandomSymbol, msg.Symbol, msg.TotalSupply, msg.SourceAddress, msg.Mintable)

	keeperErr := keeper.coinKeeper.SetCoins(ctx, msg.SourceAddress, token.TotalSupply)
	if keeperErr != nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to store new token in bank: %s", keeperErr)).Result()
	}

	err := keeper.SetToken(ctx, newRandomSymbol, token)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to store new token: '%s'", err)).Result()
	}
	return sdk.Result{} // Todo: return new symbol name?
}

// handle message to mint coins
func handleMsgMintCoins(ctx sdk.Context, keeper Keeper, msg types.MsgMintCoins) sdk.Result {
	owner, err := keeper.GetOwner(ctx, msg.Symbol)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", msg.Symbol)).Result()
	}
	if !msg.Owner.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins, err := keeper.coinKeeper.AddCoins(ctx, owner,
		sdk.NewCoins(sdk.NewInt64Coin(msg.Symbol, msg.Amount)))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to mint coins: '%s'", err)).Result()
	}

	err = keeper.SetTotalSupply(ctx, msg.Symbol, coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to set total supply when minting coins: '%s'", err)).Result()
	}
	return sdk.Result{}
}

// handle message to burn coins
func handleMsgBurnCoins(ctx sdk.Context, keeper Keeper, msg types.MsgBurnCoins) sdk.Result {
	owner, err := keeper.GetOwner(ctx, msg.Symbol)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", msg.Symbol)).Result()
	}
	if !msg.Owner.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins, err := keeper.coinKeeper.SubtractCoins(ctx, owner,
		sdk.NewCoins(sdk.NewInt64Coin(msg.Symbol, msg.Amount)))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to burn coins: '%s'", err)).Result()
	}

	err = keeper.SetTotalSupply(ctx, msg.Symbol, coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to set total supply when burning coins: '%s'", err)).Result()
	}
	return sdk.Result{}
}

// handle message to freeze coins for specific wallet
func handleMsgFreezeCoins(ctx sdk.Context, keeper Keeper, msg types.MsgFreezeCoins) sdk.Result {
	owner := msg.Owner

	// Todo: Validate you are allowed access to account?
	var customAccount types.CustomAccount = keeper.accountKeeper.GetAccount(ctx, owner)
	err := customAccount.FreezeCoins(sdk.Coins{sdk.NewInt64Coin(msg.Symbol, msg.Amount)})
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to freeze coins: '%s'", err)).Result()
	}

	// Save changes to account
	keeper.accountKeeper.SetAccount(ctx, customAccount)
	return sdk.Result{}
}

// handle message to freeze coins for specific wallet
func handleMsgUnfreezeCoins(ctx sdk.Context, keeper Keeper, msg types.MsgUnfreezeCoins) sdk.Result {
	owner := msg.Owner

	// Todo: Validate you are allowed access to account?
	var customAccount types.CustomAccount = keeper.accountKeeper.GetAccount(ctx, owner)
	err := customAccount.UnfreezeCoins(sdk.Coins{sdk.NewInt64Coin(msg.Symbol, msg.Amount)})
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to unfreeze coins: '%s'", err)).Result()
	}

	// Save changes to account
	keeper.accountKeeper.SetAccount(ctx, customAccount)
	return sdk.Result{}
}
