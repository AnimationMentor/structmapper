
# structmapper

Sometimes you want to convert a Go struct into a `map[string]string` (or the inverse). And sometimes your struct is not a simple match for that. Maybe you have a field that's a `[]string` or some other embedded structure.

Structmapper solves this by encoding non-strings as json.

# Why?

This was written to make it easy to convert Go structs to and from redis hashes. Redis hashes are precisely equivalent to `map[string]string` .


# Method used

When converting from a struct to a map (`StructToStringMap`), an entry is made for each field. The keys used and omit behaviour follows the json tag settings on the struct. String values are copied as is, other types are JSON encoded.

When converting from a map to a struct (`StringMapToStruct`), the operation is reversed. If the `strict` option is unset then some additional attempts are made to convert non-JSON inputs. String slices which don't successfully JSON decode are treated as comma separated lists. Bool values are more liberally detected in either case (even when `strict` is set).

# Example

(See `examples` directory.)

```
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

```

prints

```
from struct: main.example1{Tuna:"hello", Songs:[]string{"hi", "nice"}, FavNumber:2, Temperature:20.5, LikeCandy:true}
  to    map: map[string]string{"favnum":"2", "temp":"20.5", "candy":"true", "tuna":"hello", "songs":"[\"hi\",\"nice\"]"}
```


# Redis example

See full example at `examples/redis_example.go`

```

...

	s := example1{
		Tuna:       "this is odd",
		Songs:      []string{"one", "two", "three"},
		Foreground: color{1.0002, 0.5, 0.0},
	}

	m, _ := structmapper.StructToStringMap(&s)

...

	r.Do("hmset", redis.Args{}.Add("redis_example").AddFlat(m)...)

...

	m2, _ := redis.StringMap(r.Do("hgetall", "redis_example"))

	var s2 example1

	structmapper.StringMapToStruct(m2, &s2, true)

...

```


# Future work

String maps aren't the cheapest option when using redis in this way but it was the tidiest to implement. A slice of paired key, value strings are closer to what's actually used by redis on the wire and should be a little more efficient on time and space.

In either case, the map representation is still useful.

# etc

- Why not use the redis json module? I needed this to work on a very plain and little bit out of date redis server instance.
