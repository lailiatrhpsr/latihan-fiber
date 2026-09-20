package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"latihan-fiber/pert4/app/repository"
	"latihan-fiber/pert4/app/service"
	"latihan-fiber/pert4/config"
	"latihan-fiber/pert4/database"
	"latihan-fiber/pert4/helper"
	"latihan-fiber/pert4/route"
)

const minSecretLength = 32

func main() {
	// 1. Konfigurasi dan logger
	config.LoadEnv()
	logger := config.NewLogger()

	// 2. JWT diperiksa SEBELUM server menyala
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if len(jwtSecret) < minSecretLength {
		logger.Error("JWT_SECRET tidak diisi atau terlalu pendek",
			slog.Int("minimal_karakter", minSecretLength))
		os.Exit(1)
	}

	// 3. Database
	pool := database.NewPostgresPool()
	defer pool.Close()

	// 4. Perakitan dari dalam ke luar: repository -> service
	studentRepository := repository.NewStudentRepository(pool)
	userRepository := repository.NewUserRepository(pool)
	tokenRepository := repository.NewTokenRepository(pool)
	roleRepository := repository.NewRoleRepository(pool)

	rawPermissions, err := roleRepository.LoadPermissions(context.Background())
	if err != nil {
		logger.Error("gagal memuat permission", slog.String("error", err.Error()))
		os.Exit(1)
	}
	permissions := helper.NewPermissionSet(rawPermissions)
	logger.Info("permission dimuat", slog.Any("roles", permissions.KnownRoles()))
	
	studentService := service.NewStudentService(studentRepository, pool, permissions)

	jwtManager := helper.NewJWTManager(
		jwtSecret,
		config.GetEnv("JWT_ISSUER", "praktikum-backend"),
		time.Duration(config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15))*time.Minute,
	)
	authService := service.NewAuthService(
		userRepository, tokenRepository, jwtManager,
		time.Duration(config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7))*24*time.Hour,
	)

	// 5. Aplikasi
	app := config.NewApp(logger, route.Dependencies{
		Pool:           pool,
		JWT:            jwtManager,
		StudentService: studentService,
		AuthService:    authService,
		Permissions:    permissions,
	})
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("server berjalan", slog.String("port", port))

	// 6.  beri waktu request untuk selesai
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("sinyal berhenti diterima, menutup server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("gagal menutup server dengan rapi", slog.String("error", err.Error()))
	}
	logger.Info("server berhenti dengan rapi")
}
