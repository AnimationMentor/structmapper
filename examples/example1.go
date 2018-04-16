package main

import (
	"fmt"

	"github.com/AnimationMentor/structmapper"
)

type example1 struct {
	Tuna        string   `json:"tuna"`
	Songs       []string `json:"songs"`
	FavNumber   int      `json:"favnum"`
	Temperature float64  `json:"temp"`
	LikeCandy   bool     `json:"candy"`
}

func main() {

	s := example1{"hello", []string{"hi", "nice"}, 2, 20.5, true}

	m, err := structmapper.StructToStringMap(&s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("from struct: %#v\n", s)
	fmt.Printf("  to    map: %#v\n", m)
}
