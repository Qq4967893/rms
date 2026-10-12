package crypto

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Md5(plaintext string) string {
	data := []byte(plaintext)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func PasswordHash(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), 10)
	return string(bytes), err
}

func PasswordVerify(plaintext, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintext))
	return err == nil
}
