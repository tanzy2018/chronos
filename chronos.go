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
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

// Chronos ... chronos,a test time comsume interval.
type Chronos interface {
	// WriteTo ..
	io.WriterTo
	// Reset .. reset the chronos timer.
	Reset()
	// Add .. add a new time.Time with key,if existed then replace.
	Add(key string) error
	// Link .. link key with the histrory key.
	Link(key string, fromKey string) error
	// Consume .. the time interval between key and its link or the the init key.
	Consume(key string) (time.Duration, bool)
	// ConsumeFrom .. the time interval between key and specified key.
	ConsumeFrom(key string, fromKey string) (time.Duration, bool)
}

// Ichronos ...
type Ichronos struct {
	initKey string
	c       int32
	m       map[string]*timer
}

type timer struct {
	key     string
	linkKey string
	tm      time.Time
	ic      *Ichronos
}

var globalChronos = New()

// errors
var (
	ErrKeyDuplicated = errors.New("key is existed,duplicated")
	ErrKeyNotfound   = errors.New("key not found")
	ErrDisorderedKey = errors.New("key and fromKey is disordered")
)

var zeroDuration = time.Duration(0)
var initKey = "chronos-init"

// New ... generate a new chronos.
func New() Chronos {
	c := &Ichronos{
		initKey: wrapInitKey(),
		m:       make(map[string]*timer),
	}
	atomic.StoreInt32(&c.c, 0)
	c.m[c.initKey] = &timer{
		key:     c.initKey,
		linkKey: c.initKey,
		tm:      time.Now(),
		ic:      c,
	}
	return c
}

// WriteTo ..
func (ic *Ichronos) WriteTo(w io.Writer) (n int64, err error) {
	ic.lock()
	tmStrings := make([]string, 0, len(ic.m))
	for _, tm := range ic.m {
		tmStrings = append(tmStrings, strings.Trim(tm.String(), "\n"))

	}
	ic.unlock()

	sort.Strings(tmStrings)
	var buf bytes.Buffer
	for _, str := range tmStrings {
		buf.WriteString(str)
		buf.WriteByte('\n')
	}
	return buf.WriteTo(w)
}

// Reset .. reset the chronos timer.
func (ic *Ichronos) Reset() {
	ic.lock()
	defer ic.unlock()
	ic.initKey = wrapInitKey()
	ic.m = make(map[string]*timer)
	ic.m[ic.initKey] = &timer{
		key:     ic.initKey,
		linkKey: ic.initKey,
		tm:      time.Now(),
		ic:      ic,
	}
}

// Add .. add a new time.Time with key,if existed then replace.
func (ic *Ichronos) Add(key string) error {
	tm := time.Now()
	ic.lock()
	defer ic.unlock()
	if _, ok := ic.m[key]; ok {
		return ErrKeyDuplicated
	}
	ic.m[key] = &timer{
		key:     key,
		linkKey: ic.initKey,
		tm:      tm,
		ic:      ic,
	}
	return nil
}

// Link .. link key with the histrory key.
func (ic *Ichronos) Link(key string, fromKey string) error {
	ic.lock()
	defer ic.unlock()
	var tm0, tm1 *timer
	var ok bool
	if tm0, ok = ic.m[key]; !ok {
		return ErrKeyNotfound
	}

	if tm1, ok = ic.m[fromKey]; !ok {
		return ErrKeyNotfound
	}

	if tm0.tm.Before(tm1.tm) {
		return ErrDisorderedKey
	}

	tm0.linkKey = fromKey
	ic.m[key] = tm0
	return nil
}

// Consume .. the time interval between key and its link or the the init key.
func (ic *Ichronos) Consume(key string) (time.Duration, bool) {
	ic.lock()
	defer ic.unlock()
	tm, ok := ic.m[key]
	if !ok {
		return zeroDuration, false
	}
	return tm.tm.Sub(ic.m[tm.linkKey].tm), true
}

// ConsumeFrom .. the time interval between key and specified key.
func (ic *Ichronos) ConsumeFrom(key string, fromkey string) (time.Duration, bool) {
	ic.lock()
	defer ic.unlock()
	tm0, ok := ic.m[key]
	if !ok {
		return zeroDuration, false
	}
	tm1, ok := ic.m[fromkey]
	if !ok {
		return zeroDuration, false
	}

	if tm0.tm.Before(tm1.tm) {
		return zeroDuration, false
	}

	return tm0.tm.Sub(tm1.tm), true
}

func (ic *Ichronos) lock() {
	for {
		if atomic.CompareAndSwapInt32(&ic.c, 0, 1) {
			return
		}
	}
}

func (ic *Ichronos) unlock() {
	for {
		if atomic.CompareAndSwapInt32(&ic.c, 1, 0) {
			return
		}
	}
}

func wrapInitKey() string {
	return initKey + "-" + RandomString(4)
}

// WriteTo ..
func WriteTo(w io.Writer) (int64, error) {
	return globalChronos.WriteTo(w)
}

// Reset .. reset the chronos timer.
func Reset() {
	globalChronos.Reset()
}

// Add .. add a new time.Time with key,if existed then replace.
func Add(key string) error {
	return globalChronos.Add(key)
}

// Link .. link key with the histrory key.
func Link(key string, fromKey string) error {
	return globalChronos.Link(key, fromKey)
}

// Consume .. the time interval between key and its link or the the init key.
func Consume(key string) (time.Duration, bool) {
	return globalChronos.Consume(key)
}

// ConsumeFrom .. the time interval between key and specified key.
func ConsumeFrom(key string, fromKey string) (time.Duration, bool) {
	return globalChronos.ConsumeFrom(key, fromKey)
}

func (t *timer) String() string {
	m := make(map[string]interface{})
	m["from"] = t.linkKey
	m["to"] = t.key
	m["consume"] = t.tm.Sub(t.ic.m[t.linkKey].tm)
	m["unit"] = "ns"
	out, _ := json.Marshal(m)
	return ToString(out)
}
