package stalecucumber

func populateJumpList(jl *OpcodeJumpList) {
jl[0x49] = (*PickleMachine).opcode_INT
jl[0x4C] = (*PickleMachine).opcode_LONG
jl[0x53] = (*PickleMachine).opcode_STRING
jl[0x4E] = (*PickleMachine).opcode_NONE
jl[0x56] = (*PickleMachine).opcode_UNICODE
jl[0x46] = (*PickleMachine).opcode_FLOAT
jl[0x61] = (*PickleMachine).opcode_APPEND
jl[0x6C] = (*PickleMachine).opcode_LIST
jl[0x74] = (*PickleMachine).opcode_TUPLE
jl[0x64] = (*PickleMachine).opcode_DICT
jl[0x73] = (*PickleMachine).opcode_SETITEM
jl[0x30] = (*PickleMachine).opcode_POP
jl[0x32] = (*PickleMachine).opcode_DUP
jl[0x28] = (*PickleMachine).opcode_MARK
jl[0x67] = (*PickleMachine).opcode_GET
jl[0x70] = (*PickleMachine).opcode_PUT
jl[0x63] = (*PickleMachine).opcode_GLOBAL
jl[0x52] = (*PickleMachine).opcode_REDUCE
jl[0x62] = (*PickleMachine).opcode_BUILD
jl[0x69] = (*PickleMachine).opcode_INST
jl[0x2E] = (*PickleMachine).opcode_STOP
jl[0x50] = (*PickleMachine).opcode_PERSID
}

