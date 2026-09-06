package helper

import (
	"github.com/gofiber/fiber/v2"

	"latihan-fiber/pert4/app/model"
)

// Success mengirim response sukses 
func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

// SuccessList mengirim response sukses yang menyertakan metadata pagination.
func SuccessList(c *fiber.Ctx, message string, data interface{}, meta model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// Created mengirim 201 sekaligus memasang header Location.
func Created(c *fiber.Ctx, message string, data interface{}, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

// NoContent mengirim 204 tanpa body, dipakai setelah Delete berhasil.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Fail mengirim response gagal dengan pesan umum (tanpa detail per field).
func Fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{
		Message: message,
	})
}

// FailValidation mengirim 422 beserta daftar pesan error per field.
func FailValidation(c *fiber.Ctx, errs []string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Message: "validasi gagal",
		Errors:  errs,
	})
}