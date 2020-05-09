package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions is a list of operations
type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Opcode represents an operation code supported by the virtual machine
type Opcode byte

const (
	// OpConstant allows us to store constants throughout the execution of the program
	OpConstant Opcode = iota

	// OpAdd tells us to pop the top two items off of the vm, add them and put the result back
	OpAdd
	OpSub
	OpMul
	OpDiv

	// OpTrue tells the Vm to put an object.Boolean onto the stack
	OpTrue
	OpFalse

	OpEqual
	OpNotEqual
	OpGreaterThan

	// OpMinus is the prefix `-`` operator
	OpMinus
	OpBang

	// OpJumpNotTruthy jumps if the previous instruction was false
	OpJumpNotTruthy
	// OpJump jumps no matter what
	OpJump

	// OpNull puts a Null value on the stack
	OpNull

	// OpGetGlobal is used to set identifiers to values
	OpGetGlobal
	OpSetGlobal

	// OpGetLocal gets a locall binding (scoped locally to a function for example)
	OpGetLocal
	OpSetLocal

	// OpGetBuiltin gets a builtin function from the builtin function scope
	OpGetBuiltin

	// OpArray tells the vm how to build an array
	OpArray

	// OpHash tells the compiler and vm to handle a dictionary (hash)
	OpHash

	// OpIndex allows us to get an element from an array or hash by index
	OpIndex

	// OpCall is used to execute a function call
	OpCall

	//OpReturnValue tells the vm to return from a function with a return value
	OpReturnValue
	// OpReturn tells the vm to return with nothing on the stack - go back to where you were
	OpReturn

	// OpPop tells us to simply pop an item off the top of the stack
	// each expression statement should be followed by it in order
	// to prevent filling up the stack
	OpPop
)

// Definition provides human readable debugging information for a specific OpCode
type Definition struct {
	Name          string // human readable opcode name
	OperandWidths []int  // each entry is the size (in bytes) for each operand
}

// definitions tracks the number of bytes an instruction operates on
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2} /* 2 bytes - unit16*/},
	OpAdd:      {"OpAdd", []int{} /*takes no operands*/},
	OpSub:      {"OpSub", []int{} /*takes no operands*/},
	OpMul:      {"OpMul", []int{} /*takes no operands*/},
	OpDiv:      {"OpDiv", []int{} /*takes no operands*/},

	OpTrue:        {"OpTrue", []int{} /*takes no operands*/},
	OpFalse:       {"OpFalse", []int{} /*takes no operands*/},
	OpEqual:       {"OpEqual", []int{} /*takes no operands*/},
	OpNotEqual:    {"OpNotEqual", []int{} /*takes no operands*/},
	OpGreaterThan: {"OpGreaterThan", []int{} /*takes no operands*/},

	OpMinus: {"OpMinus", []int{} /*takes no operands*/},
	OpBang:  {"OpBang", []int{} /*takes no operands*/},

	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2} /*single operand is the offset instruction*/},
	OpJump:          {"OpJump", []int{2} /*single operand is the offset instruction*/},

	// limited to 65536 global variables
	OpGetGlobal: {"OpGetGlobal", []int{2} /*single operand is the global reference location*/},
	OpSetGlobal: {"OpSetGlobal", []int{2} /*single operand is the global reference location*/},

	// limited to 256 local variables
	OpGetLocal: {"OpGetLocal", []int{1} /*single operand is the local reference location*/},
	OpSetLocal: {"OpSetLocal", []int{1} /*single operand is the local reference location*/},

	OpArray: {"OpArray", []int{2} /*single operand which specifies the number of elements in the array*/},
	OpHash:  {"OpHash", []int{2} /*single operand specifies the number of keys/value sitting on the stack*/},
	OpIndex: {"OpIndex", []int{} /*no operands; requires 2 items on stack: the data structure and index*/},

	OpCall: {
		"OpCall",
		[]int{1}, /*operand is number of arguments; previous item on stack is the identifier for the call*/
	},

	OpGetBuiltin: {"OpGetBuiltin", []int{1} /*operand is the index of the called builtin*/},

	OpReturnValue: {"OpReturnValue", []int{} /*no operands; returned value sits at the top of the stack*/},
	OpReturn:      {"OpReturn", []int{} /*no operands; no return value, just null*/},

	OpNull: {"OpNull", []int{}},

	OpPop: {"OpPop", []int{} /*takes no operands*/},
}

// Lookup returns the Definition for the specific op and an error if none found
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// --------------- make ---------------

// Make returns the bytecode for the specific Opcode and operands
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok { // opcode not found
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)

	instruction[0] = byte(op) // set the opcode

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

// ReadOperands reverses the work of Make
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }
