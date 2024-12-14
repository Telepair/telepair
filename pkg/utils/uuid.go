package utils

import (
	"errors"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
)

const MaxTry = 8 // MaxTry is the maximum number of tries to generate a UUID.

// ErrInvalidUUID is the error returned when the UUID is invalid.
var ErrInvalidUUID = errors.New("invalid UUID")

// UUID is a wrapper around google/uuid.UUID.
type UUID uuid.UUID

// NewUUID returns a new UUID.
func NewUUID() UUID {
	return UUIDv7()
}

// Time returns the time part of the UUID.
func (u UUID) Time() time.Time {
	return time.Unix(uuid.UUID(u).Time().UnixTime())
}

// Bytes returns the byte slice representation of the UUID.
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// Bytes returns the byte slice representation of the UUID.
func (u UUID) B58() string {
	return base58.Encode(u[:])
}

// Bytes returns the byte slice representation of the UUID.
func (u *UUID) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return ErrInvalidUUID
	}
	copy(u[:], id[:])
	return nil
}

// Bytes returns the byte slice representation of the UUID.
func (u *UUID) FromB58(s string) error {
	b := base58.Decode(s)
	if len(b) != 16 {
		return ErrInvalidUUID
	}
	copy(u[:], b)
	return nil
}

// UUIDv4 returns a new UUIDv4.
func UUIDv4() UUID {
	for i := 0; i < MaxTry; i++ {
		id, err := uuid.NewRandom()
		if err == nil {
			return UUID(id)
		}
	}
	panic("failed to generate UUIDv4")
}

// UUIDv7 returns a new UUIDv7.
func UUIDv7() UUID {
	for i := 0; i < MaxTry; i++ {
		id, err := uuid.NewV7()
		if err == nil {
			return UUID(id)
		}
	}
	panic("failed to generate UUIDv7")
}
