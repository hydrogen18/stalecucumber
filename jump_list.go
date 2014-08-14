package stalecucumber

import "errors"

type OpcodeFunc func(*PickleMachine) error

var ErrOpcodeInvalid = errors.New("Opcode is invalid")

func (pm *PickleMachine) Opcode_Invalid() error {
	return ErrOpcodeInvalid
}

type OpcodeJumpList [256]OpcodeFunc

func buildEmptyJumpList() OpcodeJumpList {
	jl := OpcodeJumpList{}

	for i := range jl {
		jl[i] = (*PickleMachine).Opcode_Invalid
	}
	return jl
}
