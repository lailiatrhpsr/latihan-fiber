package service

import (
	"strings"

	"latihan-fiber/pert4/app/model"
)

// File ini berisi business rules MURNI: tidak menyentuh fiber.Ctx,
// tidak menyentuh database, dan tidak tahu apa pun tentang HTTP.

// validateFields adalah pemeriksaan bersama yang dipakai Create dan Replace,
// karena keduanya mengirim seluruh field mahasiswa sekaligus.
func validateFields(nim, name string, grade float64) []string {
	var errs []string
	if strings.TrimSpace(nim) == "" {
		errs = append(errs, "nim wajib diisi")
	}
	if strings.TrimSpace(name) == "" {
		errs = append(errs, "name wajib diisi")
	}
	if grade < 0 || grade > 100 {
		errs = append(errs, "grade harus di antara 0 dan 100")
	}
	return errs
}

// ValidateCreate memeriksa isi permintaan pembuatan mahasiswa (POST).
func ValidateCreate(req model.CreateStudentRequest) []string {
	return validateFields(req.NIM, req.Name, req.Grade)
}

// ValidateReplace memeriksa isi permintaan penggantian penuh (PUT).
func ValidateReplace(req model.ReplaceStudentRequest) []string {
	return validateFields(req.NIM, req.Name, req.Grade)
}

// ApplyPatch menyalin field yang dikirim (tidak nil) ke data yang sudah ada.
// Field yang bernilai nil dibiarkan apa adanya. Mengembalikan data hasil
// beserta daftar error validasi; slice error kosong berarti lolos.
func ApplyPatch(current model.Student, req model.PatchStudentRequest) (model.Student, []string) {
	var errs []string

	if req.NIM != nil {
		current.NIM = *req.NIM
	}
	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs = append(errs, "grade harus di antara 0 dan 100")
		} else {
			current.Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current, errs
}