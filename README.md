# skycoin-cx-forth

A small Forth VM in Go that runs a base subset of CX IR, with three-valued
(ternary) logic. It builds with Go for the host and with TinyGo for
microcontrollers.

- `forth/` is the VM: fixed-size stacks and memory, colon definitions,
  IF/ELSE/THEN, loops, variables, and tagged cells (int, f64, f32, ternary,
  string).
- `cxir/` holds the CX IR subset and its translator to Forth words.
- `cmd/cxforth` is the host REPL. `cmd/cxforth-esp32` is a UART REPL for ESP32.

## Ternary logic

Cells can be TRUE, FALSE or NULL (unknown), with Kleene AND, OR, NOT and
IMPLIES. Comparisons on NULL give NULL, and arithmetic on a ternary value is
an error. `TRI-IF ... ELSE ... NULL-ELSE ... THEN` branches three ways.

The words MAYBE, TNOT, TAND, TOR, TEQUIV, TIMPLY, TXOR and T. follow
PETREMANN's eForth/ESP32forth
[ternary.fs](https://github.com/MPETREMANN11/eForth-Windows/blob/main/INFERENCE/ternary.fs),
and accept both its integer encoding (TRUE=-1, FALSE=0, MAYBE=1) and the
tagged cells. See `forth/doc.go`.

## Usage

```
go run ./cmd/cxforth -list
go run ./cmd/cxforth -cx fib
go run ./cmd/cxforth -e 'TRUE NULL AND .'
go test ./...
tinygo build -target=esp32c3-generic ./cmd/cxforth-esp32
```

The ESP32-C3 build is checked. The classic (Xtensa) ESP32 needs a TinyGo whose
LLVM includes the Xtensa backend. Nothing has run on hardware yet.

## License

See LICENSE. `forth/testdata/ternary.fs` is a verbatim copy of PETREMANN's
file, kept as test data under its own GNU GPL header. It is not built into
any binary.
