package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	claimCounter   int64
	claimCounterMu sync.Mutex
	claimInitOnce  sync.Once
)

func GenerateClaimNumber() string {
	claimCounterMu.Lock()
	defer claimCounterMu.Unlock()
	claimCounter++
	return fmt.Sprintf("CLM-%d-%06d", time.Now().Year(), claimCounter)
}

func InitClaimCounter(start int64) {
	claimInitOnce.Do(func() {
		claimCounterMu.Lock()
		defer claimCounterMu.Unlock()
		if start > claimCounter {
			claimCounter = start
		}
	})
}

func ResetClaimCounterForCollision(start int64) {
	claimCounterMu.Lock()
	defer claimCounterMu.Unlock()
	if start > claimCounter {
		claimCounter = start
	}
}
