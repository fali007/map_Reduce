package main

import (
	"fmt"
	"time"
	// "mr/search"
	"mr/matrix"
)

func main() {
	fmt.Println("Hello, welcome to search")
	// search.Search("key", "../../sustainableComputing")
	m := matrix.Generate_matrix(2000)
	n := matrix.Generate_matrix(2000)
	// fmt.Printf("%+v\n%+v\n", m, n)

	t1 := time.Now()
	r := matrix.Multiply_S(m,n)
	fmt.Println("Time taken :", time.Now().Sub(t1))

	t1 = time.Now()
	rm := matrix.Multiply_M(m,n)
	fmt.Println("Time taken :", time.Now().Sub(t1))

	t1 = time.Now()
	rmm := matrix.Multiply_MM(m,n)
	fmt.Println("Time taken :", time.Now().Sub(t1))

	fmt.Println(len(r), len(rm), len(rmm))
	// fmt.Printf("%+v\n%+v\n%+v\n", r, rm, rmm)
}
