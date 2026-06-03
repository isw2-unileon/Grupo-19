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
	if err := database.DB.Where("name ILIKE ?", "%"+query+"%").Find(&products).Error; err != nil {
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

			precioLimpio := strings.Replace(scrapedData.Price, ",", ".", 1)
			precioFloat, err := strconv.ParseFloat(precioLimpio, 64)
			if err != nil || precioFloat == 0 {
				return
			}

			haCambiado := false

			// Check if there is a different price
			if precioFloat != products[index].LastPrice {
				products[index].LastPrice = precioFloat
				haCambiado = true
			}

			// Check minimum price
			if products[index].LowestPrice == 0 || precioFloat < products[index].LowestPrice {
				products[index].LowestPrice = precioFloat
				haCambiado = true
			}

			// Price refresh if needed
			if haCambiado {
				products[index].UpdatedAt = time.Now()
				database.DB.Save(&products[index])
			}

			var ultimoHistorial models.PriceHistory
			errHistorial := database.DB.Where("product_id = ?", products[index].ProductID).Order("register_date desc").First(&ultimoHistorial).Error

			if errHistorial != nil || time.Since(ultimoHistorial.RegisterDate) >= 12*time.Hour {
				nuevoPrecio := models.PriceHistory{
					ProductID:    products[index].ProductID,
					Price:        precioFloat,
					RegisterDate: time.Now(),
				}
				database.DB.Create(&nuevoPrecio)
			}
		}(i)
	}

	// Wait until all goroutines end
	wg.Wait()

	// Return the product list with refreshed prices
	c.JSON(http.StatusOK, products)
}
