package security

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordVerifier struct{}

func NewBcryptPasswordVerifier() *BcryptPasswordVerifier {
	return &BcryptPasswordVerifier{}
}

func (v *BcryptPasswordVerifier) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
