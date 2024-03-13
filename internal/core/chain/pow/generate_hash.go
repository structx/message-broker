// Package pow proof of work package
package pow

import (
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

// HashWithSHA3AndDifficulty generate sha3 hash with difficulty check
func HashWithSHA3AndDifficulty(timestamp, prevHash string, data []byte, difficulty int) string {

	for i := 10; ; i++ {

		h := hashSHA3(i, timestamp, prevHash, data)
		if isHashValid(h, difficulty) {
			return h
		}
	}
}

func hashSHA3(nounce int, timestamp, prevHash string, data []byte) string {
	h := sha3.New224()
	h.Write([]byte(fmt.Sprintf("%x", nounce)))
	h.Write([]byte(timestamp + prevHash))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func isHashValid(input string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(input, prefix)
}
