package rpc

import (
	"context"
	"github.com/DenrianWeiss/anvilEstimate/chain/contracts/erc20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
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
		_, err := conn.TransactionReceipt(context.Background(), common.HexToHash(hash))
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			return
		}
	}
}
