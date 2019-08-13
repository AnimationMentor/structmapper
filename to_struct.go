package structmapper

import (
	"encoding/json"
	"reflect"
	"time"
	// "unicode"

	"github.com/pkg/errors"
)

// StringMapToStruct is the inverse of StructToStringMap.
// If strict is false, some attempts are made to coerce non-json inputs:
// - bools are checked for natural values (1 or begins with 't')
// - string slices are invented by splitting the string on commas
//
// Note: bool conversion is non-strict in all cases.
//
func StringMapToStruct(m map[string]string, s interface{}, strict bool) error {

	if reflect.TypeOf(s).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(s)).Kind() != reflect.Struct {
		return errors.Errorf("s must be a pointer to a struct")
	}
	if m == nil {
		return errors.Errorf("m must not be nil")
	}

	// m2 is an intermediate form of the data with each value being a string
	// or a json decoded value
	m2 := make(map[string]interface{}, len(m))

	var walkValue func(reflect.Value) error

	// Iterate over the given struct and collect values into the map.
	// Anonymous fields cause a recursive call to walkValue().
	walkValue = func(sv reflect.Value) error {

		st := sv.Type()

		// Iterate over the struct looking for matches in the string map.
		for i := 0; i < st.NumField(); i++ {
			f := st.Field(i)
			ft := st.Field(i)

			if ft.PkgPath != "" { // unexported
				continue
			}

			if f.Anonymous {
				if err := walkValue(sv.Field(i)); err != nil {
					return err
				}
				continue
			}

			t, _, omit := getJSONTagFromField(f)
			if omit {
				continue
			}
			if value, exists := m[t]; exists {

				// a := f.Type.Kind() == reflect.Struct
				// b:= sv.Field(i).Interface()

				_, isT := isTime(sv.Field(i))

				if f.Type.Kind() == reflect.String {
					m2[t] = value
				} else if isT {
					if value != "" {
						tval, err := time.Parse(time.RFC3339, value)
						if err != nil {
							return errors.Wrap(err, "parsing time value")
						}
						m2[t] = tval
					}
				} else if f.Type.Kind() == reflect.Bool {
					m2[t] = stringToBool(value)
				} else if value != "" {
					var decodedValue interface{}
					err := json.Unmarshal([]byte(value), &decodedValue)
					if err == nil {
						m2[t] = decodedValue
					} else if !strict {
						if f.Type == reflect.SliceOf(reflect.TypeOf("")) {
							m2[t] = stringToStringSlice(value)
							err = nil
						}
					}

					if err != nil {
						return err
					}

				}
			}
		}
		return nil
	}
	if err := walkValue(reflect.ValueOf(s).Elem()); err != nil {
		return err
	}

	// we now json encode m2 to make it a form which looks more like s
	buf, err := json.Marshal(m2)
	if err != nil {
		return err
	}

	// json decode fully into s
	return json.Unmarshal([]byte(buf), &s)
}
