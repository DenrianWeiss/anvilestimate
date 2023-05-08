package anvil

import (
	"os/exec"
	"sync"
)

var forkMap = make(map[int]*exec.Cmd)
var rwLock = sync.RWMutex{}

func SetFork(pid int, cmd *exec.Cmd) {
	rwLock.Lock()
	defer rwLock.Unlock()
	forkMap[pid] = cmd
}

func GetFork(pid int) *exec.Cmd {
	rwLock.RLock()
	defer rwLock.RUnlock()
	return forkMap[pid]
}

func DeleteFork(pid int) {
	rwLock.Lock()
	defer rwLock.Unlock()
	delete(forkMap, pid)

	DeletePidPort(pid)
}
