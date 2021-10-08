package auth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/frizz925/higuchi/internal/crypto/hasher"
)

type FileAuth struct {
	hasher hasher.PasswordHasher
}

func NewFileAuth(h hasher.PasswordHasher) *FileAuth {
	return &FileAuth{h}
}

func (a *FileAuth) ReadPasswordsFile(name string) (Users, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := bufio.NewReaderSize(f, 1024)
	users := make(Users)
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

func (a *FileAuth) WritePasswordsFile(name string, users Users) error {
	f, err := os.Create(name)
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
