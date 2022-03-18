package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions will be the byte array representation of the byte code
// it can actually hold both a single instruction or multiple instructions in a series
// NOTE: we do not have a singular Instruction, Instructions caters to both
type Instructions []byte

// print a human-readable form of instructions from their byte array form
func (ins Instructions) String() string {
	var out bytes.Buffer

	offset := 0
	for offset < len(ins) {
		def, err := Lookup(ins[offset])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[offset+1:])

		fmt.Fprintf(&out, "%04d %s\n", offset, ins.fmtInstruction(def, operands))

		offset += 1 + read
	}

	return out.String()
}

// fmtInstruction will format a single instruction
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Opcode is a single byte operation code in the instruction referencing
// the operation which needs to be performed by that instruction
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

// Make will take opcode and operands of an instruction,
// create the bytecode binary representation of that instruction
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

// ReadOperands is reverse of Make, it will disassemble a binary instruction's operands
// into actual readable values
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	// for each of the operands, query the Definition for its width
	// then increment the total number of bytes read
	// the operands (2 byte long) are stored in big-endian so decode into integer of 16bits
	for i, operandWidth := range def.OperandWidths {
		switch operandWidth {
		case 2:
			operands[i] = int(ReadUInt16(ins[offset:]))
		}

		offset += operandWidth
	}

	return operands, offset
}

func ReadUInt16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
