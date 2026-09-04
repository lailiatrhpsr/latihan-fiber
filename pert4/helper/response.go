package helper

import (
	"github.com/gofiber/fiber/v2"

	"latihan-fiber/pert4/app/model"
)

func Ok(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

func OkList(c *fiber.Ctx, message string, data interface{}, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *fiber.Ctx, message string, data interface{}, location string) error {
	if location != "" {
		c.Set("Location", location)
	}
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Message: message,
		Data:    data,
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{
		Message: message,
	})
}

func FailValidation(c *fiber.Ctx, errors []string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Message: "validasi gagal",
		Errors:  errors,
	})
}