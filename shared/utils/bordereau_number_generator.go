package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	bordereauCounter   int64
	bordereauCounterMu sync.Mutex
)

func GenerateBordereauNumber() string {
	bordereauCounterMu.Lock()
	defer bordereauCounterMu.Unlock()
	bordereauCounter++
	return fmt.Sprintf("BDX-%d-%06d", time.Now().Year(), bordereauCounter)
}
