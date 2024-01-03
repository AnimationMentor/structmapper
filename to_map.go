package structmapper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// StructToStringMap converts a struct to a map[string]string by creating entries in the map
// for each field in the struct. String type values are copied as is. Other types are JSON
// encoded.
// JSON struct tags are consulted for map key names and ommit behaviour.
func StructToStringMap(s interface{}) (map[string]string, error) {
	if reflect.TypeOf(s).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(s)).Kind() != reflect.Struct {
		return nil, fmt.Errorf("s must be a pointer to a struct")
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

			if ft.PkgPath != "" { // unexported
				continue
			}

			if ft.Anonymous {
				if err := walkValue(f); err != nil {
					return err
				}

				continue
			}

			t, omitEmpty, omit := getJSONTagFromField(ft)

			if omit {
				continue
			}

			if omitEmpty && f.Kind() == reflect.Ptr && f.IsNil() {
				continue
			}

			var stringValue string

			if timeString, isT := isTime(f); isT {
				stringValue = timeString
			} else if f.Kind() == reflect.String {
				stringValue = f.String()
			} else {
				buf, err := json.Marshal(f.Interface())
				if err != nil {
					return fmt.Errorf("json encoding: %w", err)
				}
				stringValue = string(buf)
			}

			if stringValue == "" && omitEmpty {
				continue
			}
			m[t] = stringValue
		}

		return nil
	}

	return m, walkValue(reflect.ValueOf(s).Elem())
}

// isTime returns true and the string value if this is a time.Time or *time.Time
// If it's a nil pointer or a zero value time, returns an empty string.
func isTime(f reflect.Value) (string, bool) {
	switch f.Kind() {
	case reflect.Struct:
		if t, ok := f.Interface().(time.Time); ok {
			if t.IsZero() {
				return "", true
			}

			return t.Format(time.RFC3339), true
		}
	case reflect.Ptr:
		if t, ok := f.Interface().(*time.Time); ok {
			if t == nil {
				return "", true
			}

			return t.Format(time.RFC3339), true
		}
	}

	return "", false
}
