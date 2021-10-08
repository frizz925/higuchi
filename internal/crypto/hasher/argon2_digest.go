package hasher

import "bytes"

type Argon2Digest struct {
	hasher *Argon2Hasher
	digest string
	hashed []byte
	salt   []byte
	params Argon2Params
}

func NewArgon2Digest(h *Argon2Hasher, password string) (ad Argon2Digest, err error) {
	ad.hasher, ad.params = h, h.params
	ad.salt, err = h.generateSalt()
	if err != nil {
		return
	}
	ad.hashed = h.hash(password, ad.salt, ad.params)
	ad.digest = h.format(ad.hashed, ad.salt, ad.params)
	return
}

func (ad Argon2Digest) Compare(password string) int {
	hashed := ad.hasher.hash(password, ad.salt, ad.params)
	return bytes.Compare(ad.hashed, hashed)
}

func (ad Argon2Digest) Digest() string {
	return ad.digest
}
