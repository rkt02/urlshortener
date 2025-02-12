package utils

import (
	"strings"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(num int64) string {
	if num == 0 {
		return string(charset[0])
	}

	var encoded strings.Builder
	base := int64(62)

	for num > 0 {
		remainder := num % base
		encoded.WriteByte(charset[remainder])
		num /= base
	}

	return encoded.String()
}
