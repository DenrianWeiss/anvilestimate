package handler

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/DenrianWeiss/anvilEstimate/service/anvil"
	"github.com/DenrianWeiss/anvilEstimate/service/cache"
	"github.com/DenrianWeiss/anvilEstimate/service/rpc"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"math/big"
	"net/http"
	"time"
)

type SimulationRequest struct {
	From        string   `json:"from" binding:"required"`
	To          string   `json:"to" binding:"required"`
	Amount      string   `json:"amount"`
	Data        string   `json:"data"`
	TokenChange []string `json:"token_change" binding:"required"`
}

type sendPayload struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value,omitempty"`
	Data  string `json:"data,omitempty"`
}

type SendTransactionRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []sendPayload `json:"params"`
	Id      string        `json:"id"`
}

type SendTransactionResp struct {
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewSendTxRequest(from, to, value, data string) SendTransactionRequest {
	return SendTransactionRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendTransaction",
		Params: []sendPayload{{
			From:  from,
			To:    to,
			Value: value,
			Data:  data,
		}},
		Id: "1",
	}
}

type SimulationResponse struct {
	TokenChange map[string]string `json:"token_change"`
	Status      string
	Reason      string
}

func HandleSimulationRequest(ctx *gin.Context) {
	var request SimulationRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":  "invalid request",
			"detail": err.Error(),
		})
		return
	}
	// Spin up a simulation environment
	entry := NewEntry()

	cache.SetRespCache(entry, SimulationResponse{
		TokenChange: nil,
		Status:      "pending",
	})

	go AsyncSimulation(entry, request)

	ctx.JSON(200, gin.H{
		"message": "ok",
		"data":    entry,
	})
}

func GetSimulationResult(ctx *gin.Context) {
	entry := ctx.Param("entry")
	resp, _ := cache.ReadRespCache(entry)
	if resp == nil {
		ctx.JSON(200, gin.H{
			"error": "not found",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "ok",
		"data":    resp,
	})
}

func AsyncSimulation(entry string, req SimulationRequest) {
	// Launch simulation environment
	fork, port, err := anvil.StartFork()
	if err != nil {
		return
	}
	defer anvil.StopFork(fork)
	// Wait for it to start
	anvil.WaitUntilStarted(fork, 10*time.Second)
	rpc.Impersonate(port, req.From)
	balanceSaved := make(map[string]*big.Int)
	// Record balance
	for _, t := range req.TokenChange {
		balanceSaved[t] = rpc.GetBalance(t, req.From, port)
	}
	// Send Transaction
	payload := NewSendTxRequest(req.From, req.To, req.Amount, req.Data)
	payloadB, _ := json.Marshal(&payload)
	post, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d", port), "application/json", bytes.NewReader(payloadB))
	if err != nil {
		return // todo handle error
	}
	// Read post result
	resp := SendTransactionResp{}
	err = json.NewDecoder(post.Body).Decode(&resp)
	if err != nil {
		cache.SetRespCache(entry, SimulationResponse{
			TokenChange: nil,
			Status:      "err",
			Reason:      err.Error(),
		})
		return
	}
	txId := resp.Result
	if txId == "" {
		cache.SetRespCache(entry, SimulationResponse{
			TokenChange: nil,
			Status:      "err",
			Reason:      "failed to send tx",
		})
		return
	}
	// Wait for it to be mined
	rpc.WaitMined(port, txId)

	// Record balance
	balanceNew := make(map[string]*big.Int)
	for _, t := range req.TokenChange {
		balanceNew[t] = rpc.GetBalance(t, req.From, port)
	}
	// Record Balance Diff
	balanceDiff := make(map[string]string)
	for s, b := range balanceSaved {
		balanceDiff[s] = new(big.Int).Sub(balanceNew[s], b).String()
	}
	// Save to cache
	cache.SetRespCache(entry, SimulationResponse{
		TokenChange: balanceDiff,
		Status:      "ok",
	})
	go ClearCacheTimer(entry)
}

func NewEntry() string {
	return uuid()
}

func ClearCacheTimer(entry string) {
	time.Sleep(10 * time.Minute)
	cache.DelRespCache(entry)
}

func uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
