import React, { useState, useEffect } from "react";

// Same interface as go model
export interface Product {
  ProductID: number;
  Name: string;
  SourceURL: string;
  LastPrice: number;
  LowestPrice: number;
  CreatedBy: number;
  CreateAt: string;
  UpdatedAt: string;
  image_url: string;
}

// Property definition
type ResultsGridProps = {
  products?: Product[];
  isLoading?: boolean;
};

// ResultsGrid build the grid of products
export default function ResultsGrid({ products = [], isLoading = false }: ResultsGridProps) {
  // Pagination logic
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 4;

  useEffect(() => {
    setCurrentPage(1);
  }, [products]);

  const totalPages = Math.ceil(products.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentProducts = products.slice(startIndex, startIndex + itemsPerPage);

  const handlePrev = () => {
    if (currentPage > 1) setCurrentPage(currentPage - 1);
  };

  const handleNext = () => {
    if (currentPage < totalPages) setCurrentPage(currentPage + 1);
  };

  return (
    <div style={{ position: "relative" }}>
      <h3 style={styles.title}>Results</h3>

      {/* Loading layer */}
      {isLoading && (
        <div style={styles.loadingOverlay}>
          <svg width="50" height="50" viewBox="0 0 50 50" fill="none" stroke="#FACC15" strokeWidth="4">
            <circle cx="25" cy="25" r="20" strokeDasharray="31.4 31.4" strokeLinecap="round">
              <animateTransform
                attributeName="transform"
                type="rotate"
                repeatCount="indefinite"
                dur="1s"
                values="0 25 25;360 25 25"
              />
            </circle>
          </svg>
        </div>
      )}

      {/* Empty array check */}
      {!products || products.length === 0 ? (
        <div style={styles.emptyContainer}>
          <p style={styles.emptyText}>No tienes ningún producto monitorizado todavía.</p>
          <span style={styles.emptySubtext}>¡Introduce un enlace arriba para empezar a rastrear precios!</span>
        </div>
      ) : (
        <>
          <div style={styles.grid}>
            {currentProducts.map((product) => {
              // Check if prices are valid numbers
              const currentPrice = typeof product.LastPrice === 'number' ? product.LastPrice : 0;
              const minPrice = typeof product.LowestPrice === 'number' ? product.LowestPrice : 0;

              return (
                <div key={product.ProductID} style={styles.card}>
                  <div>
                    {/* --- Image --- */}
                    <div style={styles.imageContainer}>
                      {product.image_url ? (
                        <img src={product.image_url} alt={product.Name} style={styles.productImage} />
                      ) : (
                        <div style={styles.imagePlaceholder}>Sin imagen</div>
                      )}
                    </div>

                    <h4 style={styles.productName}>{product.Name || "Producto sin nombre"}</h4>

                    {/* --- URL --- */}
                    <a
                      href={product.SourceURL || "#"}
                      target="_blank"
                      rel="noopener noreferrer"
                      style={styles.productLink}
                    >
                      Enlace de compra
                    </a>

                    <div style={styles.priceTag}>
                      Precio actual: <strong style={styles.priceNumber}>{currentPrice.toFixed(2)}€</strong>
                    </div>
                    <div style={styles.lowestPriceTag}>
                      Mínimo registrado: {minPrice.toFixed(2)}€
                    </div>
                  </div>

                  <div style={styles.notificationBox}>
                    <label style={styles.notificationLabel}>Notify when price less than</label>
                    <div style={styles.inputGroup}>
                      <input type="text" placeholder="--" style={styles.alertInput} disabled />
                      <span style={styles.currencyAddon}>€</span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

          {/* --- Pagination controls --- */}
          {totalPages > 1 && (
            <div style={styles.paginationContainer}>
              <button
                onClick={handlePrev}
                disabled={currentPage === 1}
                style={{ ...styles.paginationButton, ...(currentPage === 1 ? styles.disabledButton : {}) }}
              >
                Anterior
              </button>

              <span style={styles.paginationText}>
                Página {currentPage} de {totalPages}
              </span>

              <button
                onClick={handleNext}
                disabled={currentPage === totalPages}
                style={{ ...styles.paginationButton, ...(currentPage === totalPages ? styles.disabledButton : {}) }}
              >
                Siguiente
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}

// --- Styles ---
const styles: Record<string, React.CSSProperties> = {
  loadingOverlay: {
    position: "fixed",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: "rgba(255, 255, 255, 0.7)",
    zIndex: 10,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    borderRadius: "12px",
  },
  title: {
    fontSize: "22px",
    color: "#1f2937",
    marginBottom: "24px",
    fontWeight: "bold",
    borderLeft: "5px solid #FACC15",
    paddingLeft: "12px",
    lineHeight: "1",
  },
  grid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))",
    gap: "24px",
  },
  card: {
    background: "white",
    borderRadius: "12px",
    padding: "24px",
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-between",
    boxShadow: "0 4px 15px rgba(0,0,0,0.02)",
    border: "1px solid #eee",
    borderTop: "4px solid #FACC15",
    minHeight: "230px",
  },
  imageContainer: {
    width: "100%",
    height: "160px",
    marginBottom: "16px",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#ffffff",
    borderRadius: "8px",
    overflow: "hidden",
  },
  productImage: {
    width: "100%",
    height: "100%",
    objectFit: "contain",
  },
  imagePlaceholder: {
    color: "#9ca3af",
    fontSize: "13px",
  },
  productName: {
    margin: "0 0 6px 0",
    fontSize: "18px",
    color: "#1f2937",
    fontWeight: "bold",
  },
  productLink: {
    display: "inline-block",
    fontSize: "13px",
    color: "#2563eb",
    textDecoration: "underline",
    marginBottom: "14px",
  },
  priceTag: {
    fontSize: "15px",
    color: "#374151",
    marginBottom: "4px",
  },
  priceNumber: {
    color: "#16a34a",
    fontSize: "16px",
  },
  lowestPriceTag: {
    fontSize: "12px",
    color: "#6b7280",
    marginBottom: "15px",
  },
  notificationBox: {
    borderTop: "1px dashed #e5e7eb",
    paddingTop: "14px",
  },
  notificationLabel: {
    display: "block",
    fontSize: "13px",
    color: "#6b7280",
    marginBottom: "6px",
  },
  inputGroup: {
    display: "flex",
    width: "110px",
  },
  alertInput: {
    width: "100%",
    padding: "8px",
    border: "1px solid #ddd",
    borderRadius: "6px 0 0 6px",
    fontSize: "14px",
    textAlign: "center",
    backgroundColor: "#fafafa",
  },
  currencyAddon: {
    backgroundColor: "#f3f4f6",
    border: "1px solid #ddd",
    borderLeft: "none",
    padding: "8px 12px",
    fontSize: "14px",
    color: "#4b5563",
    borderRadius: "0 6px 6px 0",
    fontWeight: "bold",
  },
  emptyContainer: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    padding: "60px 20px",
    background: "white",
    borderRadius: "12px",
    border: "1px dashed #ccc",
    textAlign: "center",
  },
  emptyText: {
    fontSize: "16px",
    fontWeight: "bold",
    color: "#4b5563",
    margin: "0 0 8px 0",
  },
  emptySubtext: {
    fontSize: "14px",
    color: "#9ca3af",
  },
  paginationContainer: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    marginTop: "40px",
    gap: "24px",
  },
  paginationButton: {
    padding: "10px 20px",
    backgroundColor: "#FACC15",
    color: "black",
    border: "none",
    borderRadius: "8px",
    fontWeight: "bold",
    cursor: "pointer",
    fontSize: "14px",
    boxShadow: "0 4px 15px rgba(250, 204, 21, 0.2)",
  },
  disabledButton: {
    backgroundColor: "#e5e7eb",
    color: "#9ca3af",
    cursor: "not-allowed",
    boxShadow: "none",
  },
  paginationText: {
    fontSize: "14px",
    fontWeight: "bold",
    color: "#374151",
  },
};