package handler

import (
	"fmt"
	"mobile-banking/internal/models"
	"mobile-banking/internal/service"
	"mobile-banking/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler{
	return &UserHandler{userService: userService}
}

type RegisterRequest struct{
	Username 	string 	`json:"username" binding:"required,min=3,max=50"`
	Email		string	`json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct{
	Username 	string 	`json:"username" binding:"required"`
	Password 	string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context){
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
	}

	fmt.Println("ini handler register")

	registeredUser, err := h.userService.Register(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, registeredUser)

}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	
	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Attempt login
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *UserHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			c.Abort()
			return
		}

		if len(authHeader) < len(BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization token",
			})
			c.Abort()
			return
		}

		token := authHeader[len(BearerSchema):]
		
		// Verify token
		claims, err := utils.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		
		c.Next()
	}
}
