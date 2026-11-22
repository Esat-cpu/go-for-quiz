package main

import (
	"crypto/rand"
	"encoding/hex"
)

// randHex returns a hex-encoded string of 16 random bytes.
func randHex() string {
	var b [16]byte
	rand.Read(b[:])
	id := hex.EncodeToString(b[:])
	return id
}
