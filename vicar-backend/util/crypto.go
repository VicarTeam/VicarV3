package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomStringBasedOnFormat(format, letters string) (string, error) {
	ret := make([]byte, 0)
	for i := 0; i < len(format); i++ {
		if format[i] == '{' || format[i] == '[' {
			i++
			num := 0
			for ; i < len(format) && format[i] != '}' && format[i] != ']'; i++ {
				num = num*10 + int(format[i]-'0')
			}
			if i == len(format) {
				return "", nil
			}
			str, err := GenerateRandomStringSpecific(num, letters)
			if err != nil {
				return "", err
			}
			ret = append(ret, []byte(str)...)
		} else {
			ret = append(ret, format[i])
		}
	}
	return string(ret), nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GenerateRandomStringSpecific(n int, letters string) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GenerateRandomNumericString(n int) (string, error) {
	const letters = "0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
