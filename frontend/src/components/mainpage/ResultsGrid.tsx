import React from "react";

// 1. Interfaz idéntica a vuestro modelo de Go
export interface Product {
  ProductID: number;
  Name: string;
  SourceURL: string;
  LastPrice: number;
  LowestPrice: number;
  CreatedBy: number;
  CreateAt: string; 
  UpdatedAt: string;
}

// 2. Definición de las propiedades que recibe el componente
type ResultsGridProps = {
  products?: Product[]; // El signo '?' evita errores si el padre no lo envía de inmediato
};

// 3. Asignamos '= []' por defecto para que nunca intente leer un undefined
export default function ResultsGrid({ products = [] }: ResultsGridProps) {
  return (
    <div>
      <h3 style={styles.title}>Results</h3>
      
      {/* Comprobamos de forma segura si el array está vacío */}
      {!products || products.length === 0 ? (
        <div style={styles.emptyContainer}>
          <p style={styles.emptyText}>No tienes ningún producto monitorizado todavía.</p>
          <span style={styles.emptySubtext}>¡Introduce un enlace arriba para empezar a rastrear precios!</span>
        </div>
      ) : (
        <div style={styles.grid}>
          {products.map((product) => {
            // Aseguramos que los precios sean números válidos antes de usar toFixed para evitar pantallas blancas
            const currentPrice = typeof product.LastPrice === 'number' ? product.LastPrice : 0;
            const minPrice = typeof product.LowestPrice === 'number' ? product.LowestPrice : 0;

            return (
              <div key={product.ProductID} style={styles.card}>
                <div>
                  <h4 style={styles.productName}>{product.Name || "Producto sin nombre"}</h4>
                  <span style={styles.productLink}>{product.SourceURL || "#"}</span>
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
      )}
    </div>
  );
}

// --- ESTILOS VISUALES ---
const styles: Record<string, React.CSSProperties> = {
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
  productName: {
    margin: "0 0 6px 0",
    fontSize: "18px",
    color: "#1f2937",
    fontWeight: "bold",
  },
  productLink: {
    display: "block",
    fontSize: "13px",
    color: "#2563eb",
    textDecoration: "none",
    marginBottom: "14px",
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
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
};