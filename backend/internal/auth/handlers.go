package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest struct representing the form received
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest struct representing the form received
type RegisterRequest struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler checks if the user exists. If true, checks if the password is correct
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

	token, err := GenerateToken(user.UserID, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token de acceso"})
		return
	}

	// ==========================================
	// Dynamic cookie configuration
	// ==========================================
	esProduccion := os.Getenv("GIN_MODE") == "release"

	domain := ""
	secure := false
	sameSite := http.SameSiteLaxMode

	if esProduccion {
		domain = ""
		secure = true
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		MaxAge:   86400,
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login exitoso",
		"username": user.Username,
	})
}

// RegisterHandler check for duplicates and save new user
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

// UpdateProfileRequest structure for the name and email form
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Update Password Request structure for the password change form
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}

// GetProfileHandler returns the current user's information (without the password)
func GetProfileHandler(c *gin.Context) {
	// We retrieve the userID that AuthMiddleware() has saved in the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UserID":   user.UserID,
		"Username": user.Username,
		"Email":    user.Email,
	})
}

// UpdateProfileHandler modifies the Username and Email, validating duplicates.
func UpdateProfileHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Validation: Verify that the new email/username is not already in use by another person
	var existingUser models.User
	result := database.DB.Where("(email = ? OR username = ?) AND user_id != ?", req.Email, req.Username, userID.(uint)).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "El correo o nombre de usuario ya está registrado"})
		return
	}

	user.Username = req.Username
	user.Email = req.Email

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente"})
}

// UpdatePasswordHandler checks the old password and saves the new hashed password
func UpdatePasswordHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La nueva contraseña debe tener mínimo 6 caracteres"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "La contraseña actual es incorrecta"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la nueva credencial"})
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la nueva contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña modificada con éxito"})
}

// LogoutHandler destroys the browser's authentication cookie
func LogoutHandler(c *gin.Context) {
	isProduction := os.Getenv("GIN_MODE") == "release"
	domain := ""
	secure := false
	sameSite := http.SameSiteLaxMode

	if isProduction {
		domain = ""
		secure = true
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada correctamente"})
}
