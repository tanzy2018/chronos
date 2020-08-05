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
	"os"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	type args struct {
		key     string
		fromKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"key1", "key0"}, false},
		{"disordered", args{"key0", "key1"}, true},
		{"left-exists", args{"key0", "notexists"}, true},
		{"right-exists", args{"notexists", "key1"}, true},
		{"no-exsist", args{"notexist0", "notexist1"}, true},
	}
	Reset()
	Add("key0")
	Add("key1")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Link(tt.args.key, tt.args.fromKey); (err != nil) != tt.wantErr {
				t.Errorf("Link() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteTo(t *testing.T) {
	outFile := "testdata/chronos.txt"
	os.Remove(outFile)
	fi, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Logf("should open file,but failed err=%v", err)
	}
	defer fi.Close()
	Reset()
	Add("key0")
	Add("key1")
	Add("key2")
	Add("key3")
	Link("key1", "key0")
	Link("key2", "key0")
	Link("key3", "key1")
	WriteTo(fi)
}

func TestReset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"reset"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"unexisted",
			args{key: "key0"},
			false,
		},

		{
			"existed",
			args{key: "key0"},
			true,
		},
	}
	Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Add(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestConsume(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Duration
		want1 bool
	}{
		{"existed", args{"key0"}, 0, true},
		{"unexisted", args{"key1"}, 0, false},
	}
	Reset()
	Add("key0")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Consume(tt.args.key)

			if got1 != tt.want1 {
				t.Errorf("Consume() got1 = %v, want %v", got1, tt.want1)
			}

			if !got1 && got != tt.want {
				t.Errorf("Consume() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsumeFrom(t *testing.T) {
	type args struct {
		key     string
		fromKey string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Duration
		want1 bool
	}{
		{"key-notexist", args{"key2", "key1"}, 0, false},
		{"fromkey-notexist", args{"key1", "key2"}, 0, false},
		{"key-fromkey-notexist", args{"key2", "key3"}, 0, false},
		{"key-fromkey-disorder", args{"key0", "key1"}, 0, false},
		{"success", args{"key1", "key0"}, 0, true},
	}
	Reset()
	Add("key0")
	Add("key1")
	Link("key1", "key0")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ConsumeFrom(tt.args.key, tt.args.fromKey)

			if got1 != tt.want1 {
				t.Errorf("ConsumeFrom() got1 = %v, want %v", got1, tt.want1)
			}

			if !got1 && got != tt.want {
				t.Errorf("ConsumeFrom() got = %v, want %v", got, tt.want)
			}

			if got1 && got == 0 {
				t.Errorf("ConsumeFrom() got = %v, want x > %v", got, tt.want)
			}
		})
	}
}
