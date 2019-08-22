package assetmanagement

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(NewKeeper(auth.AccountKeeper{}, nil, nil, nil))

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized assetmanagement Msg type"))
}
