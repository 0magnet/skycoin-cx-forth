package forth

import (
	"strings"
	"testing"
)

func run(t *testing.T, src string) string {
	t.Helper()
	out, err := Run(src)
	if err != nil {
		t.Fatalf("%q: %v (output %q)", src, err, out)
	}
	return out
}

func TestCore(t *testing.T) {
	cases := []struct{ src, want string }{
		{"1 2 + .", "3 "},
		{"7 2 - . 6 7 * . 17 5 / . 17 5 MOD .", "5 42 3 2 "},
		{"17 5 /MOD . .", "3 2 "},
		{"1 2 SWAP . .", "1 2 "},
		{"1 2 OVER . . .", "1 2 1 "},
		{"1 2 3 ROT . . .", "1 3 2 "},
		{"5 DUP * .", "25 "},
		{"1 2 3 DEPTH .", "3 "},
		{"-5 ABS . 3 NEGATE . 3 9 MIN . 3 9 MAX .", "5 -3 3 9 "},
		{"1.5 2 * .", "3 "},
		{"0.1 >F32 .", "0.1 "},
		{"7 2 / . 7.0 2 / .", "3 3.5 "},
		{"2147483647 1 + >I32 .", "-2147483648 "},
		{"6 3 AND . 6 3 OR . 6 3 XOR . 0 INVERT .", "2 7 5 -1 "},
		{": SQ DUP * ; 9 SQ .", "81 "},
		{": SIGN DUP 0< IF DROP -1 ELSE 0= IF 0 ELSE 1 THEN THEN ; -4 SIGN . 0 SIGN . 4 SIGN .", "-1 0 1 "},
		{": T 0 10 0 DO I + LOOP ; T .", "45 "},
		{": T 0 10 0 DO I + 2 +LOOP ; T .", "20 "},
		{": T 3 0 DO 2 0 DO J 10 * I + . LOOP LOOP ; T", "0 1 10 11 20 21 "},
		{": T 0 BEGIN 1+ DUP 5 = UNTIL ; T .", "5 "},
		{": T 0 BEGIN DUP 3 < WHILE 1+ REPEAT ; T .", "3 "},
		{"VARIABLE X 42 X ! X @ . 8 X +! X @ .", "42 50 "},
		{"10 CONSTANT TEN TEN TEN + .", "20 "},
		{"HERE 1 , 2 , DUP @ SWAP 1+ @ + .", "3 "},
		{`." hello" CR`, "hello\n"},
		{`: HI ." hi " ; HI HI`, "hi hi "},
		{`S" ab" S" cd" + PRINTLN`, "abcd\n"},
		{`S" a" S" b" < .`, "true "},
		{": FIB DUP 2 < IF EXIT THEN DUP 1- FIB SWAP 2 - FIB + ; 20 FIB .", "6765 "},
		{": F DUP 2 < IF EXIT THEN DUP 1- RECURSE SWAP 2 - RECURSE + ; 10 F .", "55 "},
		{"( comment ) 1 . \\ rest ignored\n2 .", "1 2 "},
		{"1 2 >R . R> .", "1 2 "},
		{"HEX FF DECIMAL .", "255 "},
		{": T 3 FRAME 7 0 L! 8 1 L! 0 L@ 1 L@ + 2 L@ . . UNFRAME ; T", "0 15 "},
		{": T 0 LABEL 1 1+ DUP 4 < JMP3 1 2 2 LABEL 2 ; T .", "4 "},
		{": T GOTO 5 99 . LABEL 5 1 . ; T", "1 "},
	}
	for _, c := range cases {
		if got := run(t, c.src); got != c.want {
			t.Errorf("%q: got %q want %q", c.src, got, c.want)
		}
	}
}

func TestErrors(t *testing.T) {
	cases := []struct{ src, want string }{
		{"DROP", "stack underflow"},
		{"FOO", "unknown word: FOO"},
		{"1 0 /", "division by zero"},
		{": T IF ;", "unbalanced control structure"},
		{"THEN", "THEN without IF"},
		{": T GOTO 9 ;", "undefined label: 9"},
		{"5000 @", "bad address"},
		{"TRUE 1 +", "type error in arithmetic"},
		{": R RECURSE ; R", "return stack overflow"},
	}
	for _, c := range cases {
		_, err := Run(c.src)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: got %v want %q", c.src, err, c.want)
		}
	}
}

func TestStateAcrossLinesAndRecovery(t *testing.T) {
	var b strings.Builder
	vm := New(&b)
	for _, line := range []string{": SQ", "DUP *", ";", "1 2 FOO", "3 SQ ."} {
		_ = vm.Interpret(line)
	}
	if b.String() != "9 " {
		t.Fatalf("got %q", b.String())
	}
	// A failed definition is rolled back.
	if err := vm.Interpret(": BAD 1 IF ;"); err == nil {
		t.Fatal("expected error")
	}
	if vm.find("BAD") >= 0 {
		t.Fatal("BAD left in dictionary")
	}
	if vm.Depth() != 0 {
		t.Fatalf("stack not cleared: %d", vm.Depth())
	}
}
