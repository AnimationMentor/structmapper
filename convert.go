package structmapper

import (
	"strconv"
	"strings"
)

// getJSONTag takes a field name and json struct tag value and returns the effective name and omitempty setting.
// If the json tag is a "-" then the empty string is returned.
//
// Examples:
//    getJSONTag("Cake", "tuna,omitempty") -> "tuna",true
//    getJSONTag("Cake", "") -> "Cake",false
//
func getJSONTag(fieldName, tag string) (string, bool) {
	if tag == "-" {
		return "", false
	}
	if tag == "-," { // per the docs, this is how you get a json field name of "-"
		return "-", false
	}
	f := strings.Split(tag, ",")
	omitEmpty := len(f) > 1 && f[1] == "omitempty"
	if len(f) < 1 || f[0] == "" {
		return fieldName, omitEmpty
	}
	return f[0], omitEmpty
}

// stringToBool returns true if the value is true looking
// if numeric and greater than zero
// if string and starts with t or y
func stringToBool(s string) bool {
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
