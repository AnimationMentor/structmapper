package structmapper_test

import (
	"github.com/AnimationMentor/structmapper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Structmapper", func() {

	type testStruct2 struct {
		cake string
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

	It("handles bad inputs to StructToStringMap", func() {
		a := "a"
		Expect(structmapper.StructToStringMap(a)).Error().To(HaveOccurred(), "passing string value to StructToStringMap should return an error")

		Expect(structmapper.StructToStringMap(&a)).Error().To(HaveOccurred(), "passing string pointer to StructToStringMap should return an error")

		type testing struct{}

		var b *testing
		Expect(structmapper.StructToStringMap(b)).Error().To(HaveOccurred(), "passing nil struct pointer to StructToStringMap should return an error")
	})

	It("handles bad inputs to StringMapToStruct", func() {
		m := map[string]string{}

		a := "a"
		Expect(structmapper.StringMapToStruct(m, a, true)).
			To(MatchError("s must be a pointer to a struct",
				"passing string value to StringMapToStruct should return an error"))

		Expect(structmapper.StringMapToStruct(m, &a, true)).
			To(MatchError("s must be a pointer to a struct",
				"passing string pointer to StringMapToStruct should return an error"))

		type testing struct{}
		var b *testing
		Expect(structmapper.StringMapToStruct(m, b, true)).To(MatchError("s must be a pointer to a struct",
			"passing nil struct pointer to StringMapToStruct should return an error"))

		b = &testing{}
		m = nil
		Expect(structmapper.StringMapToStruct(m, b, true)).To(MatchError("m must not be nil",
			"m must not be nil"))
	})

	DescribeTable("StructToStringMap",
		func(expectedToError bool, input *testStruct, expectedMap map[string]string) {
			m, err := structmapper.StructToStringMap(input)
			if expectedToError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(m).To(Equal(expectedMap))
			}
		},

		Entry(nil,
			true,
			nil,
			nil,
		),
		Entry(nil,
			false,
			&testStruct{"hello", []string{"hi", "nice"}, 2, 20.5, true, "", "", "", "", nil},
			map[string]string{"tuna": "hello", "songs": "[\"hi\",\"nice\"]", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
		),
		Entry(nil,
			false,
			&testStruct{"hello", []string{"hi", "nice"}, 2, 20.5, true, "A", "", "", "", nil},
			map[string]string{"tuna": "hello", "songs": "[\"hi\",\"nice\"]", "favnum": "2", "temp": "20.5", "candy": "true", "quiet": "A", "NoTag": ""},
		),
		Entry(nil,
			false,
			&testStruct{},
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "NoTag": ""},
		),
	)

	DescribeTable("StringMapToStruct",
		func(expectedToError bool, inputMap map[string]string, expectedStruct *testStruct) {
			s := &testStruct{}

			err := structmapper.StringMapToStruct(inputMap, s, true)

			if expectedToError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(s).To(Equal(expectedStruct))
			}
		},
		Entry(nil,
			true,
			nil,
			nil,
		),
		// Note these two are a reverse of the StructToStringMap tests.
		Entry(nil,
			false,
			map[string]string{"tuna": "hello", "songs": "[\"hi\",\"nice\"]", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", []string{"hi", "nice"}, 2, 20.5, true, "", "", "", "", nil},
		),
		Entry(nil,
			false,
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "quiet": "A", "NoTag": ""},
			&testStruct{Quiet: "A"},
		),
		Entry(nil,
			false,
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "NoTag": ""},
			&testStruct{},
		),
		// These test more crazier inputs.
		Entry(nil,
			false,
			map[string]string{"tuna": "hello", "songs": "[]", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", []string{}, 2, 20.5, true, "", "", "", "", nil},
		),
		Entry(nil,
			false,
			map[string]string{"tuna": "hello", "songs": "", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", nil, 2, 20.5, true, "", "", "", "", nil},
		),
		Entry(nil,
			false,
			map[string]string{"tuna": "hello", "favnum": "2", "temp": "20.5", "candy": "true", "NoTag": ""},
			&testStruct{"hello", nil, 2, 20.5, true, "", "", "", "", nil},
		),
		Entry(nil,
			false,
			map[string]string{"tuna": "", "songs": "null", "favnum": "0", "temp": "0", "candy": "false", "Skip": "cake", "NoTag": ""},
			&testStruct{},
		),
	)

	DescribeTable("NonStrictStringMapToStruct",
		func(expectedToError bool, inputMap map[string]string, expectedStruct *testStruct) {
			s := &testStruct{}

			err := structmapper.StringMapToStruct(inputMap, s, false)
			if expectedToError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(s).To(Equal(expectedStruct))
			}
		},
		Entry(nil,
			false,
			map[string]string{"songs": "hi,nice"},
			&testStruct{Songs: []string{"hi", "nice"}},
		),
		Entry(nil,
			false,
			map[string]string{"songs": "hi, nice"},
			&testStruct{Songs: []string{"hi", "nice"}},
		),
		Entry(nil,
			false,
			map[string]string{"songs": "hi"},
			&testStruct{Songs: []string{"hi"}},
		),
		Entry(nil,
			false,
			map[string]string{"candy": "True"},
			&testStruct{LikeCandy: true},
		),
		Entry(nil,
			false,
			map[string]string{"candy": "1"},
			&testStruct{LikeCandy: true},
		),
		Entry(nil,
			false,
			map[string]string{"candy": "99"},
			&testStruct{LikeCandy: true},
		),
		Entry(nil,
			false,
			map[string]string{"candy": "cake"},
			&testStruct{LikeCandy: false},
		),
	)

	DescribeTable("GetJSONTag",
		func(inName string, inTag string, expectedTag string, expectedOmitEmpty bool) {
			gotTag, gotOmitEmpty, _ := structmapper.GetJSONTag(inName, inTag)
			Expect(gotTag).To(Equal(expectedTag))
			Expect(gotOmitEmpty).To(Equal(expectedOmitEmpty))
		},
		Entry(nil, "", "", "", false),
		Entry(nil, "Cake", "", "Cake", false),
		Entry(nil, "Mars", "tuna", "tuna", false),
		Entry(nil, "Mars", "tuna,", "tuna", false),
		Entry(nil, "Mars", "tuna,omitempty", "tuna", true),
		Entry(nil, "Mars", ",omitempty", "Mars", true),
		Entry(nil, "Mars", "-,omitempty", "-", true),
		Entry(nil, "Mars", "-,", "-", false),
		Entry(nil, "Mars", "-", "", false),
	)

	DescribeTable("StringToBool",
		func(in string, expected bool) {
			Expect(structmapper.StringToBool(in)).To(Equal(expected))
		},

		Entry(nil, "", false),
		Entry(nil, "0", false),
		Entry(nil, "-1", false),
		Entry(nil, "False", false),
		Entry(nil, "false", false),
		Entry(nil, "cheese", false),
		Entry(nil, "1", true),
		Entry(nil, "99", true),
		Entry(nil, "t", true),
		Entry(nil, "true", true),
		Entry(nil, "True", true),
		Entry(nil, "tonsils", true),
	)

	It("handles AnonymousFields", func() {
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

		m, err := structmapper.StructToStringMap(&t1)

		expected := map[string]string{
			"f11": "11",
			"f12": "[\"f12\"]",
			"f21": "21",
			"f22": "[\"f22\"]",
		}

		Expect(err).To(Succeed())
		Expect(m).To(Equal(expected))

		t2 := T1{}
		Expect(structmapper.StringMapToStruct(m, &t2, true)).To(Succeed())
		Expect(t1).To(Equal(t2))
	})
})
