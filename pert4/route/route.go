package route

import (
	"github.com/gofiber/fiber/v2"

	"latihan-fiber/pert4/app/service"
	"latihan-fiber/pert4/middleware"
)

// Register memetakan URL ke method pada service.
//
// Perhatikan isi file ini: tidak ada logika bisnis, tidak ada query,
// tidak ada validasi. Hanya daftar alamat dan siapa yang melayaninya.
func Register(app *fiber.App, studentService *service.StudentService) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")
	api.Get("/health", studentService.HealthCheck)

	students := api.Group("/students", middleware.RequireJSON)
	students.Get("/", studentService.List)
	students.Get("/:id", studentService.Get)
	students.Post("/", studentService.Create)
	students.Put("/:id", studentService.Replace)
	students.Patch("/:id", studentService.Patch)
	students.Delete("/:id", studentService.Delete)
}
