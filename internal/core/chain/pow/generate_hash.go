// Package pow proof of work package
package pow

import (
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

// HashWithSHA3AndDifficulty generate sha3 hash with difficulty check
func HashWithSHA3AndDifficulty(timestamp string, difficulty int) string {

	for i := 0; ; i++ {

		n := fmt.Sprintf("%x", i)
		h := hashSHA3(n, timestamp)
		if isHashValid(h, difficulty) {
			return h
		}
	}
}

func hashSHA3(nounce, timestamp string) string {
	h := sha3.New224()
	h.Write([]byte(nounce + timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func isHashValid(input string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(input, prefix)
}
