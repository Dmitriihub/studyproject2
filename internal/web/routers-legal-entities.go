package web

import (
	"net/http"

	"github.com/krisch/crm-backend/internal/legalentities"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type LegalEntityHandler struct {
	service legalentities.Service
}

func NewLegalEntityHandler(service legalentities.Service) *LegalEntityHandler {
	return &LegalEntityHandler{service: service}
}

func initOpenAPILegalEntitiesRouters(a *Web, e *echo.Echo) {
	logrus.WithField("route", "legal-entities").Debug("routes initialization")

	handler := NewLegalEntityHandler(a.app.LegalEntities)
	handler.RegisterRoutes(e)
}

func (h *LegalEntityHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/legal-entities")
	group.GET("", h.getAll)
	group.POST("", h.create)
	group.PUT("/:uuid", h.update)
	group.DELETE("/:uuid", h.delete)
}

func (h *LegalEntityHandler) getAll(c echo.Context) error {
	entities, err := h.service.GetAllLegalEntities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, entities)
}

func (h *LegalEntityHandler) create(c echo.Context) error {
	var entity legalentities.LegalEntity
	if err := c.Bind(&entity); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.service.CreateLegalEntity(&entity); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, entity)
}

func (h *LegalEntityHandler) update(c echo.Context) error {
	uuid := c.Param("uuid")
	var entity legalentities.LegalEntity
	if err := c.Bind(&entity); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	entity.UUID = uuid
	if err := h.service.UpdateLegalEntity(&entity); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, entity)
}

func (h *LegalEntityHandler) delete(c echo.Context) error {
	uuid := c.Param("uuid")
	if err := h.service.DeleteLegalEntity(uuid); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
