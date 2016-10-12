package main

import(
  "github.com/arithmetric/zcashrpcclient"
  "github.com/arithmetric/zcashrpcclient/zcashjson"
  "github.com/btcsuite/btcd/chaincfg/chainhash"
  "github.com/btcsuite/btcutil"
  "bufio"
  "fmt"
  "encoding/hex"
  "log"
  "os"
  "strconv"
  "strings"
  "time"
)

const TransactionFee = 0.0001

func main() {
  fmt.Print("zsend.go - Interactive Zcash sending tool\n\n")

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

  reader := bufio.NewReader(os.Stdin)

  // Prompt for sender's address
  fmt.Print("Sender Address: ")
  addrSender, _ := reader.ReadString('\n')
  addrSender = strings.TrimSpace(addrSender)

  // Ensure sender's address is in wallet
  zaddrs, err := client.ZListAddresses()
  if err != nil {
    log.Fatal(err)
    return
  }
  zaddrfound := false
  for _, z := range zaddrs {
    if z == addrSender {
      zaddrfound = true
      break
    }
  }
  if !zaddrfound {
    fmt.Print("\nError: Cannot send from this address. It is not in the wallet.\n")
    return
  }

  // Prompt for recipient's address
  fmt.Print("Recipient Address: ")
  addrRecipient, _ := reader.ReadString('\n')
  addrRecipient = strings.TrimSpace(addrRecipient)

  // Prompt for amount to send
  fmt.Print("Amount: ")
  amountStr, _ := reader.ReadString('\n')
  amountFloat, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
  amount, _ := btcutil.NewAmount(amountFloat)
  feeamount, _ := btcutil.NewAmount(TransactionFee)

  // Ensure that sender's address has funds to cover amount
  amt, err := client.ZGetBalance(addrSender)
  if err != nil {
    log.Fatal(err)
  }
  if amt < (amount + feeamount) {
    fmt.Printf("Error: Insufficient balance. Sender's address has: %f\n", amt.ToBTC())
    return
  }

  // Prompt for transaction memo if destination is ZAddress
  memo := ""
  if addrRecipient[0] == 'z' {
    fmt.Print("Memo: ")
    memo, _ = reader.ReadString('\n')
    memo = strings.TrimSpace(memo)
  }

  // Confirm transaction before sending
  fmt.Printf("\nPreparing to send the following transaction:\n\n  From: %s\n  To: %s\n  Amount: %f\n  Memo: %s\n\n", addrSender, addrRecipient, amount.ToBTC(), memo)
  fmt.Print("Is this correct? (yes/no) ")
  confirm, _ := reader.ReadString('\n')
  if confirm != "yes\n" {
    fmt.Printf("Cancelling transaction.\n")
    return
  }

  // Send transaction
  fmt.Printf("\nSending transaction...\n")
  envelope := zcashjson.ZSendManyEntry{Address: addrRecipient, Amount: amount.ToBTC()}
  if len(memo) > 0 {
    hexmemo := hex.EncodeToString([]byte(memo))
    envelope.Memo = &hexmemo
  }
  envelopes := make([]zcashjson.ZSendManyEntry, 1)
  envelopes[0] = envelope
  opId, err := client.ZSendMany(addrSender, envelopes)
  if err != nil {
		log.Fatal(err)
    return
	}
	fmt.Printf("Transaction in queue as operation: %s\n", opId)

  // Prompt user for waiting for transaction operation status
  fmt.Print("\nWait for transaction ID? (yes/no) ")
  confirm, _ = reader.ReadString('\n')
  if confirm != "yes\n" {
    fmt.Printf("Exiting with transaction in queue.\n")
    return
  }

  // Wait for transaction operation to succeed or fail
  var txid string
  CheckOperation:
  for {
    statuses, err := client.ZGetOperationStatus()
    if err != nil {
  		log.Fatal(err)
  	}
    opfound := false
    var status zcashjson.ZGetOperationStatusResult
    for _, status = range statuses {
      if status.Id == opId {
        opfound = true
        break
      }
    }
    if !opfound {
      fmt.Printf("Error: Transaction operation not found.\n")
      return
    }
    switch status.Status {
      case "success":
        txid = statuses[0].Result["txid"]
        fmt.Printf("\nTransaction successfully created as: %s\n", txid)
        break CheckOperation
      case "error":
        fallthrough
      case "failed":
        fmt.Printf("Transaction failed: %s\n", statuses[0].Error.Message)
        return
    }
    fmt.Printf("Waiting 15 seconds...\n")
    time.Sleep(15 * time.Second)
  }

  // Prompt user for waiting for transaction confirmation
  fmt.Print("\nWait for transaction confirmation? (yes/no) ")
  confirm, _ = reader.ReadString('\n')
  if confirm != "yes\n" {
    fmt.Printf("Exiting with transaction in queue.\n")
    return
  }

  // Wait for transaction confirmation
  txhash, _ := chainhash.NewHashFromStr(txid)
  CheckTransaction:
  for {
    tx, err := client.GetTransaction(txhash)
    if err != nil {
  		log.Fatal(err)
  	}
    if tx.Confirmations > 0 {
      fmt.Printf("\nTransaction %s has %d confirmations!\n", txid, tx.Confirmations)
      break CheckTransaction
    }
    fmt.Printf("Waiting 30 seconds...\n")
    time.Sleep(30 * time.Second)
  }
}
