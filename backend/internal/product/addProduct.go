//nolint:misspell
package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/auth"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

// TrackRequest represent the url of the product we want to track
type TrackRequest struct {
	URL string `json:"url" binding:"required"`
}

// AddProduct add an item to the system
func AddProduct(c *gin.Context) {
	var req TrackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debes proporcionar una URL válida"})
		return
	}

	product, err := scraper.Extract(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo extraer la información del producto",
			"detalle": err.Error(),
		})
		return
	}

	cleanPrice := strings.Replace(product.Price, ",", ".", 1)
	floatPrice, err := strconv.ParseFloat(cleanPrice, 64)
	if err != nil {
		floatPrice = 0.0
	}

	// User ID
	var userID uint

	tokenString, errCookie := c.Cookie("auth_token")
	if errCookie == nil && tokenString != "" {
		claims, err := auth.ValidateToken(tokenString)
		if err == nil && claims != nil {
			userID = claims.UserID
		}
	}

	var producto models.Product
	resultadoBusqueda := database.DB.Where("source_url = ?", product.URL).First(&producto)

	if resultadoBusqueda.Error != nil {
		producto = models.Product{
			Name:        product.Title,
			SourceURL:   product.URL,
			LastPrice:   floatPrice,
			LowestPrice: floatPrice,
			CreatedBy:   userID,
			CreateAt:    time.Now(),
			UpdatedAt:   time.Now(),
			ImageURL:    product.ImageURL,
			Description: product.Description,
		}

		if err := database.DB.Create(&producto).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el producto"})
			return
		}
	} else {
		producto.LastPrice = floatPrice
		producto.UpdatedAt = time.Now()

		if product.Title != "" {
			producto.Name = product.Title
		}
		if product.ImageURL != "" {
			producto.ImageURL = product.ImageURL
		}
		if product.Description != "" {
			producto.Description = product.Description
		}

		if floatPrice > 0 && (producto.LowestPrice == 0 || floatPrice < producto.LowestPrice) {
			producto.LowestPrice = floatPrice
		}

		database.DB.Save(&producto)
	}

	var lastHistory models.PriceHistory
	errHistorial := database.DB.Where("product_id = ?", producto.ProductID).Order("register_date desc").First(&lastHistory).Error

	if errHistorial != nil || time.Since(lastHistory.RegisterDate) >= 12*time.Hour {
		nuevoPrecio := models.PriceHistory{
			ProductID:    producto.ProductID,
			Price:        floatPrice,
			RegisterDate: time.Now(),
		}

		if err := database.DB.Create(&nuevoPrecio).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el historial de precios"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Operación realizada con éxito en PostgreSQL",
		"data":    producto, // Return the model
	})
}

// GetPriceHistory returns all the priceHistory entries on the db for a product
func GetPriceHistory(c *gin.Context) {
	productID := c.Param("id")
	daysStr := c.DefaultQuery("days", "7")

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 7
	}

	limitDate := time.Now().AddDate(0, 0, -days)

	var history []models.PriceHistory
	database.DB.Where("product_id = ? AND register_date >= ?", productID, limitDate).
		Order("register_date asc").Find(&history)

	c.JSON(http.StatusOK, history)
}
