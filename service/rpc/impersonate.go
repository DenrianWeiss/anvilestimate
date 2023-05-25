package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		Method:  "anvil_setBalanceOf",
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
