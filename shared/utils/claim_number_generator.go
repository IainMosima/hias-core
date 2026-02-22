package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	claimCounter     int64
	claimCounterMu   sync.Mutex
)

func GenerateClaimNumber() string {
	claimCounterMu.Lock()
	defer claimCounterMu.Unlock()
	claimCounter++
	return fmt.Sprintf("CLM-%d-%06d", time.Now().Year(), claimCounter)
}
