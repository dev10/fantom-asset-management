package assetmanagement

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/dev10/fantom-asset-management/x/assetmanagement/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter
// methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the assetmanagement Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// GetToken gets the entire Token metadata struct by symbol
func (k Keeper) GetToken(ctx sdk.Context, symbol string) types.Token {
	store := ctx.KVStore(k.storeKey)
	if !k.IsSymbolPresent(ctx, symbol) {
		return types.Token{Symbol: symbol}
	}
	bz := store.Get([]byte(symbol))
	var token types.Token
	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return token
}

// SetToken sets the entire Token metadata struct by symbol. Owner must be set.
func (k Keeper) SetToken(ctx sdk.Context, symbol string, token types.Token) {
	if token.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(symbol), k.cdc.MustMarshalBinaryBare(token))
}

// Deletes the entire Token metadata struct by symbol
func (k Keeper) DeleteToken(ctx sdk.Context, symbol string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(symbol))
}

// ResolveName - returns the name string that the symbol resolves to
func (k Keeper) ResolveName(ctx sdk.Context, symbol string) string {
	return k.GetToken(ctx, symbol).Name
}

// SetName - sets the name string that a symbol resolves to
func (k Keeper) SetName(ctx sdk.Context, symbol string, name string) {
	token := k.GetToken(ctx, symbol)
	token.Name = name
	k.SetToken(ctx, symbol, token)
}

// HasOwner - returns whether or not the symbol already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, symbol string) bool {
	return !k.GetToken(ctx, symbol).Owner.Empty()
}

// GetOwner - get the current owner of a symbol
func (k Keeper) GetOwner(ctx sdk.Context, symbol string) sdk.AccAddress {
	return k.GetToken(ctx, symbol).Owner
}

// SetOwner - sets the current owner of a symbol
func (k Keeper) SetOwner(ctx sdk.Context, symbol string, owner sdk.AccAddress) {
	token := k.GetToken(ctx, symbol)
	token.Owner = owner
	k.SetToken(ctx, symbol, token)
}

// GetTotalSupply - gets the current total supply of a symbol
func (k Keeper) GetTotalSupply(ctx sdk.Context, symbol string) sdk.Coins {
	return k.GetToken(ctx, symbol).TotalSupply
}

// SetTotalSupply - sets the current total supply of a symbol
func (k Keeper) SetTotalSupply(ctx sdk.Context, symbol string, price sdk.Coins) {
	token := k.GetToken(ctx, symbol)
	token.TotalSupply = price
	k.SetToken(ctx, symbol, token)
}

// GetNamesIterator - Get an iterator over all symbols in which the keys are the symbols and the values are the token
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// IsSymbolPresent - Check if the symbol is present in the store or not
func (k Keeper) IsSymbolPresent(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(symbol))
}
