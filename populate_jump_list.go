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
jl[0x4A] = (*PickleMachine).opcode_BININT
jl[0x4B] = (*PickleMachine).opcode_BININT1
jl[0x4D] = (*PickleMachine).opcode_BININT2
jl[0x54] = (*PickleMachine).opcode_BINSTRING
jl[0x55] = (*PickleMachine).opcode_SHORT_BINSTRING
jl[0x58] = (*PickleMachine).opcode_BINUNICODE
jl[0x47] = (*PickleMachine).opcode_BINFLOAT
jl[0x5D] = (*PickleMachine).opcode_EMPTY_LIST
jl[0x65] = (*PickleMachine).opcode_APPENDS
jl[0x29] = (*PickleMachine).opcode_EMPTY_TUPLE
jl[0x7D] = (*PickleMachine).opcode_EMPTY_DICT
jl[0x75] = (*PickleMachine).opcode_SETITEMS
jl[0x31] = (*PickleMachine).opcode_POP_MARK
jl[0x68] = (*PickleMachine).opcode_BINGET
jl[0x6A] = (*PickleMachine).opcode_LONG_BINGET
jl[0x71] = (*PickleMachine).opcode_BINPUT
jl[0x72] = (*PickleMachine).opcode_LONG_BINPUT
jl[0x6F] = (*PickleMachine).opcode_OBJ
jl[0x51] = (*PickleMachine).opcode_BINPERSID
jl[0x8A] = (*PickleMachine).opcode_LONG1
jl[0x8B] = (*PickleMachine).opcode_LONG4
jl[0x88] = (*PickleMachine).opcode_NEWTRUE
jl[0x89] = (*PickleMachine).opcode_NEWFALSE
jl[0x85] = (*PickleMachine).opcode_TUPLE1
jl[0x86] = (*PickleMachine).opcode_TUPLE2
jl[0x87] = (*PickleMachine).opcode_TUPLE3
jl[0x82] = (*PickleMachine).opcode_EXT1
jl[0x83] = (*PickleMachine).opcode_EXT2
jl[0x84] = (*PickleMachine).opcode_EXT4
jl[0x81] = (*PickleMachine).opcode_NEWOBJ
jl[0x80] = (*PickleMachine).opcode_PROTO
}

