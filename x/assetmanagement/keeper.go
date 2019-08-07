package assetmanagement

import (
	"errors"
	"fmt"

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

// GetToken gets the entire Token metadata struct by symbol. False if not found, true otherwise
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (error, *types.Token) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsSymbolPresent(ctx, symbol) {
		return fmt.Errorf("could not find Token for symbol '%s'", symbol), nil
	}
	bz := store.Get([]byte(symbol))
	var token types.Token
	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return nil, &token
}

// SetToken sets the entire Token metadata struct by symbol. Owner must be set. Returns success
func (k Keeper) SetToken(ctx sdk.Context, symbol string, token *types.Token) error {
	if token == nil {
		return errors.New("unable to store nil/empty token")
	}
	if token.Owner.Empty() {
		return fmt.Errorf("unable to store token because owner for symbol '%s' is empty", symbol)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(symbol), k.cdc.MustMarshalBinaryBare(*token))
	return nil
}

// Deletes the entire Token metadata struct by symbol
func (k Keeper) DeleteToken(ctx sdk.Context, symbol string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(symbol))
}

// ResolveName - returns the name string that the symbol resolves to
func (k Keeper) ResolveName(ctx sdk.Context, symbol string) (error, string) {
	if err, found := k.GetToken(ctx, symbol); err == nil {
		return nil, found.Name
	} else {
		return fmt.Errorf("couldn't resolve name for symbol '%s' because: %s", symbol, err), ""
	}
}

// SetName - sets the name string that a symbol resolves to
func (k Keeper) SetName(ctx sdk.Context, symbol string, name string) error {
	err, token := k.GetToken(ctx, symbol)
	if err == nil {
		token.Name = name
		return k.SetToken(ctx, symbol, token)
	} else {
		return fmt.Errorf("failed to set token name for symbol '%s' because: %s", symbol, err)
	}
}

// HasOwner - returns whether or not the symbol already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, symbol string) (error, bool) {
	if err, token := k.GetToken(ctx, symbol); err == nil {
		return nil, !token.Owner.Empty()
	} else {
		return fmt.Errorf("unable to check owner for symbol '%s' because: %s", symbol, err), false
	}
}

// GetOwner - get the current owner of a symbol
func (k Keeper) GetOwner(ctx sdk.Context, symbol string) (error, sdk.AccAddress) {
	if err, token := k.GetToken(ctx, symbol); err == nil {
		return nil, token.Owner
	} else {
		return fmt.Errorf("unable to get owner for symbol '%s' because: %s", symbol, err), nil
	}
}

// SetOwner - sets the current owner of a symbol
func (k Keeper) SetOwner(ctx sdk.Context, symbol string, owner sdk.AccAddress) error {
	if err, token := k.GetToken(ctx, symbol); err == nil {
		token.Owner = owner
		return k.SetToken(ctx, symbol, token)
	} else {
		return fmt.Errorf("unable to set owner for symbol '%s' because: %s", symbol, err)
	}
}

// GetTotalSupply - gets the current total supply of a symbol
func (k Keeper) GetTotalSupply(ctx sdk.Context, symbol string) (error, sdk.Coins) {
	if err, token := k.GetToken(ctx, symbol); err == nil {
		return nil, token.TotalSupply
	} else {
		return fmt.Errorf("failed to get total supply for symbol '%s' because: %s", symbol, err), nil
	}
}

// SetTotalSupply - sets the current total supply of a symbol
func (k Keeper) SetTotalSupply(ctx sdk.Context, symbol string, price sdk.Coins) error {
	if err, token := k.GetToken(ctx, symbol); err == nil {
		token.TotalSupply = price
		return k.SetToken(ctx, symbol, token)
	} else {
		return fmt.Errorf("failed to set total supply for symbol '%s' because: %s", symbol, err)
	}
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
