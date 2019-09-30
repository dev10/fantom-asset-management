package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgIssueToken(t *testing.T) {
	var (
		name         = "Zap"
		symbol       = "ZAP"
		total  int64 = 1
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgIssueToken(owner, name, symbol, total, false)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "issue_token")
}

func TestMsgIssueTokenValidation(t *testing.T) {
	var (
		name               = "Zap"
		symbol             = "ZAP"
		total        int64 = 1
		totalInvalid int64 = 0
		acc                = sdk.AccAddress([]byte("me"))
		name2              = "a"
		total2       int64 = 2
		acc2               = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgIssueToken
	}{
		{true, NewMsgIssueToken(acc, name, symbol, total, false)},
		{true, NewMsgIssueToken(acc, name, symbol, total, false)},
		{false, NewMsgIssueToken(acc, name, symbol, totalInvalid, false)},
		{true, NewMsgIssueToken(acc2, name2, symbol, total2, false)},
		{true, NewMsgIssueToken(acc2, name2, symbol, total, false)},
		{true, NewMsgIssueToken(acc, name2, symbol, total2, false)},
		{false, NewMsgIssueToken(nil, name, symbol, total2, false)},
		{false, NewMsgIssueToken(acc2, "", symbol, total2, false)},
		{false, NewMsgIssueToken(acc2, name, symbol, totalInvalid, false)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, fmt.Sprintf("transaction [%v] failed but was marked valid", tc))
		} else {
			require.NotNil(t, err, fmt.Sprintf("transaction [%v] is valid but has an error", tc))
		}
	}
}

func TestMsgIssueTokenGetSignBytes(t *testing.T) {
	var (
		name         = "Zap"
		symbol       = "ZAP"
		total  int64 = 1
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgIssueToken(owner, name, symbol, total, false)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/IssueToken",` +
		`"value":{` +
		`"mintable":false,` +
		`"name":"Zap",` +
		`"source_address":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP",` +
		`"total_supply":"1"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgMintCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgMintCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "mint_coins")
}

func TestMsgMintCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 10
		symbol        = "ZAP-001"
		symbol2       = "MNT-500"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgMintCoins
	}{
		{true, NewMsgMintCoins(amount, symbol, owner)},
		{true, NewMsgMintCoins(amount, symbol2, owner2)},
		{false, NewMsgMintCoins(-1, symbol, owner)},
		{false, NewMsgMintCoins(0, symbol, owner)},
		{true, NewMsgMintCoins(1, symbol, owner)},
		{false, NewMsgMintCoins(amount, symbol, nil)},
		{false, NewMsgMintCoins(amount, "", owner)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, fmt.Sprintf("transaction [%v] failed but was marked valid", tc))
		} else {
			require.NotNil(t, err, fmt.Sprintf("transaction [%v] is valid but has an error", tc))
		}
	}
}

func TestMsgMintCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgMintCoins(amount, symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/MintCoins","value":{` +
		`"amount":"10",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP-001"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgBurnCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgBurnCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "burn_coins")
}

func TestMsgBurnCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 20
		symbol        = "ZAP-001"
		symbol2       = "BRN-002"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgBurnCoins
	}{
		{true, NewMsgBurnCoins(amount, symbol, owner)},
		{true, NewMsgBurnCoins(amount, symbol2, owner2)},
		{false, NewMsgBurnCoins(-1, symbol, owner)},
		{false, NewMsgBurnCoins(0, symbol, owner)},
		{true, NewMsgBurnCoins(1, symbol, owner)},
		{false, NewMsgBurnCoins(amount, symbol, nil)},
		{false, NewMsgBurnCoins(amount, "", owner)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgBurnCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "ZAP-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgBurnCoins(amount, symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/BurnCoins","value":{` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP-999"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgFreezeCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgFreezeCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "freeze_coins")
}
