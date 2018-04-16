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

	sv := reflect.ValueOf(s).Elem()
	st := sv.Type()

	m := make(map[string]string, st.NumField())

	for i := 0; i < st.NumField(); i++ {
		f := sv.Field(i)
		ft := st.Field(i)

		t, omitEmpty := getJSONTag(ft.Name, ft.Tag.Get("json"))
		if t == "" || omitEmpty && f.Len() == 0 {
			continue
		}

		if f.Kind() == reflect.String {
			m[t] = f.String()
		} else {
			buf, err := json.Marshal(f.Interface())
			if err != nil {
				return nil, err
			}
			m[t] = string(buf)
		}
	}
	return m, nil
}
