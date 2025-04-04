package legalentities

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c echo.Context) error {
	entities, err := h.service.GetAllLegalEntities()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, entities)
}

func (h *Handler) Create(c echo.Context) error {
	var entity LegalEntity
	if err := c.Bind(&entity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	if err := h.service.CreateLegalEntity(&entity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, entity)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("uuid")

	var input struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	entity := &LegalEntity{
		UUID: id,
		Name: input.Name,
	}

	if err := h.service.UpdateLegalEntity(entity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, entity)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("uuid")
	if err := h.service.DeleteLegalEntity(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
