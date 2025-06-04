package handler

import (
	"finanzas-api/internal/users/domain"
	"finanzas-api/internal/users/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(r *gin.Engine, uc *usecase.UserUseCase) *UserHandler {
	handler := &UserHandler{uc: uc}

	r.GET("/users", handler.GetAll)
	r.GET("/users/:id", handler.GetByID)
	r.POST("/users", handler.Create)
	r.PUT("/users/:id", handler.Update)
	r.DELETE("/users/:id", handler.Delete)

	return handler
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, _ := h.uc.GetAll()
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, _ := h.uc.Create(user)
	c.JSON(http.StatusCreated, created)
}

func (h *UserHandler) Update(c *gin.Context) {
	var user domain.User
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	updated, err := h.uc.Update(user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
