package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// SliceContains will return true if the target object is
// within the specified slice
func SliceContains[K comparable](slice []K, target K) bool {
	for _, obj := range slice {
		if obj == target {
			return true
		}
	}
	return false
}

// HashString will return a hex-encoded sha256 hash of the given string
func HashString(in string) string {
	h := sha256.New()
	h.Write([]byte(in))
	return hex.EncodeToString(h.Sum(nil))
}
