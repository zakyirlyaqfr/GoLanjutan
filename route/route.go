package route

import (
	// "golanjutan/app/model" // Tidak perlu lagi
	"golanjutan/app/repository"
	"golanjutan/app/service"
	"golanjutan/database"
	"golanjutan/middleware"

	// "strconv" // Tidak perlu lagi

	"github.com/gofiber/fiber/v2"
)

// Helper function tidak diperlukan lagi di file route.go
// func getUserFromContext(c *fiber.Ctx) (*model.User, error) { ... }

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// repositories
	alumniRepo := repository.NewAlumniRepository(database.DB)
	pekerjaanRepo := repository.NewPekerjaanRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)

	// services
	alumniSvc := service.NewAlumniService(alumniRepo)
	pekerjaanSvc := service.NewPekerjaanService(pekerjaanRepo)
	authService := service.NewAuthService(userRepo, alumniRepo)
	trashService := service.NewTrashService(alumniRepo, pekerjaanRepo)

	// ============================
	// AUTH ROUTES
	// ============================
	auth := api.Group("/auth")
	// Diubah: Langsung memanggil handler dari service
	auth.Post("/register", authService.HandleRegister)
	auth.Post("/login", authService.HandleLogin)

	// ============================
	// ALUMNI ROUTES
	// ============================
	alumni := api.Group("/alumni", middleware.Cors(), middleware.Protected())

	// Diubah: Semua handler inline dipindahkan ke alumniSvc
	alumni.Get("/filter", alumniSvc.HandleGetAllWithFilter)
	alumni.Get("/", alumniSvc.HandleGetAll)
	alumni.Get("/:id", alumniSvc.HandleGetByID)
	
	alumni.Post("/", middleware.RequireRole("admin"), alumniSvc.HandleCreate)
	alumni.Put("/:id", middleware.RequireRole("admin"), alumniSvc.HandleUpdate)
	
	alumni.Delete("/:id", alumniSvc.HandleSoftDelete) // Otorisasi (Superadmin) ada di dalam service
	
	// Hard delete dan restore dipisah group agar middleware RequireRole tidak bentrok
	alumniAdmin := api.Group("/alumni", middleware.Cors(), middleware.Protected())
	alumniAdmin.Delete("/harddelete/:id", alumniSvc.HandleHardDelete) // Otorisasi (Admin) ada di dalam service
	alumniAdmin.Patch("/:id/restore", alumniSvc.HandleRestore) // Otorisasi (Superadmin) ada di dalam service


	// ============================
	// PEKERJAAN ROUTES
	// ============================
	// (Ini sudah benar, tidak perlu diubah)
	pekerjaan := api.Group("/pekerjaan", middleware.Cors(), middleware.Protected())

	pekerjaan.Get("/filter", pekerjaanSvc.HandleGetAllWithFilter)
	pekerjaan.Get("/", pekerjaanSvc.HandleGetAll)
	pekerjaan.Get("/:id", pekerjaanSvc.HandleGetByID)
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanSvc.HandleGetByAlumniID)
	
	pekerjaan.Post("/", middleware.RequireRole("admin"), pekerjaanSvc.HandleCreate) // TODO: User juga harusnya bisa create/update/delete?
	pekerjaan.Put("/:id", middleware.RequireRole("admin"), pekerjaanSvc.HandleUpdate)
	
	pekerjaan.Delete("/:id", pekerjaanSvc.HandleSoftDelete) // Otorisasi (Admin/User pemilik) ada di dalam service
	pekerjaan.Delete("/harddelete/:id", pekerjaanSvc.HandleHardDelete) // Otorisasi (Admin/User pemilik) ada di dalam service
	pekerjaan.Patch("/:id/restore", pekerjaanSvc.HandleRestore) // Otorisasi (Admin/User pemilik) ada di dalam service

	// ============================
	// TRASH ROUTE
	// ============================
	// Diubah: Handler inline dipindahkan ke trashSvc
	api.Get("/trash", middleware.Protected(), trashService.HandleGetTrash)
}