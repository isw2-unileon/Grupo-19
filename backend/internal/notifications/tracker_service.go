package handlers

import (
	"fmt"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// EvaluatePriceDropAndNotify analyzes the new price and decides to generate alerts based on user tracking.
func EvaluatePriceDropAndNotify(productID uint, currentPrice float64, oldPrice float64) error {
	if currentPrice >= oldPrice {
		return nil
	}

	// Take the product data
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return fmt.Errorf("error al buscar producto: %w", err)
	}

	// Search for all user alerts configured for this product
	var trackings []models.Tracking
	err := database.DB.Where("product_id = ?", productID).Find(&trackings).Error
	if err != nil {
		return fmt.Errorf("error buscando trackings: %v", err)
	}

	for _, track := range trackings {
		shouldNotify := false
		var reasonTitle string
		var reasonDesc string
		notificationType := "price_drop"

		// Any general price drop
		if track.NotifyPriceChanges {
			shouldNotify = true
			reasonTitle = "¡Bajada de precio detectada!"
			reasonDesc = fmt.Sprintf("El producto '%s' ha bajado de %.2f€ a %.2f€ (¡Ahorras %.2f€!).",
				product.Name, oldPrice, currentPrice, oldPrice-currentPrice)
		}

		// The current price breaks the user's target price barrier
		if track.NotifyTargetPrice > 0 && currentPrice <= track.NotifyTargetPrice {
			shouldNotify = true
			notificationType = "target_price"
			reasonTitle = "¡Precio objetivo alcanzado!"
			reasonDesc = fmt.Sprintf("¡Buenas noticias! '%s' ha caído hasta los %.2f€, alcanzando o mejorando el precio objetivo de %.2f€ que configuraste.",
				product.Name, currentPrice, track.NotifyTargetPrice)
		}

		if shouldNotify {
			notification := models.Notification{
				UserID:      track.UserID,
				ProductID:   productID,
				Type:        notificationType,
				Title:       reasonTitle,
				Description: reasonDesc,
				IsRead:      false,
				CreatedAt:   time.Now(),
			}

			// We create the record in the database
			if err := database.DB.Create(&notification).Error; err != nil {
				fmt.Printf("Error al generar alerta para usuario %d en producto %d: %v\n", track.UserID, productID, err)
			}
		}
	}

	return nil
}
