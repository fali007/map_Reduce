package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Result struct {
	Position int
	Value []float64
}

type Transformer struct {
	Query [][]float64
	Key [][]float64
	Value [][]float64
}

func generate_matrix(n, m, limit int) [][]float64 {
	var res [][]float64 = make([][]float64, n)
	for row := 0; row < n; row++ {
		temp := make([]float64, m)
		for col := 0; col < m; col++ {
			temp[col] = rand.Float64() * float64(limit)
		}
		res[row] = temp
	}
	return res
}	

func multiply(a, b [][]float64) [][]float64{
	a_row_s := len(a)
	a_col_s := len(a[0])
	b_row_s := len(b)

	// fmt.Printf("Multiplying A of shape (%d, %d) and B of shape (%d, %d)\n", len(a), len(a[0]), len(b), len(b[0]))
	if a_col_s != b_row_s {
		panic(fmt.Sprintf("Multiplication not defined between matrices\nA with shape (%d, %d) and B with shape (%d, %d)\n", len(b), a_col_s, b_row_s, len(b[0])))
	}

	res_chan := make(chan Result, a_row_s)

	var res [][]float64 = make([][]float64, a_row_s)
	for row := 0; row < a_row_s; row++ {
		go get_row(a[row], b, row, res_chan)
	}
	for i := 0; i < a_row_s; i++ {
		idx := <-res_chan
		res[idx.Position] = idx.Value
	}
	// fmt.Printf("Output shape (%d, %d)\n", len(res), len(res[0]))
	return res
}

func get_row(a []float64, b[][]float64, row int, c chan Result) {
	l := len(b[0])
	res := make([]float64, l)
	for col := 0; col < l; col++ {
		for index := 0; index < len(a); index++ {
			res[col] += a[index] * b[index][col]
		}
	}
	c <- Result{row, res}
}

func transpose(i [][]float64) [][]float64 {
	m := len(i[0])
	n := len(i)
	var res [][]float64 = make([][]float64, m)
	for row := 0; row < m; row++ {
		temp := make([]float64, n)
		for col := 0; col < n; col++ {
			temp[col] = i[col][row]
		}
		res[row] = temp
	}
	return res
}

func scale(i [][]float64, s float64) [][]float64{
	m := len(i)
	res_chan := make(chan Result, m)
	for row := 0; row < m; row++ {
		go scale_row(i[row], s, row, res_chan)
	}
	for index := 0; index < m; index++ {
		idx := <-res_chan
		i[idx.Position] = idx.Value
	}
	return i
}

func scale_row(a []float64, s float64, row int, c chan Result) {
	for col := 0; col < len(a); col++ {
		a[col] = a[col] / s
		if a[col] != a[col] {
			a[col] = 0
		}
	}
	c <- Result{row, a}
}

func sum_row(i []float64) float64 {
	l := len(i)
	temp := 0.0
	for index := 0; index < l; index++ {
		temp += i[index]
	}
	return temp
}

func softmax(i [][]float64) [][]float64 {
	m := len(i)
	res_chan := make(chan Result, m)
	for row := 0; row < m; row++ {
		go scale_row(i[row], sum_row(i[row]), row, res_chan)
	}
	for index := 0; index < m; index++ {
		idx := <-res_chan
		i[idx.Position] = idx.Value
	}
	return i
}

func (t Transformer) init(input_dim, embedding_dim int) Transformer{
	t.Query = generate_matrix(embedding_dim, input_dim, 1)
	t.Key = generate_matrix(embedding_dim, input_dim, 2)
	t.Value = generate_matrix(embedding_dim, input_dim, 1)
	return t
}

func (t Transformer) forward(input [][]float64) [][]float64 {
	q := multiply(input, t.Query)
	fmt.Println(q)
	k := multiply(input, t.Key)
	v := multiply(input, t.Value)

	multi := multiply(q , transpose(k))
	attn := multiply(softmax(scale(multi, math.Sqrt(float64(len(input[0]))))), v)
	return attn
}

func main(){
	input := generate_matrix(20,30,1)

	t := Transformer{}
	t = t.init(len(input), len(input[0]))
	
	fmt.Println(t.forward(input))
}

