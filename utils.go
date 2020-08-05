package chronos

import (
	"math/rand"
	"time"
	"unsafe"
)

var utilsRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// ToString ... convert []byte into string.
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// RandomString ...
func RandomString(n int) string {
	if n <= 0 {
		return ""
	}
	tpl := []byte("23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")

	for i := 0; i < len(tpl)/2; i++ {
		idx := utilsRand.Intn(len(tpl) - i)
		tpl[idx], tpl[len(tpl)-i-1] = tpl[len(tpl)-i-1], tpl[idx]
	}

	if n > len(tpl) {
		return ToString(tpl[:])
	}

	return ToString(tpl[:n])
}
