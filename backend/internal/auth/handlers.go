package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Checks if the user exists. If true, checks if the password is correct
func LoginHandler(c *gin.Context) {
	var req LoginRequest

	// Decode json from frontend
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var user models.User

	// Search for email
	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login exitoso"})
}

// Check for duplicates and save new user
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	// Decode json from frontend
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Check if username is already used
	var existingUser models.User
	result := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "El correo o nombre de usuario ya existe"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	// If username is not used, we save the new user
	newUser := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashedPassword),
		UserType:   "user", // Default UserType
		RegisterAt: time.Now(),
	}

	// Save in DB
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente"})
}
