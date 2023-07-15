package scipher

import (
	"crypto/rand"
)

// generates a random session key for a symmetric algorithm 
func SKeyGen(size int) []byte {
	key := make([]byte, size)
	rand.Read(key)

	return key
}
