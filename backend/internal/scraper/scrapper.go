package scraper

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
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
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	// Camuflaje avanzado de Headers para intentar saltar el 403
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")

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

	// Titulo
	if title, exists := doc.Find("meta[property='og:title']").Attr("content"); exists && title != "" {
		product.Title = title
	} else if title := doc.Find("#productTitle").Text(); title != "" { // Amazon
		product.Title = title
	}

	// Imagen
	if img, exists := doc.Find("meta[property='og:image']").Attr("content"); exists && img != "" {
		product.ImageURL = img
	} else if img, exists := doc.Find("#landingImage").Attr("src"); exists { // Amazon
		product.ImageURL = img
	}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		textoScript := s.Text()

		// Si el script contiene los datos analíticos de Zara
		if strings.Contains(textoScript, "zara.analyticsData") {

			// Cazamos el título usando Regex
			reTitle := regexp.MustCompile(`"productName"\s*:\s*"([^"]+)"`)
			if matchTitle := reTitle.FindStringSubmatch(textoScript); len(matchTitle) > 1 {
				product.Title = matchTitle[1]
			}

			// Cazamos el precio usando Regex
			rePrice := regexp.MustCompile(`"mainPrice"\s*:\s*([0-9.]+)`)
			if matchPrice := rePrice.FindStringSubmatch(textoScript); len(matchPrice) > 1 {
				product.Price = matchPrice[1]
			}
		}
	})

	// Precio
	if price, exists := doc.Find("meta[property='product:price:amount']").Attr("content"); exists && price != "" {
		product.Price = price
	} else {
		// Probamos varios selectores de diferentes tiendas
		precioCrudo := doc.Find("#pdp-price-current-integer").Parent().Text() // PCComponentes

		if precioCrudo == "" { // Amazon
			precioCrudo = doc.Find(".a-price .a-offscreen").First().Text()
		}
		if precioCrudo == "" {
			precioCrudo = doc.Find(".a-price-whole").First().Text()
		}
		if precioCrudo == "" { // Mediamarkt
			precioCrudo = doc.Find("[data-test='branded-price-whole-value']").Parent().Text()
		}
		if precioCrudo == "" { // Coolmod
			euros := strings.TrimSpace(doc.Find(".product_price.int_price").First().Text())
			centimos := strings.TrimSpace(doc.Find(".dec_price").Text())

			if euros != "" {
				if centimos != "" {
					precioCrudo = euros + "." + centimos
				} else {
					precioCrudo = euros
				}
			}
		}
		if precioCrudo == "" { // Game
			euros := strings.TrimSpace(doc.Find(".buy--price .int").First().Text())
			centimos := strings.TrimSpace(doc.Find(".buy--price .decimal").First().Text())

			if euros != "" {
				centimos = strings.ReplaceAll(centimos, "'", "")

				if centimos != "" {
					precioCrudo = euros + "." + centimos
				} else {
					precioCrudo = euros
				}
			}
		}

		// Extraemos solo el número con regex
		re := regexp.MustCompile(`[0-9]+(?:[.,][0-9]+)?`)
		precioLimpio := re.FindString(precioCrudo)

		if precioLimpio != "" {
			product.Price = precioLimpio
		}
	}

	// Descripcion (Sin Uso)
	if desc, exists := doc.Find("meta[property='og:description']").Attr("content"); exists && desc != "" {
		product.Description = desc
	} else {
		product.Description = doc.Find("#feature-bullets").Text() // Amazon
	}

	product.Title = strings.TrimSpace(product.Title)
	product.Description = strings.TrimSpace(product.Description)

	return product, nil
}
