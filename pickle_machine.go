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

func Unmarshal(reader io.Reader, dest interface{}) error {
	var pm PickleMachine
	pm.Reader = reader

	err := (&pm).execute()
	if err != nil {
		return err
	}

	if len(pm.Stack) == 0 {
		return ErrNoResult
	}

	v := pm.Stack[0]

	switch dest := dest.(type) {
	case *int64:
		if vi, ok := v.(int64); ok {
			*dest = vi
			return nil
		}
	case *bool:
		if vb, ok := v.(bool); ok {
			*dest = vb
			return nil
		}
	case *string:
		if vs, ok := v.(string); ok {
			*dest = vs
			return nil
		}
	}

	return fmt.Errorf("Cannot unmarshal object of type %T into destination of type %T", v, dest)

}

var jumpList = buildEmptyJumpList()

func init() {
	populateJumpList(&jumpList)
}

type PickleMachine struct {
	Stack  []interface{}
	Memo   []interface{}
	Reader io.Reader
}

func (pm *PickleMachine) execute() error {
	for {
		var opcode uint8
		err := binary.Read(pm.Reader, binary.BigEndian, &opcode)
		if err != nil {
			return err
		}

		err = jumpList[int(opcode)](pm)
		if err == ErrOpcodeStopped {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func (pm *PickleMachine) push(v interface{}) {
	pm.Stack = append(pm.Stack, v)
}

func (pm *PickleMachine) pop() error {
	return ErrStackTooSmall
}

func (pm *PickleMachine) putMemo(index int, v interface{}) {
	for len(pm.Memo) <= index {
		pm.Memo = append(pm.Memo, nil)
	}

	pm.Memo[index] = v
}

func (pm *PickleMachine) readString() (string, error) {
	var buf bytes.Buffer

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
		(&buf).WriteByte(v[0])
	}
	//Avoid getting "<nil>"
	if (&buf).Len() == 0 {
		return "", nil
	}
	return (&buf).String(), nil
}
