// MIT License

// Copyright (c) 2020 tanzy2018

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
