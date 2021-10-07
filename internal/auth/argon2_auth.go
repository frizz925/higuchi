package auth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/frizz925/higuchi/internal/crypto/hasher"
)

type Argon2Users map[string]hasher.Argon2Digest

type Argon2Auth struct {
	hasher *hasher.Argon2Hasher
}

func NewArgon2Auth(h *hasher.Argon2Hasher) *Argon2Auth {
	return &Argon2Auth{h}
}

func (a *Argon2Auth) ReadPasswordsFile(name string) (Argon2Users, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := bufio.NewReaderSize(f, 1024)
	users := make(Argon2Users)
	for i := 0; ; i++ {
		b, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		line := string(b)
		idx := strings.Index(line, ":")
		if idx < 0 {
			return nil, fmt.Errorf("malformed user entry at line %d", i+1)
		}
		user, digest := line[:idx], line[idx+1:]
		ad, err := a.hasher.ParseDigest(digest)
		if err != nil {
			return nil, fmt.Errorf("failed to parse password digest at line %d: %+v", i+1, err)
		}
		users[user] = ad
	}

	return users, nil
}

func (a *Argon2Auth) WritePasswordsFile(name string, users Argon2Users) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	for user, ad := range users {
		line := fmt.Sprintf("%s:%s\r\n", user, ad.Digest())
		bw.WriteString(line)
		if err := bw.Flush(); err != nil {
			return err
		}
	}
	return nil
}
