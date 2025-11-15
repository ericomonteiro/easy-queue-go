package handlers

import (
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserHandler gerencia as requisições HTTP relacionadas a usuários
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler cria uma nova instância de UserHandler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ErrorResponse representa uma resposta de erro padrão
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateUser godoc
// @Summary Cria um novo usuário
// @Description Cria um novo usuário no sistema com email, senha, telefone e role
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Validar role
	if req.Role != models.RoleBusinessOwner && req.Role != models.RoleCustomer {
		log.Warn(ctx, "Invalid user role", zap.String("role", string(req.Role)))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_role",
			Message: "Role must be 'BO' (Business Owner) or 'CU' (Customer)",
		})
		return
	}

	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		log.Error(ctx, "Failed to create user", zap.Error(err))
		
		// Verificar se é erro de duplicação
		if err.Error() == "user with email "+req.Email+" already exists" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "user_already_exists",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create user",
		})
		return
	}

	log.Info(ctx, "User created successfully via HTTP",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	c.JSON(http.StatusCreated, user)
}

// GetUserByID godoc
// @Summary Busca um usuário por ID
// @Description Retorna os dados de um usuário específico pelo ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Warn(ctx, "Invalid user ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid UUID format",
		})
		return
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Failed to get user", zap.Error(err), zap.String("id", id.String()))
		
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByEmail godoc
// @Summary Busca um usuário por email
// @Description Retorna os dados de um usuário específico pelo email
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/by-email [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	ctx := c.Request.Context()

	email := c.Query("email")
	if email == "" {
		log.Warn(ctx, "Email parameter is required")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_email",
			Message: "Email query parameter is required",
		})
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(ctx, "Failed to get user by email", zap.Error(err), zap.String("email", email))
		
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
