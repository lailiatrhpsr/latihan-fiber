package main

import (
	"errors"
	"strconv"

	"latihan-fiber/pert3/app/model"
	"latihan-fiber/pert3/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentHandler struct {
	repo repository.StudentRepository
	db   *pgxpool.Pool
}

func NewStudentHandler(repo repository.StudentRepository, db *pgxpool.Pool) *StudentHandler {
	return &StudentHandler{repo: repo, db: db}
}

// HealthCheck memeriksa ketersediaan server dan koneksi PostgreSQL
func (h *StudentHandler) HealthCheck(c *fiber.Ctx) error {
	if err := h.db.Ping(c.Context()); err != nil {
		return fail(c, fiber.StatusServiceUnavailable, "layanan basis data tidak tersedia")
	}
	return ok(c, "server dan basis data berjalan", fiber.Map{
		"database": "connected",
	})
}

// GET /api/v1/students
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	students, meta, err := h.repo.FindAll(c.Context(), q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return okList(c, "daftar mahasiswa berhasil diambil", students, meta)
}

// GET /api/v1/students/:id
func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	student, err := h.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
		}
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return ok(c, "mahasiswa ditemukan", student)
}

// POST /api/v1/students
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	if errs := validateStudentFields(req.NIM, req.Name, req.Grade); len(errs) > 0 {
		return failValidation(c, errs)
	}

	newStudent := model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}

	createdStudent, err := h.repo.Create(c.Context(), newStudent)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return fail(c, fiber.StatusConflict, "nim sudah terdaftar")
		}
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	c.Set("Location", "/api/v1/students/"+strconv.Itoa(createdStudent.ID))
	return created(c, "mahasiswa berhasil dibuat", createdStudent)
}

// PUT /api/v1/students/:id (Replace seluruh field)
func (h *StudentHandler) ReplaceStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	if errs := validateStudentFields(req.NIM, req.Name, req.Grade); len(errs) > 0 {
		return failValidation(c, errs)
	}

	studentToUpdate := model.Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}

	updatedStudent, err := h.repo.Update(c.Context(), studentToUpdate)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return fail(c, fiber.StatusConflict, "nim sudah dipakai mahasiswa lain")
		}
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return ok(c, "mahasiswa berhasil diperbarui (replace)", updatedStudent)
}

// PATCH /api/v1/students/:id (Partial Update)
func (h *StudentHandler) PatchStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body tidak valid")
	}

	// Ambil data lama dari basis data terlebih dahulu
	existing, err := h.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
		}
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	// Update hanya field yang dikirim (pointer not nil)
	if req.NIM != nil {
		existing.NIM = *req.NIM
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			return failValidation(c, []string{"grade harus di antara 0 dan 100"})
		}
		existing.Grade = *req.Grade
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updatedStudent, err := h.repo.Update(c.Context(), existing)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return fail(c, fiber.StatusConflict, "nim sudah dipakai mahasiswa lain")
		}
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return ok(c, "mahasiswa berhasil diperbarui (patch)", updatedStudent)
}

// DELETE /api/v1/students/:id
func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka")
	}

	err = h.repo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
		}
		return fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return noContent(c)
}