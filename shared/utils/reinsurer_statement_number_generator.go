package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	reinsurerStatementCounter   int64
	reinsurerStatementCounterMu sync.Mutex
)

func GenerateReinsurerStatementNumber() string {
	reinsurerStatementCounterMu.Lock()
	defer reinsurerStatementCounterMu.Unlock()
	reinsurerStatementCounter++
	return fmt.Sprintf("RST-%d-%06d", time.Now().Year(), reinsurerStatementCounter)
}
