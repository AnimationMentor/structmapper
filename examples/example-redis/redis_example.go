package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gomodule/redigo/redis"

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
	s := example1{
		Tuna:       "this is odd",
		Songs:      []string{"one", "two", "three"},
		Foreground: color{1.0002, 0.5, 0.0},
	}

	fmt.Printf("\noriginal struct:\n%#v\n\n", s)

	m, err := structmapper.StructToStringMap(&s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("map created from struct:\n%#v\n\n", m)

	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	_, err = r.Do("hmset", redis.Args{}.Add("redis_example").AddFlat(m)...)
	if err != nil {
		log.Fatal(err)
	}

	m2, err := redis.StringMap(r.Do("hgetall", "redis_example"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("map read from redis:\n%#v\n\n", m2)

	var s2 example1

	err = structmapper.StringMapToStruct(m2, &s2, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("struct re-hydrated from map:\n%#v\n\n", s2)

	if reflect.DeepEqual(s, s2) {
		fmt.Println("and it matches our original struct")
	} else {
		fmt.Println("OOPS, shouldn't happen, this struct does not match our original")
	}
}
