package stalecucumber

import "reflect"
import "fmt"
import "errors"
import "strings"
import "math/big"

const PICKLE_TAG = "pickle"

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
var ErrTargetTypeOverflow = errors.New("Value overflows target type")

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
	//Check if an error occurred
	if err != nil {
		return err
	}

	//Get the value of the destination
	v := reflect.ValueOf(u.dest)

	//The destination must always be a pointer
	if v.Kind() != reflect.Ptr {
		return UnpackingError{Source: srcI,
			Destination: u.dest,
			Err:         ErrNotPointer}
	}

	//The destination can never be nil
	if v.IsNil() {
		return UnpackingError{Source: srcI,
			Destination: u.dest,
			Err:         ErrNilPointer}

	}

	//Indirect the destination. This gets the actual
	//value pointed at
	vIndirect := reflect.Indirect(v)

	//Check the input against known types
	switch s := srcI.(type) {
	case int64:
		switch vIndirect.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64, reflect.Int32:
			if vIndirect.OverflowInt(s) {
				return UnpackingError{Source: srcI,
					Destination: u.dest,
					Err:         ErrTargetTypeOverflow}
			}
			vIndirect.SetInt(s)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if s < 0 || vIndirect.OverflowUint(uint64(s)) {
				return UnpackingError{Source: srcI,
					Destination: u.dest,
					Err:         ErrTargetTypeOverflow}
			}

			vIndirect.SetUint(uint64(s))
			return nil

		}
	case string:
		switch vIndirect.Kind() {
		case reflect.String:
			vIndirect.SetString(s)
			return nil
		}
	case bool:
		switch vIndirect.Kind() {
		case reflect.Bool:
			vIndirect.SetBool(s)
			return nil
		}
	case float64:
		switch vIndirect.Kind() {
		case reflect.Float32, reflect.Float64:
			vIndirect.SetFloat(s)
			return nil
		}
	case *big.Int:
		dstBig, ok := vIndirect.Addr().Interface().(*big.Int)
		if ok {
			(dstBig).Set(s)
			return nil
		}

		if vi, err := Int(srcI, nil); err == nil {
			return unpacker{dest: v.Interface(),
				AllowMismatchedFields: u.AllowMismatchedFields,
				AllowMissingFields:    u.AllowMissingFields}.From(vi, nil)
		}

	case []interface{}:
		//Check that the destination is a slice
		if vIndirect.Kind() != reflect.Slice {
			return UnpackingError{Source: s,
				Destination: u.dest,
				Err:         fmt.Errorf("Cannot unpack slice into destination")}
		}
		//Check for exact type match
		if vIndirect.Type().Elem().Kind() == reflect.Interface {
			vIndirect.Set(reflect.ValueOf(s))
			return nil
		}

		//Build the value using reflection
		replacement := reflect.MakeSlice(vIndirect.Type(),
			len(s), len(s))

		for i, srcV := range s {
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
			//Recurse to set the value
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

	case map[interface{}]interface{}:
		//Check to see if the field is exactly
		//of the type
		if vIndirect.Kind() == reflect.Map {
			dstT := vIndirect.Type()
			if dstT.Key().Kind() == reflect.Interface &&
				dstT.Elem().Kind() == reflect.Interface {
				vIndirect.Set(reflect.ValueOf(s))
				return nil
			}
		}

		var src map[string]interface{}
		src, err = DictString(srcI, err)
		if err != nil {
			return UnpackingError{Source: srcI,
				Destination: u.dest,
				Err:         fmt.Errorf("Cannot unpack source into struct")}
		}

		if vIndirect.Kind() != reflect.Struct {
			return UnpackingError{Source: src,
				Destination: u.dest,
				Err:         fmt.Errorf("Cannot unpack into %v", v.Kind().String())}

		}

		var fieldByTag map[string]int
		vIndirectType := reflect.TypeOf(vIndirect.Interface())
		numFields := vIndirectType.NumField()

		for i := 0; i != numFields; i++ {
			fv := vIndirectType.Field(i)
			tag := fv.Tag.Get(PICKLE_TAG)

			if len(tag) != 0 {
				if fieldByTag == nil {
					fieldByTag = make(map[string]int)
				}
				fieldByTag[tag] = i
			}
		}

		for k, kv := range src {
			var fv reflect.Value

			if fieldIndex, ok := fieldByTag[k]; ok {
				fv = vIndirect.Field(fieldIndex)
			} else {
				//Try the name verbatim. This catches
				//embedded fields as well
				fv = vIndirect.FieldByName(k)

				if !fv.IsValid() {
					//Capitalize the first character. Structs
					//do not export fields with a lower case
					//first character
					capk := strings.ToUpper(k[0:1]) + k[1:]

					fv = vIndirect.FieldByName(capk)
				}

			}

			if !fv.IsValid() || !fv.CanSet() {
				if !u.AllowMismatchedFields {
					return UnpackingError{Source: src,
						Destination: u.dest,
						Err:         fmt.Errorf("Cannot find field for key %q", k)}
				}
				continue
			}

			_, valueIsNone := kv.(PickleNone)

			if fv.Kind() != reflect.Ptr {
				if valueIsNone {
					panic("foo")
				}
				fv = fv.Addr()
			} else {
				if valueIsNone {
					fv.Set(reflect.Zero(fv.Type()))
					continue
				}

				if fv.IsNil() {
					fv.Set(reflect.New(fv.Type().Elem()))
				}
			}

			err := unpacker{dest: fv.Interface(),
				AllowMismatchedFields: u.AllowMismatchedFields,
				AllowMissingFields:    u.AllowMissingFields}.From(kv, nil)

			if err != nil && !u.AllowMismatchedFields {
				return err
			}
		}

		return nil
	}

	return UnpackingError{Source: srcI,
		Destination: u.dest,
		Err:         fmt.Errorf("Cannot unpack")}
}
