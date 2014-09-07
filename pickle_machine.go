/*
This package reads pickled data written from python using the "pickle" module.
Protocols 0,1,2 are implemented. These are the versions written by the Python
2.x series. Python 3 defines newer protocol versions, but can write the older
protocol versions so they are readable by this package.

TLDR:

Read a pickled string or unicode object
	var somePickledData io.Reader
	mystring, err := stalecucumber.String(stalecucumber.Unpickle(somePickledData))

Read a pickled integer
	var somePickledData io.Reader
	myint64, err := stalecucumber.Int(stalecucumber.Unpickle(somePickledData))

Read a python dictionary into a structure
	var somePickledData io.Reader
	mystruct := struct{
		Apple int
		Banana uint
		Cat string
		Dog float32}{}

	err := stalecucumber.UnpackInto(&mystruct).From(stalecucumber.From(somePickledData))

Reading Data

The stalecucumber.Unpickle function takes a reader and attempts to read
a complete pickle program from it. This is normally the output of the function
like "pickle.dump" from Python.

The returned type is interface{} because unpickling can generate any type. Use
a helper function to convert to another type without an additional type check.

An error is returned if the underlying reader fails, the program
is invalid, or unsupported opcodes are encountered. See below for the details
of unsupported opcodes.

If the version of python you are using supports protocol version 1 or 2,
you should always specify that protocol version. By default the "pickle"
and "cPickle" modules in Python write using protocol 0. Protocol 0
requires much more space to represent the same values and is much
slower to parse.

Type Conversions

Types conversion Python types to Go types is performed as followed
	int -> int64
	string -> string
	unicode -> string
	float -> float64
	long -> big.Int from the "math/big" package
	lists -> []interface{}
	tuples -> []interface{}
	dict -> []interface{}

The following values are converted from Python to the Go types
	True & False -> bool
	None -> stalecucumber.PickleNone

Helper Functions

The following helper functions were inspired by the github.com/garyburd/redigo
package. Each function takes the result of Unpickle as its arguments. If unpickle
fails it does nothing and returns that error. Otherwise it attempts to
convert to the appropriate type. If type conversion fails it returns an error

	String - string from Python string or unicode
	Int - int64 from Python int or long
	Bool - bool from Python True or False
	Big - *big.Int from Python long
	ListOrTuple - []interface{} from Python Tuple or List
	Float - float64 from Python float
	Dict - map[interface{}]interface{} from Python dictionary
	DictString -
		map[string]interface{} from Python dictionary.
		Keys must all be of type unicode or string.

Unpacking into structures

If the pickled object is a python dictionary that has only unicode and string
objects for keys, that object can be unpickled into a struct in Go by using
the "UnpackInto" function. The "From" receiver on the return value accepts
the result of "Unpickle" as its actual parameters.

The keys of the python dictionary are assigned to fields in a structure. If
the first character of the key is not uppercase, it is uppercased. If a field
matching that name is found, the value in the python dictionary is unpacked
into the value of the field within the structure.

A nested python dictionary is unpickled into nested structures in Go. If a
field is of type map[interface{}]interface{} is of course unpacked into that
as well.

By default UnpackInto skips any missing fields and fails if a field's
type is not compatible with the object's type.

This behavior can be changed by setting "AllowMissingFields" and
"AllowMismatchedFields" on the return value of UnpackInto before calling
From.

Recursive objects

You can pickle recursive objects like so

	a = {}
	a["self"] = a
	pickle.dumps(a)

Python's pickler is intelligent enough not to emit an infinite data structure
when a recursive object is pickled.

I reccomend against pickling recursive objects in the first place, but this
library handles unpickling them without a problem. The result of unpickling
the above is map[interface{}]interface{} with a key "a" that contains
a reference to itself.

Attempting to unpack the result of the above python code into a structure
with UnpackInto would either fail or recurse forever.

Unsupported Opcodes

The pickle format is incredibly flexible and as a result has some
features that are impractical or unimportant when implementing a reader in
another language.

Each set of opcodes is listed below by protocol version with the impact.

Protocol 0

	GLOBAL

This opcode is equivalent to calling "import foo; foo.bar" in python. It is
generated whenever an object instance, class definition, or method definition
is serialized. As long as the pickled data does not contain an instance
of a python class or a reference to a python callable this opcode is not
emitted by the "pickle" module.

A few examples of what will definitely cause this opcode to be emitted

	pickle.dumps(range) #Pickling the range function
	pickle.dumps(Exception()) #Pickling an instance of a python class

This opcodes will be partially supported in a future revision to this package
that allows the unpickling of instances of Python classes.

	REDUCE
	BUILD
	INST

These opcodes are used in recreating pickled python objects. That is currently
not supported by this package.

These opcodes will be supported in a future revision to this package
that allows the unpickling of instances of Python classes.

	PERSID

This opcode is used to reference concrete definitions of objects between
a pickler and an unpickler by an ID#. The pickle protocol doesn't define
what a persistent ID means.

This opcode is unlikely to ever be supported by this package.

Protocol 1

	OBJ

This opcodes is used in recreating pickled python objects. That is currently
not supported by this package.

This opcode will supported in a future revision to this package
that allows the unpickling of instances of Python classes.


	BINPERSID

This opcode is equivalent to PERSID in protocol 0 and won't be supported
for the same reason.

Protocol 2

	NEWOBJ

This opcodes is used in recreating pickled python objects. That is currently
not supported by this package.

This opcode will supported in a future revision to this package
that allows the unpickling of instances of Python classes.

	EXT1
	EXT2
	EXT4

These opcodes allow using a registry
of popular objects that are pickled by name, typically classes.
It is envisioned that through a global negotiation and
registration process, third parties can set up a mapping between
ints and object names.

These opcodes are unlikely to ever be supported by this package.

*/
package stalecucumber

import "errors"
import "io"
import "bytes"
import "encoding/binary"
import "fmt"

var ErrOpcodeStopped = errors.New("STOP opcode found")
var ErrStackTooSmall = errors.New("Stack is too small to perform requested operation")
var ErrInputTruncated = errors.New("Input to the pickle machine was truncated")
var ErrOpcodeNotImplemented = errors.New("Input encountered opcode that is not implemented")
var ErrNoResult = errors.New("Input did not place a value onto the stack")
var ErrMarkNotFound = errors.New("Mark could not be found on the stack")

/*
Unpickle a value from a reader. This function returns an error if
the reader fails, the pickled data is invalid, or if the pickled data contains
an unsupported opcode.

See unsupported opcodes in the documentation of this package for more
information.

This function is generally not used directly, but with one of the helpers
such as string.
*/
func Unpickle(reader io.Reader) (interface{}, error) {
	var pm PickleMachine
	pm.buf = &bytes.Buffer{}
	pm.Reader = reader

	err := (&pm).execute()
	if err != nil {
		return nil, err
	}

	if len(pm.Stack) == 0 {
		return nil, ErrNoResult
	}

	return pm.Stack[0], nil

}

var jumpList = buildEmptyJumpList()

func init() {
	populateJumpList(&jumpList)
}

/*
This type is returned whenever Unpickle encounters an error in pickled data.
*/
type PickleMachineError struct {
	Err       error
	StackSize int
	MemoSize  int
	Opcode    uint8
}

/*
This struct is current exposed but not useful. It is likely to be hidden
in the near future.
*/
type PickleMachine struct {
	Stack  []interface{}
	Memo   []interface{}
	Reader io.Reader

	currentOpcode uint8
	buf           *bytes.Buffer
}

func (pme PickleMachineError) Error() string {
	return fmt.Sprintf("Pickle Machine failed on opcode:0x%x. Stack size:%d. Memo size:%d. Cause:%v", pme.Opcode, pme.StackSize, pme.MemoSize, pme.Err)
}

func (pm *PickleMachine) error(src error) error {
	return PickleMachineError{
		StackSize: len(pm.Stack),
		MemoSize:  len(pm.Memo),
		Err:       src,
		Opcode:    pm.currentOpcode,
	}
}

func (pm *PickleMachine) execute() error {
	for {
		var opcode uint8
		err := binary.Read(pm.Reader, binary.BigEndian, &opcode)
		if err != nil {
			return err
		}

		pm.currentOpcode = opcode
		err = jumpList[int(opcode)](pm)
		if err == ErrOpcodeStopped {
			return nil
		} else if err == ErrOpcodeNotImplemented {
			return pm.error(ErrOpcodeNotImplemented)
		} else if err != nil {
			return pm.error(err)
		}
	}
}

func (pm *PickleMachine) storeMemo(index int64, v interface{}) error {
	if index < 0 {
		return fmt.Errorf("Requested to write to invalid memo index:%v", index)
	}

	if int64(len(pm.Memo)) <= index {
		replacement := make([]interface{}, index+1)
		copy(replacement, pm.Memo)
		pm.Memo = replacement
	}

	pm.Memo[index] = v

	return nil
}

func (pm *PickleMachine) readFromMemo(index int64) (interface{}, error) {
	if index < 0 || index >= int64(len(pm.Memo)) {
		return nil, fmt.Errorf("Requested to read from invalid memo index %d", index)
	}

	return pm.Memo[index], nil
}

func (pm *PickleMachine) push(v interface{}) {
	pm.Stack = append(pm.Stack, v)
}

func (pm *PickleMachine) pop() (interface{}, error) {
	if len(pm.Stack) == 0 {
		return nil, ErrStackTooSmall
	}

	lastIndex := len(pm.Stack) - 1
	top := pm.Stack[lastIndex]

	pm.Stack = pm.Stack[:lastIndex]
	return top, nil
}

func (pm *PickleMachine) readFromStack(offset int) (interface{}, error) {
	return pm.readFromStackAt(len(pm.Stack) - 1 - offset)
}

func (pm *PickleMachine) readFromStackAt(position int) (interface{}, error) {

	if position < 0 {
		return nil, fmt.Errorf("Request to read from invalid stack position %d", position)
	}

	return pm.Stack[position], nil

}

func (pm *PickleMachine) readIntFromStack(offset int) (int64, error) {
	v, err := pm.readFromStack(offset)
	if err != nil {
		return 0, err
	}

	vi, ok := v.(int64)
	if !ok {
		return 0, fmt.Errorf("Type %T was requested from stack but found %v(%T)", vi, v, v)
	}

	return vi, nil
}

func (pm *PickleMachine) popAfterIndex(index int) error {
	if len(pm.Stack)-1 < index {
		return ErrStackTooSmall
	}

	pm.Stack = pm.Stack[0:index]
	return nil
}

func (pm *PickleMachine) putMemo(index int, v interface{}) {
	for len(pm.Memo) <= index {
		pm.Memo = append(pm.Memo, nil)
	}

	pm.Memo[index] = v
}

func (pm *PickleMachine) findMark() (int, error) {
	for i := len(pm.Stack) - 1; i != -1; i-- {
		if _, ok := pm.Stack[i].(PickleMark); ok {
			return i, nil
		}
	}
	return -1, ErrMarkNotFound
}

func (pm *PickleMachine) readFixedLengthRaw(l int64) ([]byte, error) {

	pm.buf.Reset()
	_, err := io.CopyN(pm.buf, pm.Reader, l)
	if err != nil {
		return nil, err
	}

	return pm.buf.Bytes(), nil
}

func (pm *PickleMachine) readFixedLengthString(l int64) (string, error) {

	//Avoid getting "<nil>"
	if l == 0 {
		return "", nil
	}

	pm.buf.Reset()
	_, err := io.CopyN(pm.buf, pm.Reader, l)
	if err != nil {
		return "", err
	}
	return pm.buf.String(), nil
}

func (pm *PickleMachine) readString() (string, error) {
	pm.buf.Reset()
	for {
		var v [1]byte
		n, err := pm.Reader.Read(v[:])
		if n != 1 {
			return "", ErrInputTruncated
		}
		if err != nil {
			return "", err
		}

		if v[0] == '\n' {
			break
		}
		pm.buf.WriteByte(v[0])
	}

	//Avoid getting "<nil>"
	if pm.buf.Len() == 0 {
		return "", nil
	}
	return pm.buf.String(), nil
}

func (pm *PickleMachine) readBinaryInto(dst interface{}, bigEndian bool) error {
	var bo binary.ByteOrder
	if bigEndian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}
	return binary.Read(pm.Reader, bo, dst)
}
