package handlers

import (
	"errors"
	"net/http"

	"menu-tree-system/backend/models"
	"menu-tree-system/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// GetMenus godoc
// @Summary      Get all menus in tree structure
// @Description  Retrieve full hierarchical tree of menu items
// @Tags         menus
// @Produce      json
// @Success      200  {array}   models.Menu
// @Failure      500  {object}  ErrorResponse
// @Router       /api/menus [get]
func (h *MenuHandler) GetMenus(c *gin.Context) {
	menus, err := h.service.GetMenuTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, menus)
}

// GetMenu godoc
// @Summary      Get a single menu item
// @Description  Get a menu item by ID with parent details
// @Tags         menus
// @Produce      json
// @Param        id   path      string  true  "Menu ID (UUID)"
// @Success      200  {object}  models.Menu
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/menus/{id} [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid UUID format"})
		return
	}

	menu, err := h.service.GetMenuItem(id)
	if err != nil {
		if errors.Is(err, services.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// CreateMenu godoc
// @Summary      Create a new menu item
// @Description  Create a new root or child menu item
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        menu  body      models.CreateMenuInput  true  "Menu Creation Input"
// @Success      201   {object}  models.Menu
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var input models.CreateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	menu, err := h.service.CreateMenuItem(input)
	if err != nil {
		if errors.Is(err, services.ErrParentNotFound) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu)
}

// UpdateMenu godoc
// @Summary      Update an existing menu item
// @Description  Update details of a menu item by ID
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "Menu ID (UUID)"
// @Param        menu  body      models.UpdateMenuInput  true  "Menu Update Input"
// @Success      200   {object}  models.Menu
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid UUID format"})
		return
	}

	var input models.UpdateMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	menu, err := h.service.UpdateMenuItem(id, input)
	if err != nil {
		if errors.Is(err, services.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// DeleteMenu godoc
// @Summary      Delete a menu item
// @Description  Delete a menu item and all its descendants (cascading)
// @Tags         menus
// @Produce      json
// @Param        id   path      string  true  "Menu ID (UUID)"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid UUID format"})
		return
	}

	if err := h.service.DeleteMenuItem(id); err != nil {
		if errors.Is(err, services.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// MoveMenu godoc
// @Summary      Move menu item to a different parent
// @Description  Reparent menu item and adjust depth of subtrees
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Menu ID (UUID)"
// @Param        body  body      models.MoveMenuInput  true  "Move Menu Input"
// @Success      200   {object}  models.Menu
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/menus/{id}/move [patch]
func (h *MenuHandler) MoveMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid UUID format"})
		return
	}

	var input models.MoveMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	menu, err := h.service.MoveMenuItem(id, input)
	if err != nil {
		if errors.Is(err, services.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, services.ErrCircularDependency) || errors.Is(err, services.ErrParentNotFound) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// ReorderMenu godoc
// @Summary      Reorder menu item within same parent level
// @Description  Change the display sort order among sibling items
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Menu ID (UUID)"
// @Param        body  body      models.ReorderMenuInput  true  "Reorder Menu Input"
// @Success      200   {object}  models.Menu
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/menus/{id}/reorder [patch]
func (h *MenuHandler) ReorderMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid UUID format"})
		return
	}

	var input models.ReorderMenuInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	menu, err := h.service.ReorderMenuItem(id, input)
	if err != nil {
		if errors.Is(err, services.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, services.ErrInvalidSortPosition) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}