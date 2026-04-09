package port

type PasswordVerifier interface {
	Compare(hash, password string) bool
}
