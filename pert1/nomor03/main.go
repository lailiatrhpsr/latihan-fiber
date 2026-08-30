package main

import "fmt"

func swap(a, b *int) {
	*a, *b = *b, *a
}

func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

func ubahNilai(x int) {
	x = 100 
}

func ubahLewatPointer(x *int) {
	*x = 100 
}

func main() {
	a, b := 10, 20
	fmt.Printf("sebelum swap: a=%d, b=%d\n", a, b)
	swap(&a, &b)
	fmt.Printf("setelah swap: a=%d, b=%d\n\n", a, b)

	daftar := []string{"Ayam", "Sapi"}
	fmt.Println("slice awal  :", daftar)
	updateSlice(&daftar, "Kucing")
	fmt.Println("slice update:", daftar, "\n")

	num := 42
	ubahNilai(num)
	fmt.Println("setelah ubahNilai  :", num) 

	ubahLewatPointer(&num)
	fmt.Println("setelah ubahLewatPointer   :", num) 
}