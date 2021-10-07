package hasher

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/argon2"
)

const PHCStringFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s"

type Argon2Hasher struct {
	pepper []byte
	params Argon2Params
}

type Argon2Params struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewArgon2Hasher(pepper []byte, params ...Argon2Params) *Argon2Hasher {
	p := Argon2Params{
		Memory:      65536,
		Time:        1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   16,
	}
	if len(params) > 0 {
		x := params[0]
		if x.Memory > 0 {
			p.Memory = x.Memory
		}
		if x.Time > 0 {
			p.Time = x.Time
		}
		if x.Parallelism > 0 {
			p.Parallelism = x.Parallelism
		}
		if x.SaltLength > 0 {
			p.SaltLength = x.SaltLength
		}
		if x.KeyLength > 0 {
			p.KeyLength = x.KeyLength
		}
	}
	return &Argon2Hasher{
		pepper: pepper,
		params: p,
	}
}

func (h *Argon2Hasher) Compare(password, digest string) (int, error) {
	hashed, salt, params, err := h.parse(digest)
	if err != nil {
		return 0, err
	}
	return bytes.Compare(h.hash(password, salt, params), hashed), nil
}

func (h *Argon2Hasher) Hash(password string) (string, error) {
	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}
	hashed := h.hash(password, salt, h.params)
	return h.format(hashed, salt, h.params), nil
}

func (h *Argon2Hasher) generateSalt() ([]byte, error) {
	salt := make([]byte, h.params.SaltLength)
	n, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt[:n], nil
}

func (h *Argon2Hasher) format(hashed, salt []byte, params Argon2Params) string {
	encsalt := base64.RawStdEncoding.EncodeToString(salt)
	enchash := base64.RawStdEncoding.EncodeToString(hashed)
	return fmt.Sprintf(
		PHCStringFormat,
		argon2.Version, params.Memory, params.Time, params.Parallelism,
		fmt.Sprintf("%s$%s", encsalt, enchash),
	)
}

func (h *Argon2Hasher) parse(digest string) (hashed, salt []byte, params Argon2Params, err error) {
	var (
		version int
		tail    string
	)
	_, err = fmt.Sscanf(
		digest, PHCStringFormat,
		&version, &params.Memory, &params.Time, &params.Parallelism, &tail,
	)
	if err != nil {
		return
	}
	parts := strings.Split(tail, "$")
	if len(parts) >= 2 {
		hashed, err = base64.RawStdEncoding.DecodeString(parts[1])
		if err != nil {
			return
		}
		salt, err = base64.RawStdEncoding.DecodeString(parts[0])
	} else {
		hashed, err = base64.RawStdEncoding.DecodeString(tail)
	}
	params.KeyLength = uint32(len(hashed))
	params.SaltLength = uint32(len(salt))
	return
}

func (h *Argon2Hasher) hash(password string, salt []byte, params Argon2Params) []byte {
	pn := len(password)
	pb := make([]byte, pn+len(h.pepper))
	copy(pb, []byte(password))
	copy(pb[pn:], h.pepper)
	return argon2.IDKey(pb, salt, params.Time, params.Memory, params.Parallelism, params.KeyLength)
}
