package handlers

import (
	"errors"
	"net/http"

	"task-manager-api/internal/auth"
	"task-manager-api/internal/middleware"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service      *services.AuthService
	jwt          *auth.JWTManager
	cookieSecure bool
}

func NewAuthHandler(service *services.AuthService, jwt *auth.JWTManager, cookieSecure bool) *AuthHandler {
	return &AuthHandler{service: service, jwt: jwt, cookieSecure: cookieSecure}
}

type signupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindingError(c, err)
		return
	}

	user, token, err := h.service.Signup(req.Name, req.Email, req.Password)
	if errors.Is(err, services.ErrEmailTaken) {
		respondError(c, http.StatusConflict, "EMAIL_TAKEN", err.Error())
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Something went wrong")
		return
	}

	h.setAuthCookie(c, token)
	respondData(c, http.StatusCreated, authResponse{User: user, Token: token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindingError(c, err)
		return
	}

	user, token, err := h.service.Login(req.Email, req.Password)
	if errors.Is(err, services.ErrInvalidCredentials) {
		respondError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Something went wrong")
		return
	}

	h.setAuthCookie(c, token)
	respondData(c, http.StatusOK, authResponse{User: user, Token: token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "/", "", h.cookieSecure, true)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.service.GetUser(middleware.CurrentUserID(c))
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Account no longer exists")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Something went wrong")
		return
	}
	respondData(c, http.StatusOK, user)
}

func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, int(h.jwt.TTL().Seconds()), "/", "", h.cookieSecure, true)
}
