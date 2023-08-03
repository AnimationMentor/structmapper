package structmapper

import (
	"reflect"
	"strconv"
	"strings"
)

// GetJSONTag takes a field name and json struct tag value and returns the effective name, omitempty setting
// and an indicator if the field value should be ommitted (if the tag value is "-").
//
//	name, omitEmpty, omit := GetJSONTag(fieldName, fieldTag)
//
// Examples:
//
//	GetJSONTag("Cake", "tuna,omitempty") -> "tuna",true
//	GetJSONTag("Cake", "") -> "Cake",false
func GetJSONTag(fieldName, tag string) (string, bool, bool) {
	if tag == "-" {
		return "", false, true
	}
	if tag == "-," { // per the docs, this is how you get a json field name of "-"
		return "-", false, false
	}
	f := strings.Split(tag, ",")
	omitEmpty := len(f) > 1 && f[1] == "omitempty"
	if len(f) < 1 || f[0] == "" {
		return fieldName, omitEmpty, false
	}
	return f[0], omitEmpty, false
}

// getJSONTagFromField gathers tag and name from a field and passes it to getJSONTag.
// If an "sm" tag is present, it is used in favor of the "json" field tag.
func getJSONTagFromField(f reflect.StructField) (string, bool, bool) {
	t := f.Tag.Get("sm")
	if t == "" {
		t = f.Tag.Get("json")
	}
	return GetJSONTag(f.Name, t)
}

// StringToBool returns true if the value is true looking
// if numeric and greater than zero
// if string and starts with t or y
func StringToBool(s string) bool {
	// optimise for normal json case
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	// now try to guess from the input
	if s == "" {
		return false
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i > 0
	}
	s = strings.ToLower(s)
	return len(s) > 0 && s[0] == 't' || s[0] == 'y'
}

// stringToStringSlice hopes the input is a comma separated string and returns it as a slice.
func stringToStringSlice(s string) []string {
	s = strings.Replace(s, ", ", ",", -1)
	return strings.Split(s, ",")
}
