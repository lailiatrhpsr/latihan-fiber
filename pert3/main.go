package main

import (
	"fmt"
	"log"
	"strings"

	"latihan-fiber/pert3/app/repository"
	"latihan-fiber/pert3/config"
	"latihan-fiber/pert3/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var metodeBerbody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func requireJSON(c *fiber.Ctx) error {
	if metodeBerbody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType,
				"Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	config.LoadEnv()

	// Inisialisasi PostgreSQL Connection Pool
	dbPool := database.NewPostgresPool()
	defer dbPool.Close()

	// Inisialisasi Repository & Handler
	studentRepo := repository.NewStudentRepository(dbPool)
	studentHandler := NewStudentHandler(studentRepo, dbPool)

	app := fiber.New(fiber.Config{
		AppName: "Tugas Mandiri - REST API Students (PostgreSQL)",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			pesan := "terjadi kesalahan pada server"
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
				pesan = e.Message
			}
			return fail(c, status, pesan)
		},
	})

	// Middleware global
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	// Endpoint Health Check terintegrasi database ping
	api.Get("/health", studentHandler.HealthCheck)

	// Routes Students dengan validasi Header Content-Type
	s := api.Group("/students", requireJSON)
	s.Get("/", studentHandler.ListStudents)
	s.Get("/:id", studentHandler.GetStudent)
	s.Post("/", studentHandler.CreateStudent)
	s.Put("/:id", studentHandler.ReplaceStudent)
	s.Patch("/:id", studentHandler.PatchStudent)
	s.Delete("/:id", studentHandler.DeleteStudent)

	// Route 404 handler untuk endpoint tak dikenal
	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	port := config.GetEnv("APP_PORT", "3000")
	fmt.Printf("Server berjalan di http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}