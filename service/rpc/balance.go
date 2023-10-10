package rpc

import (
	"context"
	"fmt"
	"github.com/DenrianWeiss/anvilEstimate/chain/contracts/erc20"
	"github.com/DenrianWeiss/anvilEstimate/service/env"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const ETHPlaceHolder = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
const ZeroStr = "0x0000000000000000000000000000000000000000"

func GetBalance(token, account string, port int) *big.Int {
	conn, _ := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(port))
	defer conn.Close()
	if strings.ToLower(token) == ETHPlaceHolder || strings.ToLower(token) == ZeroStr {
		b, _ := conn.BalanceAt(context.Background(), common.HexToAddress(account), nil)
		return b
	}
	// ERC20 token, get using contract
	tokenAddress := common.HexToAddress(token)
	// 1. Construct contract.
	tokenInstance, _ := erc20.NewErc20(tokenAddress, conn)
	b, _ := tokenInstance.BalanceOf(nil, common.HexToAddress(account))
	return b
}

func WaitMined(port int, hash string) {
	conn, _ := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(port))
	defer conn.Close()
	for i := 0; i < 10; i++ {
		rct, err := conn.TransactionReceipt(context.Background(), common.HexToHash(hash))
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			log.Printf("Receipt for tx %x: ", rct.TxHash)
			if rct.Status == 1 {
				return
			} else {
				log.Println("tx failed", rct.Status)
			}
			// Lets run cast to get tx receipt
			castCmd := exec.Command(env.GetCastPath(), "run", rct.TxHash.Hex(), "--rpc-url", fmt.Sprintf("http://127.0.0.1:%d", port))
			castCmd.Stdout = os.Stdout
			castCmd.Stderr = os.Stderr
			err = castCmd.Run()
			if err != nil {
				log.Println("error during cast run", err)
			}
			return
		}
	}
	log.Println("WaitMined timeout")
}

func GetGasCost(port int, hash string) *big.Int {
	conn, _ := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(port))
	defer conn.Close()
	rct, err := conn.TransactionReceipt(context.Background(), common.HexToHash(hash))
	// Calculate Gas Cost
	if err != nil {
		return big.NewInt(0)
	}
	gasUsed := rct.GasUsed
	gasPriceInTx := rct.EffectiveGasPrice
	gasCost := new(big.Int).Mul(big.NewInt(int64(gasUsed)), gasPriceInTx)
	return gasCost
}
