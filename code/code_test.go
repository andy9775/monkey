package code_test

import (
	"testing"

	"github.com/andy9775/monkey/code"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		expected []byte
	}{
		{
			code.OpConstant, []int{65534},
			[]byte{
				byte(code.OpConstant), // expected opcode
				255, 254,              // big endian encoding of 65534
			},
		},
		{
			code.OpAdd, []int{},
			[]byte{
				byte(code.OpAdd), // expected opcode
				// no operands
			},
		},
		{
			code.OpGetLocal, []int{255}, []byte{byte(code.OpGetLocal), 255},
		},
		{code.OpClosure, []int{65534, 255}, []byte{byte(code.OpClosure), 255, 254, 255}},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte ad pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpAdd),
		code.Make(code.OpGetLocal, 1),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
		code.Make(code.OpClosure, 65535, 255),
	}

	expected := "0000 OpAdd\n0001 OpGetLocal 1\n0003 OpConstant 2\n0006 OpConstant 65535\n0009 OpClosure 65535 255\n"

	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        code.Opcode
		operands  []int
		bytesRead int
	}{
		{code.OpConstant, []int{65535}, 2},
		{code.OpGetLocal, []int{255}, 1},
		{code.OpClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)

		def, err := code.Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := code.ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}

	}
}
