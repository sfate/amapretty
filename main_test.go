package amapretty

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	timeNow = func() time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", "2023-02-24 05:02:03")
		return t
	}

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}
}

func setOutput() func() string {
	r, w, _ := os.Pipe()
	output = w

	return func() string {
		w.Close()
		out, _ := io.ReadAll(r)
		output = os.Stdout
		return string(out)
	}
}

func TestPrint(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		args     interface{}
	}{
		{
			name:     "with simple string argument",
			args:     "test",
			expected: "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m -- [\n\t\"test\"\n]\n",
		},
		{
			name:     "with complex args",
			args:     []struct{ Name string }{{Name: "One"}, {Name: "Chosen"}},
			expected: "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m -- [\n\t[\n\t\t{\n\t\t\t\"Name\": \"One\"\n\t\t},\n\t\t{\n\t\t\t\"Name\": \"Chosen\"\n\t\t}\n\t]\n]\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			outputCallbackF := setOutput()
			Print(c.args)
			assert.Equal(t, c.expected, outputCallbackF())
		})
	}
}
