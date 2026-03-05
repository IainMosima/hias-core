package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	creditNoteCounter   int64
	creditNoteCounterMu sync.Mutex
)

func GenerateCreditNoteNumber() string {
	creditNoteCounterMu.Lock()
	defer creditNoteCounterMu.Unlock()
	creditNoteCounter++
	return fmt.Sprintf("CN-%d-%06d", time.Now().Year(), creditNoteCounter)
}
