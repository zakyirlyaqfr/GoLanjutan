// golanjutan/route/route.go
package route

import (
	"golanjutan/app/repository"
	"golanjutan/app/service"
	"golanjutan/database"
	"golanjutan/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// ============================
	// REPOSITORIES
	// ============================
	alumniRepo := repository.NewAlumniRepository(database.DB)
	pekerjaanRepo := repository.NewPekerjaanRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	// BARU: File Repository
	fileRepo := repository.NewFileRepository(database.DB)

	// ============================
	// SERVICES
	// ============================
	alumniSvc := service.NewAlumniService(alumniRepo)
	pekerjaanSvc := service.NewPekerjaanService(pekerjaanRepo)
	authService := service.NewAuthService(userRepo, alumniRepo)
	trashService := service.NewTrashService(alumniRepo, pekerjaanRepo)
	// BARU: File Service
	fileService := service.NewFileService(fileRepo, "./uploads")

	// ============================
	// AUTH ROUTES
	// ============================
	auth := api.Group("/auth")
	auth.Post("/register", authService.HandleRegister)
	auth.Post("/login", authService.HandleLogin)

	// ============================
	// ALUMNI ROUTES
	// ============================
	alumni := api.Group("/alumni", middleware.Cors(), middleware.Protected())
	alumni.Get("/filter", alumniSvc.HandleGetAllWithFilter)
	alumni.Get("/", alumniSvc.HandleGetAll)
	alumni.Get("/:id", alumniSvc.HandleGetByID)
	alumni.Post("/", middleware.RequireRole("admin"), alumniSvc.HandleCreate)
	alumni.Put("/:id", middleware.RequireRole("admin"), alumniSvc.HandleUpdate)
	alumni.Delete("/:id", alumniSvc.HandleSoftDelete)

	alumniAdmin := api.Group("/alumni", middleware.Cors(), middleware.Protected())
	alumniAdmin.Delete("/harddelete/:id", alumniSvc.HandleHardDelete)
	alumniAdmin.Patch("/:id/restore", alumniSvc.HandleRestore)

	// ============================
	// PEKERJAAN ROUTES
	// ============================
	pekerjaan := api.Group("/pekerjaan", middleware.Cors(), middleware.Protected())
	pekerjaan.Get("/filter", pekerjaanSvc.HandleGetAllWithFilter)
	pekerjaan.Get("/", pekerjaanSvc.HandleGetAll)
	pekerjaan.Get("/:id", pekerjaanSvc.HandleGetByID)
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanSvc.HandleGetByAlumniID)
	pekerjaan.Post("/", middleware.RequireRole("admin"), pekerjaanSvc.HandleCreate)
	pekerjaan.Put("/:id", middleware.RequireRole("admin"), pekerjaanSvc.HandleUpdate)
	pekerjaan.Delete("/:id", pekerjaanSvc.HandleSoftDelete)
	pekerjaan.Delete("/harddelete/:id", pekerjaanSvc.HandleHardDelete)
	pekerjaan.Patch("/:id/restore", pekerjaanSvc.HandleRestore)

	// ============================
	// TRASH ROUTE
	// ============================
	api.Get("/trash", middleware.Protected(), trashService.HandleGetTrash)

	// ============================
	// BARU: FILE UPLOAD ROUTES
	// ============================
	files := api.Group("/files", middleware.Cors(), middleware.Protected())

	// 1. Rute untuk User/Admin upload untuk dirinya sendiri
	// URL: POST /api/files/upload
	// User ID akan diambil dari token.
	files.Post("/upload", fileService.HandleUpload)

	// 2. Rute untuk Admin upload untuk user lain
	// URL: POST /api/files/upload/:id
	// ID adalah ID alumni/user yang dituju. Otorisasi (admin) dilakukan di HandleUpload.
	files.Post("/upload/:id", fileService.HandleUpload)

	// 3. Endpoint tambahan lainnya
	files.Get("/", middleware.RequireRole("admin"), fileService.GetAllFiles)
	files.Get("/:id", fileService.GetFileByID)
	files.Delete("/:id", middleware.RequireRole("admin"), fileService.HandleDeleteFile)
}