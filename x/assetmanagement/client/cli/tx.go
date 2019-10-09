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
	_from string
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txRootCmd := &cobra.Command{
		Use:                        "token",
		Short:                      "Asset Management transaction sub-commands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txRootCmd.PersistentFlags().StringVarP(&_from, "from", "f", "", "the account")
	err := txRootCmd.MarkPersistentFlagRequired("from")
	if err != nil {
		panic(fmt.Sprintf("failed to setup 'from' flag: %s", err))
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

func getAccountAddress(bech32 string) sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(bech32)
	if err != nil {
		panic(fmt.Sprintf("failed to get account address from '%s': %v", bech32, err))
	}
	return address
}

// GetCmdIssueToken is the CLI command for sending a IssueToken transaction
func GetCmdIssueToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: `issue --token-name [name] --total-supply [amount]
			--symbol [ABC] --mintable --from [account]`,
		Short: "create a new asset",
		// Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// find given account
			from := fetchStringFlag(cmd, "from")
			address := getAccountAddress(from)
			fmt.Printf("token account: %s / %v", from, address)

			name := fetchStringFlag(cmd, "token-name")
			symbol := fetchStringFlag(cmd, "symbol")
			totalSupply := fetchInt64Flag(cmd, "total-supply")
			mintable := fetchBoolFlag(cmd, "mintable")
			fmt.Printf("token is mintable? %t\n", mintable)

			msg := types.NewMsgIssueToken(address, name, symbol, totalSupply, mintable)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupBoolFlag(cmd, "mintable", "", false, "is the new token mintable", false)
	setupStringFlag(cmd, "token-name", "", "", "the name of the new token", true)
	setupInt64Flag(cmd, "total-supply", "", -1,
		"what is the total supply for the new token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC, for the new token", true)

	return cmd
}

// GetCmdMintCoins is the CLI command for sending a MintCoins transaction
func GetCmdMintCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `mint --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "mint more coins for the specified token",
		// Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// if err := cliCtx.EnsureAccountExists(); err != nil {
			// 	return err
			// }

			address, symbol, amount := getCommonParameters(cmd)

			msg := types.NewMsgMintCoins(amount, symbol, address)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to mint for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

func getCommonParameters(cmd *cobra.Command) (sdk.AccAddress, string, int64) {
	// find given account
	from := fetchStringFlag(cmd, "from")
	address := getAccountAddress(from)
	fmt.Printf("token account: %s / %v", from, address)

	symbol := fetchStringFlag(cmd, "symbol")
	amount := fetchInt64Flag(cmd, "amount")
	return address, symbol, amount
}

// GetCmdBurnCoins is the CLI command for sending a BurnCoins transaction
func GetCmdBurnCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `burn --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "destroy the given amount of token/coins, reducing the total supply",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cmd)

			msg := types.NewMsgBurnCoins(amount, symbol, address)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to burn for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

// GetCmdFreezeCoins is the CLI command for sending a FreezeCoins transaction
func GetCmdFreezeCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `freeze --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "move specified amount of token/coins into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cmd)

			msg := types.NewMsgFreezeCoins(amount, symbol, address)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to freeze for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

// GetCmdUnfreezeCoins is the CLI command for sending a FreezeCoins transaction
func GetCmdUnfreezeCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `unfreeze --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "move specified amount of token into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cmd)

			msg := types.NewMsgUnfreezeCoins(amount, symbol, address)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to unfreeze for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}
