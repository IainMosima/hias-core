package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	quotationCounter   int64
	quotationCounterMu sync.Mutex
)

func GenerateQuotationNumber() string {
	quotationCounterMu.Lock()
	defer quotationCounterMu.Unlock()
	quotationCounter++
	return fmt.Sprintf("QUO-%d-%06d", time.Now().Year(), quotationCounter)
}
