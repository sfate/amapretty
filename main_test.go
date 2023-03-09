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

func TestColors(t *testing.T) {
	t.Run("check color codes", func(t *testing.T) {
		assert.Equal(t, "\x1b[1;30mtest\x1b[0m", black("test"))
		assert.Equal(t, "\x1b[1;31mtest\x1b[0m", red("test"))
		assert.Equal(t, "\x1b[1;32mtest\x1b[0m", green("test"))
		assert.Equal(t, "\x1b[1;33mtest\x1b[0m", yellow("test"))
		assert.Equal(t, "\x1b[1;34mtest\x1b[0m", purple("test"))
		assert.Equal(t, "\x1b[1;35mtest\x1b[0m", magenta("test"))
		assert.Equal(t, "\x1b[1;36mtest\x1b[0m", teal("test"))
		assert.Equal(t, "\x1b[1;37mtest\x1b[0m", white("test"))
	})
}

func bufferOutput(args ...interface{}) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print(args)

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	return string(output)
}

func TestPrint(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		args     interface{}
	}{
		{
			name:     "with simple string argument",
			expected: "\n\x1b[1;32m[PrettyPrint] \x1b[1;34m02-24-2023 05:02:03\x1b[0m -- \x1b[0m\x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m\n\x1b[1;32m[PrettyPrint] \x1b[1;34m02-24-2023 05:02:03\x1b[0m -- \x1b[0m[\n\t[\n\t\t\"test\"\n\t]\n]\n",
			args:     "test",
		},
		{
			name:     "with complex args",
			expected: "\n\x1b[1;32m[PrettyPrint] \x1b[1;34m02-24-2023 05:02:03\x1b[0m -- \x1b[0m\x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m\n\x1b[1;32m[PrettyPrint] \x1b[1;34m02-24-2023 05:02:03\x1b[0m -- \x1b[0m[\n\t[\n\t\t[\n\t\t\t{\n\t\t\t\t\"Name\": \"One\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"Name\": \"Chosen\"\n\t\t\t}\n\t\t]\n\t]\n]\n",
			args:     []struct{ Name string }{{Name: "One"}, {Name: "Chosen"}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := bufferOutput(c.args)
			assert.Equal(t, c.expected, output)
		})
	}
}
