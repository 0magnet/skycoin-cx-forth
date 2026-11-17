package forth

import (
	"os"
	"strings"
	"testing"
)

// testdata/ternary.fs is a verbatim copy of
// https://raw.githubusercontent.com/MPETREMANN11/eForth-Windows/main/INFERENCE/ternary.fs
// (GPL, adapted from https://rosettacode.org/wiki/Ternary_logic#Forth). It is
// test data only and is not built into any binary.

// Kleene tables in T F ? letters, rows a and columns b in the order T F ?.
var kleeneLetters = map[string][]string{
	"NOT":   {"F", "T", "?"},
	"AND":   {"T F ?", "F F F", "? F ?"},
	"OR":    {"T T T", "T F ?", "T ? ?"},
	"XOR":   {"F T ?", "T F ?", "? ? ?"},
	"IMPLY": {"T F ?", "T T T", "T ? ?"},
	"EQUIV": {"T F ?", "F T ?", "? ? ?"},
}

// parseTables reads the demo output into section name -> row cells.
func parseTables(out string) map[string][]string {
	tabs := map[string][]string{}
	cur := ""
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			cur = line[1 : len(line)-1]
			continue
		}
		if row, ok := strings.CutPrefix(line, "T | "); ok {
			tabs[cur] = append(tabs[cur], strings.Join(strings.Fields(row), " "))
		} else if row, ok := strings.CutPrefix(line, "F | "); ok {
			tabs[cur] = append(tabs[cur], strings.Join(strings.Fields(row), " "))
		} else if row, ok := strings.CutPrefix(line, "? | "); ok {
			tabs[cur] = append(tabs[cur], strings.Join(strings.Fields(row), " "))
		}
	}
	return tabs
}

func checkTables(t *testing.T, out string) {
	t.Helper()
	got := parseTables(out)
	for name, want := range kleeneLetters {
		if strings.Join(got[name], "/") != strings.Join(want, "/") {
			t.Errorf("[%s] got %q want %q", name, got[name], want)
		}
	}
}

func TestTernaryFsVerbatim(t *testing.T) {
	src, err := os.ReadFile("testdata/ternary.fs")
	if err != nil {
		t.Fatal(err)
	}
	out, err := Run(string(src))
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.HasPrefix(out, "\x1b[2J\x1b[H") {
		t.Errorf("PAGE did not run first: %q", out)
	}
	checkTables(t, out)
}

// The native words give the same tables on eForth integers and on tagged
// ternary cells, without loading ternary.fs.
func TestNativeTernaryVocabulary(t *testing.T) {
	ops := []struct{ name, word string }{
		{"AND", "TAND"}, {"OR", "TOR"}, {"XOR", "TXOR"}, {"IMPLY", "TIMPLY"}, {"EQUIV", "TEQUIV"},
	}
	for _, enc := range []struct {
		name string
		vals [3]string
	}{
		{"eforth", [3]string{"-1", "0", "1"}},
		{"tagged", [3]string{"TRUE", "FALSE", "MAYBE"}},
		{"mixed", [3]string{"-1", "FALSE", "NULL"}},
	} {
		var b strings.Builder
		b.WriteString(`." [NOT]" CR `)
		for _, a := range enc.vals {
			b.WriteString(a + ` T. ."  | " ` + a + " TNOT T. CR\n")
		}
		for _, op := range ops {
			b.WriteString(`." [` + op.name + `]" CR `)
			for _, a := range enc.vals {
				b.WriteString(a + ` T. ."  | " `)
				for _, c := range enc.vals {
					b.WriteString(a + " " + c + " " + op.word + " T. SPACE ")
				}
				b.WriteString("CR\n")
			}
		}
		src := b.String()
		out, err := Run(src)
		if err != nil {
			t.Fatalf("%s: %v\n%s", enc.name, err, out)
		}
		checkTables(t, out)
	}
}

func TestEforthIntegerEncoding(t *testing.T) {
	// Integer results stay integers in the eForth encoding.
	cases := []struct{ src, want string }{
		{"1 TNOT . -1 TNOT . 0 TNOT .", "1 0 -1 "},
		{"-1 1 TAND . 0 1 TOR . 1 1 TEQUIV . 1 0 TIMPLY .", "1 1 1 1 "},
		{"TRUE MAYBE TAND .", "null "},
		{"MAYBE NULL? .", "true "},
	}
	for _, c := range cases {
		if got := run(t, c.src); got != c.want {
			t.Errorf("%q: got %q want %q", c.src, got, c.want)
		}
	}
}
