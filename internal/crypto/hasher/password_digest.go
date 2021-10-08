package hasher

type PasswordDigest interface {
	Compare(password string) int
	Digest() string
}
