package scraper

/*
Funciona en:
	Amazon
	PcComponentes
	Mediamarkt
	CoolMod
	Game
	Neobyte
	Pcbox
	Nike
	Worten
	La casa del libro
	Tradeinn
	Drunni
	Primor
	Kave Home
No Funciona en:
	(El html se genera una vez el usuario carga la pagina por lo que no existe un html de donde se pueda sacar el precio sin el uso de herramientas extenas
	como pupeteer que permited el uso de una navegador web invisible, el problema es el alto consumo de recusrsos que pide para hacer esto)
	Zara (o cualquier web del grupo inditex)
	El corte Ingles
	Mango
	(El codigo no logra acceder a la pagina por bloqueos y devulve un error 403 o una pagina en blanco)
	Inditex
	fnac
	adidas
	Decathlon
	(Webs en la que el precio que se ofrece depende por perfil)
	Aliexpress
	Temu
Puede fallar el codigo por bloqueos de paginas o conexiones (Gracias Tebas)
*/
import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ProductData contains the info about the product
type ProductData struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	ImageURL    string `json:"image_url"`
}

// Extract create or update a Product given a targetURL for the product of any of the soported sites
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

	// Title
	if title, exists := doc.Find("meta[property='og:title']").Attr("content"); exists && title != "" {
		product.Title = strings.TrimSpace(title)
	} else if title, exists := doc.Find("meta[name='og:title']").Attr("content"); exists && title != "" { // Worten
		product.Title = strings.TrimSpace(strings.ReplaceAll(title, " | Worten.es", ""))
	} else if title := doc.Find("#productTitle").Text(); title != "" { // Amazon
		product.Title = strings.TrimSpace(title)
	} else if title := doc.Find("h1").First().Text(); title != "" {
		product.Title = strings.TrimSpace(title)
	} else {
		product.Title = strings.TrimSpace(doc.Find("title").Text())
	}

	// Image
	if strings.Contains(targetURL, "worten") {
		if img, exists := doc.Find("meta[name='og:image']").Attr("content"); exists && img != "" {
			if strings.HasPrefix(img, "/") {
				product.ImageURL = "https://www.worten.es" + img
			} else {
				product.ImageURL = img
			}
		} else if img, exists := doc.Find("img.product-gallery__slider-image").Attr("src"); exists && img != "" {
			if strings.HasPrefix(img, "/") {
				product.ImageURL = "https://www.worten.es" + img
			} else {
				product.ImageURL = img
			}
		}
	} else if htmlCrudo, err := doc.Html(); err == nil { // Amazon
		reHiRes := regexp.MustCompile(`"hiRes"\s*:\s*"([^"]+)"`)
		if match := reHiRes.FindStringSubmatch(htmlCrudo); len(match) > 1 {
			product.ImageURL = match[1]
		} else {
			reOldHires := regexp.MustCompile(`data-old-hires="([^"]+)"`)
			if match2 := reOldHires.FindStringSubmatch(htmlCrudo); len(match2) > 1 {
				product.ImageURL = match2[1]
			} else if img, exists := doc.Find("#landingImage").Attr("src"); exists && img != "" {
				product.ImageURL = strings.TrimSpace(img)
			}
		}

		if product.ImageURL == "" {
			if img, exists := doc.Find("meta[property='og:image']").Attr("content"); exists {
				if strings.Contains(img, ".jpg") || strings.Contains(img, ".png") {
					product.ImageURL = strings.TrimSpace(img)
				}
			}
		}
	} else {
		if img, exists := doc.Find("meta[property='og:image']").Attr("content"); exists && img != "" {
			product.ImageURL = strings.TrimSpace(img)
		}
	}

	// Price
	if price, exists := doc.Find("meta[property='product:price:amount']").Attr("content"); exists && price != "" {
		product.Price = price
	} else {
		// Tags from different shops
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
		if precioCrudo == "" { // Nike
			precioCrudo = doc.Find("[data-testid='currentPrice-container']").First().Text()

			if precioCrudo != "" {
				precioCrudo = strings.ReplaceAll(precioCrudo, "€", "")
				precioCrudo = strings.ReplaceAll(precioCrudo, "\u00a0", "")
				precioCrudo = strings.ReplaceAll(precioCrudo, ",", ".")
				precioCrudo = strings.TrimSpace(precioCrudo)
			}
		}
		if precioCrudo == "" { // Worten
			euros := strings.TrimSpace(doc.Find(".price__numbers .value").First().Text())
			centimos := strings.TrimSpace(doc.Find(".price__numbers .decimal").First().Text())

			if euros != "" {
				euros = strings.ReplaceAll(euros, ".", "")
				if centimos != "" {
					precioCrudo = euros + "." + centimos
				} else {
					precioCrudo = euros
				}
			}
		}
		if precioCrudo == "" && strings.Contains(targetURL, "casadellibro") {
			doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
				textoScript := s.Text()

				rePrice := regexp.MustCompile(`(?i)"price"\s*:\s*"?([0-9.,]+)"?`)
				if matchPrice := rePrice.FindStringSubmatch(textoScript); len(matchPrice) > 1 {
					precioCrudo = matchPrice[1]
					precioCrudo = strings.ReplaceAll(precioCrudo, ",", ".")
				}
			})
		}
		if precioCrudo == "" { // Tradeinn
			if val, exists := doc.Find("#productFinalPrice").Attr("value"); exists && val != "" {
				precioCrudo = val
			}
		}
		if precioCrudo == "" { // Kave Home
			doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
				textoScript := s.Text()
				if strings.Contains(textoScript, `"price"`) && precioCrudo == "" {
					rePrice := regexp.MustCompile(`"price"\s*:\s*"?(\d+(?:[.,]\d+)?)"?`)
					if matchPrice := rePrice.FindStringSubmatch(textoScript); len(matchPrice) > 1 {
						precioCrudo = matchPrice[1]
					}
				}
			})
			if precioCrudo == "" {
				precioCrudo = doc.Find("main").Find("[class*='price'], [class*='Price']").First().Text()
				if precioCrudo != "" {
					precioCrudo = strings.ReplaceAll(precioCrudo, "€", "")
				}
			}
			if precioCrudo != "" {
				precioCrudo = strings.ReplaceAll(precioCrudo, ",", ".")
				precioCrudo = strings.TrimSpace(precioCrudo)
			}
		}

		// Using Regex for extracting the price
		re := regexp.MustCompile(`\d+(?:[.,]\d+)?`)
		precioLimpio := re.FindString(precioCrudo)

		if precioLimpio != "" {
			product.Price = precioLimpio
		}
	}

	// Description
	if strings.Contains(targetURL, "druni") {
		descDruni := doc.Find(".product.attribute.description .value").Text()
		if descDruni != "" {
			descDruni = strings.ReplaceAll(descDruni, "\t", "")

			reNewlines := regexp.MustCompile(`\n{3,}`)
			descDruni = reNewlines.ReplaceAllString(descDruni, "\n\n")

			product.Description = strings.TrimSpace(descDruni)
		}
	} else if desc, exists := doc.Find("meta[property='og:description']").Attr("content"); exists && desc != "" {
		product.Description = desc
	} else {
		product.Description = doc.Find("#feature-bullets").Text() // Amazon
	}

	product.Title = strings.TrimSpace(product.Title)
	product.Description = strings.TrimSpace(product.Description)

	return product, nil
}
