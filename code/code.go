package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

// Definition is a handy debugging view of the opcode and
// the number of operands an opcode will take
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// increment the instruction length to total number of
	// bytes all the operands will occupy for that opcode
	// we need the total instructionLen to construct the requisite sized
	// byte array for the instruction
	instructionLen := 1
	for _, operandWidth := range def.OperandWidths {
		instructionLen += operandWidth
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, operand := range operands {
		operandWidth := def.OperandWidths[i]
		switch operandWidth {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}

		offset += operandWidth
	}

	return instruction
}
