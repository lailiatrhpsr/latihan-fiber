package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihan-fiber/pert4/app/service"
	"latihan-fiber/pert4/helper"
	"latihan-fiber/pert4/middleware"
)

type Dependencies struct {
	Pool           *pgxpool.Pool
	JWT            *helper.JWTManager
	StudentService *service.StudentService
	AuthService    *service.AuthService
}
	
func Register(app *fiber.App, deps Dependencies) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")
	
	// --- publik ---
	api.Get("/health", deps.StudentService.HealthCheck)

	// --- autentikasi ---
	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", deps.AuthService.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
	auth.Post("/refresh", deps.AuthService.Refresh)
	auth.Post("/logout", deps.AuthService.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)


	students := api.Group("/students",
		middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	students.Get("/", deps.StudentService.List)
	students.Get("/:id", deps.StudentService.Get)
	students.Post("/", deps.StudentService.Create)
	students.Put("/:id", deps.StudentService.Replace)
	students.Patch("/:id", deps.StudentService.Patch)
	students.Delete("/:id", deps.StudentService.Delete)
}