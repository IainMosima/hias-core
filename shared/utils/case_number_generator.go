package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	caseCounter   int64
	caseCounterMu sync.Mutex
)

func GenerateCaseNumber() string {
	caseCounterMu.Lock()
	defer caseCounterMu.Unlock()
	caseCounter++
	return fmt.Sprintf("CASE-%d-%06d", time.Now().Year(), caseCounter)
}
