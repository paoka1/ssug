package utils

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

func IsURL(s string) bool {
	urlRegex := regexp.MustCompile(`^(https?|http)://[^\s/$.?#].\S*$`)
	return urlRegex.MatchString(s)
}

func IsLegalValue(v string) bool {
	urlRegex := regexp.MustCompile(`^[0-9A-Fa-f]*$`)
	return urlRegex.MatchString(v)
}

func CalculateMD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}
