import versions
from io import open

pkg_stmt = u'package stalecucumber\n\n'

def make_name(n):
	return "opcode_%s" % n

with open('protocol_0.go','w') as fout:
	fout.write(pkg_stmt) 
	
	for opcode in versions.v0:
		func_name = make_name(opcode.name)
		fout.write(u'/**\n')
		fout.write(u"Opcode: %s\n%s**\n" % (opcode.name,opcode.doc,))
		fout.write(u"Stack before: %s\n" % opcode.stack_before)
		fout.write(u"Stack after: %s\n" % opcode.stack_after)
		fout.write(u'**/\n')
		fout.write(u"func (pm *PickleMachine) %s () error {\n" % func_name)
		fout.write(u'return ErrOpcodeNotImplemented\n')
		fout.write(u'}\n\n')

with open('populate_jump_list.go','w') as fout:
	fout.write(pkg_stmt)

	fout.write(u'func populateJumpList(jl *OpcodeJumpList) {\n')
	for opcode in versions.v0:
		fout.write(u"jl[0x%X] = (*PickleMachine).%s\n" % (ord(opcode.code),make_name(opcode.name)))
	fout.write(u"}\n\n")

		
