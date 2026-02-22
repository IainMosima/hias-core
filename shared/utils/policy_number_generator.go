package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	policyCounter   int64
	policyCounterMu sync.Mutex
)

func GeneratePolicyNumber() string {
	policyCounterMu.Lock()
	defer policyCounterMu.Unlock()
	policyCounter++
	return fmt.Sprintf("POL-%d-%06d", time.Now().Year(), policyCounter)
}
