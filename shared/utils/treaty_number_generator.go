package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	treatyCounter   int64
	treatyCounterMu sync.Mutex
)

func GenerateTreatyNumber() string {
	treatyCounterMu.Lock()
	defer treatyCounterMu.Unlock()
	treatyCounter++
	return fmt.Sprintf("TRY-%d-%06d", time.Now().Year(), treatyCounter)
}
