package handlers

import (
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

// init() function runs automatically when the server starts.
func init() {
	// Avoids init when working on testing
	if strings.HasSuffix(os.Args[0], ".test") {
		return
	}

	interval := 1 * time.Hour

	go func() {
		slog.Info("[INTERNAL CRON] Waiting for the database to be ready...")

		// ACTIVE CHECK LOOP (Maximum 30 attempts, one per second)
		dbList := false
		for i := 0; i < 30; i++ {
			if database.DB != nil {
				var test int
				err := database.DB.Raw("SELECT 1").Scan(&test).Error
				if err == nil {
					dbList = true
					break
				}
			}
			time.Sleep(1 * time.Second)
		}

		if !dbList {
			slog.Error("[INTERNAL CHRON] CRITICAL: Could not connect to the DB after 30 seconds. Cron aborted.")
			return
		}

		slog.Info("[INTERNAL CHRON] Database detected! The automatic alarm clock starts right now.")

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			checkSavedProductsPrices()
			<-ticker.C
		}
	}()
}

// checkSavedProductsPrices It searches and analyzes the products that users actively follow.
func checkSavedProductsPrices() {
	slog.Info("[CRON] Periodic scan: Checking prices of active followed products...")

	if database.DB == nil {
		slog.Error("[CRON] Database connection not yet ready")
		return
	}

	var products []models.Product

	// Select unique products that have at least one associated record in the 'trackings' table
	err := database.DB.Distinct("products.*").
		Joins("JOIN trackings ON trackings.product_id = products.product_id").
		Find(&products).Error
	if err != nil {
		slog.Error("[CRON] Error when recovering active products from the DB", "error", err)
		return
	}

	slog.Info("[CRON]", "total_productos_a_scrapear", len(products))

	for i := range products {
		slog.Info("[CRON] Analyzing active item", "id", products[i].ProductID, "name", products[i].Name)

		// We run the scraper by passing it the URL
		scrapedData, err := scraper.Extract(products[i].SourceURL)
		if err != nil {
			slog.Error("[CRON] The scraper has failed for this URL", "url", products[i].SourceURL, "error", err)
			// Dejamos también un pequeño respiro de cortesía tras un error por seguridad
			time.Sleep(2 * time.Second)
			continue
		}

		// Clean up the formatting of the price text.
		cleanPrice := strings.Replace(scrapedData.Price, ",", ".", 1)
		newPrice, err := strconv.ParseFloat(cleanPrice, 64)
		if err != nil || newPrice == 0 {
			slog.Error("[CRON] Extracted price is not valid", "raw", scrapedData.Price)
			time.Sleep(2 * time.Second)
			continue
		}

		oldPrice := products[i].LastPrice
		hasChanged := false

		// If the price on the website is different from the saved price, we update it.
		if newPrice != oldPrice {
			products[i].LastPrice = newPrice
			hasChanged = true
		}

		// We check if it is the all-time low
		if products[i].LowestPrice == 0 || newPrice < products[i].LowestPrice {
			products[i].LowestPrice = newPrice
			hasChanged = true
		}

		// If there have been any structural changes, we save them in the database.
		if hasChanged {
			products[i].UpdatedAt = time.Now()
			if err := database.DB.Save(&products[i]).Error; err != nil {
				slog.Error("[CRON] Error al guardar el nuevo precio en la BBDD", "id", products[i].ProductID, "error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			// If the website price is lower than the price of our database, we trigger alerts
			if newPrice < oldPrice {
				slog.Info("Reduced price detected in the background!", "product", products[i].Name, "old", oldPrice, "new", newPrice)

				err := EvaluatePriceDropAndNotify(products[i].ProductID, newPrice, oldPrice)
				if err != nil {
					slog.Error("[CRON] Error processing user alerts", "error", err)
				}
			}
		}

		//nolint:gosec
		randomTime := 30 + rand.IntN(20)
		time.Sleep(time.Duration(randomTime) * time.Second)
	}

	slog.Info("[CRON] Periodic scan completed successfully.")
}
