package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	cessionCounter   int64
	cessionCounterMu sync.Mutex
)

func GenerateCessionNumber() string {
	cessionCounterMu.Lock()
	defer cessionCounterMu.Unlock()
	cessionCounter++
	return fmt.Sprintf("CES-%d-%06d", time.Now().Year(), cessionCounter)
}
