import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

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
  const navigate = useNavigate();

  // Pagination logic
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

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

  // Price tracking logic
  const [targetPrices, setTargetPrices] = useState<Record<number, string>>({});

  const handlePriceChange = (productID: number, value: string) => {
    setTargetPrices((prev) => ({ ...prev, [productID]: value }));
  };

  const irAFichaProducto = (productID: number) => {
    navigate(`/product/${productID}`);
  };

  const guardarTracking = async (productID: number) => {
    const precioDeseado = targetPrices[productID];

    if (!precioDeseado) {
      setModalState({
        isOpen: true,
        title: "Campo vacío",
        message: "Por favor, introduce un precio antes de guardar."
      });
      return;
    }
    const precioFloat = parseFloat(precioDeseado);

    // Check if valid number
    if (isNaN(precioFloat) || precioFloat < 0) {
      setModalState({
        isOpen: true,
        title: "Precio inválido",
        message: "Por favor, introduce un precio válido que no sea negativo."
      });
      setTargetPrices((prev) => ({ ...prev, [productID]: "" }));
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/api/tracking", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          product_id: productID,
          target_price: precioFloat
        }),
      });

      if (response.ok) {
        setModalState({
          isOpen: true,
          title: "¡Guardado!",
          message: "La alerta de precio se ha configurado con éxito."
        });
      } else {
        const data = await response.json();
        setModalState({
          isOpen: true,
          title: "Error",
          message: data.error || "Ocurrió un problema al guardar la alerta."
        });
      }
    } catch (error) {
      console.error("Error al guardar el tracking:", error);
      setModalState({
        isOpen: true,
        title: "Error de conexión",
        message: "No se pudo conectar con el servidor."
      });
    } finally {
      setTargetPrices((prev) => ({ ...prev, [productID]: "" }));
    }
  };

  const [modalState, setModalState] = useState({
    isOpen: false,
    title: "",
    message: ""
  });

  const closeModal = () => setModalState({ ...modalState, isOpen: false });

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
          {/* NUEVO: Texto cambiado */}
          <p style={styles.emptyText}>No se ha realizado ninguna búsqueda todavía.</p>
          <span style={styles.emptySubtext}>¡Introduce un enlace o nombre de un producto para empezar a rastrear precios!</span>
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
                    <button
                      style={styles.viewProductBtn}
                      onClick={() => irAFichaProducto(product.ProductID)}
                    >
                      Ver producto
                    </button>

                    <div style={styles.alertWrapper}>
                      <label style={styles.notificationLabel}>Establecer alerta de precio:</label>
                      <div style={styles.trackingGroup}>
                        <input
                          type="number"
                          placeholder="Precio objetivo..."
                          style={styles.alertInput}
                          value={targetPrices[product.ProductID] || ""}
                          onChange={(e) => handlePriceChange(product.ProductID, e.target.value)}
                        />
                        <span style={styles.currencyAddon}>€</span>
                        <button
                          style={styles.saveAlertBtn}
                          onClick={() => guardarTracking(product.ProductID)}
                        >
                          Guardar
                        </button>
                      </div>
                    </div>
                  </div>
                  {modalState.isOpen && (
                    <div style={styles.modalOverlay}>
                      <div style={styles.modalContent}>
                        <h3 style={{ marginTop: 0, color: "#1f2937" }}>{modalState.title}</h3>
                        <p style={{ color: "#4b5563", marginBottom: "24px" }}>
                          {modalState.message}
                        </p>
                        <div style={styles.modalActions}>
                          <button onClick={closeModal} style={styles.confirmButton}>
                            Aceptar
                          </button>
                        </div>
                      </div>
                    </div>
                  )}
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
  notificationBox: {
    borderTop: "1px dashed #e5e7eb",
    paddingTop: "16px",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-end",
    gap: "12px",
  },
  viewProductBtn: {
    padding: "0 14px",
    height: "36px",
    backgroundColor: "#FACC15",
    color: "black",
    border: "1px solid #e5e7eb",
    borderRadius: "8px",
    fontSize: "13px",
    fontWeight: "bold",
    cursor: "pointer",
    transition: "background-color 0.2s",
    whiteSpace: "nowrap",
  },
  alertWrapper: {
    display: "flex",
    flexDirection: "column",
    alignItems: "flex-start",
    flex: 1,
    maxWidth: "250px",
  },
  notificationLabel: {
    fontSize: "11px",
    color: "#6b7280",
    marginBottom: "6px",
  },
  trackingGroup: {
    display: "flex",
    width: "100%",
    height: "36px",
  },
  alertInput: {
    flex: 1,
    padding: "0 8px",
    border: "1px solid #ddd",
    borderRight: "none",
    borderRadius: "6px 0 0 6px",
    fontSize: "13px",
    backgroundColor: "#fafafa",
    outline: "none",
    width: "100%",
  },
  currencyAddon: {
    backgroundColor: "#fafafa",
    border: "1px solid #ddd",
    borderLeft: "none",
    borderRight: "none",
    padding: "0 10px",
    fontSize: "13px",
    color: "#4b5563",
    display: "flex",
    alignItems: "center",
  },
  saveAlertBtn: {
    padding: "0 12px",
    backgroundColor: "#FACC15",
    color: "black",
    border: "none",
    borderRadius: "0 6px 6px 0",
    fontSize: "13px",
    fontWeight: "bold",
    cursor: "pointer",
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
  emptyContainer: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    padding: "60px 20px",
    minHeight: "400px",
    background: "white",
    borderRadius: "12px",
    border: "2px dashed #e5e7eb",
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
  modalOverlay: {
    position: "fixed",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: "rgba(0, 0, 0, 0.1)",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    zIndex: 9999,
  },
  modalContent: {
    backgroundColor: "#ffffff",
    padding: "24px",
    borderRadius: "8px",
    boxShadow: "0 4px 6px rgba(0, 0, 0, 0.1)",
    maxWidth: "400px",
    width: "90%",
    textAlign: "center",
  },
  modalActions: {
    display: "flex",
    justifyContent: "center",
  },
  confirmButton: {
    padding: "10px 28px",
    border: "none",
    backgroundColor: "#FACC15",
    color: "#1f2937",
    borderRadius: "6px",
    cursor: "pointer",
    fontWeight: "bold",
  },
};