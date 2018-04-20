package commands

import (
	"fmt"
	"math/big"

	agtypes "github.com/Baptist-Publication/angine/types"
	"github.com/Baptist-Publication/chorus/src/eth/common"
	"github.com/Baptist-Publication/chorus/src/eth/crypto"
	"github.com/Baptist-Publication/chorus/src/tools"
	"github.com/Baptist-Publication/chorus/src/types"
	"gopkg.in/urfave/cli.v1"

	cl "github.com/Baptist-Publication/chorus-module/lib/go-rpc/client"
	"github.com/Baptist-Publication/chorus/src/client/commons"
)

var (
	TxCommands = cli.Command{
		Name:     "tx",
		Usage:    "operations for transaction",
		Category: "Transaction",
		Subcommands: []cli.Command{
			{
				Name:   "send",
				Usage:  "send a transaction",
				Action: sendTx,
				Flags: []cli.Flag{
					anntoolFlags.payload,
					anntoolFlags.privkey,
					anntoolFlags.nonce,
					anntoolFlags.to,
					anntoolFlags.value,
				},
			},
		},
	}
)

func sendTx(ctx *cli.Context) error {
	skbs := ctx.String("privkey")
	privkey, err := crypto.HexToECDSA(skbs)
	if err != nil {
		panic(err)
	}

	var chainID string
	if !ctx.GlobalIsSet("target") {
		chainID = "chorus"
	} else {
		chainID = ctx.GlobalString("target")
	}

	nonce := ctx.Uint64("nonce")
	to := common.HexToAddress(ctx.String("to"))
	value := ctx.Int64("value")
	payload := ctx.String("payload")
	data := common.Hex2Bytes(payload)

	bodyTx := types.TxEvmCommon{
		To:     to[:],
		Amount: big.NewInt(value),
	}
	bodyBs, err := tools.TxToBytes(bodyTx)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	from := crypto.PubkeyToAddress(privkey.PublicKey)
	tx := types.NewBlockTx(big.NewInt(90000), big.NewInt(2), nonce, from[:], data)

	if err := tx.Sign(privkey); err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	b, err := tools.TxToBytes(tx)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	tmResult := new(agtypes.RPCResult)
	clientJSON := cl.NewClientJSONRPC(logger, commons.QueryServer)
	_, err = clientJSON.Call("broadcast_tx_commit", []interface{}{chainID, agtypes.WrapTx(types.TxTagAppEvmCommon, b)}, tmResult)
	if err != nil {
		panic(err)
	}
	//res := (*tmResult).(*types.ResultBroadcastTxCommit)

	fmt.Println("tx result:", sigTx.Hash().Hex())

	return nil
}
