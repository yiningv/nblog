package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func MD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

func Sha1(value string) string {
	s := sha1.New()
	s.Write([]byte(value))
	return hex.EncodeToString(s.Sum(nil))
}
