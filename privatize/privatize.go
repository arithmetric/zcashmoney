package main

import(
  "arithmetric.com/libzcash"
  "bytes"
  "encoding/json"
  "fmt"
  "log"
  "os"
)

const AmountFactor = 1000000000;
const TxnFee = 0.0001;

/**
 *
 * 1. Send unspent to address
 *   - check `zcash-cli listunspent`
 *   - for each, run `zcash-cli z_sendmany` from unspent address to address argument
 *
 **/
func main() {
  args := os.Args[1:]
  destination := args[0]
  //TODO verify address as private
  result := libzcash.Runner("listunspent")
  items := result.([]interface{})
  for i := 0; i < len(items); i++ {
    item := items[i].(map[string]interface{})
    txn_amount := float64(int(item["amount"].(float64) * AmountFactor) - int(TxnFee * AmountFactor)) / AmountFactor
    txn := map[string]interface{}{"address": destination, "amount": txn_amount}
    txns := []interface{}{txn}
    txn_json, err := json.Marshal(txns)
    if err != nil {
  		log.Fatal(err)
  	}
    result := libzcash.Runner("z_sendmany", item["address"].(string), string(txn_json))
    txn_id := result.(bytes.Buffer)
    fmt.Printf("OK Transaction %s -- From %s (%fZTC) -- To %s (%fZTC)\n", txn_id.String(), item["address"], item["amount"], txn["address"], txn["amount"])
	}
}

//  libzcash.Runner("z_gettotalbalance", nil)
//  libzcash.Runner("z_listaddresses", nil)
//  libzcash.Runner("getinfo", nil)
