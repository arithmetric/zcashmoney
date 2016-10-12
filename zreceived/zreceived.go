package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/arithmetric/zcashrpcclient"
)

func main() {
	fmt.Print("zreceived.go - List balance and received transactions by address\n\n")

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

	// Get all addresses
	zaddrs, err := client.ZListAddresses()
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, zaddr := range zaddrs {
		fmt.Printf("[ZAddress] %s\n", zaddr)

		// Get balance for each address
		amt, err := client.ZGetBalance(zaddr)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Balance: %f\n", amt.ToBTC())

		// Get received transactions for each address
		recvs, err := client.ZListReceivedByAddress(zaddr)
		if err != nil {
			log.Fatal(err)
			return
		}
		if len(recvs) < 1 {
			fmt.Printf("\nNo transactions received to ZAddress.\n")
		}
		for _, recv := range recvs {
			memobytes, err := hex.DecodeString(recv.Memo)
			if err != nil {
				log.Fatal(err)
			}
			idx := bytes.IndexByte(memobytes, 0)
			var memo string
			if idx == 1 && memobytes[0] == 0xf6 {
				memo = "(generated)"
			} else {
				memo = string(memobytes[:idx])
			}
			fmt.Printf("\n  Tx ID: %s\n  Amount: %f\tMemo: %s\n", recv.TxID, recv.Amount, memo)
		}
		fmt.Printf("\n")
	}
}
