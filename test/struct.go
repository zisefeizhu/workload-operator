package main

import "fmt"

type a struct {
	name string `json:"name"`
}

func main() {
	reps := make([]a, 0, 10)
	fmt.Println(reps)
}
