package main

import (
	"fmt"

	"github.com/AnimationMentor/structmapper"
)

type color struct {
	Red   float32 `json:"r"`
	Green float32 `json:"g"`
	Blue  float32 `json:"b"`
}

type example1 struct {
	Tuna        string   `json:"tuna"`
	Songs       []string `json:"songs"`
	FavNumber   int      `json:"favnum"`
	Temperature float64  `json:"temp"`
	LikeCandy   bool     `json:"candy"`
	Foreground  color    `json:"color"`
}

func main() {
	m := map[string]string{
		"candy": "",
		// "candy": "99",
		// "candy": "\"false\"",
		// "candy": "false",
		// "candy": "true",
		// "songs":  "1",
		// "songs": "[1,2,3]",
		"songs":  "[\"1\",\"2\",\"3\"]",
		"temp":   "98.6",
		"favnum": "13",
		"tuna":   "100",
		"color":  "{\"r\":100}",
	}
	s := example1{}

	err := structmapper.StringMapToStruct(m, &s, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("from    map: %#v\n", m)
	fmt.Printf("  to struct: %#v\n", s)
}
