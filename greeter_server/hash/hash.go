package hash

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
)

const (
	hashString    = `123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	hashStringLen = len(hashString)
)

func GenHash(length uint) string {
	var ret []byte
	for uint(len(ret)) < length {
		ret = append(ret, hashString[rand.Intn(hashStringLen)])
	}
	return string(ret)
}

func NewSha1Hash(data string) string {
	h := sha1.New()
	io.WriteString(h, data)
	//s := sha1.Sum([]byte(data))
	s := h.Sum(nil)
	return fmt.Sprintf("%x", string(s[0:20]))
}

func NewSha1HashWithSalt(data string, saltLength uint) (hash string, salt string) {
	salt = GenHash(saltLength)
	hash = NewSha1Hash(data + salt)
	return hash, salt
}
