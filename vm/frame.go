package vm

import (
	"github.com/andy9775/monkey/code"
	"github.com/andy9775/monkey/object"
)

// Frame represents a function stack frame
type Frame struct {
	cl *object.Closure
	ip int // instruction pointer in this fram for this function

	// basePointer is the old value of the stack pointer (bottom of stack of current call frame)
	// also called  the frame pointer
	basePointer int
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
