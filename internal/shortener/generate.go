package shortener

import (
	"crypto/rand"
	"math/big"
)

const alphabet string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func generateCode(n int) (string, error) {
	res := make([]byte, n)
	alphabet_len := len(alphabet)
	max := big.NewInt(int64(alphabet_len))
	for i := 0; i < alphabet_len; i++ {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		res[i] = alphabet[idx.Int64()]
	}
	return string(res), nil
}
