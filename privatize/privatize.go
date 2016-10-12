package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/arithmetric/zcashrpcclient"
	"github.com/arithmetric/zcashrpcclient/zcashjson"
	"github.com/btcsuite/btcutil"
)

const TransactionFee = 0.0001

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Error: Argument required to specify destination address.\n")
		return
	}
	destination := args[0]
	fmt.Printf("privatize.go - Sending unspent transactions to %s\n\n", destination)

	// Connect to Zcash RPC server
	connCfg := &zcashrpcclient.ConnConfig{
		Host:         os.Getenv("ZCASH_RPC_HOST"),
		User:         os.Getenv("ZCASH_RPC_USER"),
		Pass:         os.Getenv("ZCASH_RPC_PASS"),
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := zcashrpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Shutdown()

	// Get all unspent transparent transactions
	unspent, err := client.ListUnspent()
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(unspent) < 1 {
		fmt.Printf("No unspent transactions found.\n")
		return
	}
	for _, utx := range unspent {
		if utx.Confirmations > 0 {
			utxamount, _ := btcutil.NewAmount(utx.Amount)
			feeamount, _ := btcutil.NewAmount(TransactionFee)
			txamount := utxamount - feeamount
			memo := fmt.Sprintf("privatized from %s", utx.Address)
			hexmemo := hex.EncodeToString([]byte(memo))
			amount := zcashjson.ZSendManyEntry{Address: destination, Amount: txamount.ToBTC(), Memo: &hexmemo}
			amounts := make([]zcashjson.ZSendManyEntry, 1)
			amounts[0] = amount
			result, err := client.ZSendMany(utx.Address, amounts)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Sent %f from %s to %s (%s)\n", txamount.ToBTC(), utx.Address, destination, result)
		}
	}
	fmt.Printf("\nFinished sending unspent transactions.\n")
}
