package types

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// CustomAccount is customised to allow temporary freezing of coins to exclude them from transactions
type CustomAccount struct {
	*auth.BaseAccount
	FrozenCoins sdk.Coins `json:"coins" yaml:"coins"`
}

func NewCustomAccount(address sdk.AccAddress, coins sdk.Coins, frozenCoins sdk.Coins,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *CustomAccount {

	return &CustomAccount{
		BaseAccount: &auth.BaseAccount{
			Address:       address,
			Coins:         coins,
			PubKey:        pubKey,
			AccountNumber: accountNumber,
			Sequence:      sequence,
		},
		FrozenCoins: frozenCoins,
	}
}

// String implements fmt.Stringer
func (acc CustomAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
  Address:       %s
  Pubkey:        %s
  Coins:         %s
  FrozenCoins:   %s
  AccountNumber: %d
  Sequence:      %d`,
		acc.Address, pubkey, acc.Coins, acc.FrozenCoins, acc.AccountNumber, acc.Sequence,
	)
}

// GetFrozenCoins retrieves frozen coins from account
func (acc *CustomAccount) GetFrozenCoins() sdk.Coins {
	return acc.FrozenCoins
}

// SetFrozenCoins sets frozen coins for account
func (acc *CustomAccount) SetFrozenCoins(frozen sdk.Coins) error {
	acc.FrozenCoins = frozen
	return nil
}

func AreAnyCoinsZero(coins *sdk.Coins) bool {
	for _, coin := range *coins {
		if sdk.NewInt(0).Equal(coin.Amount) {
			return true
		}
	}
	return false
}

// FreezeCoins freezes unfrozen coins for account according to input
func (acc *CustomAccount) FreezeCoins(coinsToFreeze sdk.Coins) error {
	// Have enough coins to freeze?
	if coinsToFreeze == nil || coinsToFreeze.Empty() || coinsToFreeze.IsAnyNegative() || AreAnyCoinsZero(&coinsToFreeze) {
		return sdk.ErrInvalidCoins("No coins chosen to freeze")
	}

	currentCoins := acc.GetCoins()
	if currentCoins == nil || currentCoins.IsAllLT(coinsToFreeze) {
		return sdk.ErrInvalidCoins("Not enough coins to freeze")
	}

	// Freeze coins
	if newBalance, isNegative := currentCoins.SafeSub(coinsToFreeze); !isNegative {
		if err := acc.SetCoins(newBalance); err != nil {
			return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set coins: %s", err))
		}
	} else {
		return sdk.ErrInternal("failed to subtract coins for freezing")
	}

	frozen := acc.GetFrozenCoins()
	if frozen == nil {
		frozen = coinsToFreeze
	} else {
		frozen = frozen.Add(coinsToFreeze)
	}

	if err := acc.SetFrozenCoins(frozen); err != nil {
		return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set frozen coins: %s", err))
	}

	return nil
}

// UnfreezeCoins unfreezes frozen coins for account according to input
func (acc *CustomAccount) UnfreezeCoins(coinsToUnfreeze sdk.Coins) error {
	// Have enough coins to unfreeze?
	if coinsToUnfreeze == nil || coinsToUnfreeze.Empty() || coinsToUnfreeze.IsAnyNegative() {
		return sdk.ErrInvalidCoins("No coins chosen to unfreeze")
	}

	currentlyFrozen := acc.GetFrozenCoins()
	if currentlyFrozen == nil || currentlyFrozen.IsAllLT(coinsToUnfreeze) {
		return sdk.ErrInvalidCoins("Not enough coins to unfreeze")
	}

	// Unfreeze coins
	currentCoins := acc.GetCoins()
	if currentCoins == nil {
		currentCoins = coinsToUnfreeze
	} else {
		currentCoins = currentCoins.Add(coinsToUnfreeze)
	}

	if newFrozenBalance, isNegative := currentlyFrozen.SafeSub(coinsToUnfreeze); !isNegative {
		if err := acc.SetFrozenCoins(newFrozenBalance); err != nil {
			return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set frozen coins: %s", err))
		}
	} else {
		return sdk.ErrInternal("failed to subtract coins for unfreezing")
	}

	if err := acc.SetCoins(currentCoins); err != nil {
		return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set coins: %s", err))
	}

	return nil
}
