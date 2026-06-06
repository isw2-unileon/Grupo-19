package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtract(t *testing.T) {
	// Define the structure of our test cases table
	type testCase struct {
		name           string
		urlSuffix      string // Allows us to trigger URL-specific logic (like "worten" or "druni")
		mockHTML       string
		mockStatusCode int
		expectedTitle  string
		expectedPrice  string
		expectedImage  string
		expectedDesc   string
		expectError    bool
	}

	// Fill the table with different scenarios targeting specific shop quirks
	tests := []testCase{
		{
			name:      "Case 1: Standard Open Graph (OG) tags",
			urlSuffix: "/standard-product",
			mockHTML: `
				<html>
					<head>
						<meta property="og:title" content="Zapatillas Deportivas">
						<meta property="og:image" content="https://ejemplo.com/zapatillas.jpg">
						<meta property="product:price:amount" content="59.99">
						<meta property="og:description" content="Zapatillas para correr muy cómodas.">
					</head>
				</html>
			`,
			mockStatusCode: http.StatusOK,
			expectedTitle:  "Zapatillas Deportivas",
			expectedPrice:  "59.99",
			expectedImage:  "https://ejemplo.com/zapatillas.jpg",
			expectedDesc:   "Zapatillas para correr muy cómodas.",
			expectError:    false,
		},
		{
			name:      "Case 2: Alternative Selectors (Amazon Style)",
			urlSuffix: "/dp/B08XYZ123",
			mockHTML: `
				<html>
					<body>
						<h1 id="productTitle">  Monitor Gaming 144Hz  </h1>
						<img id="landingImage" src="https://ejemplo.com/monitor.jpg">
						<span class="a-price-whole">199,50</span>
						<div id="feature-bullets">Monitor con respuesta de 1ms.</div>
					</body>
				</html>
			`,
			mockStatusCode: http.StatusOK,
			expectedTitle:  "Monitor Gaming 144Hz",
			expectedPrice:  "199,50",
			expectedImage:  "https://ejemplo.com/monitor.jpg",
			expectedDesc:   "Monitor con respuesta de 1ms.",
			expectError:    false,
		},
		{
			name:           "Case 3: Server Error 404",
			urlSuffix:      "/not-found",
			mockHTML:       `<html><body>No encontrado</body></html>`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "Case 4: Coolmod specific logic (Separated int and decimal)",
			urlSuffix: "/coolmod-item",
			mockHTML: `
				<html>
					<body>
						<h1>Tarjeta Gráfica RTX</h1>
						<div class="product_price int_price">399</div>
						<div class="dec_price">95</div>
					</body>
				</html>
			`,
			mockStatusCode: http.StatusOK,
			expectedTitle:  "Tarjeta Gráfica RTX",
			expectedPrice:  "399.95",
			expectedImage:  "",
			expectedDesc:   "",
			expectError:    false,
		},
		{
			name:      "Case 5: Game specific logic (Quotes in decimals)",
			urlSuffix: "/game-item",
			mockHTML: `
				<html>
					<body>
						<h1>Mando Inalámbrico</h1>
						<div class="buy--price">
							<span class="int">59</span>
							<span class="decimal">'99</span>
						</div>
					</body>
				</html>
			`,
			mockStatusCode: http.StatusOK,
			expectedTitle:  "Mando Inalámbrico",
			expectedPrice:  "59.99", // Scraper must remove the ' character
			expectedImage:  "",
			expectedDesc:   "",
			expectError:    false,
		},
		{
			name:      "Case 6: Worten specific logic (URL detection & absolute image formatting)",
			urlSuffix: "/worten/frigorifico", // Triggers strings.Contains(targetURL, "worten")
			mockHTML: `
				<html>
					<head>
						<meta name="og:title" content="Frigorífico Combi | Worten.es">
						<meta name="og:image" content="/images/frigo.jpg">
					</head>
					<body>
						<div class="price__numbers">
							<span class="value">499.</span>
							<span class="decimal">00</span>
						</div>
					</body>
				</html>
			`,
			mockStatusCode: http.StatusOK,
			expectedTitle:  "Frigorífico Combi", // Must clean " | Worten.es"
			expectedPrice:  "499.00",
			expectedImage:  "https://www.worten.es/images/frigo.jpg", // Must prepend domain
			expectedDesc:   "",
			expectError:    false,
		},
	}

	// Execute all table cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Start a mock server to return our test HTML
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.mockStatusCode)
				fmt.Fprint(w, tc.mockHTML)
			}))
			// Ensure the mock server is closed when the test finishes
			defer server.Close()

			// Build the target URL combining the mock server IP with our specific suffix
			targetURL := server.URL + tc.urlSuffix

			// 2. Execute the extract function against our mock server
			product, err := Extract(targetURL)

			// 3. Evaluate results
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected an error (status %d) but none occurred", tc.mockStatusCode)
				}
				return // Test passes if we expected an error and got one
			}

			if err != nil {
				t.Fatalf("Unexpected error occurred: %v", err)
			}

			// Validate all product fields
			if product.Title != tc.expectedTitle {
				t.Errorf("Title mismatch. Expected: %q, Got: %q", tc.expectedTitle, product.Title)
			}
			if product.Price != tc.expectedPrice {
				t.Errorf("Price mismatch. Expected: %q, Got: %q", tc.expectedPrice, product.Price)
			}
			if product.ImageURL != tc.expectedImage {
				t.Errorf("Image URL mismatch. Expected: %q, Got: %q", tc.expectedImage, product.ImageURL)
			}
			if product.Description != tc.expectedDesc {
				t.Errorf("Description mismatch. Expected: %q, Got: %q", tc.expectedDesc, product.Description)
			}
		})
	}
}
