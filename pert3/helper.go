package main

import (
	"strconv"
	"strings"

	"latihan-fiber/pert3/app/model"

	"github.com/gofiber/fiber/v2"
)

func ok(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

func okList(c *fiber.Ctx, message string, data interface{}, meta model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

func created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{
		Message: message,
	})
}

func failValidation(c *fiber.Ctx, errors []string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Message: "validasi gagal",
		Errors:  errors,
	})
}

var allowedSortHelper = map[string]bool{
	"id":    true,
	"name":  true,
	"grade": true,
}

func parseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Page:  1,
		Limit: 10,
		Sort:  "id",
		Order: "asc",
	}

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		q.Page = p
	}

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		q.Limit = l
	}

	q.Search = strings.TrimSpace(c.Query("search"))

	sortParam := c.Query("sort")
	if allowedSortHelper[sortParam] {
		q.Sort = sortParam
	}

	orderParam := strings.ToLower(c.Query("order"))
	if orderParam == "desc" {
		q.Order = "desc"
	}

	if raw := c.Query("is_active"); raw != "" {
		val := strings.ToLower(raw) == "true"
		q.IsActive = &val
	}

	return q
}

func validateStudentFields(nim, name string, grade float64) []string {
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