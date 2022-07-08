package utils

import (
	"math/rand"
	"time"
)

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-")

func GenerateRandomNumber(low int, high int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(src)
	return low + rand.Intn(high-low)
}

func GenerateBigString(size int) string {
	rand.Seed(time.Now().UnixNano())

	str := make([]rune, 100)

	for i := range str {
		str[i] = runes[rand.Intn(len(runes))]
	}

	return string(str)
}
