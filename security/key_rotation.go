package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"time"
)

var (
	currentKey          []byte
	keyRotationInterval = 24 * time.Hour
)

func StartKeyRotationScheduler() {
	generateNewKey()

	go func() {
		ticker := time.NewTicker(keyRotationInterval)
		for range ticker.C {
			generateNewKey()
		}
	}()
}

func generateNewKey() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	currentKey = key
}

func EncryptSecret(plaintext []byte) []byte {
	block, _ := aes.NewCipher(currentKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	return gcm.Seal(nonce, nonce, plaintext, nil)
}
