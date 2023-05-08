package anvil

import (
	"sync"
)

var port = []int{}
var pidPortMap = map[int]int{}
var pidPortMapLock = sync.Mutex{}

func init() {
	// Init Port queue
	for i := 33000; i < 33020; i++ {
		port = append(port, i)
	}
}

func GetPort() int {
	if len(port) == 0 {
		return 0
	}
	p := port[0]
	port = port[1:]
	return p
}

func ReturnPort(p int) {
	port = append(port, p)
}

func GetPidPort(pid int) int {
	pidPortMapLock.Lock()
	defer pidPortMapLock.Unlock()
	return pidPortMap[pid]
}

func SetPidPort(pid int, p int) {
	pidPortMapLock.Lock()
	defer pidPortMapLock.Unlock()
	pidPortMap[pid] = p
}

func DeletePidPort(pid int) {
	pidPortMapLock.Lock()
	defer pidPortMapLock.Unlock()
	delete(pidPortMap, pid)
}
