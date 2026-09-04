package service

import (
	"strings"

	"latihan-fiber/pert4/app/model"
)

// memvalidasi data dasar mahasiswa
func ValidateStudentFields(nim, name string, grade float64) []string {
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

// memvalidasi request POST
func ValidateCreateStudent(req model.CreateStudentRequest) []string {
	return ValidateStudentFields(req.NIM, req.Name, req.Grade)
}

// memvalidasi request PUT
func ValidateReplaceStudent(req model.ReplaceStudentRequest) []string {
	return ValidateStudentFields(req.NIM, req.Name, req.Grade)
}

// memperbarui field jika tidak bernilai nil pada PATCH
func ApplyPatchStudent(current model.Student, req model.PatchStudentRequest) (model.Student, []string) {
	var errs []string

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			errs = append(errs, "nim tidak boleh kosong")
		} else {
			current.NIM = *req.NIM
		}
	}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs = append(errs, "name tidak boleh kosong")
		} else {
			current.Name = *req.Name
		}
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

// memeriksa apakah request PATCH tidak mengirim field apa pun
func IsEmptyPatchStudent(req model.PatchStudentRequest) bool {
	return req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil
}

// menghitung total halaman pembulatan ke atas
func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}