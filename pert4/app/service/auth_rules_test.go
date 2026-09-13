package service

import (
	"testing"

	"latihan-fiber/pert4/app/model"
)

func TestValidateRegister(t *testing.T) {
	cases := []struct {
		name    string
		req     model.RegisterRequest
		wantErr bool
	}{
		{
			"data lengkap dan valid",
			model.RegisterRequest{Username: "sari", Email: "sari@example.com", Password: "rahasia123"},
			false,
		},
		{
			"username kosong",
			model.RegisterRequest{Username: "  ", Email: "sari@example.com", Password: "rahasia123"},
			true,
		},
		{
			"username terlalu pendek",
			model.RegisterRequest{Username: "ab", Email: "sari@example.com", Password: "rahasia123"},
			true,
		},
		{
			"username mengandung karakter tidak diizinkan",
			model.RegisterRequest{Username: "sari!", Email: "sari@example.com", Password: "rahasia123"},
			true,
		},
		{
			"email tidak valid",
			model.RegisterRequest{Username: "sari", Email: "bukan-email", Password: "rahasia123"},
			true,
		},
		{
			"password terlalu pendek",
			model.RegisterRequest{Username: "sari", Email: "sari@example.com", Password: "rhs123"},
			true,
		},
		{
			"password tanpa angka",
			model.RegisterRequest{Username: "sari", Email: "sari@example.com", Password: "rahasiasaja"},
			true,
		},
		{
			"password terlalu umum",
			model.RegisterRequest{Username: "sari", Email: "sari@example.com", Password: "password1"},
			true,
		},
	}

	for _, tc := range cases {
		errs := ValidateRegister(tc.req)
		if tc.wantErr && len(errs) == 0 {
			t.Errorf("%s: harap ada error, tidak dapat satu pun", tc.name)
		}
		if !tc.wantErr && len(errs) != 0 {
			t.Errorf("%s: tidak harap error, tetapi dapat %v", tc.name, errs)
		}
	}
}

func TestValidateLogin(t *testing.T) {
	cases := []struct {
		name    string
		req     model.LoginRequest
		wantErr bool
	}{
		{"data lengkap", model.LoginRequest{Username: "sari", Password: "apapun"}, false},
		{"username kosong", model.LoginRequest{Username: "", Password: "apapun"}, true},
		{"password kosong", model.LoginRequest{Username: "sari", Password: ""}, true},
		{
			// ValidateLogin sengaja tidak memberlakukan aturan kekuatan
			// password: user lama mungkin dibuat sebelum aturannya ada.
			"password lama yang lemah tetap lolos validasi kelengkapan",
			model.LoginRequest{Username: "sari", Password: "123"},
			false,
		},
	}

	for _, tc := range cases {
		errs := ValidateLogin(tc.req)
		if tc.wantErr && len(errs) == 0 {
			t.Errorf("%s: harap ada error, tidak dapat satu pun", tc.name)
		}
		if !tc.wantErr && len(errs) != 0 {
			t.Errorf("%s: tidak harap error, tetapi dapat %v", tc.name, errs)
		}
	}
}

func TestCheckPasswordStrength(t *testing.T) {
	if msg := checkPasswordStrength("rahasia123"); msg != "" {
		t.Errorf("password kuat seharusnya lolos, tetapi dapat: %s", msg)
	}
	if msg := checkPasswordStrength("pendek1"); msg == "" {
		t.Error("password di bawah 8 karakter seharusnya ditolak")
	}
	if msg := checkPasswordStrength("hanyahuruf"); msg == "" {
		t.Error("password tanpa angka seharusnya ditolak")
	}
	if msg := checkPasswordStrength("Password123"); msg == "" {
		t.Error("password yang ada di daftar lemah seharusnya ditolak (case-insensitive)")
	}
}
