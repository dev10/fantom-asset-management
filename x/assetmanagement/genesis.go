package assetmanagement

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	TokenRecords []Token `json:"token_records"`
}

func NewGenesisState(tokenRecords []Token) GenesisState {
	return GenesisState{TokenRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.TokenRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid TokenRecord: Value: %s. Error: Missing Owner", record.Symbol)
		}
		if record.Symbol == "" {
			return fmt.Errorf("invalid TokenRecord: Owner: %s. Error: Missing Symbol", record.Owner)
		}
		if record.TotalSupply == nil || record.TotalSupply.Len() == 0 {
			return fmt.Errorf("invalid TokenRecord: Symbol: %s. Error: Missing TotalSupply", record.Symbol)
		}
		if record.Name == "" {
			return fmt.Errorf("invalid TokenRecord: Symbol: %s. Error: Missing Name", record.Symbol)
		}
		if record.OriginalSymbol == "" {
			return fmt.Errorf("invalid TokenRecord: Symbol: %s. Error: Missing OriginalSymbol", record.Symbol)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		TokenRecords: []Token{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.TokenRecords {
		err := keeper.SetToken(ctx, record.Symbol, &record)
		if err != nil {
			panic(fmt.Sprintf("failed to set token for symbol: %s. Error: %s", record.Symbol, err))
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Token
	iterator := k.GetTokensIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		symbol := string(iterator.Key())
		token, err := k.GetToken(ctx, symbol)
		if err != nil {
			panic(fmt.Sprintf("failed to find token for symbol: %s. Error: %s", symbol, err))
		}
		records = append(records, *token)

	}
	return GenesisState{TokenRecords: records}
}
