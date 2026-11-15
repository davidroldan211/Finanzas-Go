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

// registrationAcceptedMessage es la única respuesta de un registro con
// datos válidos, se haya creado el usuario o no (A5).
const registrationAcceptedMessage = "Registro recibido. Revisa tu correo para confirmarlo."

// CreateUser es el registro público. Con datos válidos responde siempre 202
// con el mismo cuerpo: no devuelve el usuario ni dice si el correo existía.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	if err := h.userService.Create(c.Request.Context(), req.toCommand()); err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusAccepted, gin.H{"message": registrationAcceptedMessage})
}

// GetUser obtiene un usuario por ID.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httpx.BadRequest("ID de usuario inválido."))
		return
	}

	actor, ok := actorFrom(c)
	if !ok {
		httpx.Abort(c, httpx.Unauthorized(""))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), actor, id)
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

	actor, ok := actorFrom(c)
	if !ok {
		httpx.Abort(c, httpx.Unauthorized(""))
		return
	}

	var req UpdateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	user, err := h.userService.Update(c.Request.Context(), actor, req.toCommand(id))
	if err != nil {
		httpx.Abort(c, toAppError(err))
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    toUserResponse(user),
	})
}

// actorFrom lee la identidad que el AuthGuard dejó en el contexto. Devuelve
// false si la ruta no pasó por el guard o los datos no son válidos.
func actorFrom(c *gin.Context) (in.Actor, bool) {
	id, err := uuid.Parse(c.GetString(httpx.CtxUserID))
	role := c.GetString(httpx.CtxUserRole)
	if err != nil || id == uuid.Nil || role == "" {
		return in.Actor{}, false
	}
	return in.Actor{ID: id, Role: role}, true
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
