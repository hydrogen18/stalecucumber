package stalecucumber

import "fmt"
import "errors"

// A type to convert to a GLOBAL opcode to something meaningful in golang
type PythonResolver interface {
	Resolve(module string, name string, args interface{}) (interface{}, error)
}

var ErrUnresolvablePythonGlobal = errors.New("Unresolvable Python global value")
type UnparseablePythonGlobalError struct {
	Args interface{}
	Message string
}

func(this UnparseablePythonGlobalError) Error() string {
	return fmt.Sprintf("%s; arguments (%T): %v", this.Message, this.Args, this.Args)
}

type PythonBuiltinResolver struct {}

func (this PythonBuiltinResolver) Resolve(module string, name string, args interface{}) (interface{}, error) {
	if module != "__builtin__" {
		return nil, ErrUnresolvablePythonGlobal
	}

	if name == "set" {
		return this.handlePythonSet(args)
	}

	return nil, ErrUnresolvablePythonGlobal
}

func (this PythonBuiltinResolver) handlePythonSet(args interface{}) (interface{}, error){
	list, ok := args.([]interface{})
	if !ok {
		return nil, UnparseablePythonGlobalError{
			Args: args, 
			Message: fmt.Sprintf("Expected args to be of type %T", list),
		}
	}

	if len(list) != 1 {
		return nil, UnparseablePythonGlobalError{
			Args: args, 
			Message: "Expected args to be of length 1",
		}
	}

	tuple, ok := list[0].([]interface{})
	if !ok {
		return nil, UnparseablePythonGlobalError{
			Args: args, 
			Message: fmt.Sprintf("Expected first argument of args to be of type %T", tuple),
		}
	}

	// A map is the equivalent golang type for a python set
	set := make(map[interface{}]bool, len(tuple))
	for _, item := range tuple {
		set[item] = true
	}

	return set, nil
}
 