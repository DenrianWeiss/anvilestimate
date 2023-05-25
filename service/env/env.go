package env

import (
	"os"
	"strconv"
)

func GetUpstreamRpc() string {
	env, b := os.LookupEnv("UPSTREAM_RPC")
	if !b {
		return "http://127.0.0.1:8545"
	}
	return env
}

func GetAnvilPath() string {
	env, b := os.LookupEnv("ANVIL_PATH")
	if !b {
		return "anvil"
	}
	return env
}

func GetDelay() uint64 {
	env, b := os.LookupEnv("DELAY_BLOCK")
	if !b {
		return 0
	}
	v, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		return 30
	}
	return uint64(v)
}
