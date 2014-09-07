package stalecucumber

import "reflect"
import "fmt"
import "errors"
import "strings"
import "math/big"

type UnpackingError struct {
	Source      interface{}
	Destination interface{}
	Err         error
}

func (ue UnpackingError) Error() string {
	return fmt.Sprintf("Error unpacking %v(%T) into %v(%T):%v",
		ue.Source,
		ue.Source,
		ue.Destination,
		ue.Destination,
		ue.Err)
}

var ErrNilPointer = errors.New("Destination cannot be a nil pointer")
var ErrNotPointer = errors.New("Destination must be a pointer type")

type unpacker struct {
	dest interface{}
}

func UnpackInto(dest interface{}) unpacker {
	return unpacker{dest: dest}
}

func (u unpacker) From(srcI interface{}, err error) error {
	var src map[string]interface{}
	src, err = DictString(srcI, err)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(u.dest)

	if v.Kind() != reflect.Ptr {
		return UnpackingError{Source: src,
			Destination: u.dest,
			Err:         ErrNotPointer}
	}

	if v.IsNil() {
		return UnpackingError{Source: src,
			Destination: u.dest,
			Err:         ErrNilPointer}

	}

	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return UnpackingError{Source: src,
			Destination: u.dest,
			Err:         fmt.Errorf("Cannot unpickle into %v", v.Kind().String())}
	}

	for k, kv := range src {
		//Ignore zero length strings, a struct
		//cannot have such a field
		if len(k) == 0 {
			continue
		}
		//Capitalize the first character. Structs
		//do not export fields with a lower case
		//first character
		k = strings.ToUpper(k[0:1]) + k[1:]

		fv := v.FieldByName(k)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		if !assignTo(kv, fv) {
			return UnpackingError{Source: src,
				Destination: u.dest,
				Err:         fmt.Errorf("Cannot unpack into field %q with type %s a value %v of with type %T ", k, fv.Type(), kv, kv)}
		}
	}

	return nil
}

func assignTo(v interface{}, dst reflect.Value) bool {
	//If the destination is a pointer then
	//it cannot be assigned directly
	if dst.Kind() == reflect.Ptr {
		//Construct an instance of the type pointed at
		//if needed
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return assignTo(v, reflect.Indirect(dst))
	}
	switch v := v.(type) {
	case int64:
		switch dst.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
			dst.SetInt(v)
			return true
		}
	case string:
		switch dst.Kind() {
		case reflect.String:
			dst.SetString(v)
			return true
		}
	case bool:
		switch dst.Kind() {
		case reflect.Bool:
			dst.SetBool(v)
			return true
		}
	case float64:
		switch dst.Kind() {
		case reflect.Float32, reflect.Float64:
			dst.SetFloat(v)
			return true
		}
	case *big.Int:
		dstBig, ok := dst.Addr().Interface().(*big.Int)
		if ok {
			(dstBig).Set(v)
			return true
		}

	case []interface{}:
		if dst.Kind() == reflect.Slice &&
			dst.Type().Elem().Kind() == reflect.Interface {
			dst.Set(reflect.ValueOf(v))
			return true
		}
	}
	return false
}
