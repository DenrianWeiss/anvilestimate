package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
)

type impersonate struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      string   `json:"id"`
}

func Impersonate(port int, account string) {
	req := impersonate{
		Jsonrpc: "2.0",
		Method:  "anvil_impersonateAccount",
		Params:  []string{account},
		ID:      "i",
	}
	reqB, _ := json.Marshal(&req)
	reqI := bytes.NewReader(reqB)
	_, err := http.Post(fmt.Sprintf("http://%s:%d", "127.0.0.1", port), "application/json", reqI)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func SetBalanceOf(port int, account string) {
	req := impersonate{
		Jsonrpc: "2.0",
		Method:  "anvil_setBalance",
		Params:  []string{account, "0x21e19e0c9bab2400000"},
		ID:      "i",
	}
	reqB, _ := json.Marshal(&req)
	reqI := bytes.NewReader(reqB)
	_, err := http.Post(fmt.Sprintf("http://%s:%d", "127.0.0.1", port), "application/json", reqI)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func GetGasPrice(port int) string {
	conn, _ := ethclient.Dial(fmt.Sprintf("http://%s:%d", "127.0.0.1", port))
	price, err := conn.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return ""
	} else {
		r := hexutil.EncodeBig(price)
		log.Println("gas price: ", r)
		return r
	}
}
