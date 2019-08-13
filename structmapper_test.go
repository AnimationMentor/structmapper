package structmapper

import (
	"reflect"
	"testing"
)

func Test_BadInputsToStructToStringMap(t *testing.T) {
	a := "a"
	_, err := StructToStringMap(a)
	if err == nil {
		t.Errorf("passing string value to StructToStringMap should return an error")
	}

	_, err = StructToStringMap(&a)
	if err == nil {
		t.Errorf("passing string pointer to StructToStringMap should return an error")
	}

	type testing struct{}

	var b *testing
	_, err = StructToStringMap(b)
	if err == nil {
		t.Errorf("passing nil struct pointer to StructToStringMap should return an error")
	}
}

func Test_BadInputsToStringMapToStruct(t *testing.T) {

	m := map[string]string{}

	a := "a"
	err := StringMapToStruct(m, a, true)
	if err == nil {
		t.Errorf("passing string value to StringMapToStruct should return an error")
	}

	err = StringMapToStruct(m, &a, true)
	if err == nil {
		t.Errorf("passing string pointer to StringMapToStruct should return an error")
	}

	type testing struct{}

	var b *testing
	err = StringMapToStruct(m, b, true)
	if err == nil {
		t.Errorf("passing nil struct pointer to StringMapToStruct should return an error")
	}

	b = &testing{}
	m = nil
	err = StringMapToStruct(m, b, true)
	if err == nil {
		t.Errorf("passing nil map to StringMapToStruct should return an error")
	}
}

type testStruct struct {
	Tuna             string   `json:"tuna"`
	Songs            []string `json:"songs"`
	FavNumber        int      `json:"favnum"`
	Temperature      float64  `json:"temp"`
	LikeCandy        bool     `json:"candy"`
	Quiet            string   `json:"quiet,omitempty"`
	Skip             string   `json:"-"`
	NoTag            string
	unexportedString string
	unexportedPtr    *testStruct2
}

type testStruct2 struct {
	cake string
}

func Test_StructToStringMap(t *testing.T) {

	type testdata struct {
		expectedToError bool
		input           *testStruct
		expectedMap     map[string]string
	}

	for i, v := range []testdata{
		{
			true,
			nil,
			nil,
		},
		{
			false,
			&testStruct{"hello", []string{"hi", "nice"}, 2, 20.5, true, "", "", "", "", nil},
			map[string]string{"tuna": "hello", "songs": "[\"hi\",\"nice\"]", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
		},
		{
			false,
			&testStruct{"hello", []string{"hi", "nice"}, 2, 20.5, true, "A", "", "", "", nil},
			map[string]string{"tuna": "hello", "songs": "[\"hi\",\"nice\"]", "favnum": "2", "temp": "20.5", "candy": "true", "quiet": "A", "NoTag": ""},
		},
		{
			false,
			&testStruct{},
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "NoTag": ""},
		},
	} {
		i++ // prints prettier

		m, err := StructToStringMap(v.input)
		if err != nil {
			if !v.expectedToError {
				t.Errorf("test %d: got unexpected error - %v", i, err)
			}
			continue
		}

		if !reflect.DeepEqual(m, v.expectedMap) {
			t.Errorf("test %d: unexpected result - %#v", i, m)
		}

	}

}

func Test_StringMapToStruct(t *testing.T) {

	type testdata struct {
		expectedToError bool
		inputMap        map[string]string
		expectedStruct  *testStruct
	}

	for i, v := range []testdata{
		{
			true,
			nil,
			nil,
		},
		// Note these two are a reverse of the StructToStringMap tests.
		{
			false,
			map[string]string{"tuna": "hello", "songs": "[\"hi\",\"nice\"]", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", []string{"hi", "nice"}, 2, 20.5, true, "", "", "", "", nil},
		},
		{
			false,
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "quiet": "A", "NoTag": ""},
			&testStruct{Quiet: "A"},
		},
		{
			false,
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "NoTag": ""},
			&testStruct{},
		},
		// These test more crazier inputs.
		{
			false,
			map[string]string{"tuna": "hello", "songs": "[]", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", []string{}, 2, 20.5, true, "", "", "", "", nil},
		},
		{
			false,
			map[string]string{"tuna": "hello", "songs": "", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", nil, 2, 20.5, true, "", "", "", "", nil},
		},
		{
			false,
			map[string]string{"tuna": "hello", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", nil, 2, 20.5, true, "", "", "", "", nil},
		},
		{
			false,
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "Skip": "cake", "NoTag": ""},
			&testStruct{},
		},
	} {
		i++ // prints prettier

		s := &testStruct{}

		err := StringMapToStruct(v.inputMap, s, true)
		if err != nil {
			if !v.expectedToError {
				t.Errorf("test %d: got unexpected error - %v", i, err)
			}
			continue
		}

		if !reflect.DeepEqual(s, v.expectedStruct) {
			t.Errorf("test %d: unexpected result - %#v", i, s)
		}

	}

}

func Test_NonStrictStringMapToStruct(t *testing.T) {

	type testdata struct {
		expectedToError bool
		inputMap        map[string]string
		expectedStruct  *testStruct
	}

	for i, v := range []testdata{
		{
			false,
			map[string]string{"songs": "hi,nice"},
			&testStruct{Songs: []string{"hi", "nice"}},
		},
		{
			false,
			map[string]string{"songs": "hi, nice"},
			&testStruct{Songs: []string{"hi", "nice"}},
		},
		{
			false,
			map[string]string{"songs": "hi"},
			&testStruct{Songs: []string{"hi"}},
		},
		{
			false,
			map[string]string{"candy": "True"},
			&testStruct{LikeCandy: true},
		},
		{
			false,
			map[string]string{"candy": "1"},
			&testStruct{LikeCandy: true},
		},
		{
			false,
			map[string]string{"candy": "99"},
			&testStruct{LikeCandy: true},
		},
		{
			false,
			map[string]string{"candy": "cake"},
			&testStruct{LikeCandy: false},
		},
	} {
		i++ // prints prettier

		s := &testStruct{}

		err := StringMapToStruct(v.inputMap, s, false)
		if err != nil {
			if !v.expectedToError {
				t.Errorf("test %d: got unexpected error - %v", i, err)
			}
			continue
		}

		if !reflect.DeepEqual(s, v.expectedStruct) {
			t.Errorf("test %d: got %#v, expected %#v", i, s, v.expectedStruct)
		}

	}

}

func Test_getJSONTag(t *testing.T) {

	type testdata struct {
		inName            string
		inTag             string
		expectedTag       string
		expectedOmitEmpty bool
	}

	for i, v := range []testdata{
		{"", "", "", false},
		{"Cake", "", "Cake", false},
		{"Mars", "tuna", "tuna", false},
		{"Mars", "tuna,", "tuna", false},
		{"Mars", "tuna,omitempty", "tuna", true},
		{"Mars", ",omitempty", "Mars", true},
		{"Mars", "-,omitempty", "-", true},
		{"Mars", "-,", "-", false},
		{"Mars", "-", "", false},
	} {
		i++ // prints prettier

		gotTag, gotOmitEmpty, _ := getJSONTag(v.inName, v.inTag)

		if gotTag != v.expectedTag || gotOmitEmpty != v.expectedOmitEmpty {
			t.Errorf("test %d: expected (%q,%t) got (%q,%t)", i, v.expectedTag, v.expectedOmitEmpty, gotTag, gotOmitEmpty)
		}

	}

}

func Test_stringToBool(t *testing.T) {

	type testdata struct {
		in       string
		expected bool
	}

	for i, v := range []testdata{
		{"", false},
		{"0", false},
		{"-1", false},
		{"False", false},
		{"false", false},
		{"cheese", false},
		{"1", true},
		{"99", true},
		{"t", true},
		{"true", true},
		{"True", true},
		{"tonsils", true},
	} {
		i++ // prints prettier

		got := stringToBool(v.in)

		if got != v.expected {
			t.Errorf("test %d: expected %t got %t for %q", i, v.expected, got, v.in)
		}

	}

}

func Test_AnonymousFields(t *testing.T) {

	type T2 struct {
		F21 int      `json:"f21"`
		F22 []string `json:"f22"`
	}

	type T1 struct {
		T2
		F11 int      `json:"f11"`
		F12 []string `json:"f12"`
	}

	t1 := T1{
		F11: 11,
		F12: []string{"f12"},
	}
	t1.F21 = 21
	t1.F22 = []string{"f22"}

	m, err := StructToStringMap(&t1)

	expected := map[string]string{
		"f11": "11",
		"f12": "[\"f12\"]",
		"f21": "21",
		"f22": "[\"f22\"]",
	}

	if err != nil {
		t.Errorf("StructToStringMap - unexpected error: %v", err)
		return
	} else if !reflect.DeepEqual(m, expected) {
		t.Errorf("StructToStringMap - expected %#v got %#v", expected, m)
		return
	}

	t2 := T1{}
	err = StringMapToStruct(m, &t2, true)

	if err != nil {
		t.Errorf("StringMapToStruct - unexpected error: %v", err)
		return
	} else if !reflect.DeepEqual(t1, t2) {
		t.Errorf("StringMapToStruct - expected %#v got %#v", t1, t2)
		return
	}
}
