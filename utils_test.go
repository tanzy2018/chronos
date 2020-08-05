package chronos

import (
	"testing"
)

func TestToString(t *testing.T) {
	testData := []struct {
		name     string
		expected string
		data     []byte
	}{
		{"accii", "012abc", []byte("012abc")},
		{"中文", "中文ab", []byte("中文abc")},

		{"日本語", "こんにちは", []byte("こんにちは")},
	}

	for _, td := range testData {
		if actual := ToString(td.data); td.expected != actual {
			t.Logf("\nname:%s,\nexpected:%s,\nactual:%s\n", td.name, td.expected, actual)
		}
	}
}
