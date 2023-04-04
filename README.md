# amapretty

A basic stdout shout printer (prints values as json objects) with few colors and caller reference.

## Usage

Use `go get` to download the dependency.

```bash
go get github.com/sfate/amapretty@latest
```

Then, `import` it in your Go files:

```go
import "github.com/sfate/amapretty"
```

This lib comes with `Print` function which accepts any type (and amount) of interface(-s).

```go
amapretty.Print("string value")
amapretty.Print([]struct { Name string }{ { Name: "One" }, { Name: "Chosen" } })
```

## Preview
[![asciicast](https://asciinema.org/a/573936.svg)](https://asciinema.org/a/573936)

## License

[MIT](/LICENSE)
