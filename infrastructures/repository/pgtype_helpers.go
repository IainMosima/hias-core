package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// uuidToPgtype converts a google/uuid.UUID to pgtype.UUID.
// A zero-value UUID is treated as invalid (NULL).
func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}

// uuidPtrToPgtype converts a *uuid.UUID to pgtype.UUID.
// A nil pointer is treated as invalid (NULL).
func uuidPtrToPgtype(id *uuid.UUID) pgtype.UUID {
	if id == nil || *id == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// pgtypeToUUID converts a pgtype.UUID to a google/uuid.UUID.
// An invalid pgtype.UUID returns uuid.Nil.
func pgtypeToUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return uuid.UUID(id.Bytes)
}

// stringToPgtypeText converts a Go string to pgtype.Text.
// An empty string is treated as invalid (NULL).
func stringToPgtypeText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

// timeToPgtypeTimestamptz converts a time.Time to pgtype.Timestamptz.
// A zero time is treated as invalid (NULL).
func timeToPgtypeTimestamptz(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// timePtrToPgtypeTimestamptz converts a *time.Time to pgtype.Timestamptz.
// A nil pointer is treated as invalid (NULL).
func timePtrToPgtypeTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// pgtypeTimestamptzToTimePtr converts a pgtype.Timestamptz to *time.Time.
// An invalid timestamptz returns nil.
func pgtypeTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	result := t.Time
	return &result
}

// pgtypeTimestamptzToTime converts a pgtype.Timestamptz to time.Time.
// An invalid timestamptz returns zero time.
func pgtypeTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// timeToPgtypeDate converts a time.Time to pgtype.Date.
// A zero time is treated as invalid (NULL).
func timeToPgtypeDate(t time.Time) pgtype.Date {
	if t.IsZero() {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: t, Valid: true}
}

// pgtypeDateToTime converts a pgtype.Date to time.Time.
// An invalid date returns zero time.
func pgtypeDateToTime(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}

// int64ToPgtypeInt8 converts an int64 to pgtype.Int8.
func int64ToPgtypeInt8(v int64) pgtype.Int8 {
	return pgtype.Int8{Int64: v, Valid: true}
}

// intToPgtypeInt4 converts an int to pgtype.Int4.
func intToPgtypeInt4(v int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(v), Valid: true}
}

// pgtypeDateToTimePtr converts a pgtype.Date to *time.Time.
// An invalid date returns nil.
func pgtypeDateToTimePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	result := d.Time
	return &result
}

// timePtrToPgtypeDate converts a *time.Time to pgtype.Date.
// A nil pointer is treated as invalid (NULL).
func timePtrToPgtypeDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

// pgtypeInt8ToInt64Ptr converts a pgtype.Int8 to *int64.
// An invalid Int8 returns nil.
func pgtypeInt8ToInt64Ptr(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	result := v.Int64
	return &result
}

// int64PtrToPgtypeInt8 converts *int64 to pgtype.Int8.
// A nil pointer is treated as invalid (NULL).
func int64PtrToPgtypeInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

// boolToPgtypeBool converts a bool to pgtype.Bool.
func boolToPgtypeBool(v bool) pgtype.Bool {
	return pgtype.Bool{Bool: v, Valid: true}
}

// float64ToPgNumeric converts a float64 to pgtype.Numeric.
func float64ToPgNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(fmt.Sprintf("%f", f))
	return n
}

// pgNumericToFloat64 converts a pgtype.Numeric to float64.
func pgNumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
