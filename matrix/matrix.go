package matrix

import (
	"math/rand"
)

const (
	system_cores = 8
)

type Result struct {
	Position int
	Value []float64
}

type ResultM struct {
	PositionR int
	PositionC int
	Value float64
}

func Multiply_S (a, b [][]float64) [][]float64{
	a_row_s := len(a)
	a_col_s := len(a[0])
	b_row_s := len(b)
	b_col_s := len(b[0])

	if a_col_s != b_row_s {
		panic("Multiplication not defined between matrices")
	}

	var res [][]float64 = make([][]float64, a_row_s)

	for row := 0; row < a_row_s; row++ {
		temp_row := make([]float64, b_row_s)
		for col := 0; col < b_col_s; col++ {
			for index := 0; index < b_row_s; index++ {
				temp_row[col] += a[row][index] * b[index][col]
			}
		}
		res[row] = temp_row
	}
	return res
}

func Multiply_M(a, b [][]float64) [][]float64{
	a_row_s := len(a)
	a_col_s := len(a[0])
	b_row_s := len(b)

	if a_col_s != b_row_s {
		panic("Multiplication not defined between matrices")
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
	return res
}

func get_row(a []float64, b[][]float64, row int, c chan Result) {
	res := make([]float64, len(b))
	for col := 0; col < len(a); col++ {
		for index := 0; index < len(a); index++ {
			res[col] += a[index] * b[index][col]
		}
	}
	c <- Result{row, res}
}

// useless
func Multiply_MM (a, b [][]float64) [][]float64 {
	a_row_s := len(a)
	a_col_s := len(a[0])
	b_row_s := len(b)
	b_col_s := len(b[0])

	if a_col_s != b_row_s {
		panic("Multiplication not defined between matrices")
	}

	res_chan := make(chan ResultM, a_row_s)

	var res [][]float64 = make([][]float64, a_row_s)

	for row := 0; row < a_row_s; row++ {
		res[row] = make([]float64, b_col_s)
		go get_row_m(a[row], b, row, res_chan)
	}

	iter :=  a_row_s * b_col_s
	for i := 0; i < iter; i++ {
		idx := <-res_chan
		res[idx.PositionR][idx.PositionC] = idx.Value
	}
	return res
}

func get_row_m(a []float64, b[][]float64, row int, c chan ResultM) {
	n := len(b)
	iter := len(b[0])
	for col := 0; col < iter; col++ {
		col_mat := make([]float64, n)
		for index := 0; index < n; index++ {
			col_mat[index] = b[index][col]
		}
		go get_index_value(a, col_mat, row, col, c)
	}
}

func get_index_value(a, b []float64, row, col int, c chan ResultM) {
	n := len(a)
	res := 0.0
	for i := 0; i < n; i++ {
		res += a[i] * b[i]
	}
	c <- ResultM{row, col, res}
}

func Generate_matrix(n int) [][]float64 {
	limit := 10.0

	var res [][]float64 = make([][]float64, n)
	for row := 0; row < n; row++ {
		temp := make([]float64, n)
		for col := 0; col < n; col++ {
			temp[col] = rand.Float64() * limit
		}
		res[row] = temp
	}
	return res
}	