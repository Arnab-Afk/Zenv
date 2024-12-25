package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"time"
)

var (
	currentKey          []byte
	previousKey         []byte
	keyRotationInterval = 24 * time.Hour
)

func StartKeyRotationScheduler() {
	generateNewKey()

	go func() {
		ticker := time.NewTicker(keyRotationInterval)
		for range ticker.C {
			previousKey = currentKey
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
	block, err := aes.NewCipher(currentKey)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func DecryptSecret(ciphertext []byte) ([]byte, error) {
	// Try with current key first
	plaintext, err := decryptWithKey(ciphertext, currentKey)
	if err == nil {
		return plaintext, nil
	}

	// Try with previous key if current fails
	if previousKey != nil {
		plaintext, err = decryptWithKey(ciphertext, previousKey)
		if err == nil {
			return plaintext, nil
		}
	}

	return nil, errors.New("decryption failed")
}

func decryptWithKey(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
