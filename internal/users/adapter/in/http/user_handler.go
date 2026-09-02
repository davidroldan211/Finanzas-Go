package http

import (
	nethttp "net/http"
	"strconv"

	"finanzas-api/internal/httpx"
	"finanzas-api/internal/users/port/in"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService in.UserService
}

// NewUserHandler crea una nueva instancia del handler de usuarios.
func NewUserHandler(userService in.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser maneja la creación de nuevos usuarios.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req.toCommand())
	if err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    toUserResponse(user),
	})
}

// GetUser obtiene un usuario por ID.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httpx.BadRequest("ID de usuario inválido."))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"user": toUserResponse(user)})
}

// UpdateUser actualiza un usuario.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httpx.BadRequest("ID de usuario inválido."))
		return
	}

	var req UpdateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	user, err := h.userService.Update(c.Request.Context(), req.toCommand(id))
	if err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    toUserResponse(user),
	})
}

// DeleteUser elimina un usuario.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httpx.BadRequest("ID de usuario inválido."))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers obtiene una lista de usuarios.
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.userService.List(c.Request.Context(), in.ListUsersQuery{Limit: limit, Offset: offset})
	if err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	userResponses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, toUserResponse(user))
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"users": userResponses,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(userResponses),
		},
	})
}
