package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/dev10/fantom-asset-management/x/assetmanagement/internal/types"
)

var (
	_chainId   string
	_from      string
	_node      string
	_trustNode bool
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txRootCmd := &cobra.Command{
		Use:                        "token",
		Short:                      "Asset Management transaction sub-commands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txRootCmd.PersistentFlags().StringVarP(&_chainId, "chain-id", "ci", "", "the chain id")
	txRootCmd.PersistentFlags().StringVarP(&_from, "from", "f", "", "the account")
	txRootCmd.PersistentFlags().StringVarP(&_node, "node", "n", "", "the node URL to connect to")
	txRootCmd.PersistentFlags().BoolVarP(&_trustNode, "trust-node", "tn", false, "whether to trust the specified node")
	err := txRootCmd.MarkPersistentFlagRequired("chain-id")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'chain-id' flag: %s", err))
	}
	err = txRootCmd.MarkPersistentFlagRequired("from")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'from' flag: %s", err))
	}
	err = txRootCmd.MarkPersistentFlagRequired("node")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'node' flag: %s", err))
	}

	txRootCmd.AddCommand(client.PostCommands(
		GetCmdIssueToken(cdc),
		GetCmdMintCoins(cdc),
		GetCmdBurnCoins(cdc),
		GetCmdFreezeCoins(cdc),
		GetCmdUnfreezeCoins(cdc),
	)...)

	return txRootCmd
}

func flagError2String(name string, err error) string {
	return fmt.Sprintf("unable to find '%s' flag: %v", name, err)
}

func fetchStringFlag(cmd *cobra.Command, flagName string) string {
	flag, err := cmd.Flags().GetString(flagName)
	if err != nil {
		panic(flagError2String(flagName, err))
	}

	return flag
}

func fetchInt64Flag(cmd *cobra.Command, flagName string) int64 {
	flag, err := cmd.Flags().GetInt64(flagName)
	if err != nil {
		panic(flagError2String(flagName, err))
	}

	return flag
}

func fetchBoolFlag(cmd *cobra.Command, flagName string) bool {
	flag, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		panic(flagError2String(flagName, err))
	}

	return flag
}

func setupRequiredFlag(cmd *cobra.Command, name string) {
	err := cmd.MarkFlagRequired(name)
	if err != nil {
		panic(fmt.Sprintf("failed to setup '%s' flag: %s", name, err))
	}
}

func setupBoolFlag(cmd *cobra.Command, name string, shorthand string, value bool, usage string, required bool) {
	cmd.Flags().BoolP(name, shorthand, value, usage)
	if required {
		setupRequiredFlag(cmd, name)
	}
}

func setupStringFlag(cmd *cobra.Command, name string, shorthand string, value string, usage string, required bool) {
	cmd.Flags().StringP(name, shorthand, value, usage)
	if required {
		setupRequiredFlag(cmd, name)
	}
}

func setupInt64Flag(cmd *cobra.Command, name string, shorthand string, value int64, usage string, required bool) {
	cmd.Flags().Int64P(name, shorthand, value, usage)
	if required {
		setupRequiredFlag(cmd, name)
	}
}

// GetCmdIssueToken is the CLI command for sending a IssueToken transaction
func GetCmdIssueToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: `issue --token-name [name] --total-supply [amount]
			--symbol [ABC] --mintable --from [account] --chain-id [name]
			--node [URL] --trust-node`,
		Short: "create a new asset",
		// Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// find given account
			sourceAddress := fetchStringFlag(cmd, "from")
			address, err := sdk.AccAddressFromBech32(sourceAddress)
			fmt.Printf("token account: %s / %v", sourceAddress, address)

			name := fetchStringFlag(cmd, "token-name")
			symbol := fetchStringFlag(cmd, "symbol")
			totalSupply := fetchInt64Flag(cmd, "total-supply")
			mintable := fetchBoolFlag(cmd, "mintable")
			fmt.Printf("token is mintable? %t\n", mintable)

			msg := types.NewMsgIssueToken(address, name, symbol, totalSupply, mintable)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupBoolFlag(cmd, "mintable", "m", false, "is the new token mintable", false)
	setupStringFlag(cmd, "token-name", "tn", "", "the name of the new token", true)
	setupInt64Flag(cmd, "total-supply", "ts", -1, "what is the total supply for the new token", true)
	setupStringFlag(cmd, "symbol", "s", "",
		"what is the shorthand symbol, eg ABC / ABC-123, for the new/existing token", true)

	return cmd
}

// GetCmdMintCoins is the CLI command for sending a MintCoins transaction
func GetCmdMintCoins(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: `mint --amount [amount] --symbol [ABC-123] --from [account]
			--chain-id [name] --node [URL] --trust-node`,
		Short: "mint more coins for the specified token",
		// Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// if err := cliCtx.EnsureAccountExists(); err != nil {
			// 	return err
			// }

			msg := types.NewMsgMintCoins()
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurnCoins is the CLI command for sending a BurnCoins transaction
func GetCmdBurnCoins(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: `burn --amount [amount] --symbol [ABC-123] --from [account]
			--chain-id [name] --node [URL] --trust-node`,
		Short: "destroy the given amount of token/coins, reducing the total supply",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgBurnCoins()
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdFreezeCoins is the CLI command for sending a FreezeCoins transaction
func GetCmdFreezeCoins(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: `freeze --amount [amount] --symbol [ABC-123] --from [account]
			--chain-id [name] --node [URL] --trust-node`,
		Short: "move specified amount of token/coins into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgFreezeCoins()
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUnfreezeCoins is the CLI command for sending a FreezeCoins transaction
func GetCmdUnfreezeCoins(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: `unfreeze --amount [amount] --symbol [ABC-123] --from [account]
			--chain-id [name] --node [URL] --trust-node`,
		Short: "move specified amount of token into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgUnfreezeCoins()
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
