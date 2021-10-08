package hasher

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

var _ PasswordHasher = (*MD5Hasher)(nil)

type MD5Hasher struct {
	pepper     []byte
	saltLength int
}

func NewMD5Hasher(pepper []byte, saltLength ...int) *MD5Hasher {
	ah := &MD5Hasher{
		pepper:     pepper,
		saltLength: 16,
	}
	if len(saltLength) > 0 && saltLength[0] > 0 {
		ah.saltLength = saltLength[0]
	}
	return ah
}

func (h *MD5Hasher) Hash(password string) (string, error) {
	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}
	hashed := h.hash(password, salt)
	return h.format(hashed, salt), nil
}

func (h *MD5Hasher) Compare(password, digest string) (int, error) {
	hashed, salt, err := h.parse(digest)
	if err != nil {
		return 0, err
	}
	phashed := h.hash(password, salt)
	return bytes.Compare(phashed, hashed), nil
}

func (h *MD5Hasher) ParseDigest(digest string) (PasswordDigest, error) {
	var err error
	md := MD5Digest{
		hasher: h,
		digest: digest,
	}
	md.hashed, md.salt, err = h.parse(digest)
	return md, err
}

func (h *MD5Hasher) format(hashed, salt []byte) string {
	return fmt.Sprintf(
		"$apr1$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashed),
	)
}

func (h *MD5Hasher) parse(digest string) (hashed, salt []byte, err error) {
	var tail string
	_, err = fmt.Sscanf(digest, "$apr1$%s", &tail)
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
		if err != nil {
			return
		}
	} else {
		hashed, err = base64.RawStdEncoding.DecodeString(tail)
		if err != nil {
			return
		}
	}
	return
}

func (h *MD5Hasher) hash(password string, salt []byte) []byte {
	pn, sn := len(password), len(salt)
	pb := make([]byte, pn+sn+len(h.pepper))
	copy(pb, []byte(password))
	copy(pb[pn:], salt)
	copy(pb[pn+sn:], h.pepper)
	hashed := md5.Sum(pb)
	return hashed[:]
}

func (h *MD5Hasher) generateSalt() ([]byte, error) {
	salt := make([]byte, h.saltLength)
	n, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt[:n], nil
}
