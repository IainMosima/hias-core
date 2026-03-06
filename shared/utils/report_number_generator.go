package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	reportCounter   int64
	reportCounterMu sync.Mutex
)

func GenerateReportNumber() string {
	reportCounterMu.Lock()
	defer reportCounterMu.Unlock()
	reportCounter++
	return fmt.Sprintf("RPT-%d-%06d", time.Now().Year(), reportCounter)
}
