package structmapper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// StringMapToStruct is the inverse of StructToStringMap.
// If strict is false, some attempts are made to coerce non-json inputs:
// - bools are checked for natural values (1 or begins with 't')
// - string slices are invented by splitting the string on commas
//
// Note: bool conversion is non-strict in all cases.
func StringMapToStruct(inputMap map[string]string, str interface{}, strict bool) error {
	if reflect.TypeOf(str).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(str)).Kind() != reflect.Struct {
		return fmt.Errorf("second param must be a pointer to a struct")
	}
	if inputMap == nil {
		return fmt.Errorf("first param must not be nil")
	}

	// outputMap is an intermediate form of the data with each value being a string
	// or a json decoded value
	outputMap := make(map[string]interface{}, len(inputMap))

	var walkValue func(reflect.Value) error

	// Iterate over the given struct and collect values into the map.
	// Anonymous fields cause a recursive call to walkValue().
	walkValue = func(sv reflect.Value) error {
		st := sv.Type()

		// Iterate over the struct looking for matches in the string map.
		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			ft := st.Field(i)

			if ft.PkgPath != "" { // unexported
				continue
			}

			if field.Anonymous {
				if err := walkValue(sv.Field(i)); err != nil {
					return err
				}

				continue
			}

			tag, _, omit := getJSONTagFromField(field)
			if omit {
				continue
			}
			if value, exists := inputMap[tag]; exists {
				if err := convert(outputMap, field, sv.Field(i), tag, value, strict); err != nil {
					return err
				}
			}
		}

		return nil
	}
	if err := walkValue(reflect.ValueOf(str).Elem()); err != nil {
		return err
	}

	if err := encodeMaps(outputMap, str); err != nil {
		return err
	}

	return nil
}

func encodeMaps(outputMap map[string]interface{}, str interface{}) error {
	// we now json encode outputMap to make it a form which looks more like str
	buf, err := json.Marshal(outputMap)
	if err != nil {
		return fmt.Errorf("json encoding: %w", err)
	}
	// json decode fully into str

	if err := json.Unmarshal(buf, &str); err != nil {
		return fmt.Errorf("json decoding: %w", err)
	}

	return nil
}

func convert(
	outputMap map[string]interface{},
	field reflect.StructField,
	fieldValue reflect.Value,
	tag, value string,
	strict bool,
) error {
	if field.Type.Kind() == reflect.String {
		outputMap[tag] = value

		return nil
	}

	if _, isTimeVariable := isTime(fieldValue); isTimeVariable {
		if value != "" {
			tval, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("parsing time value: %w", err)
			}
			outputMap[tag] = tval
		}

		return nil
	}

	if field.Type.Kind() == reflect.Bool {
		outputMap[tag] = StringToBool(value)

		return nil
	}

	if value == "" {
		return nil
	}

	var decodedValue interface{}
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil && strict {
		return fmt.Errorf("json decoding: %w", err)
	} else if err == nil { // no error, so all good to return
		outputMap[tag] = decodedValue

		return nil
	}

	// not in strict mode, so try to find a value worth returning
	if field.Type == reflect.SliceOf(reflect.TypeOf("")) {
		outputMap[tag] = stringToStringSlice(value)
	}

	return nil
}
