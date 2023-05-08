package anvil

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DenrianWeiss/anvilEstimate/service/env"
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
	port = GetPort()
	if port == 0 {
		return 0, 0, errors.New("no port available")
	}
	cmd := exec.Command(env.GetAnvilPath(), "--fork-url", ethRpc, "-p", strconv.Itoa(port))

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
			if err == os.ErrProcessDone {
				DeleteFork(pid)
				ReturnPort(GetPidPort(pid))
				return nil
			}
			return err
		}
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
		if GetPidPort(pid) != 0 {
			continue
		}
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
		post, err := http.Post("http://localhost:"+strconv.Itoa(GetPidPort(pid)), "application/json", r)
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
		time.Sleep(1 * time.Second)
	}
}
