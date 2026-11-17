// Package forth is a small Forth with tagged cells and three valued logic.
//
// It is written for TinyGo on microcontrollers as well as for the host. All
// stacks, data memory and code space are fixed size arrays, there is no
// reflection, and output is formatted with strconv rather than fmt.
//
// # Cells
//
// A cell carries a tag: integer, float64, float32, ternary or string. The
// arithmetic words are polymorphic. Integer with integer stays integer, any
// float makes the result a float. >I32 and >F32 narrow a value the way CX
// i32 and f32 arithmetic does.
//
// # Ternary logic
//
// TRUE, FALSE and NULL push ternary cells. They are stored in balanced form,
// FALSE=-1 NULL=0 TRUE=1, which makes Kleene's strong logic arithmetic:
//
//	AND      min(a, b)
//	OR       max(a, b)
//	NOT      -a
//	IMPLIES  max(-a, b)          Kleene
//	LIMPLIES min(1, 1 - a + b)   Lukasiewicz, NULL LIMPLIES NULL is TRUE
//	XOR      -(a * b)
//	EQV      a * b
//
// AND, OR and XOR on two integers stay bitwise, as in standard Forth. When
// either operand is ternary, integers are read as Forth flags, zero false and
// anything else true. Comparisons always return a ternary cell, and any NULL
// operand makes the result NULL. Arithmetic on NULL also gives NULL, so an
// unknown reading stays unknown through a calculation. KNOWN? NULL? TRUE?
// and FALSE? test a value and never return NULL.
//
// # Conditionals on NULL
//
// IF runs its true branch only for TRUE. FALSE and NULL both take the ELSE
// branch, the same rule SQL uses for WHERE. UNTIL and WHILE follow the same
// rule, so a loop ends only on TRUE.
//
// TRI-IF branches three ways.
//
//	cond TRI-IF  true-part  ELSE  false-part  NULL-ELSE  null-part  THEN
//
// ELSE and NULL-ELSE are both optional. Without NULL-ELSE, NULL takes the
// ELSE branch exactly as IF does.
//
// # Words for translated code
//
// n FRAME opens n zeroed local slots, i L@ and i L! read and write them and
// UNFRAME closes the frame. LABEL n, GOTO n and JMP3 t f n give numbered
// labels inside one definition, which is what jump based IR needs.
package forth
