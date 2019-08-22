package assetmanagement

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the assetmanagement Querier
const (
	QuerySymbols = "symbols"
	QueryToken   = "token"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, path[1:], req, keeper)
		case QuerySymbols:
			return querySymbols(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown assetmanagement query endpoint")
		}
	}
}

// nolint: unparam
func queryToken(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	token, err := keeper.GetToken(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, token)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

// nolint: unparam
func querySymbols(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var symbolList QueryResultSymbol

	iterator := keeper.GetTokensIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		symbolList = append(symbolList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, symbolList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
