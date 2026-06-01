package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtract(t *testing.T) {
	// Definimos la estructura de nuestra "tabla" de pruebas
	type testCase struct {
		name           string
		mockHTML       string
		mockStatusCode int
		expectedTitle  string
		expectedPrice  string
		expectedImage  string
		expectedDesc   string
		expectError    bool
	}

	// Llenamos la tabla con los diferentes escenarios a probar
	tests := []testCase{
		{
			name: "Caso 1: Etiquetas estándar Open Graph (OG)",
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
			name: "Caso 2: Selectores alternativos (Estilo Amazon)",
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
			expectedPrice:  "199,50", // El scraper actualiza esto limpiando el HTML
			expectedImage:  "https://ejemplo.com/monitor.jpg",
			expectedDesc:   "Monitor con respuesta de 1ms.",
			expectError:    false,
		},
		{
			name:           "Caso 3: Error 404 del servidor",
			mockHTML:       `<html><body>No encontrado</body></html>`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
	}

	// Ejecutamos cada caso de la tabla
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Levantamos un servidor falso que devuelva nuestro HTML de prueba
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.mockStatusCode)
				fmt.Fprint(w, tc.mockHTML)
			}))
			// Aseguramos que el servidor falso se apague al terminar este test
			defer server.Close()

			// 2. Ejecutamos nuestra función original pasándole la URL del servidor falso
			product, err := Extract(server.URL)

			// 3. Comprobamos los resultados
			if tc.expectError {
				if err == nil {
					t.Fatalf("Se esperaba un error (código %d) pero no ocurrió", tc.mockStatusCode)
				}
				return // Si esperábamos error y dio error, el test pasa correctamente
			}

			if err != nil {
				t.Fatalf("No se esperaba error, pero ocurrió: %v", err)
			}

			// Validamos todos los campos del producto
			if product.Title != tc.expectedTitle {
				t.Errorf("Título incorrecto. Esperado: %q, Obtenido: %q", tc.expectedTitle, product.Title)
			}
			if product.Price != tc.expectedPrice {
				t.Errorf("Precio incorrecto. Esperado: %q, Obtenido: %q", tc.expectedPrice, product.Price)
			}
			if product.ImageURL != tc.expectedImage {
				t.Errorf("Imagen incorrecta. Esperado: %q, Obtenido: %q", tc.expectedImage, product.ImageURL)
			}
			if product.Description != tc.expectedDesc {
				t.Errorf("Descripción incorrecta. Esperado: %q, Obtenido: %q", tc.expectedDesc, product.Description)
			}
		})
	}
}
