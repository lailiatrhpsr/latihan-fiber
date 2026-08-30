package main

import "fmt"

type Student struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

// Value receiver
func (s Student) GetInfo() string {
	return fmt.Sprintf("ID: %d | Nama: %s | Nilai: %.2f | Aktif: %v", s.ID, s.Name, s.Grade, s.IsActive)
}

// Pointer receiver
func (s *Student) UpdateGrade(grade float64) { s.Grade = grade }
func (s *Student) Activate()                 { s.IsActive = true }
func (s *Student) Deactivate()               { s.IsActive = false }

func main() {
	mhs := Student{ID: 1, Name: "Lailia", Grade: 80.0, IsActive: false}

	fmt.Println("--- Kondisi Awal ---")
	fmt.Println(mhs.GetInfo())

	fmt.Println("\n--- Setelah Activate & UpdateGrade ---")
	mhs.Activate()
	mhs.UpdateGrade(92.5)
	fmt.Println(mhs.GetInfo())

	fmt.Println("\n--- Setelah Deactivate ---")
	mhs.Deactivate()
	fmt.Println(mhs.GetInfo())
}
