package structmapper

import (
	"encoding/json"
	"reflect"

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

	sv := reflect.ValueOf(s).Elem()
	st := sv.Type()

	// Iterate over the struct looking for matches in the string map.
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		t, _ := getJSONTag(f.Name, f.Tag.Get("json"))
		if t == "" {
			continue
		}
		if value, exists := m[t]; exists {
			if f.Type.Kind() == reflect.String {
				m2[t] = value
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

	// we now json encode m2 to make it a form which looks more like s
	buf, err := json.Marshal(m2)
	if err != nil {
		return err
	}

	// json decode fully into s
	return json.Unmarshal([]byte(buf), &s)
}
