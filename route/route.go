package route

import (
	"golanjutan/app/model"
	"golanjutan/app/repository"
	"golanjutan/app/service"
	"golanjutan/database"
	"golanjutan/middleware"

	"strconv"
	// "time"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// repositories
	alumniRepo := repository.NewAlumniRepository(database.DB)
	pekerjaanRepo := repository.NewPekerjaanRepository(database.DB)

	// services
	alumniSvc := service.NewAlumniService(alumniRepo)
	pekerjaanSvc := service.NewPekerjaanService(pekerjaanRepo)

	// alumni routes
	alumni := api.Group("/alumni", middleware.Cors())
	alumni.Get("/", func(c *fiber.Ctx) error {
		res, err := alumniSvc.GetAll()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "id invalid")
		}
		res, err := alumniSvc.GetByID(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "alumni not found")
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})
	alumni.Post("/", func(c *fiber.Ctx) error {
		var req model.CreateAlumniRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		id, err := alumniSvc.Create(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		newAlumni, _ := alumniSvc.GetByID(id)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newAlumni})
	})
	alumni.Put("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "id invalid")
		}
		var req model.UpdateAlumniRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if err := alumniSvc.Update(id, req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		updated, _ := alumniSvc.GetByID(id)
		return c.JSON(fiber.Map{"success": true, "data": updated})
	})
	alumni.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "id invalid")
		}
		if err := alumniSvc.Delete(id); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "message": "alumni deleted"})
	})

	// pekerjaan routes
	pekerjaan := api.Group("/pekerjaan", middleware.Cors())
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		res, err := pekerjaanSvc.GetAll()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "id invalid")
		}
		res, err := pekerjaanSvc.GetByID(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "pekerjaan not found")
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})
	pekerjaan.Get("/alumni/:alumni_id", func(c *fiber.Ctx) error {
		alumniID, err := strconv.Atoi(c.Params("alumni_id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "alumni_id invalid")
		}
		res, err := pekerjaanSvc.GetByAlumniID(alumniID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "data": res})
	})
	pekerjaan.Post("/", func(c *fiber.Ctx) error {
		var req model.CreatePekerjaanRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		id, err := pekerjaanSvc.Create(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		newPekerjaan, _ := pekerjaanSvc.GetByID(id)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newPekerjaan})
	})
	pekerjaan.Put("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "id invalid")
		}
		var req model.UpdatePekerjaanRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if err := pekerjaanSvc.Update(id, req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		updated, _ := pekerjaanSvc.GetByID(id)
		return c.JSON(fiber.Map{"success": true, "data": updated})
	})
	pekerjaan.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "id invalid")
		}
		if err := pekerjaanSvc.Delete(id); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"success": true, "message": "pekerjaan deleted"})
	})
}
