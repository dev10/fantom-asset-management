package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgIssueToken defines a IssueToken message
type MsgIssueToken struct {
	SourceAddress sdk.AccAddress `json:"source_address"`
	Name          string         `json:"name"`
	Symbol        string         `json:"symbol"`
	TotalSupply   string         `json:"total_supply"`
	Mintable      bool           `json:"mintable"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgIssueToken(sourceAddress sdk.AccAddress, name, symbol, totalSupply string, mintable bool) MsgIssueToken {
	return MsgIssueToken{
		SourceAddress: sourceAddress,
		Name:          name,
		Symbol:        symbol,
		TotalSupply:   totalSupply,
		Mintable:      mintable,
	}
}

// Route should return the name of the module
func (msg MsgIssueToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgIssueToken) Type() string { return "issue_token" }

// ValidateBasic runs stateless checks on the message
func (msg MsgIssueToken) ValidateBasic() sdk.Error {
	if msg.SourceAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.SourceAddress.String())
	}
	if len(msg.Name) == 0 || len(msg.Symbol) == 0 || len(msg.TotalSupply) == 0 {
		return sdk.ErrUnknownRequest("Name, Symbol and/or TotalSupply cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgIssueToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SourceAddress}
}

// MsgMintCoins defines the BuyName message
type MsgMintCoins struct {
	Amount string         `json:"amount"`
	Symbol string         `json:"symbol"`
	Minter sdk.AccAddress `json:"minter"`
}

// NewMsgMintCoins is the constructor function for MsgMintCoins
func NewMsgMintCoins(amount, symbol string, minter sdk.AccAddress) MsgMintCoins {
	return MsgMintCoins{
		Amount: amount,
		Symbol: symbol,
		Minter: minter,
	}
}

// Route should return the name of the module
func (msg MsgMintCoins) Route() string { return RouterKey }

// Type should return the action
func (msg MsgMintCoins) Type() string { return "mint_coins" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMintCoins) ValidateBasic() sdk.Error {
	if msg.Minter.Empty() {
		return sdk.ErrInvalidAddress(msg.Minter.String())
	}
	if len(msg.Symbol) == 0 || len(msg.Amount) == 0 {
		return sdk.ErrUnknownRequest("Symbol and/or Amount cannot be empty")
	}
	// Todo: make sure Amount is positive decimal
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMintCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgMintCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Minter}
}
