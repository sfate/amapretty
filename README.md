# amapretty

A small Go debug-print helper that writes values as indented JSON with a timestamp
and caller reference.

## Usage

Use `go get` to download the dependency.

```bash
go get github.com/sfate/amapretty@latest
```

Then, `import` it in your Go files:

```go
import "github.com/sfate/amapretty"
```

Use `Print` and `Printf` to write to stdout:

```go
amapretty.Print("string value")
amapretty.Print([]struct { Name string }{ { Name: "One" }, { Name: "Chosen" } })
amapretty.Printf("user_id=%d", 123)
```

Use `Fprint` and `Fprintf` when you want to choose the destination writer:

```go
var buf bytes.Buffer

_, err := amapretty.Fprint(&buf, "string value")
_, err = amapretty.Fprintf(&buf, "user_id=%d", 123)
```

Use `Sprint` and `Sprintf` when you want the formatted string:

```go
s := amapretty.Sprint("string value")
s = amapretty.Sprintf("user_id=%d", 123)
```

## Output

Output contains the fixed `amapretty` prefix, an RFC3339 timestamp, the caller
file and line, and an indented JSON array of the provided values:

```text
[amapretty] 2023-02-24T05:02:03Z main.go:101 -- [
        "string value"
]
```

Caller paths inside the current working directory are printed relative to that
directory. Paths outside it remain absolute.

ANSI colors are emitted only when the destination writer is a terminal. Set
`NO_COLOR` to disable colors.

If a value cannot be encoded as JSON, amapretty prints a JSON object containing
the marshal error and a Go-syntax fallback representation of the arguments.

## Preview
[![asciicast](https://asciinema.org/a/573936.svg)](https://asciinema.org/a/573936)

## License

[MIT](/LICENSE)
