package env

import "os"

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
