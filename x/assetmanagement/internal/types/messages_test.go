package types

import (
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
