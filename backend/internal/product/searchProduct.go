package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

// SearchProducts search products by text and refresh prices via scrapping with goroutines
func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El término de búsqueda no puede estar vacío"})
		return
	}

	var products []models.Product
	if err := database.DB.Where("LOWER(name) LIKE LOWER(?)", "%"+query+"%").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar productos"})
		return
	}

	// If there is no results, return
	if len(products) == 0 {
		c.JSON(http.StatusOK, products)
		return
	}

	// WaitGroup to sync goroutines
	var wg sync.WaitGroup

	// Iterate on products to scrap
	for i := range products {
		wg.Add(1)

		// One goroutine for each product
		go func(index int) {
			defer wg.Done()

			scrapedData, err := scraper.Extract(products[index].SourceURL)
			if err != nil {
				// If the scrap fails, we keep old price
				return
			}

			cleanPrice := strings.Replace(scrapedData.Price, ",", ".", 1)
			floatPrice, err := strconv.ParseFloat(cleanPrice, 64)
			if err != nil || floatPrice == 0 {
				return
			}

			hasChanged := false

			// Check if there is a different price
			if floatPrice != products[index].LastPrice {
				products[index].LastPrice = floatPrice
				hasChanged = true
			}

			// Check minimum price
			if products[index].LowestPrice == 0 || floatPrice < products[index].LowestPrice {
				products[index].LowestPrice = floatPrice
				hasChanged = true
			}

			// Price refresh if needed
			if hasChanged {
				products[index].UpdatedAt = time.Now()
				database.DB.Save(&products[index])
			}

			var lastHistory models.PriceHistory
			errHistory := database.DB.Where("product_id = ?", products[index].ProductID).Order("register_date desc").First(&lastHistory).Error

			if errHistory != nil || time.Since(lastHistory.RegisterDate) >= 12*time.Hour {
				newPrice := models.PriceHistory{
					ProductID:    products[index].ProductID,
					Price:        floatPrice,
					RegisterDate: time.Now(),
				}
				database.DB.Create(&newPrice)
			}
		}(i)
	}

	// Wait until all goroutines end
	wg.Wait()

	// Return the product list with refreshed prices
	c.JSON(http.StatusOK, products)
}

// GetProductByID obtiene los detalles de un producto específico
func GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, product)
}
