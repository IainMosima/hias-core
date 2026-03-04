package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	leadCounter   int64
	leadCounterMu sync.Mutex
)

func GenerateLeadNumber() string {
	leadCounterMu.Lock()
	defer leadCounterMu.Unlock()
	leadCounter++
	return fmt.Sprintf("LED-%d-%06d", time.Now().Year(), leadCounter)
}
