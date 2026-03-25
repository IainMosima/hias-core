package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

var (
	policyCounter   int64
	policyCounterMu sync.Mutex
)

func init() {
	// Seed counter from random source to avoid collisions across restarts
	var b [4]byte
	if _, err := rand.Read(b[:]); err == nil {
		policyCounter = int64(binary.BigEndian.Uint32(b[:]) % 900000)
	}
}

func GeneratePolicyNumber() string {
	policyCounterMu.Lock()
	defer policyCounterMu.Unlock()
	policyCounter++
	return fmt.Sprintf("POL-%d-%06d", time.Now().Year(), policyCounter%1000000)
}
