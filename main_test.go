package amapretty

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	timeNow = func() time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", "2023-02-24 05:02:03")
		return t
	}

	os.Exit(m.Run())
}

func setOutput(t *testing.T) func() string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	output = w

	return func() string {
		require.NoError(t, w.Close())
		out, err := io.ReadAll(r)
		require.NoError(t, err)
		output = os.Stdout
		return string(out)
	}
}

func setTerminal(t *testing.T, enabled bool) {
	t.Helper()

	previousIsTerminal := isTerminal
	previousLookupEnv := lookupEnv
	isTerminal = func(w io.Writer) bool {
		return enabled
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	t.Cleanup(func() {
		isTerminal = previousIsTerminal
		lookupEnv = previousLookupEnv
	})
}

func TestPrint(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		require.Equal(t, 3, skip)
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

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
			outputCallbackF := setOutput(t)
			Print(c.args)
			require.Equal(t, c.expected, outputCallbackF())
		})
	}
}

func TestPrintf(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		require.Equal(t, 3, skip)
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	outputCallbackF := setOutput(t)
	Printf("dime: %d, val: %s, time: %v", 123, "none", timeNow().Format(time.RFC3339))
	expected := "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m -- [\n\t\"dime: 123, val: none, time: 2023-02-24T05:02:03Z\"\n]\n"
	require.Equal(t, expected, outputCallbackF())
}

func TestFprint(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		require.Equal(t, 3, skip)
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	var buf bytes.Buffer
	n, err := Fprint(&buf, "test")
	expected := "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m -- [\n\t\"test\"\n]\n"
	require.NoError(t, err)
	require.Equal(t, len(expected), n)
	require.Equal(t, expected, buf.String())
}

func TestFprintf(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		require.Equal(t, 3, skip)
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	var buf bytes.Buffer
	n, err := Fprintf(&buf, "value: %d", 123)
	expected := "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m -- [\n\t\"value: 123\"\n]\n"
	require.NoError(t, err)
	require.Equal(t, len(expected), n)
	require.Equal(t, expected, buf.String())
}

func TestFprintOmitsColorForNonTerminalWriter(t *testing.T) {
	setTerminal(t, false)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	var buf bytes.Buffer
	_, err := Fprint(&buf, "test")
	require.NoError(t, err)
	require.Equal(t, "[amapretty] 2023-02-24T05:02:03Z /Users/username/path/project/main.go:101 -- [\n\t\"test\"\n]\n", buf.String())
}

func TestFprintRespectsNoColor(t *testing.T) {
	setTerminal(t, true)
	lookupEnv = func(key string) (string, bool) {
		if key == "NO_COLOR" {
			return "1", true
		}
		return "", false
	}

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	var buf bytes.Buffer
	_, err := Fprint(&buf, "test")
	require.NoError(t, err)
	require.Equal(t, "[amapretty] 2023-02-24T05:02:03Z /Users/username/path/project/main.go:101 -- [\n\t\"test\"\n]\n", buf.String())
}

func TestPrintUsesRelativeCallerForLocalPath(t *testing.T) {
	setTerminal(t, true)
	wd, err := os.Getwd()
	require.NoError(t, err)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), filepath.Join(wd, "main.go"), 101, true
	}

	outputCallbackF := setOutput(t)
	Print("test")
	expected := "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36mmain.go:101\x1b[0m -- [\n\t\"test\"\n]\n"
	require.Equal(t, expected, outputCallbackF())
}

func TestPrintWithMultipleArguments(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	outputCallbackF := setOutput(t)
	Print("first", 2)
	expected := "[\x1b[1;32mamapretty\x1b[0m] \x1b[1;34m2023-02-24T05:02:03Z\x1b[0m \x1b[1;36m/Users/username/path/project/main.go:101\x1b[0m -- [\n\t\"first\",\n\t2\n]\n"
	require.Equal(t, expected, outputCallbackF())
}

func TestPrintWithUnsupportedJSONValue(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	outputCallbackF := setOutput(t)
	Print(func() {})
	out := outputCallbackF()
	require.Contains(t, out, `"error": "json: unsupported type: func()"`)
	require.Contains(t, out, `"args": "[]interface {}{(func())(`)
}

func TestPrintWithoutCaller(t *testing.T) {
	setTerminal(t, true)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "", 0, false
	}

	outputCallbackF := setOutput(t)
	Print("test")
	out := outputCallbackF()
	require.NotContains(t, out, "\x1b[1;36m")
	require.True(t, strings.Contains(out, " -- [\n\t\"test\"\n]\n"))
}

func TestFprintSerializesConcurrentWrites(t *testing.T) {
	setTerminal(t, false)

	runtimeCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return uintptr(0), "/Users/username/path/project/main.go", 101, true
	}

	w := &overlapDetectingWriter{}
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := Fprint(w, i)
			if err != nil {
				t.Errorf("Fprint() error = %v", err)
			}
		}(i)
	}
	wg.Wait()

	require.Zero(t, w.overlaps.Load())
}

type overlapDetectingWriter struct {
	active   atomic.Int32
	overlaps atomic.Int32
}

func (w *overlapDetectingWriter) Write(p []byte) (int, error) {
	if w.active.Add(1) > 1 {
		w.overlaps.Add(1)
	}
	time.Sleep(time.Millisecond)
	w.active.Add(-1)
	return len(p), nil
}
