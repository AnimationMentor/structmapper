package main

// This tests anonymous fields in a struct.

import (
	"fmt"

	"github.com/AnimationMentor/structmapper"
)

type RecordCommon struct {
	ID        string `json:"id"`
	Writeable bool   `json:"writeable"`
}

type color struct {
	Red   float32 `json:"r"`
	Green float32 `json:"g"`
	Blue  float32 `json:"b"`
}

type example1 struct {
	RecordCommon
	Tuna        string   `json:"tuna"`
	Songs       []string `json:"songs"`
	FavNumber   int      `json:"favnum"`
	Temperature float64  `json:"temp"`
	LikeCandy   bool     `json:"candy"`
	Foreground  color    `json:"color"`
}

func main() {

	s := example1{
		Tuna:        "hello",
		Songs:       []string{"hi", "nice"},
		FavNumber:   2,
		Temperature: 20.5,
		LikeCandy:   true,
	}
	s.ID = "12345"

	m, err := structmapper.StructToStringMap(&s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("from struct: %#v\n", s)
	fmt.Printf("  to    map: %#v\n", m)

}
