package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	invoiceCounter   int64
	invoiceCounterMu sync.Mutex
)

func GenerateInvoiceNumber() string {
	invoiceCounterMu.Lock()
	defer invoiceCounterMu.Unlock()
	invoiceCounter++
	return fmt.Sprintf("INV-%d-%06d", time.Now().Year(), invoiceCounter)
}
