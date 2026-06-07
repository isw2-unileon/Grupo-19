package scraper

/*
Works on:
	amazon
	PcComponents
	Mediamarkt
	CoolMod
	Games
	Neobyte
	pcbox
	Nike
	Worten
	The book house
	Tradeinn
	Drunni
	Primor
	Kave Home
Does not work on:
	(The html is generated once the user loads the page, so there is no html from which the price can be obtained without the use of external tools.
	like a pupeteer that allows the use of an invisible web browser, the problem is the high consumption of resources it requires to do this)
	Zara (or any Inditex group website)
	The English court
	Handle
	(The code cannot access the page due to blockages and returns a 403 error or a blank page)
	Inditex
	fnac
	adidas
	Decathlon
	(Websites in which the price offered depends on the profile)
	AliExpress
	Temu
The code may fail due to page or connection blocking (Thanks Tebas)
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

	// Advanced Headers camouflage to try to skip the 403
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
		return nil, fmt.Errorf("server error: %d", res.StatusCode)
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
	} else if rawHTML, err := doc.Html(); err == nil { // Amazon
		reHiRes := regexp.MustCompile(`"hiRes"\s*:\s*"([^"]+)"`)
		if match := reHiRes.FindStringSubmatch(rawHTML); len(match) > 1 {
			product.ImageURL = match[1]
		} else {
			reOldHires := regexp.MustCompile(`data-old-hires="([^"]+)"`)
			if match2 := reOldHires.FindStringSubmatch(rawHTML); len(match2) > 1 {
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
		rawPrice := doc.Find("#pdp-price-current-integer").Parent().Text() // PCComponentes

		if rawPrice == "" { // Amazon
			rawPrice = doc.Find(".a-price .a-offscreen").First().Text()
		}
		if rawPrice == "" {
			rawPrice = doc.Find(".a-price-whole").First().Text()
		}
		if rawPrice == "" { // Mediamarkt
			rawPrice = doc.Find("[data-test='branded-price-whole-value']").Parent().Text()
		}
		if rawPrice == "" { // Coolmod
			euros := strings.TrimSpace(doc.Find(".product_price.int_price").First().Text())
			cents := strings.TrimSpace(doc.Find(".dec_price").Text())

			if euros != "" {
				if cents != "" {
					rawPrice = euros + "." + cents
				} else {
					rawPrice = euros
				}
			}
		}
		if rawPrice == "" { // Game
			euros := strings.TrimSpace(doc.Find(".buy--price .int").First().Text())
			cents := strings.TrimSpace(doc.Find(".buy--price .decimal").First().Text())

			if euros != "" {
				cents = strings.ReplaceAll(cents, "'", "")

				if cents != "" {
					rawPrice = euros + "." + cents
				} else {
					rawPrice = euros
				}
			}
		}
		if rawPrice == "" { // Nike
			rawPrice = doc.Find("[data-testid='currentPrice-container']").First().Text()

			if rawPrice != "" {
				rawPrice = strings.ReplaceAll(rawPrice, "€", "")
				rawPrice = strings.ReplaceAll(rawPrice, "\u00a0", "")
				rawPrice = strings.ReplaceAll(rawPrice, ",", ".")
				rawPrice = strings.TrimSpace(rawPrice)
			}
		}
		if rawPrice == "" { // Worten
			euros := strings.TrimSpace(doc.Find(".price__numbers .value").First().Text())
			cents := strings.TrimSpace(doc.Find(".price__numbers .decimal").First().Text())

			if euros != "" {
				euros = strings.ReplaceAll(euros, ".", "")
				if cents != "" {
					rawPrice = euros + "." + cents
				} else {
					rawPrice = euros
				}
			}
		}
		if rawPrice == "" && strings.Contains(targetURL, "casadellibro") {
			doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
				textoScript := s.Text()

				rePrice := regexp.MustCompile(`(?i)"price"\s*:\s*"?([0-9.,]+)"?`)
				if matchPrice := rePrice.FindStringSubmatch(textoScript); len(matchPrice) > 1 {
					rawPrice = matchPrice[1]
					rawPrice = strings.ReplaceAll(rawPrice, ",", ".")
				}
			})
		}
		if rawPrice == "" { // Tradeinn
			if val, exists := doc.Find("#productFinalPrice").Attr("value"); exists && val != "" {
				rawPrice = val
			}
		}
		if rawPrice == "" { // Kave Home
			doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
				textoScript := s.Text()
				if strings.Contains(textoScript, `"price"`) && rawPrice == "" {
					rePrice := regexp.MustCompile(`"price"\s*:\s*"?(\d+(?:[.,]\d+)?)"?`)
					if matchPrice := rePrice.FindStringSubmatch(textoScript); len(matchPrice) > 1 {
						rawPrice = matchPrice[1]
					}
				}
			})
			if rawPrice == "" {
				rawPrice = doc.Find("main").Find("[class*='price'], [class*='Price']").First().Text()
				if rawPrice != "" {
					rawPrice = strings.ReplaceAll(rawPrice, "€", "")
				}
			}
			if rawPrice != "" {
				rawPrice = strings.ReplaceAll(rawPrice, ",", ".")
				rawPrice = strings.TrimSpace(rawPrice)
			}
		}

		// Using Regex for extracting the price
		re := regexp.MustCompile(`\d+(?:[.,]\d+)?`)
		cleanPrice := re.FindString(rawPrice)

		if cleanPrice != "" {
			product.Price = cleanPrice
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
