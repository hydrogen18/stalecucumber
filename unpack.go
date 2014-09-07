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

/*
This type is returned when a call to From() fails.
Setting "AllowMissingFields" and "AllowMismatchedFields"
on the result of "UnpackInto" controls if this error is
returned or not.
*/
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
	dest                  interface{}
	AllowMissingFields    bool
	AllowMismatchedFields bool
}

func UnpackInto(dest interface{}) unpacker {
	return unpacker{dest: dest,
		AllowMissingFields:    true,
		AllowMismatchedFields: false}
}

func (u unpacker) From(srcI interface{}, err error) error {
	if err != nil {
		return err
	}

	v := reflect.ValueOf(u.dest)
	if v.Kind() != reflect.Ptr {
		return UnpackingError{Source: srcI,
			Destination: u.dest,
			Err:         ErrNotPointer}
	}

	if v.IsNil() {
		return UnpackingError{Source: srcI,
			Destination: u.dest,
			Err:         ErrNilPointer}

	}

	vIndirect := reflect.Indirect(v)
	if srcList, ok := srcI.([]interface{}); ok {
		if vIndirect.Kind() != reflect.Slice {
			return UnpackingError{Source: srcI,
				Destination: u.dest,
				Err:         fmt.Errorf("Cannot unpack slice into destination")}
		}

		replacement := reflect.MakeSlice(vIndirect.Type(),
			len(srcList), len(srcList))

		for i, srcV := range srcList {
			dstV := replacement.Index(i)

			//Check if the slice element type is
			//a pointer.
			if dstV.Kind() != reflect.Ptr {
				//If not a pointer, then indirect the
				//value here
				dstV = dstV.Addr()
			} else {
				//If it is a pointer, check for being nil.
				//Allocate a structure if it is nil
				if dstV.IsNil() {
					dstV.Set(reflect.New(dstV.Type().Elem()))
				}
			}
			err := unpacker{dest: dstV.Interface(),
				AllowMissingFields:    u.AllowMissingFields,
				AllowMismatchedFields: u.AllowMismatchedFields}.
				From(srcV, nil)
			if err != nil {
				return err
			}
		}
		vIndirect.Set(replacement)
		return nil
	}

	v = vIndirect

	var src map[string]interface{}
	src, err = DictString(srcI, err)
	if err != nil {
		return UnpackingError{Source: srcI,
			Destination: u.dest,
			Err:         fmt.Errorf("Cannot unpack source into struct")}
	}

	if v.Kind() != reflect.Struct {
		return UnpackingError{Source: src,
			Destination: u.dest,
			Err:         fmt.Errorf("Cannot unpack into %v", v.Kind().String())}

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
			if !u.AllowMismatchedFields {
				return UnpackingError{Source: src,
					Destination: u.dest,
					Err:         fmt.Errorf("Cannot find field for key %q", k)}
			}
			continue
		}
		err := u.assignTo(k, kv, fv)
		if err != nil && !u.AllowMismatchedFields {
			return err
		}
	}

	return nil
}

func (u unpacker) assignTo(fieldName string, v interface{}, dst reflect.Value) error {
	//If the destination is a pointer then
	//it cannot be assigned directly
	if dst.Kind() == reflect.Ptr {
		//Construct an instance of the type pointed at
		//if needed
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return u.assignTo(fieldName, v, reflect.Indirect(dst))
	}
	switch v := v.(type) {
	case int64:
		switch dst.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
			dst.SetInt(v)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64:
			if v >= 0 {
				dst.SetUint(uint64(v))
				return nil
			}
		}
	case string:
		switch dst.Kind() {
		case reflect.String:
			dst.SetString(v)
			return nil
		}
	case bool:
		switch dst.Kind() {
		case reflect.Bool:
			dst.SetBool(v)
			return nil
		}
	case float64:
		switch dst.Kind() {
		case reflect.Float32, reflect.Float64:
			dst.SetFloat(v)
			return nil
		}
	case *big.Int:
		dstBig, ok := dst.Addr().Interface().(*big.Int)
		if ok {
			(dstBig).Set(v)
			return nil
		}

		if vi, err := Int(v, nil); err == nil {
			return u.assignTo(fieldName, vi, dst)
		}

	case []interface{}:
		if dst.Kind() == reflect.Slice &&
			dst.Type().Elem().Kind() == reflect.Interface {
			dst.Set(reflect.ValueOf(v))
			return nil
		}
		return unpacker{dest: dst.Addr().Interface(),
			AllowMismatchedFields: u.AllowMismatchedFields,
			AllowMissingFields:    u.AllowMissingFields}.From(v, nil)

	case map[interface{}]interface{}:
		//Check to see if the field is exactly
		//of the type.
		if dst.Kind() == reflect.Map {
			dstT := dst.Type()
			if dstT.Key().Kind() == reflect.Interface &&
				dstT.Elem().Kind() == reflect.Interface {
				dst.Set(reflect.ValueOf(v))
				return nil
			}
		}

		//Try to assign this recursively
		return unpacker{dest: dst.Addr().Interface(),
			AllowMismatchedFields: u.AllowMismatchedFields,
			AllowMissingFields:    u.AllowMissingFields}.From(v, nil)
	}
	return UnpackingError{Source: v,
		Destination: dst.Interface(),
		Err:         fmt.Errorf("For field %q source type doesn't match destination field", fieldName)}
}
