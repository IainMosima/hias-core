package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	memberCounter   int64
	memberCounterMu sync.Mutex
)

func GenerateMemberNumber() string {
	memberCounterMu.Lock()
	defer memberCounterMu.Unlock()
	memberCounter++
	return fmt.Sprintf("MBR-%d-%06d", time.Now().Year(), memberCounter)
}
