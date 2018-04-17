package structmapper

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// StructToStringMap converts a struct to a map[string]string by creating entries in the map
// for each field in the struct. String type values are copied as is. Other types are JSON
// encoded.
// JSON struct tags are consulted for map key names and ommit behaviour.
func StructToStringMap(s interface{}) (map[string]string, error) {

	if reflect.TypeOf(s).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(s)).Kind() != reflect.Struct {
		return nil, errors.Errorf("s must be a pointer to a struct")
	}
	m := make(map[string]string, 20)

	var walkValue func(reflect.Value) error

	// Iterate over the given struct and collect values into the map.
	// Anonymous fields cause a recursive call to walkValue().
	walkValue = func(sv reflect.Value) error {
		st := sv.Type()

		for i := 0; i < st.NumField(); i++ {
			f := sv.Field(i)
			ft := st.Field(i)

			if ft.Anonymous {
				if err := walkValue(f); err != nil {
					return err
				}
				continue
			}

			t, omitEmpty := getJSONTag(ft.Name, ft.Tag.Get("json"))
			if t == "" || omitEmpty && f.Len() == 0 {
				continue
			}

			if f.Kind() == reflect.String {
				m[t] = f.String()
			} else {
				buf, err := json.Marshal(f.Interface())
				if err != nil {
					return err
				}
				m[t] = string(buf)
			}
		}
		return nil
	}

	return m, walkValue(reflect.ValueOf(s).Elem())
}
