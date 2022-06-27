//package main
//
//import (
//	"fmt"
//	"time"
//)

//func main() {
//	to, _ := time.Parse("2006-01-02T15:04:05Z", "2021-01-27T10:10:10.294Z")
//	fmt.Println(to)
//	stamp := to.Format("2006-01-02 15:04:05")
//	fmt.Println(stamp)
//}

package main

import (
	"fmt"
	"time"
)

func main() {
	d, err := time.ParseDuration("1h15m30.918273645s")
	if err != nil {
		panic(err)
	}

	round := []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
		time.Hour,
	}

	for _, r := range round {
		fmt.Printf("d.Round(%6s) = %s\n", r, d.Round(r).String())
	}
}
