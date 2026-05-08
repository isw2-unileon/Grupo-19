package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProductData struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	ImageURL    string `json:"image_url"`
}

func Extract(targetURL string) (*ProductData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	// Amazon bloquea agresivamente. A veces añadir el encabezado "Accept-Language" ayuda.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error de servidor: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	product := &ProductData{URL: targetURL}

	// --- 1. TÍTULO ---
	if title, exists := doc.Find("meta[property='og:title']").Attr("content"); exists && title != "" {
		product.Title = title
	} else if title := doc.Find("#productTitle").Text(); title != "" { // Amazon
		product.Title = title
	} else {
		product.Title = doc.Find("title").Text()
	}

	// --- 2. IMAGEN ---
	if img, exists := doc.Find("meta[property='og:image']").Attr("content"); exists && img != "" {
		product.ImageURL = img
	} else if img, exists := doc.Find("#landingImage").Attr("src"); exists { // Amazon
		product.ImageURL = img
	}

	// --- 3. PRECIO ---
	if price, exists := doc.Find("meta[property='product:price:amount']").Attr("content"); exists && price != "" {
		product.Price = price
	} else {
		// Probamos varios selectores de diferentes tiendas
		precioCrudo := doc.Find("#pdp-price-current-integer").Parent().Text() // PCComponentes

		if precioCrudo == "" {
			precioCrudo = doc.Find(".a-price .a-offscreen").First().Text() // Amazon (precio oculto limpio)
		}
		if precioCrudo == "" {
			precioCrudo = doc.Find(".a-price-whole").First().Text() // Amazon (solo números grandes)
		}

		// Extraemos solo el número con regex
		re := regexp.MustCompile(`[0-9]+(?:[.,][0-9]+)?`)
		precioLimpio := re.FindString(precioCrudo)

		if precioLimpio != "" {
			product.Price = precioLimpio
		}
	}

	// --- 4. DESCRIPCIÓN ---
	if desc, exists := doc.Find("meta[property='og:description']").Attr("content"); exists && desc != "" {
		product.Description = desc
	} else {
		product.Description = doc.Find("#feature-bullets").Text() // Amazon viñetas
	}

	product.Title = strings.TrimSpace(product.Title)
	product.Description = strings.TrimSpace(product.Description)

	return product, nil
}
