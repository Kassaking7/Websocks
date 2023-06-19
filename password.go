package Websocks

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
)

const passwordLength = 256

type password [passwordLength]byte

func (password *password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}

func ParsePassword(passwordString string) (*password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != passwordLength {
		return nil, errors.New("不合法的密码")
	}
	password := password{}
	copy(password[:], bs)
	bs = nil
	return &password, nil
}

func RandPassword() string {
	intArr := rand.Perm(passwordLength)
	password := &password{}
	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			// 确保不会出现如何一个byte位出现重复
			return RandPassword()
		}
	}
	return password.String()
}
