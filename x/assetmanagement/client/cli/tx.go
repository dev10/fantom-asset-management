package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
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
			sourceAddress, err := cmd.Flags().GetString("from")
			if err != nil {
				panic(fmt.Sprintf("unable to find 'from' flag: %v", err))
			}
			fmt.Printf("token account: %s", sourceAddress)

			// mintable
			mintable, err := cmd.LocalFlags().GetBool("mintable")
			if err != nil {
				panic(fmt.Sprintf("unable to find 'mintable' flag: %v", err))
			}
			fmt.Printf("token is mintable? %t\n", mintable)

			msg := types.NewMsgIssueToken(cliCtx.GetFromAddress(), name, symbol, totalSupply, mintable)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().BoolP("mintable", "m", false, "is the new token mintable")
	cmd.Flags().StringP("token-name", "tn", "", "the name of the new token")
	cmd.Flags().Int64P("total-supply", "ts", -1, "what is the total supply for the new token")
	cmd.Flags().StringP("symbol", "s", "",
		"what is the shorthand symbol, eg ABC / ABC-123, for the new/existing token")

	err := cmd.MarkFlagRequired("token-name")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'token-name' flag: %s", err))
	}

	err = cmd.MarkFlagRequired("total-supply")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'total-supply' flag: %s", err))
	}

	err = cmd.MarkFlagRequired("symbol")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'symbol' flag: %s", err))
	}

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
