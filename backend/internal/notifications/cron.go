package handlers

import (
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/scraper"
)

// init() function runs automatically when the server starts.
func init() {
	// Avoid big work charge for the server just when booting up
	time.Sleep(3 * time.Minute)
	intervalo := 1 * time.Hour

	go func() {
		slog.Info("[CRON INTERNO] Esperando a que la base de datos esté lista...")

		// ACTIVE CHECK LOOP (Maximum 30 attempts, one per second)
		dbLista := false
		for i := 0; i < 30; i++ {
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

// checkSavedProductsPrices It searches and analyzes the products that users actively follow.
func checkSavedProductsPrices() {
	slog.Info("[CRON] Pasada horaria: Comprobando precios de productos activos en seguimiento...")

	if database.DB == nil {
		slog.Error("[CRON] La conexión a la base de datos aún no está lista")
		return
	}

	var productos []models.Product

	// Select unique products that have at least one associated record in the 'trackings' table
	err := database.DB.Distinct("products.*").
		Joins("JOIN trackings ON trackings.product_id = products.product_id").
		Find(&productos).Error
	if err != nil {
		slog.Error("[CRON] Error al recuperar productos activos de la BBDD", "error", err)
		return
	}

	slog.Info("[CRON]", "total_productos_a_scrapear", len(productos))

	for i := range productos {
		slog.Info("[CRON] Analizando artículo activo", "id", productos[i].ProductID, "name", productos[i].Name)

		// We run the scraper by passing it the URL
		scrapedData, err := scraper.Extract(productos[i].SourceURL)
		if err != nil {
			slog.Error("[CRON] El scraper ha fallado para esta URL", "url", productos[i].SourceURL, "error", err)
			// Dejamos también un pequeño respiro de cortesía tras un error por seguridad
			time.Sleep(2 * time.Second)
			continue
		}

		// Clean up the formatting of the price text.
		precioLimpio := strings.Replace(scrapedData.Price, ",", ".", 1)
		nuevoPrecio, err := strconv.ParseFloat(precioLimpio, 64)
		if err != nil || nuevoPrecio == 0 {
			slog.Error("[CRON] El precio extraído no es válido", "raw", scrapedData.Price)
			time.Sleep(2 * time.Second)
			continue
		}

		precioViejo := productos[i].LastPrice
		haCambiado := false

		// If the price on the website is different from the saved price, we update it.
		if nuevoPrecio != precioViejo {
			productos[i].LastPrice = nuevoPrecio
			haCambiado = true
		}

		// We check if it is the all-time low
		if productos[i].LowestPrice == 0 || nuevoPrecio < productos[i].LowestPrice {
			productos[i].LowestPrice = nuevoPrecio
			haCambiado = true
		}

		// If there have been any structural changes, we save them in the database.
		if haCambiado {
			productos[i].UpdatedAt = time.Now()
			if err := database.DB.Save(&productos[i]).Error; err != nil {
				slog.Error("[CRON] Error al guardar el nuevo precio en la BBDD", "id", productos[i].ProductID, "error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			// If the website price is lower than the price of our database, we trigger alerts
			if nuevoPrecio < precioViejo {
				slog.Info("¡Precio reducido detectado en segundo plano!", "product", productos[i].Name, "old", precioViejo, "new", nuevoPrecio)

				err := EvaluatePriceDropAndNotify(productos[i].ProductID, nuevoPrecio, precioViejo)
				if err != nil {
					slog.Error("[CRON] Error al procesar alertas de usuarios", "error", err)
				}
			}
		}

		//nolint:gosec
		tiempoAleatorio := 3 + rand.IntN(6)
		time.Sleep(time.Duration(tiempoAleatorio) * time.Second)
	}

	slog.Info("[CRON] Pasada horaria finalizada con éxito.")
}
