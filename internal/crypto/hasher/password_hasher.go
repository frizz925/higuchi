package hasher

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hashed string) (bool, error)
}
