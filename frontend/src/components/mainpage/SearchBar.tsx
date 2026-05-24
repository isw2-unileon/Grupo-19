import React from "react";

export default function SearchBar() {
  return (
    <div style={styles.container}>
      {/* --- Contenedor de la barra de búsqueda con la lupa --- */}
      <div style={styles.searchWrapper}>
        {/* Icono de Lupa (SVG Nativo) */}
        <div style={styles.iconContainer}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={2.5}
            stroke="#9ca3af"
            style={styles.magnifier}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              pathLength="1"
              d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
            />
          </svg>
        </div>
        
        {/* Input de texto */}
        <input
          type="text"
          placeholder="Busca un producto..."
          style={styles.input}
          disabled
        />
      </div>

      {/* --- Botón Añadir (Ahora separado a la derecha) --- */}
      <button type="button" style={styles.addButton}>
        Añadir
      </button>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: "flex",
    width: "100%",
    maxWidth: "750px", // Ampliado un poco para dar más espacio a la separación
    gap: "16px",       // Este hueco separa visualmente la barra del botón
    alignItems: "center",
  },
  searchWrapper: {
    display: "flex",
    flex: 1,
    boxShadow: "0 4px 20px rgba(0,0,0,0.05)",
    borderRadius: "8px",
    border: "1px solid #ddd",
    backgroundColor: "white",
    overflow: "hidden",
    alignItems: "center",
  },
  iconContainer: {
    paddingLeft: "16px",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  magnifier: {
    width: "20px",
    height: "20px",
  },
  input: {
    flex: 1,
    padding: "16px 16px 16px 12px", // Ajustado el padding izquierdo para que no pegue al icono
    fontSize: "15px",
    border: "none",
    backgroundColor: "transparent",
    outline: "none",
  },
  addButton: {
    padding: "16px 32px", // Ajustado el padding para que tenga la misma altura que la barra
    border: "none",
    background: "#FACC15",
    color: "black",
    fontWeight: "bold",
    fontSize: "15px",
    borderRadius: "8px",  // Ahora lleva bordes redondeados propios al estar suelto
    cursor: "pointer",
    boxShadow: "0 4px 15px rgba(250, 204, 21, 0.2)", // Sutil sombra amarilla para que destaque
  },
};