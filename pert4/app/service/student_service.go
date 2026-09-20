package service

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihan-fiber/pert4/app/model"
	"latihan-fiber/pert4/app/repository"
	"latihan-fiber/pert4/helper"
)

// StudentService memegang dua tanggung jawab sekaligus pada struktur baku
// mata kuliah ini: menerima *fiber.Ctx (peran controller) dan menjalankan
// business rules (peran use case). Pemisahan murninya ada di student_rules.go.
type StudentService struct {
	repo  repository.StudentRepository
	db    *pgxpool.Pool
	perms *helper.PermissionSet
}

// NewStudentService menerima INTERFACE repository, bukan struct konkret.
func NewStudentService(
	repo repository.StudentRepository, db *pgxpool.Pool, perms *helper.PermissionSet,
) *StudentService {
	return &StudentService{repo: repo, db: db, perms: perms}
}

// HealthCheck memeriksa ketersediaan server dan koneksi PostgreSQL.
func (s *StudentService) HealthCheck(c *fiber.Ctx) error {
	if err := s.db.Ping(c.Context()); err != nil {
		return helper.Fail(c, fiber.StatusServiceUnavailable, "layanan basis data tidak tersedia")
	}
	return helper.Success(c, fiber.StatusOK, "server dan basis data berjalan", fiber.Map{
		"database": "connected",
	})
}

// GET /api/v1/students
func (s *StudentService) List(c *fiber.Ctx) error {
	q := helper.ParseListQuery(c)

	students, meta, err := s.repo.FindAll(c.Context(), q)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return helper.SuccessList(c, "daftar mahasiswa berhasil diambil", students, meta)
}

// GET /api/v1/students/:id
func (s *StudentService) Get(c *fiber.Ctx) error {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return translateError(c, err, "terjadi kesalahan pada server")
	}

	if !CanAccessStudent(current, student.OwnerID, s.perms, "student:read:any") {
		return helper.Fail(c, fiber.StatusForbidden,
			"tidak berhak mengakses data mahasiswa ini")
	}

	return helper.Success(c, fiber.StatusOK, "mahasiswa ditemukan", student)
}

// POST /api/v1/students
// Dijaga middleware dengan permission student:create (lihat route/route.go).
func (s *StudentService) Create(c *fiber.Ctx) error {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	// Business rules-nya dipanggil, bukan ditulis ulang di sini.
	if errs := ValidateCreate(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	newStudent := model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
		OwnerID:  current.UserID,
	}

	createdStudent, err := s.repo.Create(c.Context(), newStudent)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, "nim sudah terdaftar")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return helper.Created(c, "mahasiswa berhasil dibuat", createdStudent,
		"/api/v1/students/"+strconv.Itoa(createdStudent.ID))
}

// PUT /api/v1/students/:id (Replace seluruh field)
func (s *StudentService) Replace(c *fiber.Ctx) error {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	if errs := ValidateReplace(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	existing, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return translateError(c, err, "terjadi kesalahan pada server")
	}

	if !CanAccessStudent(current, existing.OwnerID, s.perms, "student:update:any") {
		return helper.Fail(c, fiber.StatusForbidden,
			"tidak berhak mengubah data mahasiswa ini")
	}

	studentToUpdate := model.Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}

	updatedStudent, err := s.repo.Update(c.Context(), studentToUpdate)
	if err != nil {
		return translateError(c, err, "terjadi kesalahan pada server")
	}

	return helper.Success(c, fiber.StatusOK, "mahasiswa berhasil diperbarui (replace)", updatedStudent)
}

// PATCH /api/v1/students/:id (Partial Update)
func (s *StudentService) Patch(c *fiber.Ctx) error {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	// Ambil data lama dari basis data terlebih dahulu.
	existing, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return translateError(c, err, "terjadi kesalahan pada server")
	}

	if !CanAccessStudent(current, existing.OwnerID, s.perms, "student:update:any") {
		return helper.Fail(c, fiber.StatusForbidden,
			"tidak berhak mengubah data mahasiswa ini")
	}

	updated, errs := ApplyPatch(existing, req)
	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	result, err := s.repo.Update(c.Context(), updated)
	if err != nil {
		return translateError(c, err, "terjadi kesalahan pada server")
	}

	return helper.Success(c, fiber.StatusOK, "mahasiswa berhasil diperbarui (patch)", result)
}

// DELETE /api/v1/students/:id
func (s *StudentService) Delete(c *fiber.Ctx) error {
	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		return translateError(c, err, "terjadi kesalahan pada server")
	}

	return helper.NoContent(c)
}

func translateError(c *fiber.Ctx, err error, generalMessage string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.Fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Fail(c, fiber.StatusConflict, "nim sudah dipakai mahasiswa lain")
	default:
		return helper.Fail(c, fiber.StatusInternalServerError, generalMessage)
	}
}
