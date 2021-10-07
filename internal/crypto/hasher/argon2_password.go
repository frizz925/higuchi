package hasher

import "bytes"

type Argon2Password struct {
	Hasher *Argon2Hasher
	Hashed []byte
	Salt   []byte
	Params Argon2Params
}

func ParseArgon2DigestToPassword(h *Argon2Hasher, digest string) (pw Argon2Password, err error) {
	pw.Hasher = h
	pw.Hashed, pw.Salt, pw.Params, err = h.parse(digest)
	return
}

func (pw *Argon2Password) Compare(password string) int {
	hashed := pw.Hasher.hash(password, pw.Salt, pw.Params)
	return bytes.Compare(hashed, pw.Hashed)
}
