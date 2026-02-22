package utils

import (
	"time"

	"github.com/google/uuid"
)

func UUIDToPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func PtrToUUID(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func PtrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func PtrToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func BoolPtr(b bool) *bool {
	return &b
}
