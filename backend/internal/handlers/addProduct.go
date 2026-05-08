package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

type TrackRequest struct {
	// Usamos binding:"required" de Gin para validar que siempre envíen la URL
	URL string `json:"url" binding:"required"`
}

func AddProduct(c *gin.Context) {
	var req TrackRequest

	// BindJSON lee el body y comprueba que traiga la URL
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debes proporcionar una URL válida"})
		return
	}

	// Llamamos a nuestro nuevo paquete scraper
	producto, err := scraper.Extract(req.URL)
	if err != nil {
		// ¡Cambiamos esta línea para que nos devuelva el motivo exacto del fallo!
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo extraer la información del producto",
			"detalle": err.Error(),
		})
		return
	}

	// Aquí en el futuro llamaremos a la base de datos usando tus modelos de GORM
	// ej: database.DB.Create(&models.Product{Title: producto.Title, ...})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Producto añadido y rastreado correctamente",
		"data":    producto,
	})
}
