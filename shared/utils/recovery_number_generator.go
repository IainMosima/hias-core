package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	recoveryCounter   int64
	recoveryCounterMu sync.Mutex
)

func GenerateRecoveryNumber() string {
	recoveryCounterMu.Lock()
	defer recoveryCounterMu.Unlock()
	recoveryCounter++
	return fmt.Sprintf("REC-%d-%06d", time.Now().Year(), recoveryCounter)
}
