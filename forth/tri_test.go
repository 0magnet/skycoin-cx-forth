package forth

import (
	"strings"
	"testing"
)

var triWords = []string{"FALSE", "NULL", "TRUE"}

// Expected tables written out by hand, rows a = F N T, columns b = F N T.
var binTables = map[string][3][3]string{
	"AND": {
		{"false", "false", "false"},
		{"false", "null", "null"},
		{"false", "null", "true"},
	},
	"OR": {
		{"false", "null", "true"},
		{"null", "null", "true"},
		{"true", "true", "true"},
	},
	"IMPLIES": {
		{"true", "true", "true"},
		{"null", "null", "true"},
		{"false", "null", "true"},
	},
	"LIMPLIES": {
		{"true", "true", "true"},
		{"null", "true", "true"},
		{"false", "null", "true"},
	},
	"XOR": {
		{"false", "null", "true"},
		{"null", "null", "null"},
		{"true", "null", "false"},
	},
	"EQV": {
		{"true", "null", "false"},
		{"null", "null", "null"},
		{"false", "null", "true"},
	},
	"=": {
		{"true", "null", "false"},
		{"null", "null", "null"},
		{"false", "null", "true"},
	},
}

func TestTruthTables(t *testing.T) {
	for op, tab := range binTables {
		for i, a := range triWords {
			for j, b := range triWords {
				got := strings.TrimSpace(run(t, a+" "+b+" "+op+" ."))
				if got != tab[i][j] {
					t.Errorf("%s %s %s = %s, want %s", a, b, op, got, tab[i][j])
				}
			}
		}
	}
	not := map[string]string{"FALSE": "true", "NULL": "null", "TRUE": "false"}
	for a, want := range not {
		if got := strings.TrimSpace(run(t, a+" NOT .")); got != want {
			t.Errorf("%s NOT = %s, want %s", a, got, want)
		}
	}
}

func TestNullPropagation(t *testing.T) {
	cases := []struct{ src, want string }{
		{"NULL 3 < .", "null "},
		{"3 NULL = .", "null "},
		{"NULL 0= .", "null "},
		{"NULL >I32 . NULL >INT . NULL >F32 . NULL >F64 .", "null null null null "},
		{"3 4 < . 4 3 < . 3 3 <= .", "true false true "},
		{"NULL KNOWN? . NULL NULL? . 5 KNOWN? .", "false true true "},
		{"0 >TRI . 7 >TRI .", "false true "},
		// Forth integer flags mix with ternary values.
		{"-1 NULL AND . 0 NULL AND . 0 NULL OR .", "null false null "},
	}
	for _, c := range cases {
		if got := run(t, c.src); got != c.want {
			t.Errorf("%q: got %q want %q", c.src, got, c.want)
		}
	}
}

func TestConditionals(t *testing.T) {
	def := `: PLAIN IF ." yes" ELSE ." no" THEN ; ` +
		`: TRI TRI-IF ." yes" ELSE ." no" NULL-ELSE ." unknown" THEN ; ` +
		`: TRI2 TRI-IF ." yes" ELSE ." no" THEN ; ` +
		`: TRI3 TRI-IF ." yes" NULL-ELSE ." unknown" THEN ." ." ; `
	cases := []struct{ src, want string }{
		{"TRUE PLAIN FALSE PLAIN NULL PLAIN", "yesnono"},
		{"TRUE TRI FALSE TRI NULL TRI", "yesnounknown"},
		{"TRUE TRI2 FALSE TRI2 NULL TRI2", "yesnono"},
		{"TRUE TRI3 FALSE TRI3 NULL TRI3", "yes..unknown."},
		{"1 TRI 0 TRI", "yesno"},
	}
	for _, c := range cases {
		if got := run(t, def+c.src); got != c.want {
			t.Errorf("%q: got %q want %q", c.src, got, c.want)
		}
	}
}
