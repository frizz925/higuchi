package hasher

import "bytes"

type MD5Digest struct {
	hasher *MD5Hasher
	hashed []byte
	salt   []byte
	digest string
}

func NewMD5Digest(h *MD5Hasher, password string) (md MD5Digest, err error) {
	md.hasher = h
	md.salt, err = h.generateSalt()
	if err != nil {
		return
	}
	md.hashed = h.hash(password, md.salt)
	md.digest = h.format(md.hashed, md.salt)
	return
}

func (d MD5Digest) Compare(password string) int {
	hashed := d.hasher.hash(password, d.salt)
	return bytes.Compare(d.hashed, hashed)
}

func (d MD5Digest) Digest() string {
	return d.digest
}
