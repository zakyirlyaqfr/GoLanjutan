package route

import (
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"golanjutan/app/service"
	"golanjutan/database"
	"golanjutan/middleware"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// repositories
	alumniRepo := repository.NewAlumniRepository(database.DB)
	pekerjaanRepo := repository.NewPekerjaanRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)

	// services
	alumniSvc := service.NewAlumniService(alumniRepo)
	pekerjaanSvc := service.NewPekerjaanService(pekerjaanRepo)
	authService := service.NewAuthService(*userRepo)

	// ============================
	// AUTH ROUTES
	// ============================
	auth := api.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		res, err := authService.Login(req.Username, req.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})

	// ============================
	// ALUMNI ROUTES
	// ============================
	alumni := api.Group("/alumni", middleware.Cors())

	// FILTER, PAGINATION, SEARCH, SORT (harus didefinisikan duluan)
	alumni.Get("/filter", middleware.Protected(), func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		sortBy := c.Query("sortBy", "created_at")
		sortOrder := c.Query("sortOrder", "DESC")
		search := c.Query("search", "")

		res, err := alumniSvc.GetAllWithFilter(page, limit, sortBy, sortOrder, search)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true, "data": res.Data, "meta": res.Meta})
	})

	alumni.Get("/", middleware.Protected(), func(c *fiber.Ctx) error {
		res, err := alumniSvc.GetAll()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})

	alumni.Get("/:id", middleware.Protected(), func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
		}
		res, err := alumniSvc.GetByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "alumni not found"})
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})

	alumni.Post("/", middleware.Protected(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		var req model.CreateAlumniRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
		}
		id, err := alumniSvc.Create(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		newAlumni, _ := alumniSvc.GetByID(id)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newAlumni})
	})

	alumni.Put("/:id", middleware.Protected(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
		}
		var req model.UpdateAlumniRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
		}
		if err := alumniSvc.Update(id, req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		updated, _ := alumniSvc.GetByID(id)
		return c.JSON(fiber.Map{"success": true, "data": updated})
	})

	alumni.Delete("/:id", middleware.Protected(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
		}
		if err := alumniSvc.Delete(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true, "message": "alumni deleted"})
	})

	// ============================
	// PEKERJAAN ROUTES
	// ============================
	pekerjaan := api.Group("/pekerjaan", middleware.Cors())

	// FILTER, PAGINATION, SEARCH, SORT (harus didefinisikan duluan)
	pekerjaan.Get("/filter", middleware.Protected(), func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		sortBy := c.Query("sortBy", "created_at")
		sortOrder := c.Query("sortOrder", "DESC")
		search := c.Query("search", "")

		res, err := pekerjaanSvc.GetAllWithFilter(page, limit, sortBy, sortOrder, search)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true, "data": res.Data, "meta": res.Meta})
	})

	pekerjaan.Get("/", middleware.Protected(), func(c *fiber.Ctx) error {
		res, err := pekerjaanSvc.GetAll()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})

	pekerjaan.Get("/:id", middleware.Protected(), func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
		}
		res, err := pekerjaanSvc.GetByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "pekerjaan not found"})
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})

	pekerjaan.Get("/alumni/:alumni_id", middleware.Protected(), func(c *fiber.Ctx) error {
		alumniID, err := strconv.Atoi(c.Params("alumni_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "alumni_id invalid"})
		}
		res, err := pekerjaanSvc.GetByAlumniID(alumniID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})

	pekerjaan.Post("/", middleware.Protected(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		var req model.CreatePekerjaanRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
		}
		id, err := pekerjaanSvc.Create(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		newPekerjaan, _ := pekerjaanSvc.GetByID(id)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newPekerjaan})
	})

	pekerjaan.Put("/:id", middleware.Protected(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
		}
		var req model.UpdatePekerjaanRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid body"})
		}
		if err := pekerjaanSvc.Update(id, req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		updated, _ := pekerjaanSvc.GetByID(id)
		return c.JSON(fiber.Map{"success": true, "data": updated})
	})

	pekerjaan.Delete("/:id", middleware.Protected(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "id invalid"})
		}
		if err := pekerjaanSvc.Delete(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true, "message": "pekerjaan deleted"})
	})
}
