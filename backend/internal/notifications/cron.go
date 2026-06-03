package handlers

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

// init() function runs automatically when the server starts.
func init() {
	intervalo := 30 * time.Minute

	go func() {
		slog.Info("[CRON INTERNO] Esperando a que la base de datos esté lista...")

		// ACTIVE CHECK LOOP (Maximum 30 attempts, one per second)
		dbLista := false
		for i := 0; i < 30; i++ {
			// We check if the DB object has already been initialized by main.go and attempt a quick check.
			if database.DB != nil {
				var dePrueba int
				err := database.DB.Raw("SELECT 1").Scan(&dePrueba).Error
				if err == nil {
					dbLista = true
					break
				}
			}

			time.Sleep(1 * time.Second)
		}

		// After 30 seconds, we'll cancel for security reasons.
		if !dbLista {
			slog.Error("[CRON INTERNO] CRÍTICO: No se pudo conectar a la BBDD tras 30 segundos. Cron abortado.")
			return
		}

		slog.Info("[CRON INTERNO] ¡Base de datos detectada! El despertador automático se inicia ahora mismo.")

		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()

		for {
			checkSavedProductsPrices()
			<-ticker.C
		}
	}()
}

// checkSavedProductsPrices searches for saved products and checks if their prices have dropped.
func checkSavedProductsPrices() {
	slog.Info("[CRON] Pasada periódica: Comprobando precios de productos guardados...")

	if database.DB == nil {
		slog.Error("[CRON] La conexión a la base de datos aún no está lista")
		return
	}

	var productos []models.Product
	// searching all products that are registered in the system
	if err := database.DB.Find(&productos).Error; err != nil {
		slog.Error("[CRON] Error al recuperar productos de la BBDD", "error", err)
		return
	}

	for i := range productos {
		slog.Info("[CRON] Analizando artículo", "id", productos[i].ProductID, "name", productos[i].Name)

		// Run the scraper by passing it the URL
		scrapedData, err := scraper.Extract(productos[i].SourceURL)
		scrapedData = &scraper.ProductData{Price: "750.00"}
		if err != nil {
			slog.Error("[CRON] El scraper ha fallado para esta URL", "url", productos[i].SourceURL, "error", err)
			continue
		}

		// clean the price text.
		precioLimpio := strings.Replace(scrapedData.Price, ",", ".", 1)
		nuevoPrecio, err := strconv.ParseFloat(precioLimpio, 64)
		if err != nil || nuevoPrecio == 0 {
			slog.Error("[CRON] El precio extraído no es válido", "raw", scrapedData.Price)
			continue
		}

		precioViejo := productos[i].LastPrice
		haCambiado := false

		// If the scraper price is different, update the record
		if nuevoPrecio != precioViejo {
			productos[i].LastPrice = nuevoPrecio
			haCambiado = true
		}

		// We check if it's the all-time low
		if productos[i].LowestPrice == 0 || nuevoPrecio < productos[i].LowestPrice {
			productos[i].LowestPrice = nuevoPrecio
			haCambiado = true
		}

		// If the price has changed we save it
		if haCambiado {
			productos[i].UpdatedAt = time.Now()
			if err := database.DB.Save(&productos[i]).Error; err != nil {
				slog.Error("[CRON] Error al guardar el nuevo precio en la BBDD", "id", productos[i].ProductID, "error", err)
				continue
			}

			// If the price has dropped generate the notification
			if nuevoPrecio < precioViejo {
				slog.Info("⬇️ ¡Precio reducido detectado en segundo plano!", "product", productos[i].Name, "old", precioViejo, "new", nuevoPrecio)

				err := EvaluatePriceDropAndNotify(productos[i].ProductID, nuevoPrecio, precioViejo)
				if err != nil {
					slog.Error("[CRON] Error al procesar alertas de usuarios", "error", err)
				}
			}
		}

		time.Sleep(2 * time.Second)
	}

	slog.Info("[CRON] Pasada periódica finalizada con éxito.")
}
