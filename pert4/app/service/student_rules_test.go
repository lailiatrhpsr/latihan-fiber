package service

import (
	"testing"

	"latihan-fiber/pert4/app/model"
)

// Perhatikan: pengujian di file ini tidak menyalakan server, tidak
// menyentuh database, dan tidak membuat fiber.Ctx sama sekali.

func TestValidateCreate(t *testing.T) {
	cases := []struct {
		name    string
		req     model.CreateStudentRequest
		wantErr bool
	}{
		{"data lengkap dan valid", model.CreateStudentRequest{NIM: "12345", Name: "Budi", Grade: 80}, false},
		{"nim kosong", model.CreateStudentRequest{NIM: "   ", Name: "Budi", Grade: 80}, true},
		{"name kosong", model.CreateStudentRequest{NIM: "12345", Name: "", Grade: 80}, true},
		{"grade di atas 100", model.CreateStudentRequest{NIM: "12345", Name: "Budi", Grade: 150}, true},
		{"grade negatif", model.CreateStudentRequest{NIM: "12345", Name: "Budi", Grade: -1}, true},
	}
	for _, tc := range cases {
		errs := ValidateCreate(tc.req)
		if tc.wantErr && len(errs) == 0 {
			t.Errorf("%s: harap ada error, tidak dapat satu pun", tc.name)
		}
		if !tc.wantErr && len(errs) != 0 {
			t.Errorf("%s: tidak harap error, tetapi dapat %v", tc.name, errs)
		}
	}
}

func TestValidateReplace(t *testing.T) {
	cases := []struct {
		name    string
		req     model.ReplaceStudentRequest
		wantErr bool
	}{
		{"data lengkap dan valid", model.ReplaceStudentRequest{NIM: "12345", Name: "Sari", Grade: 90}, false},
		{"nim kosong pada PUT", model.ReplaceStudentRequest{NIM: "", Name: "Sari", Grade: 90}, true},
	}
	for _, tc := range cases {
		errs := ValidateReplace(tc.req)
		if tc.wantErr && len(errs) == 0 {
			t.Errorf("%s: harap ada error, tidak dapat satu pun", tc.name)
		}
		if !tc.wantErr && len(errs) != 0 {
			t.Errorf("%s: tidak harap error, tetapi dapat %v", tc.name, errs)
		}
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{ID: 1, NIM: "111", Name: "Sari", Grade: 70, IsActive: true}
	newGrade := 95.0

	result, errs := ApplyPatch(initial, model.PatchStudentRequest{Grade: &newGrade})
	if len(errs) != 0 {
		t.Fatalf("tidak seharusnya ada error: %v", errs)
	}
	if result.Grade != 95.0 {
		t.Error("grade seharusnya berubah menjadi 95")
	}
	if result.Name != "Sari" {
		t.Error("field yang tidak dikirim seharusnya tidak berubah")
	}
}

func TestApplyPatchGradeDiLuarBatas(t *testing.T) {
	initial := model.Student{ID: 1, NIM: "111", Name: "Sari", Grade: 70}
	badGrade := 200.0

	result, errs := ApplyPatch(initial, model.PatchStudentRequest{Grade: &badGrade})
	if len(errs) == 0 {
		t.Error("grade di luar 0-100 seharusnya menghasilkan error")
	}
	if result.Grade != 70 {
		t.Error("grade seharusnya tidak berubah ketika validasi gagal")
	}
}