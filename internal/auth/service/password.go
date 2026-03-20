package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

const (
	passwordSaltSize = 16
	passwordKeySize  = 32
	passwordIters    = 100_000
	passwordFormat   = "pbkdf2-sha256"
)

func hashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := deriveKey([]byte(password), salt, passwordIters, passwordKeySize)

	return fmt.Sprintf(
		"%s$%d$%s$%s",
		passwordFormat,
		passwordIters,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func verifyPassword(password string, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 {
		return false, errors.New("invalid password hash format")
	}

	if parts[0] != passwordFormat {
		return false, errors.New("unsupported password hash format")
	}

	var iterations int
	if _, err := fmt.Sscanf(parts[1], "%d", &iterations); err != nil || iterations <= 0 {
		return false, errors.New("invalid password hash iterations")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	actualHash := deriveKey([]byte(password), salt, iterations, len(expectedHash))

	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1, nil
}

func deriveKey(password []byte, salt []byte, iter int, keyLen int) []byte {
	hashLen := sha256.Size
	numBlocks := (keyLen + hashLen - 1) / hashLen
	derived := make([]byte, 0, numBlocks*hashLen)

	for block := 1; block <= numBlocks; block++ {
		u := pbkdf2Block(password, salt, iter, block)
		derived = append(derived, u...)
	}

	return derived[:keyLen]
}

func pbkdf2Block(password []byte, salt []byte, iter int, block int) []byte {
	mac := hmac.New(sha256.New, password)
	mac.Write(salt)

	blockIndex := make([]byte, 4)
	binary.BigEndian.PutUint32(blockIndex, uint32(block))
	mac.Write(blockIndex)

	sum := mac.Sum(nil)
	result := make([]byte, len(sum))
	copy(result, sum)

	for i := 1; i < iter; i++ {
		mac = hmac.New(sha256.New, password)
		mac.Write(sum)
		sum = mac.Sum(nil)

		for j := range result {
			result[j] ^= sum[j]
		}
	}

	return result
}
