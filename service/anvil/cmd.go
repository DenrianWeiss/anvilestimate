package anvil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/DenrianWeiss/anvilEstimate/service/env"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type JsonRpcReq struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

type JsonResp struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  string `json:"result"`
}

func StartFork() (pid int, port int, err error) {
	// Get RPC url
	ethRpc := env.GetUpstreamRpc()
	// Get Upstream Block Number
	dial, err := ethclient.Dial(ethRpc)
	if err != nil {
		return 0, 0, err
	}
	number, err := dial.BlockNumber(context.Background())
	if err != nil {
		return 0, 0, err
	}
	port = GetPort()
	if port == 0 {
		return 0, 0, errors.New("no port available")
	}
	cmd := exec.Command(env.GetAnvilPath(), "--fork-url", ethRpc, "-p", strconv.Itoa(port), "--fork-block-number", strconv.FormatUint(number-env.GetDelay(), 10))

	err = cmd.Start()
	if err != nil {
		return 0, 0, err
	}
	SetFork(cmd.Process.Pid, cmd)
	SetPidPort(cmd.Process.Pid, port)

	WaitUntilStarted(cmd.Process.Pid, 10*time.Second)

	go TimedDeletion(cmd.Process.Pid, 2*time.Minute)

	return cmd.Process.Pid, port, nil

}

func StopFork(pid int) error {
	cmd := GetFork(pid)
	if cmd == nil {
		return errors.New("fork not found")
	} else {
		err := cmd.Process.Kill()
		if err != nil {
			if errors.Is(err, os.ErrProcessDone) {
				log.Println("process already done")
				DeleteFork(pid)
				ReturnPort(GetPidPort(pid))
				return nil
			}
			return err
		}
		_ = cmd.Wait()
		DeleteFork(pid)
		ReturnPort(GetPidPort(pid))
		return nil
	}
}

func TimedDeletion(pid int, t time.Duration) {
	time.Sleep(t)
	StopFork(pid)
}

func WaitUntilStarted(pid int, t time.Duration) {
	for i := 0; i < int(t.Seconds()); i++ {
		if GetPidPort(pid) == 0 {
			continue
		}

		time.Sleep(1 * time.Second)
		// Try to get block number
		// If block number is 0, continue
		req := JsonRpcReq{
			Jsonrpc: "2.0",
			Method:  "eth_blockNumber",
			Params:  []interface{}{},
			Id:      "dontcare",
		}
		reqS, _ := json.Marshal(&req)
		r := bytes.NewReader(reqS)
		post, err := http.Post("http://127.0.0.1:"+strconv.Itoa(GetPidPort(pid)), "application/json", r)
		if err != nil {
			continue
		}
		resp := JsonResp{}
		err = json.NewDecoder(post.Body).Decode(&resp)
		if err != nil {
			continue
		} else {
			return
		}
	}
}
