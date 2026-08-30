package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"latihan-fiber/pert3/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound  = errors.New("mahasiswa tidak ditemukan")
	ErrDuplicate = errors.New("nim sudah terdaftar")
)

var allowedSort = map[string]string{
	"id":         "id",
	"name":       "name",
	"grade":      "grade",
	"created_at": "created_at",
}

type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, model.Meta, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

type studentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, model.Meta, error) {
	var (
		conditions []string
		args       []interface{}
		argIdx     = 1
	)

	// 1. Pencarian ILIKE
	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	// 2. Filter IsActive
	if q.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *q.IsActive)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 3. SELECT COUNT(*) untuk Total Data
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, model.Meta{}, err
	}

	// 4. Pengurutan ORDER BY dengan Whitelist
	sortCol, ok := allowedSort[strings.ToLower(q.Sort)]
	if !ok {
		sortCol = "id"
	}

	orderDir := "ASC"
	if strings.ToLower(q.Order) == "desc" {
		orderDir = "DESC"
	}

	// 5. Paginasi LIMIT & OFFSET
	limit := q.Limit
	if limit < 1 {
		limit = 10
	}
	page := q.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	dataQuery := fmt.Sprintf(
		"SELECT id, nim, name, grade, is_active, created_at FROM students %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sortCol, orderDir, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, model.Meta{}, err
	}
	defer rows.Close()

	students := make([]model.Student, 0)
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, model.Meta{}, err
		}
		students = append(students, s)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	meta := model.Meta{
		Page:       page,
		Limit:      limit,
		TotalData:  total,
		TotalPages: totalPages,
	}

	return students, meta, nil
}

func (r *studentRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	var s model.Student
	query := "SELECT id, nim, name, grade, is_active, created_at FROM students WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return s, ErrNotFound
		}
		return s, err
	}
	return s, nil
}

func (r *studentRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	query := `
		INSERT INTO students (nim, name, grade, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nim, name, grade, is_active, created_at
	`
	var created model.Student
	err := r.db.QueryRow(ctx, query, s.NIM, s.Name, s.Grade, s.IsActive).Scan(
		&created.ID, &created.NIM, &created.Name, &created.Grade, &created.IsActive, &created.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return created, ErrDuplicate
		}
		return created, err
	}
	return created, nil
}

func (r *studentRepository) Update(ctx context.Context, s model.Student) (model.Student, error) {
	query := `
		UPDATE students
		SET nim = $1, name = $2, grade = $3, is_active = $4
		WHERE id = $5
		RETURNING id, nim, name, grade, is_active, created_at
	`
	var updated model.Student
	err := r.db.QueryRow(ctx, query, s.NIM, s.Name, s.Grade, s.IsActive, s.ID).Scan(
		&updated.ID, &updated.NIM, &updated.Name, &updated.Grade, &updated.IsActive, &updated.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return updated, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return updated, ErrDuplicate
		}
		return updated, err
	}
	return updated, nil
}

func (r *studentRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM students WHERE id = $1"
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}