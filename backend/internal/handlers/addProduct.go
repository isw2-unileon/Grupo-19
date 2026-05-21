package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

type TrackRequest struct {
	URL string `json:"url" binding:"required"`
}

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

	precioLimpio := strings.Replace(product.Price, ",", ".", 1)
	precioFloat, err := strconv.ParseFloat(precioLimpio, 64)
	if err != nil {
		precioFloat = 0.0
	}

	// ID de usuario falso, a cambiar mas adelante
	var usuarioID uint = 1

	var producto models.Product
	resultadoBusqueda := database.DB.Where("source_url = ?", product.URL).First(&producto)

	if resultadoBusqueda.Error != nil {
		producto = models.Product{
			Name:        product.Title,
			SourceURL:   product.URL,
			LastPrice:   precioFloat,
			LowestPrice: precioFloat,
			CreatedBy:   usuarioID,
			CreateAt:    time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := database.DB.Create(&producto).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el producto"})
			return
		}
	} else {
		producto.LastPrice = precioFloat
		producto.UpdatedAt = time.Now()

		if precioFloat > 0 && (producto.LowestPrice == 0 || precioFloat < producto.LowestPrice) {
			producto.LowestPrice = precioFloat
		}

		database.DB.Save(&producto)
	}

	nuevoPrecio := models.PriceHistory{
		ProductID:    producto.ProductID,
		Price:        precioFloat,
		RegisterDate: time.Now(),
	}

	if err := database.DB.Create(&nuevoPrecio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el historial de precios"})
		return
	}

	var tracking models.Tracking
	errTracking := database.DB.Where("user_id = ? AND product_id = ?", usuarioID, producto.ProductID).First(&tracking).Error

	if errTracking != nil {
		tracking = models.Tracking{
			UserID:             usuarioID,
			ProductID:          producto.ProductID,
			NotifyPriceChanges: true,
			NotifyTargetPrice:  0.0,
			TrackingStartDate:  time.Now(),
		}
		database.DB.Create(&tracking)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Operación realizada con éxito en PostgreSQL",
		"data": gin.H{
			"product_id":   producto.ProductID,
			"name":         producto.Name,
			"source_url":   producto.SourceURL,
			"last_price":   producto.LastPrice,
			"lowest_price": producto.LowestPrice,
			"updated_at":   producto.UpdatedAt,
		},
	})
}
