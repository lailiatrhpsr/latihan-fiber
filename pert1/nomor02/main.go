package main

import "fmt"

func main() {
	fmt.Println("=== DEKLARASI 5 VARIABEL ===")
	var nama string = "Lailia Trihapsari"
	var umur int = 20
	var ipk float64 = 3.85
	var isAktif bool = true
	hobi := []string{"Coding,", "Desain,", "Membaca"}

	fmt.Printf("Nama   : %s \n", nama)
	fmt.Printf("Umur   : %d \n", umur)
	fmt.Printf("IPK    : %.2f \n", ipk)
	fmt.Printf("Aktif  : %t \n", isAktif)
	fmt.Printf("Hobi   : %v \n\n", hobi)

	fmt.Println("=== OPERASI MAP DATA MAHASISWA ===")
	nilaiMahasiswa := make(map[string]int)

	// Menambah 
	nilaiMahasiswa["Valerina"] = 85
	nilaiMahasiswa["Gracia"] = 92
	nilaiMahasiswa["Alisya"] = 78
	nilaiMahasiswa["Lailia"] = 95
	fmt.Println("Data mahasiswa berhasil ditambahkan.")

	// Membaca
	namaCek := "Alisya"
	if nilai, exists := nilaiMahasiswa[namaCek]; exists {
		fmt.Printf("Data ditemukan -> %s: %d\n", namaCek, nilai)
	} else {
		fmt.Printf("Data %s tidak ditemukan.\n", namaCek)
	}

	// Menghapus 
	fmt.Println("\nMenghapus data Lailia dari map...")
	delete(nilaiMahasiswa, "Lailia")

	// Menelusuri seluruh isi
	fmt.Println("\nSeluruh data mahasiswa setelah dihapus:")
	for namaMhs, nilai := range nilaiMahasiswa {
		fmt.Printf("- %-8s : %d\n", namaMhs, nilai)
	}
}