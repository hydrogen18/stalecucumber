package stalecucumber

import "errors"
import "io"
import "bytes"
import "encoding/binary"
import "fmt"
import "math/big"

var ErrOpcodeStopped = errors.New("STOP opcode found")
var ErrStackTooSmall = errors.New("Stack is too small to perform requested operation")
var ErrInputTruncated = errors.New("Input to the pickle machine was truncated")
var ErrOpcodeNotImplemented = errors.New("Input encountered opcode that is not implemented")
var ErrNoResult = errors.New("Input did not place a value onto the stack")
var ErrMarkNotFound = errors.New("Mark could not be found on the stack")

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
	case **big.Int:
		if vbi, ok := v.(*big.Int); ok {
			*dest = vbi
			return nil
		}
	case *float64:
		if vf, ok := v.(float64); ok {
			*dest = vf
			return nil
		}
	case *[]interface{}:
		if vsi, ok := v.([]interface{}); ok {
			*dest = vsi
			return nil
		}
	case *map[interface{}]interface{}:
		if vsm, ok := v.(map[interface{}]interface{}); ok {
			*dest = vsm
			return nil
		}
	}

	//TODO handle the case of v.(PickleMark{})
	//& dest is a pointer type. Just set equal to nil
	//and return

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
		} else if err == ErrOpcodeNotImplemented {
			return fmt.Errorf("Opcode 0x%X not impleneted", opcode)
		} else if err != nil {
			return err
		}
	}
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
