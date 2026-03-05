package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	statementCounter   int64
	statementCounterMu sync.Mutex
)

func GenerateStatementNumber() string {
	statementCounterMu.Lock()
	defer statementCounterMu.Unlock()
	statementCounter++
	return fmt.Sprintf("STMT-%d-%06d", time.Now().Year(), statementCounter)
}
